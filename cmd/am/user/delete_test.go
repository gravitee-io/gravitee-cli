package user

import (
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestDeleteUser(t *testing.T) {
	t.Run("deletes a user", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			DeleteUserFunc: func(d, userID string) error {
				if d != domainID {
					t.Errorf("expected domain %q, got %q", domainID, d)
				}
				if userID != "user-1" {
					t.Errorf("expected userID %q, got %q", "user-1", userID)
				}

				return nil
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestUser(tc, mock, domainID)

		err := testutil.Execute(newDeleteCmd(tc.Factory, &domainID), "user-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "User 'user-1' deleted.")
	})

	t.Run("returns error on API failure", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			DeleteUserFunc: func(_, _ string) error {
				return &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestUser(tc, mock, domainID)

		err := testutil.Execute(newDeleteCmd(tc.Factory, &domainID), "user-1")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("requires user ID argument", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)

		err := testutil.Execute(newDeleteCmd(tc.Factory, &domainID))

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newDeleteCmd(tc.Factory, &domainID), "user-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
