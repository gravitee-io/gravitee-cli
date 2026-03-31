package environment

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func TestListSuccess(t *testing.T) {
	envs := []map[string]string{
		{"id": "dev-1111", "name": "Development", "description": "Development environment"},
		{"id": "prod-2222", "name": "Production", "description": "Production environment"},
	}

	data, _ := json.Marshal(envs)

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if path != "/management/organizations/DEFAULT/environments" {
				t.Errorf("unexpected path: %s", path)
			}

			return data, nil
		},
	}

	f, out := newTestFactory(fake)

	cmd := newListCmd(f)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()

	if !strings.Contains(output, "Development") {
		t.Errorf("expected 'Development' in output, got: %s", output)
	}

	if !strings.Contains(output, "Production") {
		t.Errorf("expected 'Production' in output, got: %s", output)
	}
}

func TestListError(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(_ string) ([]byte, error) {
			return nil, &client.APIError{Status: 401, Message: "authentication failed (HTTP 401)"}
		},
	}

	f, _ := newTestFactory(fake)

	cmd := newListCmd(f)
	cmd.SetArgs([]string{})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "authentication failed") {
		t.Errorf("expected auth error, got: %v", err)
	}
}

func TestGetSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]string{
		"id": "prod-2222", "name": "Production", "description": "Production environment",
	})

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if path != "/management/organizations/DEFAULT/environments/prod-2222" {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake)

	cmd := newGetCmd(f)
	cmd.SetArgs([]string{"prod-2222"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()

	if !strings.Contains(output, "Production") {
		t.Errorf("expected 'Production' in output, got: %s", output)
	}

	if !strings.Contains(output, "prod-2222") {
		t.Errorf("expected 'prod-2222' in output, got: %s", output)
	}
}

func TestGetNotFound(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(_ string) ([]byte, error) {
			return nil, &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
		},
	}

	f, _ := newTestFactory(fake)

	cmd := newGetCmd(f)
	cmd.SetArgs([]string{"env-999"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected not found error, got: %v", err)
	}
}
