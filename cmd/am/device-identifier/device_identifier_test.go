package deviceidentifier

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

func TestListDeviceIdentifiers(t *testing.T) {
	t.Run("returns device identifiers", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			ListDeviceIdentifiersFunc: func(domainID string) ([]json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return []json.RawMessage{
					json.RawMessage(`{"id":"di-1","name":"My Device","type":"fingerprintjs"}`),
					json.RawMessage(`{"id":"di-2","name":"Other","type":"custom"}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewDeviceIdentifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Device")
		testutil.AssertOutputContains(t, tc.Out, "Other")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			ListDeviceIdentifiersFunc: func(_ string) ([]json.RawMessage, error) {
				return []json.RawMessage{
					json.RawMessage(`{"id":"di-1","name":"Test"}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewDeviceIdentifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		cmd := NewDeviceIdentifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertErrorContains(t, err, "no context configured")
	})

	t.Run("requires domain flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewDeviceIdentifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "list")

		testutil.AssertErrorContains(t, err, "required")
	})
}

// --- Get ---

func TestGetDeviceIdentifier(t *testing.T) {
	t.Run("returns device identifier details", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			GetDeviceIdentifierFunc: func(domainID, deviceIdentifierID string) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if deviceIdentifierID != "di-1" {
					t.Errorf("expected deviceIdentifierID 'di-1', got %q", deviceIdentifierID)
				}

				return json.Marshal(map[string]any{
					"id": "di-1", "name": "My Device", "type": "fingerprintjs",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewDeviceIdentifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "di-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Device")
		testutil.AssertOutputContains(t, tc.Out, "di-1")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			GetDeviceIdentifierFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "di-1", "name": "Test"})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewDeviceIdentifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "get", "di-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires device identifier ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewDeviceIdentifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		cmd := NewDeviceIdentifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "di-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Create ---

func TestCreateDeviceIdentifier(t *testing.T) {
	t.Run("creates a device identifier from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			CreateDeviceIdentifierFunc: func(domainID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return json.Marshal(map[string]any{
					"id": "new-di", "name": "My Device", "type": "fingerprintjs",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"name":"My Device","type":"fingerprintjs"}`)

		cmd := NewDeviceIdentifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Device")
		testutil.AssertOutputContains(t, tc.Out, "new-di")
	})

	t.Run("requires file flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewDeviceIdentifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewDeviceIdentifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Update ---

func TestUpdateDeviceIdentifier(t *testing.T) {
	t.Run("updates a device identifier from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			UpdateDeviceIdentifierFunc: func(domainID, deviceIdentifierID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if deviceIdentifierID != "di-1" {
					t.Errorf("expected deviceIdentifierID 'di-1', got %q", deviceIdentifierID)
				}

				return json.Marshal(map[string]any{
					"id": "di-1", "name": "Updated", "type": "fingerprintjs",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"name":"Updated"}`)

		cmd := NewDeviceIdentifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "di-1", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Updated")
	})

	t.Run("requires file flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewDeviceIdentifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "di-1")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires device identifier ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewDeviceIdentifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewDeviceIdentifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "di-1", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Delete ---

func TestDeleteDeviceIdentifier(t *testing.T) {
	t.Run("deletes a device identifier", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			DeleteDeviceIdentifierFunc: func(domainID, deviceIdentifierID string) error {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if deviceIdentifierID != "di-1" {
					t.Errorf("expected deviceIdentifierID 'di-1', got %q", deviceIdentifierID)
				}

				return nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewDeviceIdentifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "di-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Device identifier 'di-1' deleted.")
	})

	t.Run("returns error on API failure", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			DeleteDeviceIdentifierFunc: func(_, _ string) error {
				return &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewDeviceIdentifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "di-1")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("requires device identifier ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewDeviceIdentifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		cmd := NewDeviceIdentifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "di-1")

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
