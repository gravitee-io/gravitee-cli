package am

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func TestLoginWithToken(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.json")
	f := &factory.Factory{
		Config:     &config.Config{Contexts: make(map[string]*config.Context)},
		ConfigPath: cfgPath,
		IOStreams:  factory.IOStreams{Out: &discardWriter{}, Err: &discardWriter{}},
	}
	opts := &loginOptions{
		factory: f, url: "https://am.example.com", token: "my-token",
		contextName: "test-am", org: "DEFAULT", envID: "DEFAULT",
	}
	if err := opts.run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}
	ctx, ok := cfg.Contexts["test-am"]
	if !ok {
		t.Fatal("context 'test-am' not found")
	}
	if ctx.Type != "am" {
		t.Errorf("expected type 'am', got %q", ctx.Type)
	}
	if ctx.AM == nil {
		t.Fatal("expected AM config to be set")
	}
	if ctx.AM.Token != "my-token" {
		t.Errorf("expected token 'my-token', got %q", ctx.AM.Token)
	}
	if ctx.AM.URL != "https://am.example.com" {
		t.Errorf("expected URL, got %q", ctx.AM.URL)
	}
	if cfg.Current != "test-am" {
		t.Errorf("expected current context 'test-am', got %q", cfg.Current)
	}
}

func TestLoginRequiresCredentials(t *testing.T) {
	f := &factory.Factory{
		Config:   &config.Config{Contexts: make(map[string]*config.Context)},
		IOStreams: factory.IOStreams{Out: &discardWriter{}, Err: &discardWriter{}},
	}
	opts := &loginOptions{factory: f, url: "https://am.example.com"}
	err := opts.run()
	if err == nil {
		t.Fatal("expected error when no credentials")
	}
}

func TestDeriveContextName(t *testing.T) {
	tests := []struct {
		url, want string
	}{
		{"https://am.example.com", "am-example-com-am"},
		{"http://localhost:8093", "localhost-am"},
		{"invalid", "am"},
	}
	for _, tt := range tests {
		got := deriveContextName(tt.url)
		if got != tt.want {
			t.Errorf("deriveContextName(%q) = %q, want %q", tt.url, got, tt.want)
		}
	}
}

type discardWriter struct{}

func (d *discardWriter) Write(p []byte) (int, error) { return len(p), nil }
func (d *discardWriter) Read(p []byte) (int, error)  { return 0, nil }
