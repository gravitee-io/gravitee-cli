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

package shell

import (
	"testing"
)

func TestSplitArgs(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{"domain list", []string{"domain", "list"}},
		{`role get "my role"`, []string{"role", "get", "my role"}},
		{"user list --page 1", []string{"user", "list", "--page", "1"}},
		{"", []string{}},
	}
	for _, tc := range cases {
		got := splitArgs(tc.input)
		if len(got) != len(tc.expected) {
			t.Errorf("splitArgs(%q): expected %v, got %v", tc.input, tc.expected, got)
			continue
		}
		for i, v := range got {
			if v != tc.expected[i] {
				t.Errorf("splitArgs(%q)[%d]: expected %q, got %q", tc.input, i, tc.expected[i], v)
			}
		}
	}
}

func TestBuildPromptNoContext(t *testing.T) {
	p := buildPrompt("", "")
	expected := "[not-configured] am> "
	if p != expected {
		t.Errorf("buildPrompt(\"\", \"\"): expected %q, got %q", expected, p)
	}
}

func TestBuildPromptWithContext(t *testing.T) {
	p := buildPrompt("myws", "dom1")
	expected := "[myws:dom1] am> "
	if p != expected {
		t.Errorf("buildPrompt(\"myws\", \"dom1\"): expected %q, got %q", expected, p)
	}
}
