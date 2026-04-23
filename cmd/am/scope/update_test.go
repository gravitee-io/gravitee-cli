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

package scope

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestUpdateScope(t *testing.T) {
	t.Run("updates scope name", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			PatchScopeFunc: func(d, scopeID string, _ json.RawMessage) (json.RawMessage, error) {
				if d != domainID {
					t.Errorf("expected domain %q, got %q", domainID, d)
				}
				if scopeID != "scope-1" {
					t.Errorf("expected scopeID %q, got %q", "scope-1", scopeID)
				}

				return json.Marshal(map[string]any{
					"id": "scope-1", "key": "openid", "name": "Updated",
				})
			},
		}
		tc := testutil.NewFactory(nil)
		newTestScope(tc, mock, domainID)

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID), "scope-1", "--name", "Updated")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Updated")
	})

	t.Run("updates scope description", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			PatchScopeFunc: func(_, _ string, _ json.RawMessage) (json.RawMessage, error) {
				return json.Marshal(map[string]any{
					"id": "scope-1", "key": "openid", "name": "OpenID", "description": "New desc",
				})
			},
		}
		tc := testutil.NewFactory(nil)
		newTestScope(tc, mock, domainID)

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID), "scope-1", "--description", "New desc")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "New desc")
	})

	t.Run("requires at least one flag", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil)

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID), "scope-1")

		testutil.AssertErrorContains(t, err, "at least one flag")
	})

	t.Run("requires scope ID argument", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil)

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID))

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newUpdateCmd(tc.Factory, &domainID), "scope-1", "--name", "Test")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
