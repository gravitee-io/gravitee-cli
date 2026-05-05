package supportdump

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
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

func TestSupportDumpSingleDomain(t *testing.T) {
	domain := map[string]interface{}{"id": "test-domain", "name": "Test Domain", "enabled": true}
	emptyPaginated := map[string]interface{}{"data": []interface{}{}, "totalCount": 0}
	emptyArr := []interface{}{}

	domainBytes, _ := json.Marshal(domain)
	emptyPaginatedBytes, _ := json.Marshal(emptyPaginated)
	emptyArrBytes, _ := json.Marshal(emptyArr)

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			base := strings.Split(path, "?")[0]
			trimmed := strings.TrimRight(base, "/")
			switch {
			case strings.HasSuffix(trimmed, "/domains/test-domain"):
				return domainBytes, nil
			case strings.Contains(path, "/applications"),
				strings.Contains(path, "/roles"),
				strings.Contains(path, "/scopes"),
				strings.Contains(path, "/groups"),
				strings.Contains(path, "/audits"):
				return emptyPaginatedBytes, nil
			case strings.Contains(path, "/identities"),
				strings.Contains(path, "/certificates"),
				strings.Contains(path, "/flows"),
				strings.Contains(path, "/factors"),
				strings.Contains(path, "/members"):
				return emptyArrBytes, nil
			default:
				return nil, fmt.Errorf("unexpected path in fake client: %s", path)
			}
		},
	}
	f, out := newTestFactory(fake)
	cmd := NewSupportDumpCmd(f)
	cmd.SetArgs([]string{"--no-redact"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "_metadata") {
		t.Errorf("expected '_metadata' in JSON output, got: %s", output)
	}
	if !strings.Contains(output, "test-domain") {
		t.Errorf("expected domain ID in JSON output, got: %s", output)
	}
	if !strings.Contains(output, "Test Domain") {
		t.Errorf("expected domain name in JSON output, got: %s", output)
	}
	if strings.Contains(output, `"_errors"`) {
		t.Errorf("unexpected _errors in output: %s", output)
	}
}
