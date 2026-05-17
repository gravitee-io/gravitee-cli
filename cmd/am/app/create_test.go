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

func TestCreateApplication(t *testing.T) {
	t.Run("creates an application with name and type", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			CreateApplicationFunc: func(domainID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return json.Marshal(map[string]any{
					"id": "new-app", "name": "My App", "type": "web", "enabled": false,
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--name", "My App", "--type", "web")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My App")
		testutil.AssertOutputContains(t, tc.Out, "new-app")
	})

	t.Run("creates with description and redirect URIs", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			CreateApplicationFunc: func(_ string, body json.RawMessage) (json.RawMessage, error) {
				var m map[string]any
				_ = json.Unmarshal(body, &m)

				if m["description"] != "Desc" {
					t.Errorf("expected description 'Desc', got %v", m["description"])
				}

				uris, ok := m["redirectUris"].([]any)
				if !ok || len(uris) != 2 {
					t.Errorf("expected 2 redirect URIs, got %v", m["redirectUris"])
				}

				return json.Marshal(map[string]any{
					"id": "new-app", "name": "My App", "type": "web", "description": "Desc",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--name", "My App", "--type", "web",
			"--description", "Desc", "--redirect-uris", "http://a.com,http://b.com")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Desc")
	})

	t.Run("returns JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			CreateApplicationFunc: func(_ string, _ json.RawMessage) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "new-app", "name": "Test"})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "create", "--name", "Test", "--type", "web")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires name flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--type", "web")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires type flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--name", "Test")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("rejects invalid app type", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--name", "Test", "--type", "invalid")

		testutil.AssertErrorContains(t, err, "invalid value 'invalid' for flag --type")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--name", "Test", "--type", "web")

		testutil.AssertErrorContains(t, err, "no context configured")
	})

	t.Run("surfaces initial client secret for service apps", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			CreateApplicationFunc: func(_ string, _ json.RawMessage) (json.RawMessage, error) {
				return json.Marshal(map[string]any{
					"id":   "svc-1",
					"name": "Service",
					"type": "service",
					"settings": map[string]any{
						"oauth": map[string]any{
							"clientId":     "svc-client",
							"clientSecret": "supersecret",
						},
					},
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--name", "Service", "--type", "service")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "supersecret")
		testutil.AssertOutputContains(t, tc.Out, "svc-client")
	})
}
