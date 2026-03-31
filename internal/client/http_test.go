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
