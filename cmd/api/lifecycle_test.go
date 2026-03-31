package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func TestDeleteSuccess(t *testing.T) {
	fake := &client.FakeClient{
		DeleteFunc: func(path string) error {
			if !strings.Contains(path, "/apis/api-1") {
				t.Errorf("unexpected path: %s", path)
			}

			return nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newDeleteCmd(f)
	cmd.SetArgs([]string{"api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "API 'api-1' deleted.") {
		t.Errorf("unexpected output: %s", out.String())
	}
}

func TestDeleteReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)

	cmd := newDeleteCmd(f)
	cmd.SetArgs([]string{"api-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}

func TestDeleteClosePlans(t *testing.T) {
	fake := &client.FakeClient{
		DeleteFunc: func(path string) error {
			if !strings.Contains(path, "closePlans=true") {
				t.Errorf("expected closePlans param, got: %s", path)
			}

			return nil
		},
	}

	f, _ := newTestFactory(fake, false)

	cmd := newDeleteCmd(f)
	cmd.SetArgs([]string{"api-1", "--close-plans"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStartSuccess(t *testing.T) {
	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/_start") {
				t.Errorf("unexpected path: %s", path)
			}

			return nil, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newStartCmd(f)
	cmd.SetArgs([]string{"api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "API 'api-1' started.") {
		t.Errorf("unexpected output: %s", out.String())
	}
}

func TestStartReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)

	cmd := newStartCmd(f)
	cmd.SetArgs([]string{"api-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}

func TestStopSuccess(t *testing.T) {
	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/_stop") {
				t.Errorf("unexpected path: %s", path)
			}

			return nil, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newStopCmd(f)
	cmd.SetArgs([]string{"api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "API 'api-1' stopped.") {
		t.Errorf("unexpected output: %s", out.String())
	}
}

func TestDeploySuccess(t *testing.T) {
	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/deployments") {
				t.Errorf("unexpected path: %s", path)
			}

			return nil, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newDeployCmd(f)
	cmd.SetArgs([]string{"api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "deployment requested") {
		t.Errorf("unexpected output: %s", out.String())
	}
}

func TestDeployLabelTooLong(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, false)

	cmd := newDeployCmd(f)
	cmd.SetArgs([]string{"api-1", "--label", strings.Repeat("x", 33)})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "exceeds 32 characters") {
		t.Errorf("expected label error, got: %v", err)
	}
}

func TestUpdateSuccess(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "api.json")
	_ = os.WriteFile(file, []byte(`{"name":"Updated"}`), 0600)

	resp, _ := json.Marshal(map[string]string{"id": "api-1", "name": "Updated"})

	fake := &client.FakeClient{
		PutFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newUpdateCmd(f)
	cmd.SetArgs([]string{"api-1", "-f", file})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Updated") {
		t.Errorf("expected 'Updated' in output, got: %s", out.String())
	}
}

func TestImportSuccess(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "export.json")
	_ = os.WriteFile(file, []byte(`{"api":{"name":"Imported"}}`), 0600)

	resp, _ := json.Marshal(map[string]string{
		"id": "new-id", "name": "Imported", "state": "STOPPED",
		"definitionVersion": "V4", "type": "PROXY",
	})

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/_import/definition") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newImportCmd(f)
	cmd.SetArgs([]string{"-f", file})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Imported") {
		t.Errorf("expected 'Imported' in output, got: %s", out.String())
	}
}

func TestImportReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)

	cmd := newImportCmd(f)
	cmd.SetArgs([]string{"-f", "test.json"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}

func TestExportSuccess(t *testing.T) {
	resp := []byte(`{"api":{"id":"api-1","name":"Weather API"},"plans":[]}`)

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/_export/definition") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newExportCmd(f)
	cmd.SetArgs([]string{"api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Weather API") {
		t.Errorf("expected 'Weather API' in output, got: %s", out.String())
	}
}

func TestExportWithExclude(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "excludeAdditionalData=members,pages") {
				t.Errorf("expected exclude params, got: %s", path)
			}

			return []byte(`{}`), nil
		},
	}

	f, _ := newTestFactory(fake, false)

	cmd := newExportCmd(f)
	cmd.SetArgs([]string{"api-1", "--exclude", "members", "--exclude", "pages"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRollbackSuccess(t *testing.T) {
	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/_rollback") {
				t.Errorf("unexpected path: %s", path)
			}

			return nil, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newRollbackCmd(f)
	cmd.SetArgs([]string{"api-1", "--event-id", "evt-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "rolled back to event 'evt-1'") {
		t.Errorf("unexpected output: %s", out.String())
	}
}

func TestRollbackReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)

	cmd := newRollbackCmd(f)
	cmd.SetArgs([]string{"api-1", "--event-id", "evt-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}
