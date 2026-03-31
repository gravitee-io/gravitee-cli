package api

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func TestGetSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]string{
		"id": "api-1", "name": "Weather API", "state": "STARTED",
		"definitionVersion": "V4", "type": "PROXY",
	})

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newGetCmd(f)
	cmd.SetArgs([]string{"api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Weather API") {
		t.Errorf("expected 'Weather API' in output, got: %s", out.String())
	}
}

func TestGetNotFound(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(_ string) ([]byte, error) {
			return nil, &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
		},
	}

	f, _ := newTestFactory(fake, false)

	cmd := newGetCmd(f)
	cmd.SetArgs([]string{"api-999"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestGetJSON(t *testing.T) {
	resp, _ := json.Marshal(map[string]string{"id": "api-1", "state": "STARTED"})

	fake := &client.FakeClient{
		GetFunc: func(_ string) ([]byte, error) { return resp, nil },
	}

	f, out := newTestFactory(fake, false)
	f.OutputFormat = "json"

	cmd := newGetCmd(f)
	cmd.SetArgs([]string{"api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), `"state"`) {
		t.Errorf("expected JSON state field, got: %s", out.String())
	}
}
