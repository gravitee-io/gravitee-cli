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

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestUpdateDomain(t *testing.T) {
	t.Run("updates domain name", func(t *testing.T) {
		fake := &client.FakeClient{
			PatchFunc: func(path string, body any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "domains/dom-1")

				return json.Marshal(map[string]any{
					"id": "dom-1", "name": "Updated", "enabled": true,
				})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newUpdateCmd(tc.Factory), "dom-1", "--name", "Updated")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Updated")
	})

	t.Run("updates domain description", func(t *testing.T) {
		fake := &client.FakeClient{
			PatchFunc: func(_ string, _ any) ([]byte, error) {
				return json.Marshal(map[string]any{
					"id": "dom-1", "name": "Test", "description": "New desc",
				})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newUpdateCmd(tc.Factory), "dom-1", "--description", "New desc")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "New desc")
	})

	t.Run("requires at least one flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		err := testutil.Execute(newUpdateCmd(tc.Factory), "dom-1")

		testutil.AssertErrorContains(t, err, "at least one flag")
	})

	t.Run("requires domain ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		err := testutil.Execute(newUpdateCmd(tc.Factory))

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newUpdateCmd(tc.Factory), "dom-1", "--name", "Test")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
