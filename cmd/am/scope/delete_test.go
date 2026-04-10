package scope

import (
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestDeleteScope(t *testing.T) {
	t.Run("deletes a scope", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			DeleteScopeFunc: func(d, scopeID string) error {
				if d != domainID {
					t.Errorf("expected domain %q, got %q", domainID, d)
				}
				if scopeID != "scope-1" {
					t.Errorf("expected scopeID %q, got %q", "scope-1", scopeID)
				}

				return nil
			},
		}
		tc := testutil.NewFactory(nil)
		newTestScope(tc, mock, domainID)

		err := testutil.Execute(newDeleteCmd(tc.Factory, &domainID), "scope-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Scope 'scope-1' deleted.")
	})

	t.Run("returns error on API failure", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			DeleteScopeFunc: func(_, _ string) error {
				return &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
			},
		}
		tc := testutil.NewFactory(nil)
		newTestScope(tc, mock, domainID)

		err := testutil.Execute(newDeleteCmd(tc.Factory, &domainID), "scope-1")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("requires scope ID argument", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil)

		err := testutil.Execute(newDeleteCmd(tc.Factory, &domainID))

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newDeleteCmd(tc.Factory, &domainID), "scope-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
