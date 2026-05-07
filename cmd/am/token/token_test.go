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

package token

import (
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func TestTokenList(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/users/user-1/tokens") {
				t.Errorf("unexpected path: %s", path)
			}
			return []byte(`[{"id":"token-1","token":"abc"}]`), nil
		},
	}
	f, out := newTestFactory(fake, false)
	cmd := newListCmd(f)
	cmd.SetArgs([]string{"--user", "user-1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "token-1") {
		t.Errorf("expected 'token-1' in output, got: %s", out.String())
	}
}

func TestTokenCreate(t *testing.T) {
	fake := &client.FakeClient{
		PostFunc: func(path string, body interface{}) ([]byte, error) {
			if !strings.Contains(path, "/users/user-1/tokens") {
				t.Errorf("unexpected path: %s", path)
			}
			return []byte(`{"id":"token-new","token":"xyz"}`), nil
		},
	}
	f, out := newTestFactory(fake, false)
	cmd := newCreateCmd(f)
	cmd.SetArgs([]string{"--user", "user-1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "token-new") {
		t.Errorf("expected 'token-new' in output, got: %s", out.String())
	}
}

func TestTokenRevoke(t *testing.T) {
	revoked := false
	fake := &client.FakeClient{
		DeleteFunc: func(path string) error {
			if !strings.Contains(path, "/users/user-1/tokens/token-1") {
				t.Errorf("unexpected path: %s", path)
			}
			revoked = true
			return nil
		},
	}
	f, _ := newTestFactory(fake, false)
	cmd := newRevokeCmd(f)
	cmd.SetArgs([]string{"token-1", "--user", "user-1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !revoked {
		t.Error("expected Delete to be called")
	}
}
