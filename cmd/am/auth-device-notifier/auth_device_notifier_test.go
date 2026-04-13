package authdevicenotifier

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

func TestListAuthDeviceNotifiers(t *testing.T) {
	t.Run("returns auth device notifiers", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			ListAuthDeviceNotifiersFunc: func(domainID string) ([]json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return []json.RawMessage{
					json.RawMessage(`{"id":"adn-1","name":"My Notifier","type":"http"}`),
					json.RawMessage(`{"id":"adn-2","name":"Other","type":"smtp"}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAuthDeviceNotifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Notifier")
		testutil.AssertOutputContains(t, tc.Out, "Other")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			ListAuthDeviceNotifiersFunc: func(_ string) ([]json.RawMessage, error) {
				return []json.RawMessage{
					json.RawMessage(`{"id":"adn-1","name":"Test"}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAuthDeviceNotifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		cmd := NewAuthDeviceNotifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertErrorContains(t, err, "no context configured")
	})

	t.Run("requires domain flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewAuthDeviceNotifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "list")

		testutil.AssertErrorContains(t, err, "required")
	})
}

// --- Get ---

func TestGetAuthDeviceNotifier(t *testing.T) {
	t.Run("returns auth device notifier details", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			GetAuthDeviceNotifierFunc: func(domainID, adnID string) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if adnID != "adn-1" {
					t.Errorf("expected authDeviceNotifierID 'adn-1', got %q", adnID)
				}

				return json.Marshal(map[string]any{
					"id": "adn-1", "name": "My Notifier", "type": "http",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAuthDeviceNotifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "adn-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Notifier")
		testutil.AssertOutputContains(t, tc.Out, "adn-1")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			GetAuthDeviceNotifierFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "adn-1", "name": "Test"})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAuthDeviceNotifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "get", "adn-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires auth device notifier ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewAuthDeviceNotifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		cmd := NewAuthDeviceNotifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "adn-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Create ---

func TestCreateAuthDeviceNotifier(t *testing.T) {
	t.Run("creates an auth device notifier from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			CreateAuthDeviceNotifierFunc: func(domainID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return json.Marshal(map[string]any{
					"id": "new-adn", "name": "My Notifier", "type": "http",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"name":"My Notifier","type":"http"}`)

		cmd := NewAuthDeviceNotifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Notifier")
		testutil.AssertOutputContains(t, tc.Out, "new-adn")
	})

	t.Run("requires file flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewAuthDeviceNotifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewAuthDeviceNotifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Update ---

func TestUpdateAuthDeviceNotifier(t *testing.T) {
	t.Run("updates an auth device notifier from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			UpdateAuthDeviceNotifierFunc: func(domainID, adnID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if adnID != "adn-1" {
					t.Errorf("expected authDeviceNotifierID 'adn-1', got %q", adnID)
				}

				return json.Marshal(map[string]any{
					"id": "adn-1", "name": "Updated", "type": "http",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"name":"Updated"}`)

		cmd := NewAuthDeviceNotifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "adn-1", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Updated")
	})

	t.Run("requires file flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewAuthDeviceNotifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "adn-1")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires auth device notifier ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewAuthDeviceNotifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewAuthDeviceNotifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "adn-1", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Delete ---

func TestDeleteAuthDeviceNotifier(t *testing.T) {
	t.Run("deletes an auth device notifier", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			DeleteAuthDeviceNotifierFunc: func(domainID, adnID string) error {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if adnID != "adn-1" {
					t.Errorf("expected authDeviceNotifierID 'adn-1', got %q", adnID)
				}

				return nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAuthDeviceNotifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "adn-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Auth device notifier 'adn-1' deleted.")
	})

	t.Run("returns error on API failure", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			DeleteAuthDeviceNotifierFunc: func(_, _ string) error {
				return &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAuthDeviceNotifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "adn-1")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("requires auth device notifier ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewAuthDeviceNotifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		cmd := NewAuthDeviceNotifierCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "adn-1")

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
