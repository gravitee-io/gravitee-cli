package api

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func TestListSuccess(t *testing.T) {
	resp := map[string]interface{}{
		"data": []map[string]string{
			{"id": "api-1", "name": "Weather API", "state": "STARTED", "definitionVersion": "V4"},
		},
		"pagination": map[string]int{"page": 1, "perPage": 10, "pageCount": 1, "totalCount": 1, "pageItemsCount": 1},
	}

	data, _ := json.Marshal(resp)

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/apis?") {
				t.Errorf("unexpected path: %s", path)
			}

			return data, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newListCmd(f)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Weather API") {
		t.Errorf("expected 'Weather API' in output, got: %s", out.String())
	}
}

func TestListAPIError(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(_ string) ([]byte, error) {
			return nil, &client.APIError{Status: 401, Message: "authentication failed (HTTP 401)"}
		},
	}

	f, _ := newTestFactory(fake, false)

	cmd := newListCmd(f)
	cmd.SetArgs([]string{})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "authentication failed") {
		t.Errorf("expected auth error, got: %v", err)
	}
}
