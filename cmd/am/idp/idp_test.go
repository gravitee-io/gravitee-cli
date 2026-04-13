package idp

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

func TestListIdentityProviders(t *testing.T) {
	t.Run("returns identity providers", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			ListIdentityProvidersFunc: func(domainID string, userProvider bool) ([]json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if userProvider {
					t.Error("expected userProvider=false")
				}

				return []json.RawMessage{
					json.RawMessage(`{"id":"idp-1","name":"My IdP","type":"inline","enabled":true}`),
					json.RawMessage(`{"id":"idp-2","name":"Other","type":"ldap","enabled":false}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My IdP")
		testutil.AssertOutputContains(t, tc.Out, "Other")
	})

	t.Run("passes user-provider flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			ListIdentityProvidersFunc: func(_ string, userProvider bool) ([]json.RawMessage, error) {
				if !userProvider {
					t.Error("expected userProvider=true")
				}

				return []json.RawMessage{}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list", "--user-provider")

		testutil.AssertNoError(t, err)
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			ListIdentityProvidersFunc: func(_ string, _ bool) ([]json.RawMessage, error) {
				return []json.RawMessage{
					json.RawMessage(`{"id":"idp-1","name":"Test"}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertErrorContains(t, err, "no context configured")
	})

	t.Run("requires domain flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "list")

		testutil.AssertErrorContains(t, err, "required")
	})
}

// --- Get ---

func TestGetIdentityProvider(t *testing.T) {
	t.Run("returns identity provider details", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			GetIdentityProviderFunc: func(domainID, idpID string) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if idpID != "idp-1" {
					t.Errorf("expected idpID 'idp-1', got %q", idpID)
				}

				return json.Marshal(map[string]any{
					"id": "idp-1", "name": "My IdP", "type": "inline", "enabled": true,
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "idp-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My IdP")
		testutil.AssertOutputContains(t, tc.Out, "idp-1")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			GetIdentityProviderFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "idp-1", "name": "Test"})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "get", "idp-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires idp ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "idp-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Create ---

func TestCreateIdentityProvider(t *testing.T) {
	t.Run("creates an identity provider from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			CreateIdentityProviderFunc: func(domainID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return json.Marshal(map[string]any{
					"id": "new-idp", "name": "My IdP", "type": "inline", "enabled": false,
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"name":"My IdP","type":"inline"}`)

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My IdP")
		testutil.AssertOutputContains(t, tc.Out, "new-idp")
	})

	t.Run("requires file flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Update ---

func TestUpdateIdentityProvider(t *testing.T) {
	t.Run("updates an identity provider from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			UpdateIdentityProviderFunc: func(domainID, idpID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if idpID != "idp-1" {
					t.Errorf("expected idpID 'idp-1', got %q", idpID)
				}

				return json.Marshal(map[string]any{
					"id": "idp-1", "name": "Updated", "type": "inline", "enabled": true,
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"name":"Updated"}`)

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "idp-1", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Updated")
	})

	t.Run("requires file flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "idp-1")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires idp ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "idp-1", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Delete ---

func TestDeleteIdentityProvider(t *testing.T) {
	t.Run("deletes an identity provider", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			DeleteIdentityProviderFunc: func(domainID, idpID string) error {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if idpID != "idp-1" {
					t.Errorf("expected idpID 'idp-1', got %q", idpID)
				}

				return nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "idp-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Identity provider 'idp-1' deleted.")
	})

	t.Run("returns error on API failure", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			DeleteIdentityProviderFunc: func(_, _ string) error {
				return &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "idp-1")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("requires idp ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		cmd := NewIDPCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "idp-1")

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
