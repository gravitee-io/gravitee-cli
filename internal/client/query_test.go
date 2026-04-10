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
