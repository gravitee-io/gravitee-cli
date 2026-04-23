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

func TestCreateDomain(t *testing.T) {
	t.Run("creates a domain with name", func(t *testing.T) {
		fake := &client.FakeClient{
			PostFunc: func(path string, body any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "domains")

				return json.Marshal(map[string]any{
					"id": "new-dom", "name": "My Domain", "enabled": false,
				})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newCreateCmd(tc.Factory), "--name", "My Domain")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Domain")
		testutil.AssertOutputContains(t, tc.Out, "new-dom")
	})

	t.Run("creates a domain with name and description", func(t *testing.T) {
		fake := &client.FakeClient{
			PostFunc: func(_ string, _ any) ([]byte, error) {
				return json.Marshal(map[string]any{
					"id": "new-dom", "name": "My Domain", "description": "Desc",
				})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newCreateCmd(tc.Factory), "--name", "My Domain", "--description", "Desc")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Desc")
	})

	t.Run("returns JSON with -o json", func(t *testing.T) {
		fake := &client.FakeClient{
			PostFunc: func(_ string, _ any) ([]byte, error) {
				return json.Marshal(map[string]any{"id": "new-dom", "name": "Test"})
			},
		}
		tc := testutil.NewFactory(fake)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newCreateCmd(tc.Factory), "--name", "Test")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("includes dataPlaneId in request body", func(t *testing.T) {
		var sentBody map[string]any
		fake := &client.FakeClient{
			PostFunc: func(_ string, body any) ([]byte, error) {
				raw, _ := body.(json.RawMessage)
				_ = json.Unmarshal(raw, &sentBody)

				return json.Marshal(map[string]any{"id": "new", "name": "Test"})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newCreateCmd(tc.Factory), "--name", "Test")

		testutil.AssertNoError(t, err)
		if sentBody["dataPlaneId"] != "default" {
			t.Errorf("expected dataPlaneId 'default', got %v", sentBody["dataPlaneId"])
		}
	})

	t.Run("requires name flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		err := testutil.Execute(newCreateCmd(tc.Factory))

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newCreateCmd(tc.Factory), "--name", "Test")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
