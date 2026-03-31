package page

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func pageJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":         "page-1",
		"name":       "Getting Started",
		"apiId":      "api-1",
		"type":       "MARKDOWN",
		"visibility": "PUBLIC",
		"published":  true,
		"parentId":   "folder-1",
		"createdAt":  "2026-03-15T10:00:00Z",
		"updatedAt":  "2026-03-25T14:30:00Z",
	}
}

func TestListSuccess(t *testing.T) {
	resp := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"id": "page-1", "name": "Getting Started", "type": "MARKDOWN",
				"visibility": "PUBLIC", "published": true,
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
			if !strings.Contains(path, "/apis/api-1/pages?") {
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

	if !strings.Contains(out.String(), "Getting Started") {
		t.Errorf("expected 'Getting Started' in output, got: %s", out.String())
	}

	if !strings.Contains(out.String(), "MARKDOWN") {
		t.Errorf("expected 'MARKDOWN' in output, got: %s", out.String())
	}
}

func TestGetSuccess(t *testing.T) {
	resp, _ := json.Marshal(pageJSON())

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/pages/page-1") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newGetCmd(f)
	cmd.SetArgs([]string{"page-1", "--api", "api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Getting Started") {
		t.Errorf("expected 'Getting Started' in output, got: %s", output)
	}

	if !strings.Contains(output, "MARKDOWN") {
		t.Errorf("expected 'MARKDOWN' in output, got: %s", output)
	}

	if !strings.Contains(output, "true") {
		t.Errorf("expected 'true' (published) in output, got: %s", output)
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
	cmd.SetArgs([]string{"page-999", "--api", "api-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestCreateSuccess(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "page.json")
	_ = os.WriteFile(file, []byte(`{"name":"Getting Started","type":"MARKDOWN"}`), 0600)

	resp, _ := json.Marshal(pageJSON())

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.HasSuffix(path, "/apis/api-1/pages") {
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

	if !strings.Contains(out.String(), "Getting Started") {
		t.Errorf("expected 'Getting Started' in output, got: %s", out.String())
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

func TestDeleteSuccess(t *testing.T) {
	fake := &client.FakeClient{
		DeleteFunc: func(path string) error {
			if !strings.Contains(path, "/apis/api-1/pages/page-1") {
				t.Errorf("unexpected path: %s", path)
			}

			return nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newDeleteCmd(f)
	cmd.SetArgs([]string{"page-1", "--api", "api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Page 'page-1' deleted.") {
		t.Errorf("unexpected output: %s", out.String())
	}
}

func TestDeleteReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)

	cmd := newDeleteCmd(f)
	cmd.SetArgs([]string{"page-1", "--api", "api-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}

func TestPublishSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]interface{}{
		"id": "page-1", "name": "Getting Started", "apiId": "api-1",
		"type": "MARKDOWN", "visibility": "PUBLIC", "published": true,
		"createdAt": "2026-03-15T10:00:00Z", "updatedAt": "2026-03-25T14:30:00Z",
	})

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/pages/page-1/_publish") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newPublishCmd(f)
	cmd.SetArgs([]string{"page-1", "--api", "api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "true") {
		t.Errorf("expected 'true' (published) in output, got: %s", output)
	}

	if !strings.Contains(output, "Getting Started") {
		t.Errorf("expected 'Getting Started' in output, got: %s", output)
	}
}

func TestPublishReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)

	cmd := newPublishCmd(f)
	cmd.SetArgs([]string{"page-1", "--api", "api-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}

func TestUnpublishSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]interface{}{
		"id": "page-1", "name": "Getting Started", "apiId": "api-1",
		"type": "MARKDOWN", "visibility": "PUBLIC", "published": false,
		"createdAt": "2026-03-15T10:00:00Z", "updatedAt": "2026-03-30T09:00:00Z",
	})

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/pages/page-1/_unpublish") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newUnpublishCmd(f)
	cmd.SetArgs([]string{"page-1", "--api", "api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "false") {
		t.Errorf("expected 'false' (published) in output, got: %s", output)
	}

	if !strings.Contains(output, "Getting Started") {
		t.Errorf("expected 'Getting Started' in output, got: %s", output)
	}
}

func TestUnpublishReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)

	cmd := newUnpublishCmd(f)
	cmd.SetArgs([]string{"page-1", "--api", "api-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}
