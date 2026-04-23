// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
		tc := testutil.NewFactory(nil)
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
		tc := testutil.NewFactory(nil)
		newTestUser(tc, mock, domainID)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newListCmd(tc.Factory, &domainID))

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"totalCount"`)
	})

	t.Run("requires a configured context", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newListCmd(tc.Factory, &domainID))

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
