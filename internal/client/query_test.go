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

package client

import (
	"testing"
)

func TestBuildQuery(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]string
	}{
		{
			name:   "empty params",
			params: map[string]string{},
		},
		{
			name:   "skips empty values",
			params: map[string]string{"page": "1", "query": ""},
		},
		{
			name:   "single param",
			params: map[string]string{"page": "1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildQuery(tt.params)

			for k, v := range tt.params {
				if v == "" {
					if containsStr(result, k+"=") {
						t.Errorf("expected empty param %q to be skipped, got: %s", k, result)
					}
				} else {
					if !containsStr(result, k+"="+v) {
						t.Errorf("expected param %s=%s in query, got: %s", k, v, result)
					}
				}
			}
		})
	}
}

func TestItoa(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{0, "0"},
		{1, "1"},
		{42, "42"},
		{-1, "-1"},
	}

	for _, tt := range tests {
		got := Itoa(tt.input)
		if got != tt.want {
			t.Errorf("Itoa(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
