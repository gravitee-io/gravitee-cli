// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package alert

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

// --- Notifier List ---

func TestListAlertNotifiers(t *testing.T) {
	t.Run("returns notifiers", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			ListAlertNotifiersFunc: func(domainID string) ([]json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return []json.RawMessage{
					json.RawMessage(`{"id":"n-1","name":"Slack","type":"slack-notifier","enabled":true}`),
					json.RawMessage(`{"id":"n-2","name":"Email","type":"email-notifier","enabled":false}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "notifier", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Slack")
		testutil.AssertOutputContains(t, tc.Out, "Email")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			ListAlertNotifiersFunc: func(_ string) ([]json.RawMessage, error) {
				return []json.RawMessage{
					json.RawMessage(`{"id":"n-1","name":"Test"}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "notifier", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "notifier", "list")

		testutil.AssertErrorContains(t, err, "no context configured")
	})

	t.Run("requires domain flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "notifier", "list")

		testutil.AssertErrorContains(t, err, "required")
	})
}

// --- Notifier Get ---

func TestGetAlertNotifier(t *testing.T) {
	t.Run("returns notifier details", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetAlertNotifierFunc: func(domainID, notifierID string) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if notifierID != "n-1" {
					t.Errorf("expected notifierID 'n-1', got %q", notifierID)
				}

				return json.Marshal(map[string]any{
					"id": "n-1", "name": "Slack", "type": "slack-notifier", "enabled": true,
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "notifier", "get", "n-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Slack")
		testutil.AssertOutputContains(t, tc.Out, "n-1")
	})

	t.Run("requires notifier ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "notifier", "get")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "notifier", "get", "n-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Notifier Create ---

func TestCreateAlertNotifier(t *testing.T) {
	t.Run("creates a notifier from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			CreateAlertNotifierFunc: func(domainID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return json.Marshal(map[string]any{
					"id": "new-n", "name": "Slack", "type": "slack-notifier", "enabled": false,
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"name":"Slack","type":"slack-notifier"}`)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "notifier", "create", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Slack")
		testutil.AssertOutputContains(t, tc.Out, "new-n")
	})

	t.Run("requires file flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "notifier", "create")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "notifier", "create", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Notifier Update ---

func TestUpdateAlertNotifier(t *testing.T) {
	t.Run("updates a notifier from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			UpdateAlertNotifierFunc: func(domainID, notifierID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if notifierID != "n-1" {
					t.Errorf("expected notifierID 'n-1', got %q", notifierID)
				}

				return json.Marshal(map[string]any{
					"id": "n-1", "name": "Updated", "type": "slack-notifier", "enabled": true,
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"name":"Updated"}`)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "notifier", "update", "n-1", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Updated")
	})

	t.Run("requires file flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "notifier", "update", "n-1")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires notifier ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "notifier", "update", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "notifier", "update", "n-1", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Notifier Delete ---

func TestDeleteAlertNotifier(t *testing.T) {
	t.Run("deletes a notifier", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			DeleteAlertNotifierFunc: func(domainID, notifierID string) error {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if notifierID != "n-1" {
					t.Errorf("expected notifierID 'n-1', got %q", notifierID)
				}

				return nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "notifier", "delete", "n-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Alert notifier 'n-1' deleted.")
	})

	t.Run("returns error on API failure", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			DeleteAlertNotifierFunc: func(_, _ string) error {
				return &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "notifier", "delete", "n-1")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("requires notifier ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "notifier", "delete")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "notifier", "delete", "n-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Trigger Get ---

func TestGetAlertTriggers(t *testing.T) {
	t.Run("returns trigger data", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetAlertTriggersFunc: func(domainID string) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return json.Marshal([]map[string]any{
					{"type": "too_many_login_failures", "enabled": true},
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "trigger", "get")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "too_many_login_failures")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "trigger", "get")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Trigger Update ---

func TestUpdateAlertTriggers(t *testing.T) {
	t.Run("updates triggers from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			UpdateAlertTriggersFunc: func(domainID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return json.Marshal([]map[string]any{
					{"type": "too_many_login_failures", "enabled": true},
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `[{"type":"too_many_login_failures","enabled":true}]`)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "trigger", "update", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Alert triggers updated successfully.")
	})

	t.Run("returns JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			UpdateAlertTriggersFunc: func(_ string, _ json.RawMessage) (json.RawMessage, error) {
				return json.Marshal([]map[string]any{{"type": "test"}})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `[{"type":"test"}]`)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "trigger", "update", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"type"`)
	})

	t.Run("requires file flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "trigger", "update")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `[{"type":"test"}]`)

		cmd := NewAlertCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "trigger", "update", "--file", tmpFile)

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
