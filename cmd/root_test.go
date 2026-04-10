package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootHelp(t *testing.T) {
	cmd := NewRootCmd("0.1.0-test")

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Gravitee platform") {
		t.Errorf("expected help to contain 'Gravitee platform', got: %s", output)
	}
}

func TestVersionFlag(t *testing.T) {
	cmd := NewRootCmd("1.2.3")

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--version"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := strings.TrimSpace(out.String())
	if !strings.Contains(output, "1.2.3") {
		t.Errorf("expected version '1.2.3' in output, got %q", output)
	}
}

func TestGlobalFlags(t *testing.T) {
	cmd := NewRootCmd("dev")

	flags := cmd.PersistentFlags()

	tests := []struct {
		name string
	}{
		{"context"},
		{"org"},
		{"env"},
		{"debug"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if flags.Lookup(tt.name) == nil {
				t.Errorf("expected global flag %q to be registered", tt.name)
			}
		})
	}
}
