package scope

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestGetScope(t *testing.T) {
	t.Run("returns scope details", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			GetScopeFunc: func(d, scopeID string) (json.RawMessage, error) {
				if d != domainID {
					t.Errorf("expected domain %q, got %q", domainID, d)
				}
				if scopeID != "scope-1" {
					t.Errorf("expected scopeID %q, got %q", "scope-1", scopeID)
				}

				return json.Marshal(map[string]any{
					"id": "scope-1", "key": "openid", "name": "OpenID", "description": "OpenID scope",
				})
			},
		}
		tc := testutil.NewFactory(nil)
		newTestScope(tc, mock, domainID)

		err := testutil.Execute(newGetCmd(tc.Factory, &domainID), "scope-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "OpenID")
		testutil.AssertOutputContains(t, tc.Out, "scope-1")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			GetScopeFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "scope-1", "key": "openid", "name": "OpenID"})
			},
		}
		tc := testutil.NewFactory(nil)
		newTestScope(tc, mock, domainID)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newGetCmd(tc.Factory, &domainID), "scope-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires scope ID argument", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil)

		err := testutil.Execute(newGetCmd(tc.Factory, &domainID))

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newGetCmd(tc.Factory, &domainID), "scope-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
