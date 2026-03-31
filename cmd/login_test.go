package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newLoginTestFactory(t *testing.T) (*factory.Factory, *bytes.Buffer, string) {
	t.Helper()

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")
	out := &bytes.Buffer{}

	f := &factory.Factory{
		Config:     &config.Config{Contexts: make(map[string]config.Context)},
		ConfigPath: cfgPath,
		IOStreams:  factory.IOStreams{Out: out, Err: &bytes.Buffer{}},
	}

	return f, out, cfgPath
}

func TestLoginSuccess(t *testing.T) {
	f, out, cfgPath := newLoginTestFactory(t)

	cmd := newLoginCmd(f)
	cmd.SetArgs([]string{
		"--url", "https://apim.example.com",
		"--token", "gioat_abc123",
		"--context", "myctx",
		"--org", "MY_ORG",
		"--env-id", "staging",
		"--read-only",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Context 'myctx' saved and set as current.") {
		t.Errorf("unexpected output: %s", out.String())
	}

	cfg, err := config.LoadFrom(cfgPath)
	if err != nil {
		t.Fatalf("failed to load saved config: %v", err)
	}

	if cfg.CurrentContext != "myctx" {
		t.Errorf("expected current context 'myctx', got %q", cfg.CurrentContext)
	}

	ctx := cfg.Contexts["myctx"]
	if ctx.URL != "https://apim.example.com" {
		t.Errorf("unexpected URL: %s", ctx.URL)
	}

	if ctx.Org != "MY_ORG" {
		t.Errorf("unexpected org: %s", ctx.Org)
	}

	if !ctx.ReadOnly {
		t.Error("expected read-only to be true")
	}
}

func TestLoginDefaults(t *testing.T) {
	f, _, cfgPath := newLoginTestFactory(t)

	cmd := newLoginCmd(f)
	cmd.SetArgs([]string{
		"--url", "https://apim.example.com",
		"--token", "gioat_abc123",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg, err := config.LoadFrom(cfgPath)
	if err != nil {
		t.Fatalf("failed to load saved config: %v", err)
	}

	if cfg.CurrentContext != "default" {
		t.Errorf("expected current context 'default', got %q", cfg.CurrentContext)
	}

	ctx := cfg.Contexts["default"]
	if ctx.Org != "DEFAULT" {
		t.Errorf("expected org DEFAULT, got %q", ctx.Org)
	}

	if ctx.Env != "DEFAULT" {
		t.Errorf("expected env DEFAULT, got %q", ctx.Env)
	}
}

func TestLoginMissingFlags(t *testing.T) {
	f := &factory.Factory{
		Config:    &config.Config{Contexts: make(map[string]config.Context)},
		IOStreams: factory.IOStreams{Out: &bytes.Buffer{}, Err: &bytes.Buffer{}},
	}

	cmd := newLoginCmd(f)
	cmd.SetArgs([]string{})
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing flags")
	}

	if !strings.Contains(err.Error(), "required") {
		t.Errorf("expected 'required' in error, got: %v", err)
	}
}

func TestLoginOverwritesExisting(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")

	var out bytes.Buffer

	f := &factory.Factory{
		Config: &config.Config{
			CurrentContext: "old",
			Contexts: map[string]config.Context{
				"old": {URL: "https://old.com", Token: "old_tok"},
			},
		},
		ConfigPath: cfgPath,
		IOStreams:  factory.IOStreams{Out: &out, Err: &bytes.Buffer{}},
	}

	cmd := newLoginCmd(f)
	cmd.SetArgs([]string{
		"--url", "https://new.com",
		"--token", "new_tok",
		"--context", "old",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg, err := config.LoadFrom(cfgPath)
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	if cfg.Contexts["old"].URL != "https://new.com" {
		t.Errorf("expected overwritten URL, got %s", cfg.Contexts["old"].URL)
	}
}
