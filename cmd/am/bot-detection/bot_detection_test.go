package botdetection

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

func TestListBotDetections(t *testing.T) {
	t.Run("returns bot detections", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			ListBotDetectionsFunc: func(domainID string) ([]json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return []json.RawMessage{
					json.RawMessage(`{"id":"bd-1","name":"My Bot Detection","type":"recaptcha"}`),
					json.RawMessage(`{"id":"bd-2","name":"Other","type":"hcaptcha"}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewBotDetectionCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Bot Detection")
		testutil.AssertOutputContains(t, tc.Out, "Other")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			ListBotDetectionsFunc: func(_ string) ([]json.RawMessage, error) {
				return []json.RawMessage{
					json.RawMessage(`{"id":"bd-1","name":"Test"}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewBotDetectionCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		cmd := NewBotDetectionCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertErrorContains(t, err, "no context configured")
	})

	t.Run("requires domain flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewBotDetectionCmd(tc.Factory)
		err := testutil.Execute(cmd, "list")

		testutil.AssertErrorContains(t, err, "required")
	})
}

// --- Get ---

func TestGetBotDetection(t *testing.T) {
	t.Run("returns bot detection details", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			GetBotDetectionFunc: func(domainID, botDetectionID string) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if botDetectionID != "bd-1" {
					t.Errorf("expected botDetectionID 'bd-1', got %q", botDetectionID)
				}

				return json.Marshal(map[string]any{
					"id": "bd-1", "name": "My Bot Detection", "type": "recaptcha",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewBotDetectionCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "bd-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Bot Detection")
		testutil.AssertOutputContains(t, tc.Out, "bd-1")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			GetBotDetectionFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "bd-1", "name": "Test"})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewBotDetectionCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "get", "bd-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires bot detection ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewBotDetectionCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		cmd := NewBotDetectionCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "bd-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Create ---

func TestCreateBotDetection(t *testing.T) {
	t.Run("creates a bot detection from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			CreateBotDetectionFunc: func(domainID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return json.Marshal(map[string]any{
					"id": "new-bd", "name": "My Bot Detection", "type": "recaptcha",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"name":"My Bot Detection","type":"recaptcha"}`)

		cmd := NewBotDetectionCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Bot Detection")
		testutil.AssertOutputContains(t, tc.Out, "new-bd")
	})

	t.Run("requires file flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewBotDetectionCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewBotDetectionCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Update ---

func TestUpdateBotDetection(t *testing.T) {
	t.Run("updates a bot detection from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			UpdateBotDetectionFunc: func(domainID, botDetectionID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if botDetectionID != "bd-1" {
					t.Errorf("expected botDetectionID 'bd-1', got %q", botDetectionID)
				}

				return json.Marshal(map[string]any{
					"id": "bd-1", "name": "Updated", "type": "recaptcha",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"name":"Updated"}`)

		cmd := NewBotDetectionCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "bd-1", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Updated")
	})

	t.Run("requires file flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewBotDetectionCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "bd-1")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires bot detection ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewBotDetectionCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewBotDetectionCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "bd-1", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Delete ---

func TestDeleteBotDetection(t *testing.T) {
	t.Run("deletes a bot detection", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			DeleteBotDetectionFunc: func(domainID, botDetectionID string) error {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if botDetectionID != "bd-1" {
					t.Errorf("expected botDetectionID 'bd-1', got %q", botDetectionID)
				}

				return nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewBotDetectionCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "bd-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Bot detection 'bd-1' deleted.")
	})

	t.Run("returns error on API failure", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			DeleteBotDetectionFunc: func(_, _ string) error {
				return &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewBotDetectionCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "bd-1")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("requires bot detection ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewBotDetectionCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		cmd := NewBotDetectionCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "bd-1")

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
