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

package environment

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestListEnvironments(t *testing.T) {
	t.Run("returns environments from the organization", func(t *testing.T) {
		fake := environmentList(
			map[string]string{"id": "dev-1111", "name": "Development", "description": "Development environment"},
			map[string]string{"id": "prod-2222", "name": "Production", "description": "Production environment"},
		)
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Development")
		testutil.AssertOutputContains(t, tc.Out, "Production")
	})

	t.Run("calls the correct API path", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/management/organizations/DEFAULT/environments")

				data, _ := json.Marshal([]map[string]string{})

				return data, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertNoError(t, err)
	})

	t.Run("rejects invalid token with hint", func(t *testing.T) {
		fake := testutil.APIFailingWith(401, "authentication failed (HTTP 401)")
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertErrorContains(t, err, "authentication failed")
	})
}

func TestGetEnvironment(t *testing.T) {
	t.Run("returns environment details", func(t *testing.T) {
		fake := testutil.APIReturningItem(map[string]any{
			"id": "prod-2222", "name": "Production", "description": "Production environment",
		})
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newGetCmd(tc.Factory), "prod-2222")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Production")
		testutil.AssertOutputContains(t, tc.Out, "prod-2222")
	})

	t.Run("calls the correct API path", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/management/organizations/DEFAULT/environments/prod-2222")

				resp, _ := json.Marshal(map[string]string{"id": "prod-2222", "name": "Production"})

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newGetCmd(tc.Factory), "prod-2222")

		testutil.AssertNoError(t, err)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		fake := testutil.APIFailingWith(404, "resource not found (HTTP 404)")
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newGetCmd(tc.Factory), "env-999")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("rejects empty envId before calling the API", func(t *testing.T) {
		called := false
		fake := &client.FakeClient{
			GetFunc: func(_ string) ([]byte, error) {
				called = true

				return nil, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newGetCmd(tc.Factory), "")

		testutil.AssertErrorContains(t, err, "envId cannot be empty")

		if called {
			t.Fatal("expected no API call for empty envId")
		}
	})
}
