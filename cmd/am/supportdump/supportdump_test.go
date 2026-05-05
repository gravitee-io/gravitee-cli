package supportdump

import (
	"testing"
)

func TestRedactSecrets(t *testing.T) {
	input := map[string]interface{}{
		"name":         "my-cert",
		"clientSecret": "super-secret",
		"publicKey":    "pk-value",
	}
	result := redactSecrets(input)
	m := result.(map[string]interface{})
	if m["clientSecret"] != "[REDACTED]" {
		t.Errorf("expected clientSecret to be redacted, got %v", m["clientSecret"])
	}
	if m["publicKey"] != "pk-value" {
		t.Errorf("expected publicKey to be preserved, got %v", m["publicKey"])
	}
	if m["name"] != "my-cert" {
		t.Errorf("expected name to be preserved, got %v", m["name"])
	}
}

func TestShouldRedactKey(t *testing.T) {
	cases := []struct {
		key      string
		expected bool
	}{
		{"clientSecret", true},
		{"password", true},
		{"privateKey", true},
		{"apiKey", true},
		{"publicKey", false},
		{"tokenEndpoint", false},
		{"passwordPolicy", false},
		{"name", false},
	}
	for _, tc := range cases {
		got := shouldRedactKey(tc.key)
		if got != tc.expected {
			t.Errorf("shouldRedactKey(%q) = %v, want %v", tc.key, got, tc.expected)
		}
	}
}

func TestRedactSecretsNested(t *testing.T) {
	input := map[string]interface{}{
		"configuration": map[string]interface{}{
			"clientId":     "abc",
			"clientSecret": "secret123",
		},
	}
	result := redactSecrets(input)
	m := result.(map[string]interface{})
	conf := m["configuration"].(map[string]interface{})
	if conf["clientSecret"] != "[REDACTED]" {
		t.Errorf("expected nested secret to be redacted")
	}
	if conf["clientId"] != "abc" {
		t.Errorf("expected clientId to be preserved")
	}
}
