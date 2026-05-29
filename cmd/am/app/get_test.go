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

	"gravitee.io/gctl/internal/am"
	"gravitee.io/gctl/internal/testutil"
)

func TestGetApplication(t *testing.T) {
	t.Run("returns application details", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetApplicationFunc: func(domainID, appID string) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if appID != "app-1" {
					t.Errorf("expected appID 'app-1', got %q", appID)
				}

				return json.Marshal(map[string]any{
					"id": "app-1", "name": "My App", "type": "web", "enabled": true,
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "app-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My App")
		testutil.AssertOutputContains(t, tc.Out, "app-1")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetApplicationFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "app-1", "name": "Test"})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "get", "app-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires app ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "app-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
