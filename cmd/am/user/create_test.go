package user

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestCreateUser(t *testing.T) {
	t.Run("creates a user with username and email", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			CreateUserFunc: func(d string, body json.RawMessage) (json.RawMessage, error) {
				if d != domainID {
					t.Errorf("expected domain %q, got %q", domainID, d)
				}

				return json.Marshal(map[string]any{
					"id": "new-user", "username": "alice", "email": "alice@example.com",
				})
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestUser(tc, mock, domainID)

		err := testutil.Execute(newCreateCmd(tc.Factory, &domainID), "--username", "alice", "--email", "alice@example.com")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "alice")
		testutil.AssertOutputContains(t, tc.Out, "new-user")
	})

	t.Run("creates a user with all fields", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			CreateUserFunc: func(_ string, _ json.RawMessage) (json.RawMessage, error) {
				return json.Marshal(map[string]any{
					"id": "new-user", "username": "alice", "email": "alice@example.com",
					"firstName": "Alice", "lastName": "Smith",
				})
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestUser(tc, mock, domainID)

		err := testutil.Execute(newCreateCmd(tc.Factory, &domainID),
			"--username", "alice", "--email", "alice@example.com",
			"--password", "secret", "--firstName", "Alice", "--lastName", "Smith", "--preRegistration")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Alice")
	})

	t.Run("returns JSON with -o json", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			CreateUserFunc: func(_ string, _ json.RawMessage) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "new-user", "username": "alice"})
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestUser(tc, mock, domainID)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newCreateCmd(tc.Factory, &domainID), "--username", "alice", "--email", "alice@example.com")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires username flag", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)

		err := testutil.Execute(newCreateCmd(tc.Factory, &domainID), "--email", "alice@example.com")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires email flag", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)

		err := testutil.Execute(newCreateCmd(tc.Factory, &domainID), "--username", "alice")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newCreateCmd(tc.Factory, &domainID), "--username", "alice", "--email", "alice@example.com")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
