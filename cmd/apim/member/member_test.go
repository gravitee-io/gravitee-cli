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

package member

import (
	"encoding/json"
	"testing"

	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/testutil"
)

func TestListMembers(t *testing.T) {
	t.Run("returns members from the API", func(t *testing.T) {
		fake := paginatedMembers(
			map[string]any{
				"id":          "aaaa1111-2222-3333-4444-555566667777",
				"displayName": "Alice Martin",
				"roles":       []map[string]any{{"name": "PRIMARY_OWNER", "scope": "API"}},
				"type":        "USER",
			},
			map[string]any{
				"id":          "bbbb1111-2222-3333-4444-555566667777",
				"displayName": "Bob Dupont",
				"roles":       []map[string]any{{"name": "OWNER", "scope": "API"}},
				"type":        "USER",
			},
		)
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newListCmd(tc.Factory), "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Alice Martin")
		testutil.AssertOutputContains(t, tc.Out, "PRIMARY_OWNER")
		testutil.AssertOutputContains(t, tc.Out, "Bob Dupont")
	})

	t.Run("rejects invalid token with hint", func(t *testing.T) {
		fake := testutil.APIFailingWith(401, "authentication failed")
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newListCmd(tc.Factory), "--api", "api-1")

		testutil.AssertErrorContains(t, err, "authentication failed")
	})
}

func TestAddMember(t *testing.T) {
	t.Run("adds a member to the API", func(t *testing.T) {
		resp, _ := json.Marshal(map[string]any{
			"id":          "bbbb1111-2222-3333-4444-555566667777",
			"displayName": "Bob Dupont",
			"roles":       []map[string]any{{"name": "OWNER", "scope": "API"}},
			"type":        "USER",
		})
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/members")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newAddCmd(tc.Factory), "--api", "api-1", "--user", "bbbb1111-2222-3333-4444-555566667777", "--role", "OWNER")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Bob Dupont")
		testutil.AssertOutputContains(t, tc.Out, "OWNER")
	})
}

func TestRemoveMember(t *testing.T) {
	t.Run("removes a member from the API", func(t *testing.T) {
		fake := &client.FakeClient{
			DeleteFunc: func(path string) error {
				testutil.AssertPathCalled(t, path, "/apis/api-1/members/member-1")

				return nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newRemoveCmd(tc.Factory), "member-1", "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Member 'member-1' removed.")
	})
}
