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
