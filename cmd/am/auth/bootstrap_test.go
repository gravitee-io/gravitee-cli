// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// fakeJWT builds a JWT whose payload carries the given sub claim. The
// signature is a placeholder — bootstrap never verifies it.
func fakeJWT(sub string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"typ":"JWT","alg":"none"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"` + sub + `"}`))
	return header + "." + payload + ".sig"
}

// fakeAM returns a stub AM server implementing the two endpoints the
// bootstrap flow touches: login (issues a Bearer JWT cookie) and tokens
// (mints a PAT, auth via Authorization: Bearer).
func fakeAM(t *testing.T) *httptest.Server {
	t.Helper()

	sessionValue := "Bearer " + fakeJWT("user-1")

	mux := http.NewServeMux()

	mux.HandleFunc("/management/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad form", http.StatusBadRequest)
			return
		}
		if r.PostForm.Get("username") != "admin" || r.PostForm.Get("password") != "adminadmin" {
			http.Error(w, "bad creds", http.StatusUnauthorized)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: sessionCookieName, Value: sessionValue})
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/management/organizations/DEFAULT/users/user-1/tokens", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != sessionValue {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"tokenId":"tok-1","token":"gioat_abc123"}`))
	})

	return httptest.NewServer(mux)
}

func newTestFactory(t *testing.T) (*factory.Factory, *bytes.Buffer, string) {
	t.Helper()

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")

	out := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	in := strings.NewReader("")

	f := &factory.Factory{
		Config:     &config.Config{Contexts: map[string]*config.Context{}},
		ConfigPath: cfgPath,
		IOStreams:  factory.IOStreams{Out: out, Err: errBuf, In: in},
	}
	return f, out, cfgPath
}

func TestBootstrapMintsToken(t *testing.T) {
	srv := fakeAM(t)
	defer srv.Close()

	f, out, _ := newTestFactory(t)

	opts := &bootstrapOptions{
		factory:   f,
		amURL:     srv.URL,
		username:  "admin",
		password:  "adminadmin",
		tokenName: "gio-cli",
		org:       "DEFAULT",
	}

	if err := opts.run(srv.Client()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "gioat_abc123") {
		t.Errorf("expected token value in output, got: %s", out.String())
	}
}

func TestBootstrapSavesConfig(t *testing.T) {
	srv := fakeAM(t)
	defer srv.Close()

	f, _, cfgPath := newTestFactory(t)

	opts := &bootstrapOptions{
		factory:     f,
		amURL:       srv.URL,
		username:    "admin",
		password:    "adminadmin",
		tokenName:   "gio-cli",
		org:         "DEFAULT",
		save:        true,
		contextName: "local",
	}

	if err := opts.run(srv.Client()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("expected config file: %v", err)
	}
	body := string(data)
	if !strings.Contains(body, "gioat_abc123") {
		t.Errorf("expected token persisted, got: %s", body)
	}
	if !strings.Contains(body, "local:") {
		t.Errorf("expected context 'local' in config, got: %s", body)
	}
}

func TestBootstrapBadCredentials(t *testing.T) {
	srv := fakeAM(t)
	defer srv.Close()

	f, _, _ := newTestFactory(t)

	opts := &bootstrapOptions{
		factory:   f,
		amURL:     srv.URL,
		username:  "admin",
		password:  "wrong",
		tokenName: "gio-cli",
		org:       "DEFAULT",
	}

	err := opts.run(srv.Client())
	if err == nil || !strings.Contains(err.Error(), "login failed") {
		t.Errorf("expected login failure, got: %v", err)
	}
}

func TestBootstrapURLRequired(t *testing.T) {
	f, _, _ := newTestFactory(t)
	opts := &bootstrapOptions{
		factory:  f,
		username: "admin",
		password: "x",
		org:      "DEFAULT",
	}

	err := opts.run(http.DefaultClient)
	if err == nil || !strings.Contains(err.Error(), "no AM URL") {
		t.Errorf("expected URL required error, got: %v", err)
	}
}

func TestBootstrapFallsBackToResolvedURL(t *testing.T) {
	srv := fakeAM(t)
	defer srv.Close()

	f, out, _ := newTestFactory(t)
	f.Resolved = &config.ResolvedContext{URL: srv.URL}

	opts := &bootstrapOptions{
		factory:   f,
		username:  "admin",
		password:  "adminadmin",
		tokenName: "gio-cli",
		org:       "DEFAULT",
	}

	if err := opts.run(srv.Client()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "gioat_abc123") {
		t.Errorf("expected token in output, got: %s", out.String())
	}
}
