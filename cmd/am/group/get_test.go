package group

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestGetGroup(t *testing.T) {
	t.Run("returns group details", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			GetGroupFunc: func(d, groupID string) (json.RawMessage, error) {
				if d != domainID {
					t.Errorf("expected domain %q, got %q", domainID, d)
				}
				if groupID != "grp-1" {
					t.Errorf("expected groupID %q, got %q", "grp-1", groupID)
				}

				return json.Marshal(map[string]any{
					"id": "grp-1", "name": "Admins", "description": "Admin group",
				})
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestGroup(tc, mock, domainID)

		err := testutil.Execute(newGetCmd(tc.Factory, &domainID), "grp-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Admins")
		testutil.AssertOutputContains(t, tc.Out, "grp-1")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			GetGroupFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "grp-1", "name": "Admins"})
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestGroup(tc, mock, domainID)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newGetCmd(tc.Factory, &domainID), "grp-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires group ID argument", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)

		err := testutil.Execute(newGetCmd(tc.Factory, &domainID))

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newGetCmd(tc.Factory, &domainID), "grp-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
