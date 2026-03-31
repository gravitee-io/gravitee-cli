package plugin

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func TestListWithType(t *testing.T) {
	plugins := []map[string]string{
		{"id": "rate-limit", "name": "Rate Limiting", "version": "4.5.0", "description": "Rate limiting policy"},
		{"id": "api-key", "name": "API Key", "version": "4.5.0", "description": "API key validation policy"},
	}

	data, _ := json.Marshal(plugins)

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if path != "/management/v2/organizations/DEFAULT/plugins/policies" {
				t.Errorf("unexpected path: %s", path)
			}

			return data, nil
		},
	}

	f, out := newTestFactory(fake)

	cmd := newListCmd(f)
	cmd.SetArgs([]string{"--type", "policies"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Rate Limiting") {
		t.Errorf("expected 'Rate Limiting' in output, got: %s", output)
	}

	if !strings.Contains(output, "API Key") {
		t.Errorf("expected 'API Key' in output, got: %s", output)
	}

	// TYPE column should not appear when --type is set.
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 0 && strings.Contains(lines[0], "TYPE") {
		t.Errorf("TYPE column should not appear when --type is set, got header: %s", lines[0])
	}
}

func TestListWithoutType(t *testing.T) {
	endpoints := []map[string]string{
		{"id": "kafka", "name": "Kafka", "version": "4.5.0", "description": "Kafka endpoint connector"},
	}

	entrypoints := []map[string]string{
		{"id": "http-proxy", "name": "HTTP Proxy", "version": "4.5.0", "description": "HTTP proxy entrypoint"},
	}

	policies := []map[string]string{
		{"id": "rate-limit", "name": "Rate Limiting", "version": "4.5.0", "description": "Rate limiting policy"},
	}

	endpointsData, _ := json.Marshal(endpoints)
	entrypointsData, _ := json.Marshal(entrypoints)
	policiesData, _ := json.Marshal(policies)

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			switch path {
			case "/management/v2/organizations/DEFAULT/plugins/endpoints":
				return endpointsData, nil
			case "/management/v2/organizations/DEFAULT/plugins/entrypoints":
				return entrypointsData, nil
			case "/management/v2/organizations/DEFAULT/plugins/policies":
				return policiesData, nil
			default:
				t.Errorf("unexpected path: %s", path)

				return nil, nil
			}
		},
	}

	f, out := newTestFactory(fake)

	cmd := newListCmd(f)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()

	// All three types should appear.
	if !strings.Contains(output, "Kafka") {
		t.Errorf("expected 'Kafka' in output, got: %s", output)
	}

	if !strings.Contains(output, "HTTP Proxy") {
		t.Errorf("expected 'HTTP Proxy' in output, got: %s", output)
	}

	if !strings.Contains(output, "Rate Limiting") {
		t.Errorf("expected 'Rate Limiting' in output, got: %s", output)
	}

	// TYPE column should appear when --type is not set.
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 0 && !strings.Contains(lines[0], "TYPE") {
		t.Errorf("TYPE column should appear when --type is not set, got header: %s", lines[0])
	}

	// Check type labels are singular.
	if !strings.Contains(output, "endpoint") {
		t.Errorf("expected 'endpoint' type label in output, got: %s", output)
	}

	if !strings.Contains(output, "entrypoint") {
		t.Errorf("expected 'entrypoint' type label in output, got: %s", output)
	}

	if !strings.Contains(output, "policy") {
		t.Errorf("expected 'policy' type label in output, got: %s", output)
	}
}

func TestListAPIError(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(_ string) ([]byte, error) {
			return nil, &client.APIError{Status: 401, Message: "authentication failed (HTTP 401)"}
		},
	}

	f, _ := newTestFactory(fake)

	cmd := newListCmd(f)
	cmd.SetArgs([]string{"--type", "policies"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "authentication failed") {
		t.Errorf("expected auth error, got: %v", err)
	}
}

func TestListInvalidType(t *testing.T) {
	fake := &client.FakeClient{}

	f, _ := newTestFactory(fake)

	cmd := newListCmd(f)
	cmd.SetArgs([]string{"--type", "connectors"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid type")
	}

	if !strings.Contains(err.Error(), "invalid value 'connectors' for flag --type") {
		t.Errorf("expected invalid type error, got: %v", err)
	}
}
