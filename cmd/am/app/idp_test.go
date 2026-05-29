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

func TestAppIdpAdd(t *testing.T) {
	t.Run("appends a new binding preserving existing ones", func(t *testing.T) {
		var captured map[string]any
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetApplicationFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{
					"id": "app-1",
					"identityProviders": []any{
						map[string]any{"identity": "existing", "priority": 0.0},
					},
				})
			},
			PatchApplicationFunc: func(_, _ string, body json.RawMessage) (json.RawMessage, error) {
				_ = json.Unmarshal(body, &captured)
				return json.Marshal(map[string]any{
					"identityProviders": []any{
						map[string]any{"identity": "existing", "priority": 0.0},
						map[string]any{"identity": "new-idp", "priority": 10.0},
					},
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "idp", "add", "app-1", "new-idp", "--priority", "10")

		testutil.AssertNoError(t, err)
		idps, ok := captured["identityProviders"].([]any)
		if !ok || len(idps) != 2 {
			t.Fatalf("expected 2 identityProviders, got %v", captured)
		}
		second, _ := idps[1].(map[string]any)
		if second["identity"] != "new-idp" {
			t.Errorf("expected appended binding identity=new-idp, got %v", second["identity"])
		}
		if prio, _ := second["priority"].(float64); prio != 10 {
			t.Errorf("expected priority=10, got %v", second["priority"])
		}
		testutil.AssertOutputContains(t, tc.Out, "new-idp")
	})

	t.Run("updates an existing binding rather than duplicating", func(t *testing.T) {
		var captured map[string]any
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetApplicationFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{
					"identityProviders": []any{
						map[string]any{"identity": "idp-1", "priority": 0.0},
					},
				})
			},
			PatchApplicationFunc: func(_, _ string, body json.RawMessage) (json.RawMessage, error) {
				_ = json.Unmarshal(body, &captured)
				return body, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "idp", "add", "app-1", "idp-1",
			"--priority", "5", "--selection-rule", "{#context.attributes['foo'] == 'bar'}")

		testutil.AssertNoError(t, err)
		idps, _ := captured["identityProviders"].([]any)
		if len(idps) != 1 {
			t.Fatalf("expected 1 binding, got %d", len(idps))
		}
		b, _ := idps[0].(map[string]any)
		if prio, _ := b["priority"].(float64); prio != 5 {
			t.Errorf("expected priority=5, got %v", b["priority"])
		}
		if b["selectionRule"] != "{#context.attributes['foo'] == 'bar'}" {
			t.Errorf("expected selectionRule preserved, got %v", b["selectionRule"])
		}
	})
}

func TestAppIdpRemove(t *testing.T) {
	t.Run("removes a binding", func(t *testing.T) {
		var captured map[string]any
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetApplicationFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{
					"identityProviders": []any{
						map[string]any{"identity": "idp-1"},
						map[string]any{"identity": "idp-2"},
					},
				})
			},
			PatchApplicationFunc: func(_, _ string, body json.RawMessage) (json.RawMessage, error) {
				_ = json.Unmarshal(body, &captured)
				return body, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "idp", "remove", "app-1", "idp-1")

		testutil.AssertNoError(t, err)
		idps, _ := captured["identityProviders"].([]any)
		if len(idps) != 1 {
			t.Fatalf("expected 1 binding remaining, got %d", len(idps))
		}
		first, _ := idps[0].(map[string]any)
		if first["identity"] != "idp-2" {
			t.Errorf("expected idp-2 to remain, got %v", idps[0])
		}
	})

	t.Run("errors when binding not present", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetApplicationFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"identityProviders": []any{}})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "idp", "remove", "app-1", "missing")

		testutil.AssertErrorContains(t, err, "is not bound")
	})
}

func TestAppIdpList(t *testing.T) {
	t.Run("prints bindings table", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetApplicationFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{
					"identityProviders": []any{
						map[string]any{"identity": "idp-1", "priority": 0.0, "selectionRule": "rule"},
					},
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "idp", "list", "app-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "idp-1")
		testutil.AssertOutputContains(t, tc.Out, "rule")
	})

	t.Run("prints empty message when none bound", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetApplicationFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "idp", "list", "app-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "No identity providers bound")
	})
}
