package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFromFileNotFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.json")

	cfg, err := LoadFrom(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Contexts) != 0 {
		t.Errorf("expected 0 contexts, got %d", len(cfg.Contexts))
	}
}

func TestLoadFromValidConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	content := `{
		"currentContext": "dev",
		"contexts": {
			"dev": {
				"url": "https://apim-dev.company.com",
				"token": "gioat_abc",
				"org": "DEFAULT",
				"env": "DEFAULT"
			}
		}
	}`

	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadFrom(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Contexts) != 1 {
		t.Errorf("expected 1 context, got %d", len(cfg.Contexts))
	}
}

func TestLoadFromInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	if err := os.WriteFile(path, []byte(`{invalid}`), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := LoadFrom(path)
	if err == nil || !strings.Contains(err.Error(), "failed to parse config file") {
		t.Errorf("expected parse error, got: %v", err)
	}
}

func TestSaveTo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "config.json")

	cfg := &Config{
		CurrentContext: "test",
		Contexts: map[string]Context{
			"test": {
				URL:   "https://apim.example.com",
				Token: "tok_123",
				Org:   "DEFAULT",
				Env:   "DEFAULT",
			},
		},
	}

	if err := cfg.SaveTo(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, err := LoadFrom(path)
	if err != nil {
		t.Fatalf("failed to reload: %v", err)
	}

	if loaded.CurrentContext != "test" {
		t.Errorf("expected currentContext 'test', got %q", loaded.CurrentContext)
	}

	if loaded.Contexts["test"].URL != "https://apim.example.com" {
		t.Errorf("unexpected URL: %s", loaded.Contexts["test"].URL)
	}
}

func TestResolve(t *testing.T) {
	baseCfg := &Config{
		CurrentContext: "dev",
		Contexts: map[string]Context{
			"dev": {
				URL:   "https://apim-dev.company.com",
				Token: "tok_dev",
				Org:   "MY_ORG",
				Env:   "staging",
			},
			"prod": {
				URL:      "https://apim-prod.company.com",
				Token:    "tok_prod",
				ReadOnly: true,
			},
		},
	}

	t.Run("uses current context", func(t *testing.T) {
		resolved, err := baseCfg.Resolve(Overrides{})
		assertNoError(t, err)
		assertEqual(t, "dev", resolved.Name, "name")
		assertEqual(t, "MY_ORG", resolved.Org, "org")
		assertEqual(t, "staging", resolved.Env, "env")
	})

	t.Run("override context", func(t *testing.T) {
		resolved, err := baseCfg.Resolve(Overrides{Context: "prod"})
		assertNoError(t, err)
		assertEqual(t, "prod", resolved.Name, "name")
		assertEqual(t, "DEFAULT", resolved.Org, "org")
		assertEqual(t, "DEFAULT", resolved.Env, "env")

		if !resolved.ReadOnly {
			t.Error("expected readOnly=true")
		}
	})

	t.Run("override org and env", func(t *testing.T) {
		resolved, err := baseCfg.Resolve(Overrides{Org: "OTHER", EnvID: "production"})
		assertNoError(t, err)
		assertEqual(t, "OTHER", resolved.Org, "org")
		assertEqual(t, "production", resolved.Env, "env")
	})

	t.Run("no context configured", func(t *testing.T) {
		cfg := &Config{Contexts: map[string]Context{}}
		_, err := cfg.Resolve(Overrides{})
		assertErrorContains(t, err, "no context configured")
	})

	t.Run("context not found", func(t *testing.T) {
		_, err := baseCfg.Resolve(Overrides{Context: "nonexistent"})
		assertErrorContains(t, err, "context 'nonexistent' not found")
	})

	t.Run("empty org/env defaults to DEFAULT", func(t *testing.T) {
		cfg := &Config{
			CurrentContext: "bare",
			Contexts: map[string]Context{
				"bare": {URL: "https://example.com", Token: "tok"},
			},
		}

		resolved, err := cfg.Resolve(Overrides{})
		assertNoError(t, err)
		assertEqual(t, "DEFAULT", resolved.Org, "org")
		assertEqual(t, "DEFAULT", resolved.Env, "env")
	})
}

func assertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertEqual(t *testing.T, want, got, field string) {
	t.Helper()

	if got != want {
		t.Errorf("expected %s %q, got %q", field, want, got)
	}
}

func assertErrorContains(t *testing.T, err error, substr string) {
	t.Helper()

	if err == nil || !strings.Contains(err.Error(), substr) {
		t.Errorf("expected error containing %q, got: %v", substr, err)
	}
}
