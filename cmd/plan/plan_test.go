package plan

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func planJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":         "plan-1",
		"name":       "Gold Plan",
		"apiId":      "api-1",
		"status":     "PUBLISHED",
		"security":   map[string]string{"type": "API_KEY"},
		"validation": "AUTO",
		"mode":       "STANDARD",
		"createdAt":  "2026-03-15T10:00:00Z",
		"updatedAt":  "2026-03-25T14:30:00Z",
	}
}

func TestListSuccess(t *testing.T) {
	resp := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"id": "plan-1", "name": "Gold Plan", "status": "PUBLISHED",
				"security": map[string]string{"type": "API_KEY"}, "validation": "AUTO",
				"updatedAt": "2026-03-25T14:30:00Z",
			},
		},
		"pagination": map[string]int{
			"page": 1, "perPage": 10, "pageCount": 1, "totalCount": 1, "pageItemsCount": 1,
		},
	}

	data, _ := json.Marshal(resp)

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/plans?") {
				t.Errorf("unexpected path: %s", path)
			}

			return data, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newListCmd(f)
	cmd.SetArgs([]string{"--api", "api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Gold Plan") {
		t.Errorf("expected 'Gold Plan' in output, got: %s", out.String())
	}

	if !strings.Contains(out.String(), "API_KEY") {
		t.Errorf("expected 'API_KEY' in output, got: %s", out.String())
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
	cmd.SetArgs([]string{"--api", "api-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "authentication failed") {
		t.Errorf("expected auth error, got: %v", err)
	}
}

func TestGetSuccess(t *testing.T) {
	resp, _ := json.Marshal(planJSON())

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/plans/plan-1") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newGetCmd(f)
	cmd.SetArgs([]string{"plan-1", "--api", "api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Gold Plan") {
		t.Errorf("expected 'Gold Plan' in output, got: %s", output)
	}

	if !strings.Contains(output, "API_KEY") {
		t.Errorf("expected 'API_KEY' in output, got: %s", output)
	}

	if !strings.Contains(output, "PUBLISHED") {
		t.Errorf("expected 'PUBLISHED' in output, got: %s", output)
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
	cmd.SetArgs([]string{"plan-999", "--api", "api-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestCreateSuccess(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "plan.json")
	_ = os.WriteFile(file, []byte(`{"name":"Gold Plan","security":{"type":"API_KEY"}}`), 0600)

	resp, _ := json.Marshal(planJSON())

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.HasSuffix(path, "/apis/api-1/plans") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newCreateCmd(f)
	cmd.SetArgs([]string{"--api", "api-1", "-f", file})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Gold Plan") {
		t.Errorf("expected 'Gold Plan' in output, got: %s", out.String())
	}
}

func TestCreateReadOnly(t *testing.T) {
	fake := &client.FakeClient{}
	f, _ := newTestFactory(fake, true)

	cmd := newCreateCmd(f)
	cmd.SetArgs([]string{"--api", "api-1", "-f", "test.json"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}

func TestUpdateSuccess(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "plan.json")
	_ = os.WriteFile(file, []byte(`{"name":"Gold Plan v2"}`), 0600)

	resp, _ := json.Marshal(map[string]interface{}{
		"id": "plan-1", "name": "Gold Plan v2", "apiId": "api-1",
		"status": "PUBLISHED", "security": map[string]string{"type": "API_KEY"},
	})

	fake := &client.FakeClient{
		PutFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/plans/plan-1") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newUpdateCmd(f)
	cmd.SetArgs([]string{"plan-1", "--api", "api-1", "-f", file})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Gold Plan v2") {
		t.Errorf("expected 'Gold Plan v2' in output, got: %s", out.String())
	}
}

func TestUpdateReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)

	cmd := newUpdateCmd(f)
	cmd.SetArgs([]string{"plan-1", "--api", "api-1", "-f", "test.json"})
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
			if !strings.Contains(path, "/apis/api-1/plans/plan-1") {
				t.Errorf("unexpected path: %s", path)
			}

			return nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newDeleteCmd(f)
	cmd.SetArgs([]string{"plan-1", "--api", "api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Plan 'plan-1' deleted.") {
		t.Errorf("unexpected output: %s", out.String())
	}
}

func TestDeleteReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)

	cmd := newDeleteCmd(f)
	cmd.SetArgs([]string{"plan-1", "--api", "api-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}

func TestPublishSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]interface{}{
		"id": "plan-1", "name": "Gold Plan", "apiId": "api-1",
		"status": "PUBLISHED", "security": map[string]string{"type": "API_KEY"},
		"validation": "AUTO", "mode": "STANDARD",
	})

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/plans/plan-1/_publish") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newPublishCmd(f)
	cmd.SetArgs([]string{"plan-1", "--api", "api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "PUBLISHED") {
		t.Errorf("expected 'PUBLISHED' in output, got: %s", output)
	}

	if !strings.Contains(output, "Gold Plan") {
		t.Errorf("expected 'Gold Plan' in output, got: %s", output)
	}
}

func TestPublishReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)

	cmd := newPublishCmd(f)
	cmd.SetArgs([]string{"plan-1", "--api", "api-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}

func TestPublishAPIError(t *testing.T) {
	fake := &client.FakeClient{
		PostFunc: func(_ string, _ interface{}) ([]byte, error) {
			return nil, &client.APIError{Status: 400, Message: "invalid request (HTTP 400): plan is already published"}
		},
	}

	f, _ := newTestFactory(fake, false)

	cmd := newPublishCmd(f)
	cmd.SetArgs([]string{"plan-1", "--api", "api-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "already published") {
		t.Errorf("expected publish error, got: %v", err)
	}
}

func TestDeprecateSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]interface{}{
		"id": "plan-1", "name": "Gold Plan", "apiId": "api-1",
		"status": "DEPRECATED", "security": map[string]string{"type": "API_KEY"},
		"validation": "AUTO", "mode": "STANDARD",
	})

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/plans/plan-1/_deprecate") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newDeprecateCmd(f)
	cmd.SetArgs([]string{"plan-1", "--api", "api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "DEPRECATED") {
		t.Errorf("expected 'DEPRECATED' in output, got: %s", output)
	}
}

func TestDeprecateReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)

	cmd := newDeprecateCmd(f)
	cmd.SetArgs([]string{"plan-1", "--api", "api-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}

func TestCloseSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]interface{}{
		"id": "plan-1", "name": "Gold Plan", "apiId": "api-1",
		"status": "CLOSED", "security": map[string]string{"type": "API_KEY"},
		"validation": "AUTO", "mode": "STANDARD",
		"closedAt": "2026-03-27T15:00:00Z",
	})

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/plans/plan-1/_close") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newCloseCmd(f)
	cmd.SetArgs([]string{"plan-1", "--api", "api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "CLOSED") {
		t.Errorf("expected 'CLOSED' in output, got: %s", output)
	}

	if !strings.Contains(output, "2026-03-27T15:00:00Z") {
		t.Errorf("expected closedAt in output, got: %s", output)
	}
}

func TestCloseReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)

	cmd := newCloseCmd(f)
	cmd.SetArgs([]string{"plan-1", "--api", "api-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}
