package am

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func TestSetDomain(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.json")
	cfg := &config.Config{
		Current: "am-test",
		Contexts: map[string]*config.Context{
			"am-test": {Type: "am", AM: &config.ProductConfig{URL: "https://am.example.com", Token: "tok"}},
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
	_ = yaml.Unmarshal(data, &saved)
	ctx := saved.Contexts["am-test"]
	if ctx.Domain != "my-domain-123" {
		t.Errorf("expected domain 'my-domain-123', got %q", ctx.Domain)
	}
}

func TestSetDomainClear(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.json")
	cfg := &config.Config{
		Current: "am-test",
		Contexts: map[string]*config.Context{
			"am-test": {Type: "am", Domain: "old-domain", AM: &config.ProductConfig{URL: "https://am.example.com", Token: "tok"}},
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
	_ = yaml.Unmarshal(data, &saved)
	if saved.Contexts["am-test"].Domain != "" {
		t.Errorf("expected empty domain after clear, got %q", saved.Contexts["am-test"].Domain)
	}
}
