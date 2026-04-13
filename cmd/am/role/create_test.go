package role

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestCreateRole(t *testing.T) {
	t.Run("creates a role with name", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			CreateRoleFunc: func(d string, body json.RawMessage) (json.RawMessage, error) {
				if d != domainID {
					t.Errorf("expected domain %q, got %q", domainID, d)
				}

				return json.Marshal(map[string]any{
					"id": "new-role", "name": "Admin",
				})
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestRole(tc, mock, domainID)

		err := testutil.Execute(newCreateCmd(tc.Factory, &domainID), "--name", "Admin")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Admin")
		testutil.AssertOutputContains(t, tc.Out, "new-role")
	})

	t.Run("creates a role with name and description", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			CreateRoleFunc: func(_ string, _ json.RawMessage) (json.RawMessage, error) {
				return json.Marshal(map[string]any{
					"id": "new-role", "name": "Admin", "description": "Admin role",
				})
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestRole(tc, mock, domainID)

		err := testutil.Execute(newCreateCmd(tc.Factory, &domainID), "--name", "Admin", "--description", "Admin role")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Admin role")
	})

	t.Run("returns JSON with -o json", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			CreateRoleFunc: func(_ string, _ json.RawMessage) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "new-role", "name": "Admin"})
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestRole(tc, mock, domainID)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newCreateCmd(tc.Factory, &domainID), "--name", "Admin")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires name flag", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)

		err := testutil.Execute(newCreateCmd(tc.Factory, &domainID))

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newCreateCmd(tc.Factory, &domainID), "--name", "Admin")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
