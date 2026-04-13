package role

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

func TestListRoles(t *testing.T) {
	t.Run("returns roles from the domain", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			ListRolesFunc: func(d string, _ am.ListRolesParams) (*am.PaginatedResponse, error) {
				if d != domainID {
					t.Errorf("expected domain %q, got %q", domainID, d)
				}

				data := toRaw(
					map[string]any{"id": "role-1", "name": "Admin", "type": "DOMAIN", "description": "Admin role"},
					map[string]any{"id": "role-2", "name": "User", "type": "DOMAIN", "description": "User role"},
				)

				return &am.PaginatedResponse{Data: data, TotalCount: 2, CurrentPage: 0}, nil
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestRole(tc, mock, domainID)

		err := testutil.Execute(newListCmd(tc.Factory, &domainID))

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Admin")
		testutil.AssertOutputContains(t, tc.Out, "User")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			ListRolesFunc: func(_ string, _ am.ListRolesParams) (*am.PaginatedResponse, error) {
				data := toRaw(map[string]any{"id": "role-1", "name": "Admin"})

				return &am.PaginatedResponse{Data: data, TotalCount: 1, CurrentPage: 0}, nil
			},
		}
		tc := testutil.NewFactory(nil, false)
		newTestRole(tc, mock, domainID)
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
