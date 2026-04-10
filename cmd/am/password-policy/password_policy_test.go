package passwordpolicy

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

func TestListPasswordPolicies(t *testing.T) {
	t.Run("returns password policies", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			ListPasswordPoliciesFunc: func(domainID string) ([]json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return []json.RawMessage{
					json.RawMessage(`{"id":"pp-1","name":"Strong","minLength":8,"maxLength":64}`),
					json.RawMessage(`{"id":"pp-2","name":"Weak","minLength":4,"maxLength":32}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewPasswordPolicyCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Strong")
		testutil.AssertOutputContains(t, tc.Out, "Weak")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			ListPasswordPoliciesFunc: func(_ string) ([]json.RawMessage, error) {
				return []json.RawMessage{
					json.RawMessage(`{"id":"pp-1","name":"Strong"}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewPasswordPolicyCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewPasswordPolicyCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertErrorContains(t, err, "no context configured")
	})

	t.Run("requires domain flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewPasswordPolicyCmd(tc.Factory)
		err := testutil.Execute(cmd, "list")

		testutil.AssertErrorContains(t, err, "required")
	})
}

// --- Get ---

func TestGetPasswordPolicy(t *testing.T) {
	t.Run("returns password policy details", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetPasswordPolicyFunc: func(domainID, policyID string) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if policyID != "pp-1" {
					t.Errorf("expected policyID 'pp-1', got %q", policyID)
				}

				return json.Marshal(map[string]any{
					"id": "pp-1", "name": "Strong", "minLength": 8, "maxLength": 64,
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewPasswordPolicyCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "pp-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Strong")
		testutil.AssertOutputContains(t, tc.Out, "pp-1")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetPasswordPolicyFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "pp-1", "name": "Strong"})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewPasswordPolicyCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "get", "pp-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires policy ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewPasswordPolicyCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewPasswordPolicyCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "pp-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Create ---

func TestCreatePasswordPolicy(t *testing.T) {
	t.Run("creates a password policy from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			CreatePasswordPolicyFunc: func(domainID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return json.Marshal(map[string]any{
					"id": "new-pp", "name": "Custom", "minLength": 10, "maxLength": 128,
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"name":"Custom","minLength":10,"maxLength":128}`)

		cmd := NewPasswordPolicyCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Custom")
		testutil.AssertOutputContains(t, tc.Out, "new-pp")
	})

	t.Run("requires file flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewPasswordPolicyCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"name":"Custom"}`)

		cmd := NewPasswordPolicyCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Update ---

func TestUpdatePasswordPolicy(t *testing.T) {
	t.Run("updates a password policy from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			UpdatePasswordPolicyFunc: func(domainID, policyID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if policyID != "pp-1" {
					t.Errorf("expected policyID 'pp-1', got %q", policyID)
				}

				return json.Marshal(map[string]any{
					"id": "pp-1", "name": "Updated", "minLength": 12, "maxLength": 256,
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"name":"Updated","minLength":12}`)

		cmd := NewPasswordPolicyCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "pp-1", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Updated")
	})

	t.Run("requires file flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewPasswordPolicyCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "pp-1")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires policy ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		tmpFile := writeTempJSON(t, `{"name":"Updated"}`)

		cmd := NewPasswordPolicyCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"name":"Updated"}`)

		cmd := NewPasswordPolicyCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "pp-1", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Delete ---

func TestDeletePasswordPolicy(t *testing.T) {
	t.Run("deletes a password policy", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			DeletePasswordPolicyFunc: func(domainID, policyID string) error {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if policyID != "pp-1" {
					t.Errorf("expected policyID 'pp-1', got %q", policyID)
				}

				return nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewPasswordPolicyCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "pp-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Password policy 'pp-1' deleted.")
	})

	t.Run("returns error on API failure", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			DeletePasswordPolicyFunc: func(_, _ string) error {
				return &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewPasswordPolicyCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "pp-1")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("requires policy ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewPasswordPolicyCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewPasswordPolicyCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "pp-1")

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
