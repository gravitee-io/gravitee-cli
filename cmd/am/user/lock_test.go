package user

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestLockUser(t *testing.T) {
	t.Run("locks a user", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			UpdateUserStatusFunc: func(d, userID string, body json.RawMessage) (json.RawMessage, error) {
				if d != domainID {
					t.Errorf("expected domain %q, got %q", domainID, d)
				}
				if userID != "user-1" {
					t.Errorf("expected userID %q, got %q", "user-1", userID)
				}

				return json.Marshal(map[string]any{"id": "user-1", "enabled": false})
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestUser(tc, mock, domainID)

		err := testutil.Execute(newLockCmd(tc.Factory, &domainID), "user-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "User 'user-1' locked.")
	})

	t.Run("returns error on API failure", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			UpdateUserStatusFunc: func(_, _ string, _ json.RawMessage) (json.RawMessage, error) {
				return nil, &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestUser(tc, mock, domainID)

		err := testutil.Execute(newLockCmd(tc.Factory, &domainID), "user-1")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("requires user ID argument", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)

		err := testutil.Execute(newLockCmd(tc.Factory, &domainID))

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newLockCmd(tc.Factory, &domainID), "user-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

func TestUnlockUser(t *testing.T) {
	t.Run("unlocks a user", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			UpdateUserStatusFunc: func(d, userID string, _ json.RawMessage) (json.RawMessage, error) {
				if d != domainID {
					t.Errorf("expected domain %q, got %q", domainID, d)
				}
				if userID != "user-1" {
					t.Errorf("expected userID %q, got %q", "user-1", userID)
				}

				return json.Marshal(map[string]any{"id": "user-1", "enabled": true})
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestUser(tc, mock, domainID)

		err := testutil.Execute(newUnlockCmd(tc.Factory, &domainID), "user-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "User 'user-1' unlocked.")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newUnlockCmd(tc.Factory, &domainID), "user-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
