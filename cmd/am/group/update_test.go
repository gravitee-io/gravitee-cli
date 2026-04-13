package group

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestUpdateGroup(t *testing.T) {
	t.Run("updates group name", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			UpdateGroupFunc: func(d, groupID string, _ json.RawMessage) (json.RawMessage, error) {
				if d != domainID {
					t.Errorf("expected domain %q, got %q", domainID, d)
				}
				if groupID != "grp-1" {
					t.Errorf("expected groupID %q, got %q", "grp-1", groupID)
				}

				return json.Marshal(map[string]any{
					"id": "grp-1", "name": "Updated",
				})
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestGroup(tc, mock, domainID)

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID), "grp-1", "--name", "Updated")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Updated")
	})

	t.Run("updates group description", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			UpdateGroupFunc: func(_, _ string, _ json.RawMessage) (json.RawMessage, error) {
				return json.Marshal(map[string]any{
					"id": "grp-1", "name": "Admins", "description": "New desc",
				})
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestGroup(tc, mock, domainID)

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID), "grp-1", "--description", "New desc")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "New desc")
	})

	t.Run("requires at least one flag", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID), "grp-1")

		testutil.AssertErrorContains(t, err, "at least one flag")
	})

	t.Run("requires group ID argument", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID))

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID), "grp-1", "--name", "Test")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
