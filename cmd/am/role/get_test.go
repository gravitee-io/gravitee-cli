package role

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestGetRole(t *testing.T) {
	t.Run("returns role details", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			GetRoleFunc: func(d, roleID string) (json.RawMessage, error) {
				if d != domainID {
					t.Errorf("expected domain %q, got %q", domainID, d)
				}
				if roleID != "role-1" {
					t.Errorf("expected roleID %q, got %q", "role-1", roleID)
				}

				return json.Marshal(map[string]any{
					"id": "role-1", "name": "Admin", "type": "DOMAIN", "description": "Admin role",
				})
			},
		}
		tc := testutil.NewFactory(nil)
		newTestRole(tc, mock, domainID)

		err := testutil.Execute(newGetCmd(tc.Factory, &domainID), "role-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Admin")
		testutil.AssertOutputContains(t, tc.Out, "role-1")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			GetRoleFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "role-1", "name": "Admin"})
			},
		}
		tc := testutil.NewFactory(nil)
		newTestRole(tc, mock, domainID)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newGetCmd(tc.Factory, &domainID), "role-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires role ID argument", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil)

		err := testutil.Execute(newGetCmd(tc.Factory, &domainID))

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newGetCmd(tc.Factory, &domainID), "role-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
