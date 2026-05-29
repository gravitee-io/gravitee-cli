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

func TestGetDomain(t *testing.T) {
	t.Run("returns domain details", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "domains/dom-1")

				return json.Marshal(map[string]any{
					"id": "dom-1", "name": "My Domain", "enabled": true, "description": "A domain",
				})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newGetCmd(tc.Factory), "dom-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Domain")
		testutil.AssertOutputContains(t, tc.Out, "dom-1")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(_ string) ([]byte, error) {
				return json.Marshal(map[string]any{"id": "dom-1", "name": "Test"})
			},
		}
		tc := testutil.NewFactory(fake)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newGetCmd(tc.Factory), "dom-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires domain ID", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		err := testutil.Execute(newGetCmd(tc.Factory))

		testutil.AssertErrorContains(t, err, "domain ID or --hrid is required")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newGetCmd(tc.Factory), "dom-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
