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

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFromFileNotFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.yaml")

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
	path := filepath.Join(dir, "config.yaml")

	content := `
current: dev
contexts:
  dev:
    org: MY_ORG
    env: staging
    apim:
      url: https://apim-dev.company.com
      token: gioat_abc
`

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

	ctx := cfg.Contexts["dev"]
	if ctx == nil {
		t.Fatal("expected 'dev' context")
	}

	if ctx.APIM == nil || ctx.APIM.URL != "https://apim-dev.company.com" {
		t.Errorf("unexpected APIM config: %+v", ctx.APIM)
	}

	if ctx.AM != nil {
		t.Error("expected nil AM config")
	}
}

func TestLoadFromInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	if err := os.WriteFile(path, []byte(`{invalid: [}`), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := LoadFrom(path)
	if err == nil || !strings.Contains(err.Error(), "failed to parse config file") {
		t.Errorf("expected parse error, got: %v", err)
	}
}

func TestSaveToRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "config.yaml")

	cfg := &Config{
		Current: "prod",
		Contexts: map[string]*Context{
			"prod": {
				Org: "ACME",
				Env: "production",
				APIM: &ProductConfig{
					URL:   "https://apim.example.com",
					Token: "tok_apim",
				},
				AM: &ProductConfig{
					URL:   "https://am.example.com",
					Token: "tok_am",
				},
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

	assertEqual(t, "prod", loaded.Current, "current")

	ctx := loaded.Contexts["prod"]
	if ctx == nil {
		t.Fatal("expected 'prod' context")
	}

	assertEqual(t, "ACME", ctx.Org, "org")
	assertEqual(t, "production", ctx.Env, "env")
	assertEqual(t, "https://apim.example.com", ctx.APIM.URL, "apim url")
	assertEqual(t, "tok_apim", ctx.APIM.Token, "apim token")
	assertEqual(t, "https://am.example.com", ctx.AM.URL, "am url")
	assertEqual(t, "tok_am", ctx.AM.Token, "am token")
}

func TestResolve(t *testing.T) {
	baseCfg := &Config{
		Current: "dev",
		Contexts: map[string]*Context{
			"dev": {
				Org: "MY_ORG",
				Env: "staging",
				APIM: &ProductConfig{
					URL:   "https://apim-dev.company.com",
					Token: "tok_dev",
				},
			},
			"prod": {
				APIM: &ProductConfig{
					URL:   "https://apim-prod.company.com",
					Token: "tok_prod",
				},
				AM: &ProductConfig{
					URL:   "https://am-prod.company.com",
					Token: "tok_am",
				},
			},
		},
	}

	t.Run("uses current context", func(t *testing.T) {
		resolved, err := baseCfg.Resolve(Overrides{}, "apim")
		assertNoError(t, err)
		assertEqual(t, "dev", resolved.Name, "name")
		assertEqual(t, "MY_ORG", resolved.Org, "org")
		assertEqual(t, "staging", resolved.Env, "env")
		assertEqual(t, "https://apim-dev.company.com", resolved.URL, "url")
	})

	t.Run("override context", func(t *testing.T) {
		resolved, err := baseCfg.Resolve(Overrides{Context: "prod"}, "apim")
		assertNoError(t, err)
		assertEqual(t, "prod", resolved.Name, "name")
		assertEqual(t, "DEFAULT", resolved.Org, "org")
		assertEqual(t, "DEFAULT", resolved.Env, "env")
	})

	t.Run("override org and env", func(t *testing.T) {
		resolved, err := baseCfg.Resolve(Overrides{Org: "OTHER", EnvID: "production"}, "apim")
		assertNoError(t, err)
		assertEqual(t, "OTHER", resolved.Org, "org")
		assertEqual(t, "production", resolved.Env, "env")
	})

	t.Run("resolves AM product", func(t *testing.T) {
		resolved, err := baseCfg.Resolve(Overrides{Context: "prod"}, "am")
		assertNoError(t, err)
		assertEqual(t, "https://am-prod.company.com", resolved.URL, "url")
		assertEqual(t, "tok_am", resolved.Token, "token")
	})

	t.Run("no context configured", func(t *testing.T) {
		cfg := &Config{Contexts: map[string]*Context{}}
		_, err := cfg.Resolve(Overrides{}, "apim")
		assertErrorContains(t, err, "no context configured")
	})

	t.Run("context not found", func(t *testing.T) {
		_, err := baseCfg.Resolve(Overrides{Context: "nonexistent"}, "apim")
		assertErrorContains(t, err, "context 'nonexistent' not found")
	})

	t.Run("product not configured in context", func(t *testing.T) {
		_, err := baseCfg.Resolve(Overrides{}, "am")
		assertErrorContains(t, err, "AM not configured for context 'dev'")
	})

	t.Run("empty org/env defaults to DEFAULT", func(t *testing.T) {
		resolved, err := baseCfg.Resolve(Overrides{Context: "prod"}, "apim")
		assertNoError(t, err)
		assertEqual(t, "DEFAULT", resolved.Org, "org")
		assertEqual(t, "DEFAULT", resolved.Env, "env")
	})
}

func TestEnsureContext(t *testing.T) {
	t.Run("creates new context", func(t *testing.T) {
		cfg := &Config{Contexts: map[string]*Context{}}
		ctx := cfg.EnsureContext("new")

		if ctx == nil {
			t.Fatal("expected non-nil context")
		}

		if _, ok := cfg.Contexts["new"]; !ok {
			t.Error("expected context to be stored in config")
		}
	})

	t.Run("returns existing context", func(t *testing.T) {
		existing := &Context{Org: "ACME"}
		cfg := &Config{Contexts: map[string]*Context{"prod": existing}}

		ctx := cfg.EnsureContext("prod")

		if ctx.Org != "ACME" {
			t.Error("expected existing context to be returned")
		}
	})
}

func TestDeleteContext(t *testing.T) {
	t.Run("deletes context", func(t *testing.T) {
		cfg := &Config{
			Current:  "prod",
			Contexts: map[string]*Context{"prod": {}, "dev": {}},
		}

		err := cfg.DeleteContext("dev")
		assertNoError(t, err)

		if _, ok := cfg.Contexts["dev"]; ok {
			t.Error("expected context to be deleted")
		}
	})

	t.Run("clears current if deleting active context", func(t *testing.T) {
		cfg := &Config{
			Current:  "prod",
			Contexts: map[string]*Context{"prod": {}},
		}

		err := cfg.DeleteContext("prod")
		assertNoError(t, err)

		if cfg.Current != "" {
			t.Errorf("expected current to be cleared, got %q", cfg.Current)
		}
	})

	t.Run("error on nonexistent context", func(t *testing.T) {
		cfg := &Config{Contexts: map[string]*Context{}}
		err := cfg.DeleteContext("nope")
		assertErrorContains(t, err, "not found")
	})
}

func TestContextNames(t *testing.T) {
	cfg := &Config{
		Contexts: map[string]*Context{
			"prod": {},
			"dev":  {},
			"beta": {},
		},
	}

	names := cfg.ContextNames()

	if len(names) != 3 {
		t.Fatalf("expected 3 names, got %d", len(names))
	}

	if names[0] != "beta" || names[1] != "dev" || names[2] != "prod" {
		t.Errorf("expected sorted names [beta dev prod], got %v", names)
	}
}

func TestNormalizeContextName(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"local-master", "local-master"},
		{"Local Master", "local-master"},
		{"  PROD  ", "prod"},
		{"My Test Context", "my-test-context"},
		{"already-clean", "already-clean"},
		{"", ""},
		{"UPPER", "upper"},
		{"a b c", "a-b-c"},
	}

	for _, tt := range tests {
		got := NormalizeContextName(tt.in)
		if got != tt.want {
			t.Errorf("NormalizeContextName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

// Resolve must normalize overrides.Context so a user passing
// --context "Local Master" matches a stored "local-master" context.
func TestResolve_NormalizesOverrideContext(t *testing.T) {
	cfg := &Config{
		Current: "other",
		Contexts: map[string]*Context{
			"local-master": {
				Org:  "ACME",
				Env:  "prod",
				APIM: &ProductConfig{URL: "https://x", Token: "tok"},
			},
			"other": {
				APIM: &ProductConfig{URL: "https://other", Token: "tok2"},
			},
		},
	}

	resolved, err := cfg.Resolve(Overrides{Context: "Local Master"}, "apim")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resolved.Name != "local-master" {
		t.Errorf("expected resolved name 'local-master', got %q", resolved.Name)
	}
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
