package user

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func toRaw(items ...map[string]any) []json.RawMessage {
	var result []json.RawMessage
	for _, item := range items {
		raw, _ := json.Marshal(item)
		result = append(result, raw)
	}

	return result
}

func TestListUsers(t *testing.T) {
	t.Run("returns users from the domain", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			ListUsersFunc: func(d string, _ am.ListUsersParams) (*am.PaginatedResponse, error) {
				if d != domainID {
					t.Errorf("expected domain %q, got %q", domainID, d)
				}

				data := toRaw(
					map[string]any{"id": "user-1", "username": "alice", "email": "alice@example.com", "enabled": true},
					map[string]any{"id": "user-2", "username": "bob", "email": "bob@example.com", "enabled": false},
				)

				return &am.PaginatedResponse{Data: data, TotalCount: 2, CurrentPage: 0}, nil
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestUser(tc, mock, domainID)

		err := testutil.Execute(newListCmd(tc.Factory, &domainID))

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "alice")
		testutil.AssertOutputContains(t, tc.Out, "bob")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			ListUsersFunc: func(_ string, _ am.ListUsersParams) (*am.PaginatedResponse, error) {
				data := toRaw(map[string]any{"id": "user-1", "username": "alice"})

				return &am.PaginatedResponse{Data: data, TotalCount: 1, CurrentPage: 0}, nil
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestUser(tc, mock, domainID)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newListCmd(tc.Factory, &domainID))

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"totalCount"`)
	})

	t.Run("requires a configured context", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil, false)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newListCmd(tc.Factory, &domainID))

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
