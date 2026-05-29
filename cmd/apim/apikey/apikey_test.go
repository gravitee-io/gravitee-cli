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

package apikey

import (
	"encoding/json"
	"fmt"
	"testing"

	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/testutil"
)

func TestListAPIKeys(t *testing.T) {
	t.Run("returns API keys for the subscription", func(t *testing.T) {
		fake := paginatedAPIKeys(
			map[string]any{"key": "key-1", "revoked": false, "expired": false, "createdAt": "2026-03-20T10:00:00Z"},
			map[string]any{"key": "key-2", "revoked": true, "expired": false, "createdAt": "2026-03-15T08:00:00Z"},
		)
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newListCmd(tc.Factory), "--api", "api-1", "--subscription", "sub-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "key-1")
		testutil.AssertOutputContains(t, tc.Out, "key-2")
	})
}

func TestRenewAPIKey(t *testing.T) {
	t.Run("renews the API key", func(t *testing.T) {
		resp, _ := json.Marshal(map[string]any{
			"key": "new-key-1", "subscription": "sub-1", "api": "api-1",
			"revoked": false, "expired": false, "createdAt": "2026-03-27T10:00:00Z",
		})
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/subscriptions/sub-1/api-keys/_renew")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newRenewCmd(tc.Factory), "--api", "api-1", "--subscription", "sub-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "new-key-1")
		testutil.AssertOutputContains(t, tc.Out, "sub-1")
	})
}

func TestRevokeAPIKey(t *testing.T) {
	t.Run("revokes the API key", func(t *testing.T) {
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/subscriptions/sub-1/api-keys/key-1/_revoke")

				return nil, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newRevokeCmd(tc.Factory), "key-1", "--api", "api-1", "--subscription", "sub-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "revoked")
	})
}

func TestReactivateAPIKey(t *testing.T) {
	t.Run("reactivates the API key", func(t *testing.T) {
		resp, _ := json.Marshal(map[string]any{
			"key": "key-1", "subscription": "sub-1", "api": "api-1",
			"revoked": false, "expired": false, "createdAt": "2026-03-20T10:00:00Z",
		})
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/subscriptions/sub-1/api-keys/key-1/_reactivate")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newReactivateCmd(tc.Factory), "key-1", "--api", "api-1", "--subscription", "sub-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "key-1")
		testutil.AssertOutputContains(t, tc.Out, "sub-1")
	})

	t.Run("reports not found error", func(t *testing.T) {
		fake := &client.FakeClient{
			PostFunc: func(_ string, _ any) ([]byte, error) {
				return nil, fmt.Errorf("resource not found (HTTP 404)")
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newReactivateCmd(tc.Factory), "key-999", "--api", "api-1", "--subscription", "sub-1")

		testutil.AssertErrorContains(t, err, "not found")
	})
}
