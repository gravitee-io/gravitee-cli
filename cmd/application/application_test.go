package application

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func appJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":           "app-1",
		"name":         "My Mobile App",
		"description":  "Mobile client for the Weather API",
		"type":         "SIMPLE",
		"status":       "ACTIVE",
		"owner":        map[string]interface{}{"id": "user-1234", "displayName": "john.doe"},
		"api_key_mode": "UNSPECIFIED",
		"domain":       "https://my-app.com",
		"created_at":   "2026-03-15T10:00:00Z",
		"updated_at":   "2026-03-25T14:30:00Z",
	}
}

func TestListSuccess(t *testing.T) {
	resp := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"id": "app-1", "name": "My Mobile App", "type": "SIMPLE",
				"status": "ACTIVE", "owner": map[string]interface{}{"displayName": "john.doe"},
				"updated_at": "2026-03-25T14:30:00Z",
			},
		},
		"pagination": map[string]int{
			"page": 1, "perPage": 10, "pageCount": 1, "totalCount": 1, "pageItemsCount": 1,
		},
	}

	data, _ := json.Marshal(resp)

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/applications/_paged?") {
				t.Errorf("unexpected path: %s", path)
			}

			if !strings.Contains(path, "organizations/DEFAULT/environments/DEFAULT") {
				t.Errorf("expected V1 path with org/env, got: %s", path)
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

	output := out.String()
	if !strings.Contains(output, "My Mobile App") {
		t.Errorf("expected 'My Mobile App' in output, got: %s", output)
	}

	if !strings.Contains(output, "john.doe") {
		t.Errorf("expected 'john.doe' in output, got: %s", output)
	}
}

func TestGetSuccess(t *testing.T) {
	resp, _ := json.Marshal(appJSON())

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/applications/app-1") {
				t.Errorf("unexpected path: %s", path)
			}

			if !strings.Contains(path, "organizations/DEFAULT/environments/DEFAULT") {
				t.Errorf("expected V1 path with org/env, got: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newGetCmd(f)
	cmd.SetArgs([]string{"app-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "My Mobile App") {
		t.Errorf("expected 'My Mobile App' in output, got: %s", output)
	}

	if !strings.Contains(output, "john.doe") {
		t.Errorf("expected 'john.doe' in output, got: %s", output)
	}

	if !strings.Contains(output, "ACTIVE") {
		t.Errorf("expected 'ACTIVE' in output, got: %s", output)
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
	cmd.SetArgs([]string{"app-999"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestCreateSuccess(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "app.json")
	_ = os.WriteFile(file, []byte(`{"name":"My Mobile App","description":"Mobile client"}`), 0600)

	resp, _ := json.Marshal(appJSON())

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.HasSuffix(path, "/applications") {
				t.Errorf("unexpected path: %s", path)
			}

			if !strings.Contains(path, "organizations/DEFAULT/environments/DEFAULT") {
				t.Errorf("expected V1 path with org/env, got: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newCreateCmd(f)
	cmd.SetArgs([]string{"-f", file})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "My Mobile App") {
		t.Errorf("expected 'My Mobile App' in output, got: %s", out.String())
	}
}

func TestCreateReadOnly(t *testing.T) {
	fake := &client.FakeClient{}
	f, _ := newTestFactory(fake, true)

	cmd := newCreateCmd(f)
	cmd.SetArgs([]string{"-f", "test.json"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}

func TestDeleteSuccess(t *testing.T) {
	fake := &client.FakeClient{
		DeleteFunc: func(path string) error {
			if !strings.Contains(path, "/applications/app-1") {
				t.Errorf("unexpected path: %s", path)
			}

			if !strings.Contains(path, "organizations/DEFAULT/environments/DEFAULT") {
				t.Errorf("expected V1 path with org/env, got: %s", path)
			}

			return nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newDeleteCmd(f)
	cmd.SetArgs([]string{"app-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Application 'app-1' deleted.") {
		t.Errorf("unexpected output: %s", out.String())
	}
}

func TestDeleteReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)

	cmd := newDeleteCmd(f)
	cmd.SetArgs([]string{"app-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}
