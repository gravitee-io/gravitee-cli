package form

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

func TestListForms(t *testing.T) {
	t.Run("returns forms", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			ListFormsFunc: func(domainID string) ([]json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return []json.RawMessage{
					json.RawMessage(`{"id":"form-1","template":"LOGIN","enabled":true}`),
					json.RawMessage(`{"id":"form-2","template":"REGISTRATION","enabled":false}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewFormCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "LOGIN")
		testutil.AssertOutputContains(t, tc.Out, "REGISTRATION")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			ListFormsFunc: func(_ string) ([]json.RawMessage, error) {
				return []json.RawMessage{
					json.RawMessage(`{"id":"form-1","template":"LOGIN"}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewFormCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		cmd := NewFormCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertErrorContains(t, err, "no context configured")
	})

	t.Run("requires domain flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewFormCmd(tc.Factory)
		err := testutil.Execute(cmd, "list")

		testutil.AssertErrorContains(t, err, "required")
	})
}

// --- Get ---

func TestGetForm(t *testing.T) {
	t.Run("returns form details", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			GetFormFunc: func(domainID, formID string) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if formID != "form-1" {
					t.Errorf("expected formID 'form-1', got %q", formID)
				}

				return json.Marshal(map[string]any{
					"id": "form-1", "template": "LOGIN", "enabled": true,
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewFormCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "form-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "LOGIN")
		testutil.AssertOutputContains(t, tc.Out, "form-1")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			GetFormFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "form-1", "template": "LOGIN"})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewFormCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "get", "form-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires form ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewFormCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		cmd := NewFormCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "form-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Create ---

func TestCreateForm(t *testing.T) {
	t.Run("creates a form from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			CreateFormFunc: func(domainID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return json.Marshal(map[string]any{
					"id": "new-form", "template": "LOGIN", "enabled": true,
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"template":"LOGIN","enabled":true}`)

		cmd := NewFormCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "LOGIN")
		testutil.AssertOutputContains(t, tc.Out, "new-form")
	})

	t.Run("requires file flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewFormCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"template":"LOGIN"}`)

		cmd := NewFormCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Update ---

func TestUpdateForm(t *testing.T) {
	t.Run("updates a form from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			UpdateFormFunc: func(domainID, formID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if formID != "form-1" {
					t.Errorf("expected formID 'form-1', got %q", formID)
				}

				return json.Marshal(map[string]any{
					"id": "form-1", "template": "LOGIN", "enabled": false,
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"enabled":false}`)

		cmd := NewFormCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "form-1", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "LOGIN")
	})

	t.Run("requires file flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewFormCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "form-1")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires form ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		tmpFile := writeTempJSON(t, `{"enabled":false}`)

		cmd := NewFormCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"enabled":false}`)

		cmd := NewFormCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "form-1", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Delete ---

func TestDeleteForm(t *testing.T) {
	t.Run("deletes a form", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			DeleteFormFunc: func(domainID, formID string) error {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if formID != "form-1" {
					t.Errorf("expected formID 'form-1', got %q", formID)
				}

				return nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewFormCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "form-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Form 'form-1' deleted.")
	})

	t.Run("returns error on API failure", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			DeleteFormFunc: func(_, _ string) error {
				return &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewFormCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "form-1")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("requires form ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewFormCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		cmd := NewFormCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "form-1")

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
