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
	"testing"

	"gravitee.io/gctl/internal/am"
	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/testutil"
)

func TestDeleteScope(t *testing.T) {
	t.Run("deletes a scope", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			DeleteScopeFunc: func(d, scopeID string) error {
				if d != domainID {
					t.Errorf("expected domain %q, got %q", domainID, d)
				}
				if scopeID != "scope-1" {
					t.Errorf("expected scopeID %q, got %q", "scope-1", scopeID)
				}

				return nil
			},
		}
		tc := testutil.NewFactory(nil)
		newTestScope(tc, mock, domainID)

		err := testutil.Execute(newDeleteCmd(tc.Factory, &domainID), "scope-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Scope 'scope-1' deleted.")
	})

	t.Run("returns error on API failure", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			DeleteScopeFunc: func(_, _ string) error {
				return &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
			},
		}
		tc := testutil.NewFactory(nil)
		newTestScope(tc, mock, domainID)

		err := testutil.Execute(newDeleteCmd(tc.Factory, &domainID), "scope-1")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("requires scope ID argument", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil)

		err := testutil.Execute(newDeleteCmd(tc.Factory, &domainID))

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newDeleteCmd(tc.Factory, &domainID), "scope-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
