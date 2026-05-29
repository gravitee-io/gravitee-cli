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

package domain

import (
	"encoding/json"
	"testing"

	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/testutil"
)

func TestEnableDomain(t *testing.T) {
	t.Run("enables a domain", func(t *testing.T) {
		fake := &client.FakeClient{
			PatchFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "domains/dom-1")

				return json.Marshal(map[string]any{"id": "dom-1", "enabled": true})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newEnableCmd(tc.Factory), "dom-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Domain 'dom-1' enabled.")
	})

	t.Run("requires domain ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		err := testutil.Execute(newEnableCmd(tc.Factory))

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newEnableCmd(tc.Factory), "dom-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

func TestDisableDomain(t *testing.T) {
	t.Run("disables a domain", func(t *testing.T) {
		fake := &client.FakeClient{
			PatchFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "domains/dom-1")

				return json.Marshal(map[string]any{"id": "dom-1", "enabled": false})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newDisableCmd(tc.Factory), "dom-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Domain 'dom-1' disabled.")
	})

	t.Run("requires domain ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		err := testutil.Execute(newDisableCmd(tc.Factory))

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})
}
