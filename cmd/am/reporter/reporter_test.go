package reporter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

// --- List ---

func TestListReporters(t *testing.T) {
	t.Run("returns reporters", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			ListReportersFunc: func(domainID string) ([]json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return []json.RawMessage{
					json.RawMessage(`{"id":"rep-1","name":"My Reporter","type":"mongodb","enabled":true}`),
					json.RawMessage(`{"id":"rep-2","name":"Other","type":"file","enabled":false}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewReporterCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Reporter")
		testutil.AssertOutputContains(t, tc.Out, "Other")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			ListReportersFunc: func(_ string) ([]json.RawMessage, error) {
				return []json.RawMessage{
					json.RawMessage(`{"id":"rep-1","name":"Test"}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewReporterCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		cmd := NewReporterCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertErrorContains(t, err, "no context configured")
	})

	t.Run("requires domain flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewReporterCmd(tc.Factory)
		err := testutil.Execute(cmd, "list")

		testutil.AssertErrorContains(t, err, "required")
	})
}

// --- Get ---

func TestGetReporter(t *testing.T) {
	t.Run("returns reporter details", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			GetReporterFunc: func(domainID, reporterID string) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if reporterID != "rep-1" {
					t.Errorf("expected reporterID 'rep-1', got %q", reporterID)
				}

				return json.Marshal(map[string]any{
					"id": "rep-1", "name": "My Reporter", "type": "mongodb", "enabled": true,
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewReporterCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "rep-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Reporter")
		testutil.AssertOutputContains(t, tc.Out, "rep-1")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			GetReporterFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "rep-1", "name": "Test"})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewReporterCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "get", "rep-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires reporter ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewReporterCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		cmd := NewReporterCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "rep-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Create ---

func TestCreateReporter(t *testing.T) {
	t.Run("creates a reporter from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			CreateReporterFunc: func(domainID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return json.Marshal(map[string]any{
					"id": "new-rep", "name": "My Reporter", "type": "mongodb", "enabled": false,
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"name":"My Reporter","type":"mongodb"}`)

		cmd := NewReporterCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Reporter")
		testutil.AssertOutputContains(t, tc.Out, "new-rep")
	})

	t.Run("requires file flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewReporterCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewReporterCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Update ---

func TestUpdateReporter(t *testing.T) {
	t.Run("updates a reporter from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			UpdateReporterFunc: func(domainID, reporterID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if reporterID != "rep-1" {
					t.Errorf("expected reporterID 'rep-1', got %q", reporterID)
				}

				return json.Marshal(map[string]any{
					"id": "rep-1", "name": "Updated", "type": "mongodb", "enabled": true,
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"name":"Updated"}`)

		cmd := NewReporterCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "rep-1", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Updated")
	})

	t.Run("requires file flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewReporterCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "rep-1")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires reporter ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewReporterCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewReporterCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "rep-1", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Delete ---

func TestDeleteReporter(t *testing.T) {
	t.Run("deletes a reporter", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			DeleteReporterFunc: func(domainID, reporterID string) error {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if reporterID != "rep-1" {
					t.Errorf("expected reporterID 'rep-1', got %q", reporterID)
				}

				return nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewReporterCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "rep-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Reporter 'rep-1' deleted.")
	})

	t.Run("returns error on API failure", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			DeleteReporterFunc: func(_, _ string) error {
				return &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewReporterCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "rep-1")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("requires reporter ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewReporterCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		cmd := NewReporterCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "rep-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Helper ---

func writeTempJSON(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.json")

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	return path
}
