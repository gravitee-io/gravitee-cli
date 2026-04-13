package user

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestUpdateUser(t *testing.T) {
	t.Run("updates user email", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			UpdateUserFunc: func(d, userID string, _ json.RawMessage) (json.RawMessage, error) {
				if d != domainID {
					t.Errorf("expected domain %q, got %q", domainID, d)
				}
				if userID != "user-1" {
					t.Errorf("expected userID %q, got %q", "user-1", userID)
				}

				return json.Marshal(map[string]any{
					"id": "user-1", "username": "alice", "email": "new@example.com",
				})
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestUser(tc, mock, domainID)

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID), "user-1", "--email", "new@example.com")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "new@example.com")
	})

	t.Run("updates user firstName", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			UpdateUserFunc: func(_, _ string, _ json.RawMessage) (json.RawMessage, error) {
				return json.Marshal(map[string]any{
					"id": "user-1", "username": "alice", "firstName": "Alice",
				})
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestUser(tc, mock, domainID)

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID), "user-1", "--firstName", "Alice")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Alice")
	})

	t.Run("rejects invalid enabled value", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)
		newTestUser(tc, &am.MockService{}, domainID)

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID), "user-1", "--enabled", "maybe")

		testutil.AssertErrorContains(t, err, "invalid value")
	})

	t.Run("requires at least one flag", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID), "user-1")

		testutil.AssertErrorContains(t, err, "at least one flag")
	})

	t.Run("requires user ID argument", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID))

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID), "user-1", "--email", "test@example.com")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
