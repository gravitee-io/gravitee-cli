package cmdutil

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
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
