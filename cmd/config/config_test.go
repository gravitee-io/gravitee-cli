package config

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	iconfig "github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(t *testing.T) (*factory.Factory, *bytes.Buffer) {
	t.Helper()

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")
	out := &bytes.Buffer{}

	cfg := &iconfig.Config{
		CurrentContext: "dev",
		Contexts: map[string]iconfig.Context{
			"dev": {
				URL:   "https://apim-dev.company.com",
				Token: "gioat_devtoken",
				Org:   "DEFAULT",
				Env:   "DEFAULT",
			},
			"prod": {
				URL:      "https://apim-prod.company.com",
				Token:    "gioat_prodtoken",
				Org:      "DEFAULT",
				Env:      "production",
				ReadOnly: true,
			},
		},
	}

	f := &factory.Factory{
		Config:       cfg,
		ConfigPath:   cfgPath,
		OutputFormat: "table",
		IOStreams:    factory.IOStreams{Out: out, Err: &bytes.Buffer{}},
	}

	return f, out
}

func TestSetContext(t *testing.T) {
	t.Run("success creates new context", func(t *testing.T) {
		f, out := newTestFactory(t)

		cmd := newSetContextCmd(f)
		cmd.SetArgs([]string{
			"staging",
			"--url", "https://apim-staging.company.com",
			"--token", "gioat_stg",
			"--env-id", "staging",
			"--read-only",
		})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.Contains(out.String(), "Context 'staging' saved.") {
			t.Errorf("unexpected output: %s", out.String())
		}

		// Current context should NOT change.
		cfg, err := iconfig.LoadFrom(f.ConfigPath)
		if err != nil {
			t.Fatalf("failed to load: %v", err)
		}

		if cfg.CurrentContext != "dev" {
			t.Errorf("expected current context 'dev', got %q", cfg.CurrentContext)
		}

		ctx := cfg.Contexts["staging"]
		if ctx.URL != "https://apim-staging.company.com" {
			t.Errorf("unexpected URL: %s", ctx.URL)
		}

		if !ctx.ReadOnly {
			t.Error("expected read-only to be true")
		}
	})

	t.Run("missing required flags", func(t *testing.T) {
		f, _ := newTestFactory(t)

		cmd := newSetContextCmd(f)
		cmd.SetArgs([]string{"staging"})
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true

		err := cmd.Execute()
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestUseContext(t *testing.T) {
	t.Run("success switches context", func(t *testing.T) {
		f, out := newTestFactory(t)

		cmd := newUseContextCmd(f)
		cmd.SetArgs([]string{"prod"})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.Contains(out.String(), "Switched to context 'prod'.") {
			t.Errorf("unexpected output: %s", out.String())
		}

		cfg, err := iconfig.LoadFrom(f.ConfigPath)
		if err != nil {
			t.Fatalf("failed to load: %v", err)
		}

		if cfg.CurrentContext != "prod" {
			t.Errorf("expected current context 'prod', got %q", cfg.CurrentContext)
		}
	})

	t.Run("context not found", func(t *testing.T) {
		f, _ := newTestFactory(t)

		cmd := newUseContextCmd(f)
		cmd.SetArgs([]string{"nonexistent"})
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true

		err := cmd.Execute()
		if err == nil {
			t.Fatal("expected error")
		}

		if !strings.Contains(err.Error(), "context 'nonexistent' not found") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestCurrentContext(t *testing.T) {
	t.Run("success shows current", func(t *testing.T) {
		f, out := newTestFactory(t)

		cmd := newCurrentContextCmd(f)

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if strings.TrimSpace(out.String()) != "dev" {
			t.Errorf("expected 'dev', got %q", out.String())
		}
	})

	t.Run("no context configured", func(t *testing.T) {
		out := &bytes.Buffer{}

		f := &factory.Factory{
			Config:    &iconfig.Config{Contexts: make(map[string]iconfig.Context)},
			IOStreams: factory.IOStreams{Out: out, Err: &bytes.Buffer{}},
		}

		cmd := newCurrentContextCmd(f)
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true

		err := cmd.Execute()
		if err == nil {
			t.Fatal("expected error")
		}

		if !strings.Contains(err.Error(), "no context configured") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestGetContexts(t *testing.T) {
	t.Run("success table output", func(t *testing.T) {
		f, out := newTestFactory(t)

		cmd := newGetContextsCmd(f)

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := out.String()

		if !strings.Contains(output, "dev") {
			t.Error("expected 'dev' in output")
		}

		if !strings.Contains(output, "prod") {
			t.Error("expected 'prod' in output")
		}

		if !strings.Contains(output, "*") {
			t.Error("expected '*' marker for current context")
		}

		if !strings.Contains(output, "production") {
			t.Error("expected 'production' env in output")
		}
	})

	t.Run("json output", func(t *testing.T) {
		f, out := newTestFactory(t)
		f.OutputFormat = "json"

		cmd := newGetContextsCmd(f)

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := out.String()

		if !strings.Contains(output, `"currentContext"`) {
			t.Error("expected currentContext in JSON output")
		}

		if !strings.Contains(output, `"contexts"`) {
			t.Error("expected contexts in JSON output")
		}
	})
}

func TestViewSuccess(t *testing.T) {
	f, out := newTestFactory(t)

	cmd := newViewCmd(f)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()

	if !strings.Contains(output, "Context:    dev") {
		t.Errorf("expected 'Context:    dev' in output, got: %s", output)
	}

	if !strings.Contains(output, "URL:        https://apim-dev.company.com") {
		t.Error("expected URL in output")
	}

	if !strings.Contains(output, "Read-only:  no") {
		t.Error("expected 'Read-only:  no' in output")
	}
}

func TestViewTokenMasked(t *testing.T) {
	f, out := newTestFactory(t)

	cmd := newViewCmd(f)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()

	if strings.Contains(output, "gioat_devtoken") {
		t.Error("full token should not appear in output")
	}

	if !strings.Contains(output, "ken") {
		t.Error("last 3 chars of token should appear")
	}
}

func TestViewReadOnly(t *testing.T) {
	f, out := newTestFactory(t)
	f.Config.CurrentContext = "prod"

	cmd := newViewCmd(f)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Read-only:  yes") {
		t.Error("expected 'Read-only:  yes' in output")
	}
}

func TestViewNoContext(t *testing.T) {
	f := &factory.Factory{
		Config:    &iconfig.Config{Contexts: make(map[string]iconfig.Context)},
		IOStreams: factory.IOStreams{Out: &bytes.Buffer{}, Err: &bytes.Buffer{}},
	}

	cmd := newViewCmd(f)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "no context configured") {
		t.Errorf("unexpected error: %v", err)
	}
}
