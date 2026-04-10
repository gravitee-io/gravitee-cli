package cmdutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func TestRequireContext_ReturnsResolveError(t *testing.T) {
	f := &factory.Factory{
		ContextResolveErr: fmt.Errorf("context 'prod2' not found"),
	}

	err := RequireContext(f)

	if err == nil || err.Error() != "context 'prod2' not found" {
		t.Errorf("expected resolve error, got: %v", err)
	}
}

func TestRequireContext_ProductAwareHint(t *testing.T) {
	tests := []struct {
		product  string
		wantHint string
	}{
		{"apim", "gio login apim"},
		{"am", "gio login am"},
		{"", "gio login apim"},
	}

	for _, tt := range tests {
		t.Run(tt.product, func(t *testing.T) {
			f := &factory.Factory{Product: tt.product}

			err := RequireContext(f)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !contains(err.Error(), tt.wantHint) {
				t.Errorf("expected hint containing %q, got: %s", tt.wantHint, err.Error())
			}
		})
	}
}

func TestNewPrinter_ValidFormats(t *testing.T) {
	for _, format := range []string{"table", "json", "yaml"} {
		t.Run(format, func(t *testing.T) {
			f := &factory.Factory{
				OutputFormat: format,
				IOStreams:    factory.IOStreams{Out: &bytes.Buffer{}},
			}

			p, err := NewPrinter(f)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if p == nil {
				t.Fatal("expected printer, got nil")
			}
		})
	}
}

func TestNewPrinter_InvalidFormat(t *testing.T) {
	f := &factory.Factory{
		OutputFormat: "xml",
		IOStreams:    factory.IOStreams{Out: &bytes.Buffer{}, Err: &bytes.Buffer{}},
	}

	_, err := NewPrinter(f)

	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}

	if !contains(err.Error(), "xml") || !contains(err.Error(), "table, json, yaml") {
		t.Errorf("expected error mentioning format and allowed values, got: %s", err.Error())
	}
}

func TestPrintActionResult(t *testing.T) {
	tests := []struct {
		format   string
		check    func(t *testing.T, out string)
		humanMsg string
		id       string
		status   string
		name     string
	}{
		{
			name:   "table emits the human message",
			format: printer.FormatTable,
			id:     "api-1", status: "deleted",
			humanMsg: "API 'api-1' deleted.",
			check: func(t *testing.T, out string) {
				t.Helper()
				if !contains(out, "API 'api-1' deleted.") {
					t.Errorf("expected human message, got: %s", out)
				}
			},
		},
		{
			name:   "json emits {id, status} envelope",
			format: printer.FormatJSON,
			id:     "api-1", status: "deleted",
			humanMsg: "unused",
			check: func(t *testing.T, out string) {
				t.Helper()
				var got map[string]string
				if err := json.Unmarshal([]byte(out), &got); err != nil {
					t.Fatalf("expected valid JSON, got: %s", out)
				}
				if got["id"] != "api-1" || got["status"] != "deleted" {
					t.Errorf("unexpected json: %+v", got)
				}
			},
		},
		{
			name:   "yaml emits id and status",
			format: printer.FormatYAML,
			id:     "api-1", status: "started",
			humanMsg: "unused",
			check: func(t *testing.T, out string) {
				t.Helper()
				if !contains(out, "id: api-1") || !contains(out, "status: started") {
					t.Errorf("expected yaml id/status, got: %s", out)
				}
			},
		},
		{
			name:   "id emits the id alone",
			format: printer.FormatID,
			id:     "api-1", status: "deleted",
			humanMsg: "unused",
			check: func(t *testing.T, out string) {
				t.Helper()
				if strings.TrimSpace(out) != "api-1" {
					t.Errorf("expected 'api-1' alone, got: %q", out)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			p := printer.New(tt.format, out, &bytes.Buffer{}, false, false)

			if err := PrintActionResult(p, tt.id, tt.status, tt.humanMsg); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			tt.check(t, out.String())
		})
	}
}

func TestRequireNonEmpty(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		wantErr  bool
		wantHint string
	}{
		{"valid id", "abc-123", false, ""},
		{"empty string", "", true, "envId cannot be empty"},
		{"whitespace only", "   ", true, "envId cannot be empty"},
		{"tab and newline only", "\t\n", true, "envId cannot be empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RequireNonEmpty("envId", tt.value)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if !contains(err.Error(), tt.wantHint) {
					t.Errorf("expected error containing %q, got: %s", tt.wantHint, err.Error())
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
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
		{"abc", "***"},
		{"abcd", "*bcd"},
	}

	for _, tt := range tests {
		got := MaskToken(tt.input)
		if got != tt.want {
			t.Errorf("MaskToken(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestValidatePagination(t *testing.T) {
	tests := []struct {
		name    string
		page    int
		perPage int
		wantErr string
	}{
		{"valid", 1, 10, ""},
		{"page zero", 0, 10, "--page must be >= 1"},
		{"page negative", -1, 10, "--page must be >= 1"},
		{"per-page zero", 1, 0, "--per-page must be >= 1"},
		{"per-page negative", 1, -5, "--per-page must be >= 1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePagination(tt.page, tt.perPage)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			} else {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !contains(err.Error(), tt.wantErr) {
					t.Errorf("expected %q, got %q", tt.wantErr, err.Error())
				}
			}
		})
	}
}

func TestResolveProductContext_EnvOrgAppliesOnEnvBypass(t *testing.T) {
	// Product env vars (URL+Token) trigger the bypass branch; GIO_ORG/GIO_ENV
	// must still apply on that branch, not just on the config-file branch.
	t.Setenv("GIO_APIM_URL", "https://envhost")
	t.Setenv("GIO_APIM_TOKEN", "env_tok")
	t.Setenv(EnvOrg, "ENV_ORG")
	t.Setenv(EnvEnv, "ENV_ENV")

	f := &factory.Factory{
		IOStreams: factory.IOStreams{Out: &bytes.Buffer{}, Err: &bytes.Buffer{}},
	}

	if err := ResolveProductContext(f, "apim"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if f.Resolved == nil || f.Resolved.Name != "env" {
		t.Fatalf("expected env-bypass resolution, got %+v", f.Resolved)
	}

	if f.Resolved.Org != "ENV_ORG" {
		t.Errorf("Org: got %q, want ENV_ORG", f.Resolved.Org)
	}

	if f.Resolved.Env != "ENV_ENV" {
		t.Errorf("Env: got %q, want ENV_ENV", f.Resolved.Env)
	}
}

func TestResolveProductContext_EnvOrgFallsBackToDefaultOnBypass(t *testing.T) {
	// On the env-bypass branch, when GIO_ORG/GIO_ENV are not set, fall back
	// to the package-level defaults.
	t.Setenv("GIO_APIM_URL", "https://envhost")
	t.Setenv("GIO_APIM_TOKEN", "env_tok")
	t.Setenv(EnvOrg, "")
	t.Setenv(EnvEnv, "")

	f := &factory.Factory{
		IOStreams: factory.IOStreams{Out: &bytes.Buffer{}, Err: &bytes.Buffer{}},
	}

	if err := ResolveProductContext(f, "apim"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if f.Resolved.Org != config.DefaultOrg {
		t.Errorf("Org: got %q, want %s", f.Resolved.Org, config.DefaultOrg)
	}

	if f.Resolved.Env != config.DefaultEnv {
		t.Errorf("Env: got %q, want %s", f.Resolved.Env, config.DefaultEnv)
	}
}

// Priority: env var > flag > config. Sets all three, expects env var to win.
func TestResolveProductContext_EnvOrgBeatsFlagBeatsConfig(t *testing.T) {
	t.Setenv(EnvOrg, "ENV_ORG")
	t.Setenv(EnvEnv, "ENV_ENV")
	t.Setenv("GIO_APIM_URL", "")
	t.Setenv("GIO_APIM_TOKEN", "")

	f := &factory.Factory{
		Config: &config.Config{
			Current: "prod",
			Contexts: map[string]*config.Context{
				"prod": {
					Org:  "CONFIG_ORG",
					Env:  "CONFIG_ENV",
					APIM: &config.ProductConfig{URL: "https://x", Token: "tok"},
				},
			},
		},
		Overrides: config.Overrides{Org: "FLAG_ORG", EnvID: "FLAG_ENV"},
		IOStreams: factory.IOStreams{Out: &bytes.Buffer{}, Err: &bytes.Buffer{}},
	}

	if err := ResolveProductContext(f, "apim"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if f.Resolved.Org != "ENV_ORG" {
		t.Errorf("env var must win over flag and config: got Org=%q, want ENV_ORG", f.Resolved.Org)
	}

	if f.Resolved.Env != "ENV_ENV" {
		t.Errorf("env var must win over flag and config: got Env=%q, want ENV_ENV", f.Resolved.Env)
	}
}

// No env var, flag set, config set -> flag wins.
func TestResolveProductContext_FlagBeatsConfig(t *testing.T) {
	t.Setenv(EnvOrg, "")
	t.Setenv(EnvEnv, "")
	t.Setenv("GIO_APIM_URL", "")
	t.Setenv("GIO_APIM_TOKEN", "")

	f := &factory.Factory{
		Config: &config.Config{
			Current: "prod",
			Contexts: map[string]*config.Context{
				"prod": {
					Org:  "CONFIG_ORG",
					Env:  "CONFIG_ENV",
					APIM: &config.ProductConfig{URL: "https://x", Token: "tok"},
				},
			},
		},
		Overrides: config.Overrides{Org: "FLAG_ORG", EnvID: "FLAG_ENV"},
		IOStreams: factory.IOStreams{Out: &bytes.Buffer{}, Err: &bytes.Buffer{}},
	}

	if err := ResolveProductContext(f, "apim"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if f.Resolved.Org != "FLAG_ORG" {
		t.Errorf("flag must win over config: got Org=%q, want FLAG_ORG", f.Resolved.Org)
	}

	if f.Resolved.Env != "FLAG_ENV" {
		t.Errorf("flag must win over config: got Env=%q, want FLAG_ENV", f.Resolved.Env)
	}
}

func TestResolveProductContext_StoresError(t *testing.T) {
	f := &factory.Factory{
		Config: &config.Config{
			Current:  "nonexistent",
			Contexts: map[string]*config.Context{},
		},
		IOStreams: factory.IOStreams{Out: &bytes.Buffer{}, Err: &bytes.Buffer{}},
	}

	// ResolveProductContext should store the error, not return it
	err := ResolveProductContext(f, "apim")
	if err != nil {
		t.Fatalf("expected nil return, got: %v", err)
	}

	// RequireContext should surface the stored error
	err = RequireContext(f)
	if err == nil || !contains(err.Error(), "nonexistent") {
		t.Errorf("expected error about nonexistent context, got: %v", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}

func TestParseLoginURL(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		wantBase string
		wantOrg  string
		wantEnv  string
		wantErr  string
	}{
		{
			name:     "full apim path",
			in:       "https://apim.example.com/management/organizations/DEFAULT/environments/DEFAULT",
			wantBase: "https://apim.example.com",
			wantOrg:  "DEFAULT",
			wantEnv:  "DEFAULT",
		},
		{
			name:     "uuid org and env",
			in:       "https://apim.example.com/management/organizations/11111111-1111-1111-1111-111111111111/environments/22222222-2222-2222-2222-222222222222",
			wantBase: "https://apim.example.com",
			wantOrg:  "11111111-1111-1111-1111-111111111111",
			wantEnv:  "22222222-2222-2222-2222-222222222222",
		},
		{
			name:     "am environments list (no env id)",
			in:       "https://am.example.com/management/organizations/DEFAULT/environments",
			wantBase: "https://am.example.com",
			wantOrg:  "DEFAULT",
			wantEnv:  "",
		},
		{
			name:     "org only",
			in:       "https://host/management/organizations/ACME",
			wantBase: "https://host",
			wantOrg:  "ACME",
		},
		{
			name:     "trailing resource path stripped",
			in:       "https://host/management/organizations/ACME/environments/prod/apis/abc",
			wantBase: "https://host",
			wantOrg:  "ACME",
			wantEnv:  "prod",
		},
		{
			name:     "bare base",
			in:       "https://host",
			wantBase: "https://host",
		},
		{
			name:     "bare base with trailing slash",
			in:       "https://host/",
			wantBase: "https://host",
		},
		{
			name:    "missing scheme",
			in:      "localhost:8083",
			wantErr: "must start with http",
		},
		{
			name:    "empty",
			in:      "",
			wantErr: "URL is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base, org, env, err := ParseLoginURL(tt.in)

			if tt.wantErr != "" {
				if err == nil || !contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got: %v", tt.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if base != tt.wantBase {
				t.Errorf("base: got %q, want %q", base, tt.wantBase)
			}

			if org != tt.wantOrg {
				t.Errorf("org: got %q, want %q", org, tt.wantOrg)
			}

			if env != tt.wantEnv {
				t.Errorf("env: got %q, want %q", env, tt.wantEnv)
			}
		})
	}
}

func TestParseCurl(t *testing.T) {
	tests := []struct {
		name      string
		in        string
		wantURL   string
		wantToken string
		wantErr   string
	}{
		{
			name:      "apim devtools: -H before URL, double quotes",
			in:        `curl -H "Authorization: Bearer tok_apim_fake123" "https://apim.example.com/management/organizations/DEFAULT/environments/DEFAULT"`,
			wantURL:   "https://apim.example.com/management/organizations/DEFAULT/environments/DEFAULT",
			wantToken: "tok_apim_fake123",
		},
		{
			name:      "am devtools: URL before -H, single quotes",
			in:        `curl https://am.example.com/management/organizations/DEFAULT/environments -H 'Authorization: Bearer tok_am_fake456=='`,
			wantURL:   "https://am.example.com/management/organizations/DEFAULT/environments",
			wantToken: "tok_am_fake456==",
		},
		{
			name:      "lowercase bearer keyword",
			in:        `curl http://x -H "authorization: bearer abc"`,
			wantURL:   "http://x",
			wantToken: "abc",
		},
		{
			name:      "token with base64 padding",
			in:        `curl http://x -H "Authorization: Bearer abc=="`,
			wantURL:   "http://x",
			wantToken: "abc==",
		},
		{
			name:      "unknown flags ignored",
			in:        `curl -k -s --compressed http://x -H "Authorization: Bearer tok"`,
			wantURL:   "http://x",
			wantToken: "tok",
		},
		{
			name:      "header joined form --header=VALUE",
			in:        `curl http://x --header=Authorization:Bearer\ tok`,
			wantURL:   "http://x",
			wantToken: "tok",
		},
		{
			name:    "missing URL",
			in:      `curl -H "Authorization: Bearer tok"`,
			wantErr: "missing URL",
		},
		{
			name:    "missing auth header",
			in:      `curl http://x`,
			wantErr: "missing Authorization",
		},
		{
			name:    "non-bearer auth rejected",
			in:      `curl http://x -H "Authorization: Basic xxx"`,
			wantErr: "missing Authorization",
		},
		{
			name:    "unterminated quote",
			in:      `curl "http://x`,
			wantErr: "unterminated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, gotToken, err := ParseCurl(tt.in)

			if tt.wantErr != "" {
				if err == nil || !contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got: %v", tt.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if gotURL != tt.wantURL {
				t.Errorf("url: got %q, want %q", gotURL, tt.wantURL)
			}

			if gotToken != tt.wantToken {
				t.Errorf("token: got %q, want %q", gotToken, tt.wantToken)
			}
		})
	}
}
