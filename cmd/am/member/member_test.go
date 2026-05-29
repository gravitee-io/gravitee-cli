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

	"gravitee.io/gctl/internal/am"
	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/testutil"
)

// --- List ---

func TestListMembers(t *testing.T) {
	t.Run("returns members", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			ListMembersFunc: func(domainID string) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return json.Marshal(map[string]any{
					"memberships": []map[string]any{
						{"memberId": "user-1", "role": "role-1", "memberType": "USER"},
						{"memberId": "user-2", "role": "role-2", "memberType": "USER"},
					},
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewMemberCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "user-1")
		testutil.AssertOutputContains(t, tc.Out, "user-2")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			ListMembersFunc: func(_ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{
					"memberships": []map[string]any{
						{"memberId": "user-1", "role": "role-1"},
					},
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewMemberCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"memberships"`)
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewMemberCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertErrorContains(t, err, "no context configured")
	})

	t.Run("requires domain flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewMemberCmd(tc.Factory)
		err := testutil.Execute(cmd, "list")

		testutil.AssertErrorContains(t, err, "required")
	})
}

// --- Add ---

func TestAddMember(t *testing.T) {
	t.Run("adds a member", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			AddMemberFunc: func(domainID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				var m map[string]string
				if err := json.Unmarshal(body, &m); err != nil {
					t.Fatalf("failed to unmarshal body: %v", err)
				}

				if m["memberId"] != "user-123" {
					t.Errorf("expected memberId 'user-123', got %q", m["memberId"])
				}

				if m["role"] != "role-456" {
					t.Errorf("expected role 'role-456', got %q", m["role"])
				}

				if m["memberType"] != "USER" {
					t.Errorf("expected memberType 'USER', got %q", m["memberType"])
				}

				return json.RawMessage(`{}`), nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewMemberCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "add", "--member-id", "user-123", "--role", "role-456")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Member 'user-123' added.")
	})

	t.Run("requires member-id flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewMemberCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "add", "--role", "role-456")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires role flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewMemberCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "add", "--member-id", "user-123")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewMemberCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "add", "--member-id", "user-123", "--role", "role-456")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Remove ---

func TestRemoveMember(t *testing.T) {
	t.Run("removes a member", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			RemoveMemberFunc: func(domainID, memberID string) error {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if memberID != "member-1" {
					t.Errorf("expected memberID 'member-1', got %q", memberID)
				}

				return nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewMemberCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "remove", "member-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Member 'member-1' removed.")
	})

	t.Run("returns error on API failure", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			RemoveMemberFunc: func(_, _ string) error {
				return &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewMemberCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "remove", "member-1")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("requires member ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewMemberCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "remove")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewMemberCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "remove", "member-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
