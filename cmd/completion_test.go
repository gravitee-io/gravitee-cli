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
