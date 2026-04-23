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

package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestV2Path(t *testing.T) {
	tests := []struct {
		name  string
		envID string
		path  string
		want  string
	}{
		{
			name:  "basic path",
			envID: "DEFAULT",
			path:  "apis",
			want:  "/management/v2/environments/DEFAULT/apis",
		},
		{
			name:  "nested path with leading slash",
			envID: "production",
			path:  "/apis/123/plans",
			want:  "/management/v2/environments/production/apis/123/plans",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := V2Path(tt.envID, tt.path)
			if got != tt.want {
				t.Errorf("V2Path(%q, %q) = %q, want %q", tt.envID, tt.path, got, tt.want)
			}
		})
	}
}

func TestV1Path(t *testing.T) {
	tests := []struct {
		name  string
		orgID string
		envID string
		path  string
		want  string
	}{
		{
			name:  "basic path",
			orgID: "DEFAULT",
			envID: "DEFAULT",
			path:  "applications",
			want:  "/management/organizations/DEFAULT/environments/DEFAULT/applications",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := V1Path(tt.orgID, tt.envID, tt.path)
			if got != tt.want {
				t.Errorf("V1Path(%q, %q, %q) = %q, want %q", tt.orgID, tt.envID, tt.path, got, tt.want)
			}
		})
	}
}

func TestHTTPClientGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			t.Error("missing Bearer token")
		}

		w.WriteHeader(http.StatusOK)

		resp, _ := json.Marshal(map[string]string{"id": "123"})
		_, _ = w.Write(resp)
	}))
	defer server.Close()

	c := NewHTTPClient(HTTPClientConfig{BaseURL: server.URL, Token: "test-token"})

	data, err := c.Get("/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(string(data), `"id"`) {
		t.Errorf("unexpected response: %s", string(data))
	}
}

func TestHTTPClientErrorMapping(t *testing.T) {
	tests := []struct {
		name       string
		wantErr    string
		statusCode int
	}{
		{name: "401", statusCode: 401, wantErr: "authentication failed"},
		{name: "403", statusCode: 403, wantErr: "access denied"},
		{name: "404", statusCode: 404, wantErr: "resource not found"},
		{name: "500", statusCode: 500, wantErr: "server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			c := NewHTTPClient(HTTPClientConfig{BaseURL: server.URL, Token: "tok"})

			_, err := c.Get("/test")
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestHTTPClientDebugMasksToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	var debugBuf bytes.Buffer

	c := NewHTTPClient(HTTPClientConfig{
		BaseURL:  server.URL,
		Token:    "gioat_secrettoken",
		Debug:    true,
		DebugOut: &debugBuf,
	})

	_, err := c.Get("/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := debugBuf.String()
	if strings.Contains(output, "gioat_secrettoken") {
		t.Error("debug output should not contain the full token")
	}

	if !strings.Contains(output, "ken") {
		t.Error("debug output should contain the last 3 chars of the token")
	}
}

func TestMapHTTPError(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantMsg string
		status  int
	}{
		{name: "400", status: 400, body: "bad field", wantMsg: "invalid request"},
		{name: "409", status: 409, body: "already exists", wantMsg: "conflict"},
		{name: "502", status: 502, wantMsg: "server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MapHTTPError(tt.status, []byte(tt.body))
			if !strings.Contains(err.Error(), tt.wantMsg) {
				t.Errorf("expected %q in error, got: %s", tt.wantMsg, err.Error())
			}
		})
	}
}

func TestMapHTTPError_IncludesServerBody(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		wantLabel string
		wantHint  string
		status    int
	}{
		{name: "401 with body", status: 401, body: `{"message":"Token expired"}`, wantLabel: "authentication failed", wantHint: "gio login"},
		{name: "403 with body", status: 403, body: `{"message":"forbidden"}`, wantLabel: "access denied", wantHint: "token permissions"},
		{name: "404 with body", status: 404, body: `{"message":"app not found"}`, wantLabel: "resource not found"},
		{name: "500 with body", status: 500, body: "stack trace here", wantLabel: "server error", wantHint: "APIM server status"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MapHTTPError(tt.status, []byte(tt.body))

			if !strings.Contains(err.Error(), tt.wantLabel) {
				t.Errorf("expected label %q, got: %s", tt.wantLabel, err.Error())
			}

			if !strings.Contains(err.Error(), tt.body) {
				t.Errorf("expected body %q in error, got: %s", tt.body, err.Error())
			}

			if tt.wantHint != "" && !strings.Contains(err.Error(), tt.wantHint) {
				t.Errorf("expected hint %q in error, got: %s", tt.wantHint, err.Error())
			}
		})
	}
}

func TestMapHTTPError_EmptyBodyOmitsColon(t *testing.T) {
	// With an empty body, we should not emit a dangling ": " before the newline.
	tests := []struct {
		name   string
		status int
	}{
		{"401 no body", 401},
		{"404 no body", 404},
		{"400 no body", 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MapHTTPError(tt.status, nil)
			firstLine := strings.SplitN(err.Error(), "\n", 2)[0]

			if strings.HasSuffix(firstLine, ": ") || strings.HasSuffix(firstLine, ":") {
				t.Errorf("expected no trailing colon on empty body, got: %q", firstLine)
			}
		})
	}
}

func TestMaskToken(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"gioat_abc123xyz", "************xyz"},
		{"ab", "***"},
		{"", "***"},
	}

	for _, tt := range tests {
		got := maskToken(tt.input)
		if got != tt.want {
			t.Errorf("maskToken(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
