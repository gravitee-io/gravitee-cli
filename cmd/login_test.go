package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func setupInteractiveTest(t *testing.T, cfg *config.Config, input string) (*factory.Factory, string, *bytes.Buffer) {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	if cfg == nil {
		cfg = &config.Config{}
	}

	if cfg.Contexts == nil {
		cfg.Contexts = map[string]*config.Context{}
	}

	if err := cfg.SaveTo(path); err != nil {
		t.Fatalf("failed to seed config: %v", err)
	}

	out := &bytes.Buffer{}
	f := &factory.Factory{
		Config:     cfg,
		ConfigPath: path,
		IOStreams: factory.IOStreams{
			In:  strings.NewReader(input),
			Out: out,
			Err: &bytes.Buffer{},
		},
	}

	return f, path, out
}

func TestInteractive_CurlPaste_NewContext(t *testing.T) {
	input := "master\n" +
		`curl -H "Authorization: Bearer tok_abc" "https://apim.example.com/management/organizations/ACME/environments/prod"` + "\n"

	f, _, _ := setupInteractiveTest(t, &config.Config{}, input)

	if err := runInteractiveLogin(f, "apim"); err != nil {
		t.Fatalf("runInteractiveLogin: %v", err)
	}

	ctx := f.Config.Contexts["master"]
	if ctx == nil {
		t.Fatal("expected context 'master' to be created")
	}

	if ctx.Org != "ACME" {
		t.Errorf("Org: got %q, want ACME", ctx.Org)
	}

	if ctx.Env != "prod" {
		t.Errorf("Env: got %q, want prod", ctx.Env)
	}

	if ctx.APIM == nil || ctx.APIM.URL != "https://apim.example.com" || ctx.APIM.Token != "tok_abc" {
		t.Errorf("APIM config: got %+v, want URL=https://apim.example.com Token=tok_abc", ctx.APIM)
	}

	if f.Config.Current != "master" {
		t.Errorf("Current: got %q, want master", f.Config.Current)
	}
}

func TestInteractive_CurlPaste_OverwritesExistingContext(t *testing.T) {
	seed := &config.Config{
		Current: "prod",
		Contexts: map[string]*config.Context{
			"prod": {Org: "OLD_ORG", Env: "OLD_ENV"},
		},
	}

	input := "prod\n" +
		`curl -H "Authorization: Bearer new_tok" "https://new.example.com/management/organizations/NEW_ORG/environments/NEW_ENV"` + "\n"

	f, _, _ := setupInteractiveTest(t, seed, input)

	if err := runInteractiveLogin(f, "apim"); err != nil {
		t.Fatalf("runInteractiveLogin: %v", err)
	}

	ctx := f.Config.Contexts["prod"]

	if ctx.Org != "NEW_ORG" {
		t.Errorf("Org should be overwritten by curl: got %q, want NEW_ORG", ctx.Org)
	}

	if ctx.Env != "NEW_ENV" {
		t.Errorf("Env should be overwritten by curl: got %q, want NEW_ENV", ctx.Env)
	}
}

func TestInteractive_URLBare_NewContext_ExplicitOrgEnv(t *testing.T) {
	input := "dev\nhttps://apim.example.com\ntok_xyz\nMYORG\nMYENV\n"

	f, _, _ := setupInteractiveTest(t, &config.Config{}, input)

	if err := runInteractiveLogin(f, "apim"); err != nil {
		t.Fatalf("runInteractiveLogin: %v", err)
	}

	ctx := f.Config.Contexts["dev"]
	if ctx == nil {
		t.Fatal("expected context 'dev' to be created")
	}

	if ctx.Org != "MYORG" || ctx.Env != "MYENV" {
		t.Errorf("expected Org=MYORG Env=MYENV, got Org=%q Env=%q", ctx.Org, ctx.Env)
	}
}

func TestInteractive_URLBare_NewContext_EnterDefaults(t *testing.T) {
	input := "dev\nhttps://apim.example.com\ntok_xyz\n\n\n"

	f, _, _ := setupInteractiveTest(t, &config.Config{}, input)

	if err := runInteractiveLogin(f, "apim"); err != nil {
		t.Fatalf("runInteractiveLogin: %v", err)
	}

	ctx := f.Config.Contexts["dev"]
	if ctx.Org != config.DefaultOrg {
		t.Errorf("Org: got %q, want %s", ctx.Org, config.DefaultOrg)
	}

	if ctx.Env != config.DefaultEnv {
		t.Errorf("Env: got %q, want %s", ctx.Env, config.DefaultEnv)
	}
}

// URL bare + Enter on org/env while reusing an existing context must preserve
// the context's existing Org/Env (don't downgrade to DEFAULT).
func TestInteractive_URLBare_ReuseContext_EnterPreservesExistingOrgEnv(t *testing.T) {
	seed := &config.Config{
		Current: "prod",
		Contexts: map[string]*config.Context{
			"prod": {Org: "ACME", Env: "production"},
		},
	}

	input := "prod\nhttps://new.example.com\ntok_new\n\n\n"

	f, _, _ := setupInteractiveTest(t, seed, input)

	if err := runInteractiveLogin(f, "apim"); err != nil {
		t.Fatalf("runInteractiveLogin: %v", err)
	}

	ctx := f.Config.Contexts["prod"]

	if ctx.Org != "ACME" {
		t.Errorf("Org must be preserved on Enter: got %q, want ACME", ctx.Org)
	}

	if ctx.Env != "production" {
		t.Errorf("Env must be preserved on Enter: got %q, want production", ctx.Env)
	}

	// Token/URL should still update.
	if ctx.APIM == nil || ctx.APIM.Token != "tok_new" {
		t.Errorf("token should update, got %+v", ctx.APIM)
	}
}

// When reusing an existing context via URL-bare input, all three prompts
// (context name / org / env) must advertise the existing values as defaults
// so Enter preserves them instead of silently switching to DEFAULT.
func TestInteractive_ReuseContext_PromptsShowExistingValuesAsDefaults(t *testing.T) {
	seed := &config.Config{
		Current: "prod",
		Contexts: map[string]*config.Context{
			"prod": {Org: "ACME", Env: "production"},
		},
	}

	input := "\nhttps://x.example.com\ntok\n\n\n"

	f, _, out := setupInteractiveTest(t, seed, input)

	if err := runInteractiveLogin(f, "apim"); err != nil {
		t.Fatalf("runInteractiveLogin: %v", err)
	}

	if _, exists := f.Config.Contexts["default"]; exists {
		t.Error("Enter on context prompt must not create a 'default' context when a current context exists")
	}

	if f.Config.Current != "prod" {
		t.Errorf("Current: got %q, want prod", f.Config.Current)
	}

	for _, expected := range []string{
		"Context name (prod):",
		"Organization ID (ACME):",
		"Environment ID (production):",
	} {
		if !strings.Contains(out.String(), expected) {
			t.Errorf("prompt output missing %q, got:\n%s", expected, out.String())
		}
	}
}

func TestInteractive_ContextPrompt_DefaultsToDefault_WhenNoCurrent(t *testing.T) {
	input := "\nhttps://x.example.com\ntok\n\n\n"

	f, _, out := setupInteractiveTest(t, &config.Config{}, input)

	if err := runInteractiveLogin(f, "apim"); err != nil {
		t.Fatalf("runInteractiveLogin: %v", err)
	}

	if _, exists := f.Config.Contexts["default"]; !exists {
		t.Error("Enter on empty-current config should create a 'default' context")
	}

	if !strings.Contains(out.String(), "Context name (default):") {
		t.Errorf("prompt should advertise 'default' as default, got:\n%s", out.String())
	}
}

// User typing "Local Master" at the context prompt must be stored as
// "local-master" - no spaces, lowercase - so the context is easy to reuse.
func TestInteractive_ContextPrompt_NormalizesName(t *testing.T) {
	input := "Local Master\n" +
		`curl -H "Authorization: Bearer tok" "https://apim.example.com/management/organizations/ACME/environments/prod"` + "\n"

	f, _, _ := setupInteractiveTest(t, &config.Config{}, input)

	if err := runInteractiveLogin(f, "apim"); err != nil {
		t.Fatalf("runInteractiveLogin: %v", err)
	}

	if _, ok := f.Config.Contexts["local-master"]; !ok {
		t.Errorf("expected normalized context 'local-master', got keys: %v", keys(f.Config.Contexts))
	}

	if _, ok := f.Config.Contexts["Local Master"]; ok {
		t.Error("raw 'Local Master' key should not be stored")
	}

	if f.Config.Current != "local-master" {
		t.Errorf("Current: got %q, want local-master", f.Config.Current)
	}
}

func keys(m map[string]*config.Context) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}

	return out
}
