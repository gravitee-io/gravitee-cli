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

package cmdutil

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestAddAPIFlag_RejectsEmptyValue(t *testing.T) {
	ran := false

	var target string

	cmd := &cobra.Command{
		Use: "test",
		RunE: func(_ *cobra.Command, _ []string) error {
			ran = true

			return nil
		},
	}
	AddAPIFlag(cmd, &target)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"--api", ""})

	err := cmd.Execute()

	if err == nil {
		t.Fatal("expected error for empty --api, got nil")
	}

	if !contains(err.Error(), "--api cannot be empty") {
		t.Errorf("expected error mentioning empty --api, got: %s", err.Error())
	}

	if ran {
		t.Fatal("RunE should not execute when --api is empty")
	}
}

func TestAddAPIFlag_AcceptsNonEmptyValue(t *testing.T) {
	ran := false

	var target string

	cmd := &cobra.Command{
		Use: "test",
		RunE: func(_ *cobra.Command, _ []string) error {
			ran = true

			return nil
		},
	}
	AddAPIFlag(cmd, &target)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs([]string{"--api", "/my/api"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ran {
		t.Fatal("expected RunE to execute with valid --api")
	}

	if target != "/my/api" {
		t.Errorf("expected target=/my/api, got %q", target)
	}
}
