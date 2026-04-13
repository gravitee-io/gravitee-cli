package user

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestGetUser(t *testing.T) {
	t.Run("returns user details", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			GetUserFunc: func(d, userID string) (json.RawMessage, error) {
				if d != domainID {
					t.Errorf("expected domain %q, got %q", domainID, d)
				}
				if userID != "user-1" {
					t.Errorf("expected userID %q, got %q", "user-1", userID)
				}

				return json.Marshal(map[string]any{
					"id": "user-1", "username": "alice", "email": "alice@example.com", "enabled": true,
				})
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestUser(tc, mock, domainID)

		err := testutil.Execute(newGetCmd(tc.Factory, &domainID), "user-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "alice")
		testutil.AssertOutputContains(t, tc.Out, "user-1")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			GetUserFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "user-1", "username": "alice"})
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestUser(tc, mock, domainID)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newGetCmd(tc.Factory, &domainID), "user-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires user ID argument", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)

		err := testutil.Execute(newGetCmd(tc.Factory, &domainID))

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newGetCmd(tc.Factory, &domainID), "user-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
