package cmd

import (
	"bytes"
	"testing"
)

func TestCompletion(t *testing.T) {
	shells := []string{"bash", "zsh", "fish", "powershell"}

	for _, shell := range shells {
		t.Run(shell, func(t *testing.T) {
			cmd := NewRootCmd("dev")

			var out bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetErr(&out)
			cmd.SetArgs([]string{"completion", shell})

			if err := cmd.Execute(); err != nil {
				t.Fatalf("unexpected error for %s: %v", shell, err)
			}

			if out.Len() == 0 {
				t.Errorf("expected completion output for %s, got empty", shell)
			}
		})
	}
}

func TestCompletion_InvalidShell(t *testing.T) {
	cmd := NewRootCmd("dev")

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"completion", "invalid"})

	if err := cmd.Execute(); err == nil {
		t.Error("expected error for invalid shell")
	}
}
