package role

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestUpdateRole(t *testing.T) {
	t.Run("updates role name", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			UpdateRoleFunc: func(d, roleID string, _ json.RawMessage) (json.RawMessage, error) {
				if d != domainID {
					t.Errorf("expected domain %q, got %q", domainID, d)
				}
				if roleID != "role-1" {
					t.Errorf("expected roleID %q, got %q", "role-1", roleID)
				}

				return json.Marshal(map[string]any{
					"id": "role-1", "name": "Updated",
				})
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestRole(tc, mock, domainID)

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID), "role-1", "--name", "Updated")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Updated")
	})

	t.Run("updates role description", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			UpdateRoleFunc: func(_, _ string, _ json.RawMessage) (json.RawMessage, error) {
				return json.Marshal(map[string]any{
					"id": "role-1", "name": "Admin", "description": "New desc",
				})
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestRole(tc, mock, domainID)

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID), "role-1", "--description", "New desc")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "New desc")
	})

	t.Run("requires at least one flag", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID), "role-1")

		testutil.AssertErrorContains(t, err, "at least one flag")
	})

	t.Run("requires role ID argument", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID))

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID), "role-1", "--name", "Test")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
