package am

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func TestSetDomain(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.json")
	cfg := &config.Config{
		CurrentContext: "am-test",
		Contexts: map[string]config.Context{
			"am-test": {URL: "https://am.example.com", Token: "tok", Type: "am"},
		},
	}
	f := &factory.Factory{
		Config: cfg, ConfigPath: cfgPath,
		Resolved: &config.ResolvedContext{
			Name: "am-test", URL: "https://am.example.com", Token: "tok",
			Org: "DEFAULT", Env: "DEFAULT", Type: "am",
		},
		IOStreams: factory.IOStreams{Out: &discardWriter{}, Err: &discardWriter{}},
	}
	opts := &setDomainOptions{factory: f, domainID: "my-domain-123"}
	if err := opts.run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(cfgPath)
	var saved config.Config
	_ = json.Unmarshal(data, &saved)
	ctx := saved.Contexts["am-test"]
	if ctx.Domain != "my-domain-123" {
		t.Errorf("expected domain 'my-domain-123', got %q", ctx.Domain)
	}
}

func TestSetDomainClear(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.json")
	cfg := &config.Config{
		CurrentContext: "am-test",
		Contexts: map[string]config.Context{
			"am-test": {URL: "https://am.example.com", Token: "tok", Type: "am", Domain: "old-domain"},
		},
	}
	f := &factory.Factory{
		Config: cfg, ConfigPath: cfgPath,
		Resolved: &config.ResolvedContext{
			Name: "am-test", URL: "https://am.example.com", Token: "tok",
			Org: "DEFAULT", Env: "DEFAULT", Type: "am", Domain: "old-domain",
		},
		IOStreams: factory.IOStreams{Out: &discardWriter{}, Err: &discardWriter{}},
	}
	opts := &setDomainOptions{factory: f, clear: true}
	if err := opts.run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(cfgPath)
	var saved config.Config
	_ = json.Unmarshal(data, &saved)
	if saved.Contexts["am-test"].Domain != "" {
		t.Errorf("expected empty domain after clear, got %q", saved.Contexts["am-test"].Domain)
	}
}
