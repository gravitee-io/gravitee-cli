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

package email

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"gravitee.io/gctl/internal/am"
	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/testutil"
)

// --- List ---

func TestListEmails(t *testing.T) {
	t.Run("returns emails", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			ListEmailsFunc: func(domainID string) ([]json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return []json.RawMessage{
					json.RawMessage(`{"id":"email-1","template":"RESET_PASSWORD","enabled":true,"subject":"Reset"}`),
					json.RawMessage(`{"id":"email-2","template":"REGISTRATION_CONFIRMATION","enabled":false,"subject":"Welcome"}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewEmailCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "RESET_PASSWORD")
		testutil.AssertOutputContains(t, tc.Out, "REGISTRATION_CONFIRMATION")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			ListEmailsFunc: func(_ string) ([]json.RawMessage, error) {
				return []json.RawMessage{
					json.RawMessage(`{"id":"email-1","template":"RESET_PASSWORD"}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewEmailCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewEmailCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertErrorContains(t, err, "no context configured")
	})

	t.Run("requires domain flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewEmailCmd(tc.Factory)
		err := testutil.Execute(cmd, "list")

		testutil.AssertErrorContains(t, err, "required")
	})
}

// --- Get ---

func TestGetEmail(t *testing.T) {
	t.Run("returns email details", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetEmailFunc: func(domainID, emailID string) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if emailID != "email-1" {
					t.Errorf("expected emailID 'email-1', got %q", emailID)
				}

				return json.Marshal(map[string]any{
					"id": "email-1", "template": "RESET_PASSWORD", "enabled": true, "subject": "Reset",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewEmailCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "email-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "RESET_PASSWORD")
		testutil.AssertOutputContains(t, tc.Out, "email-1")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetEmailFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "email-1", "template": "RESET_PASSWORD"})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewEmailCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "get", "email-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires email ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewEmailCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewEmailCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "email-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Create ---

func TestCreateEmail(t *testing.T) {
	t.Run("creates an email from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			CreateEmailFunc: func(domainID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return json.Marshal(map[string]any{
					"id": "new-email", "template": "RESET_PASSWORD", "subject": "Reset",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"template":"RESET_PASSWORD","subject":"Reset"}`)

		cmd := NewEmailCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "RESET_PASSWORD")
		testutil.AssertOutputContains(t, tc.Out, "new-email")
	})

	t.Run("requires json input", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewEmailCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create")

		testutil.AssertErrorContains(t, err, "input")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"template":"RESET_PASSWORD"}`)

		cmd := NewEmailCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Update ---

func TestUpdateEmail(t *testing.T) {
	t.Run("updates an email from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			UpdateEmailFunc: func(domainID, emailID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if emailID != "email-1" {
					t.Errorf("expected emailID 'email-1', got %q", emailID)
				}

				return json.Marshal(map[string]any{
					"id": "email-1", "template": "RESET_PASSWORD", "subject": "Updated",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"subject":"Updated"}`)

		cmd := NewEmailCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "email-1", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Updated")
	})

	t.Run("requires json input", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewEmailCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "email-1")

		testutil.AssertErrorContains(t, err, "input")
	})

	t.Run("requires email ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		tmpFile := writeTempJSON(t, `{"subject":"Updated"}`)

		cmd := NewEmailCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"subject":"Updated"}`)

		cmd := NewEmailCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "email-1", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Delete ---

func TestDeleteEmail(t *testing.T) {
	t.Run("deletes an email", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			DeleteEmailFunc: func(domainID, emailID string) error {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if emailID != "email-1" {
					t.Errorf("expected emailID 'email-1', got %q", emailID)
				}

				return nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewEmailCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "email-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Email 'email-1' deleted.")
	})

	t.Run("returns error on API failure", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			DeleteEmailFunc: func(_, _ string) error {
				return &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewEmailCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "email-1")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("requires email ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewEmailCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewEmailCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "email-1")

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
