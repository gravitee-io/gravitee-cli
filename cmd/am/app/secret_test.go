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

package app

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestSecretRenew(t *testing.T) {
	t.Run("renews a secret and prints the new value", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			RenewAppSecretFunc: func(domainID, appID, secretID string) (json.RawMessage, error) {
				if domainID != "dom-1" || appID != "app-1" || secretID != "sec-1" {
					t.Errorf("unexpected args: %s/%s/%s", domainID, appID, secretID)
				}
				return json.Marshal(map[string]any{"id": "sec-1", "secret": "rotatedvalue"})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "secret", "--app-id", "app-1", "renew", "sec-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "rotatedvalue")
		testutil.AssertOutputContains(t, tc.Out, "sec-1")
	})

	t.Run("falls back to clientSecret key", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			RenewAppSecretFunc: func(_, _, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "sec-1", "clientSecret": "alt-value"})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "secret", "--app-id", "app-1", "renew", "sec-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "alt-value")
	})

	t.Run("emits JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			RenewAppSecretFunc: func(_, _, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "sec-1", "secret": "x"})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "secret", "--app-id", "app-1", "renew", "sec-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"secret"`)
	})
}
