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

	"gravitee.io/gctl/internal/am"
	"gravitee.io/gctl/internal/testutil"
)

func TestCreateScope(t *testing.T) {
	t.Run("creates a scope with key and name", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			CreateScopeFunc: func(d string, body json.RawMessage) (json.RawMessage, error) {
				if d != domainID {
					t.Errorf("expected domain %q, got %q", domainID, d)
				}

				return json.Marshal(map[string]any{
					"id": "new-scope", "key": "openid", "name": "OpenID",
				})
			},
		}
		tc := testutil.NewFactory(nil)
		newTestScope(tc, mock, domainID)

		err := testutil.Execute(newCreateCmd(tc.Factory, &domainID), "--key", "openid", "--name", "OpenID")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "OpenID")
		testutil.AssertOutputContains(t, tc.Out, "new-scope")
	})

	t.Run("creates a scope with description", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			CreateScopeFunc: func(_ string, _ json.RawMessage) (json.RawMessage, error) {
				return json.Marshal(map[string]any{
					"id": "new-scope", "key": "openid", "name": "OpenID", "description": "OpenID scope",
				})
			},
		}
		tc := testutil.NewFactory(nil)
		newTestScope(tc, mock, domainID)

		err := testutil.Execute(newCreateCmd(tc.Factory, &domainID), "--key", "openid", "--name", "OpenID", "--description", "OpenID scope")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "OpenID scope")
	})

	t.Run("returns JSON with -o json", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			CreateScopeFunc: func(_ string, _ json.RawMessage) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "new-scope", "key": "openid", "name": "OpenID"})
			},
		}
		tc := testutil.NewFactory(nil)
		newTestScope(tc, mock, domainID)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newCreateCmd(tc.Factory, &domainID), "--key", "openid", "--name", "OpenID")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires key flag", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil)

		err := testutil.Execute(newCreateCmd(tc.Factory, &domainID), "--name", "OpenID")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires name flag", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil)

		err := testutil.Execute(newCreateCmd(tc.Factory, &domainID), "--key", "openid")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newCreateCmd(tc.Factory, &domainID), "--key", "openid", "--name", "OpenID")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
