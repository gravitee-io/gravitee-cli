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

func TestCIMDGet(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(_ string) ([]byte, error) {
			return json.Marshal(map[string]any{
				"id": "dom-1",
				"oidc": map[string]any{
					"cimdSettings": map[string]any{
						"enabled":    true,
						"templateId": "tpl-1",
					},
				},
			})
		},
	}
	tc := testutil.NewFactory(fake)

	err := testutil.Execute(newCIMDGetCmd(tc.Factory), "dom-1")

	testutil.AssertNoError(t, err)
	testutil.AssertOutputContains(t, tc.Out, "tpl-1")
	testutil.AssertOutputContains(t, tc.Out, "true")
}

func TestCIMDEnable(t *testing.T) {
	t.Run("preserves existing oidc fields and merges settings", func(t *testing.T) {
		var captured map[string]any
		fake := &client.FakeClient{
			GetFunc: func(_ string) ([]byte, error) {
				return json.Marshal(map[string]any{
					"id": "dom-1",
					"oidc": map[string]any{
						"clientRegistrationSettings": map[string]any{"isAllowLocalhostRedirectUri": true},
						"cimdSettings":               map[string]any{"cacheTtlSeconds": 60.0},
					},
				})
			},
			PatchFunc: func(_ string, body any) ([]byte, error) {
				raw, _ := json.Marshal(body)
				_ = json.Unmarshal(raw, &captured)
				return raw, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newCIMDEnableCmd(tc.Factory),
			"dom-1",
			"--template-id", "tpl-1",
			"--allowed-domains", "a.com, b.com",
			"--allow-private",
		)

		testutil.AssertNoError(t, err)

		oidc, ok := captured["oidc"].(map[string]any)
		if !ok {
			t.Fatalf("expected oidc block, got %v", captured)
		}

		if _, found := oidc["clientRegistrationSettings"]; !found {
			t.Error("expected clientRegistrationSettings to be preserved on the PATCH body")
		}

		cimd, ok := oidc["cimdSettings"].(map[string]any)
		if !ok {
			t.Fatalf("expected oidc.cimdSettings map, got %v", oidc["cimdSettings"])
		}

		if cimd["enabled"] != true {
			t.Errorf("expected enabled=true, got %v", cimd["enabled"])
		}
		if cimd["templateId"] != "tpl-1" {
			t.Errorf("expected templateId=tpl-1, got %v", cimd["templateId"])
		}
		if cimd["allowPrivateIpAddress"] != true {
			t.Errorf("expected allowPrivateIpAddress=true, got %v", cimd["allowPrivateIpAddress"])
		}
		if ttl, _ := cimd["cacheTtlSeconds"].(float64); ttl != 60 {
			t.Errorf("expected cacheTtlSeconds=60 (preserved), got %v", cimd["cacheTtlSeconds"])
		}

		domains, _ := cimd["allowedDomains"].([]any)
		if len(domains) != 2 || domains[0] != "a.com" || domains[1] != "b.com" {
			t.Errorf("expected allowedDomains=[a.com b.com], got %v", domains)
		}
	})
}

func TestCIMDDisable(t *testing.T) {
	var captured map[string]any
	fake := &client.FakeClient{
		GetFunc: func(_ string) ([]byte, error) {
			return json.Marshal(map[string]any{
				"oidc": map[string]any{
					"cimdSettings": map[string]any{"enabled": true, "templateId": "tpl-1"},
				},
			})
		},
		PatchFunc: func(_ string, body any) ([]byte, error) {
			raw, _ := json.Marshal(body)
			_ = json.Unmarshal(raw, &captured)
			return raw, nil
		},
	}
	tc := testutil.NewFactory(fake)

	err := testutil.Execute(newCIMDDisableCmd(tc.Factory), "dom-1")

	testutil.AssertNoError(t, err)
	oidc, ok := captured["oidc"].(map[string]any)
	if !ok {
		t.Fatalf("expected oidc map, got %v", captured["oidc"])
	}
	cimd, ok := oidc["cimdSettings"].(map[string]any)
	if !ok {
		t.Fatalf("expected cimdSettings map, got %v", oidc["cimdSettings"])
	}
	if cimd["enabled"] != false {
		t.Errorf("expected enabled=false, got %v", cimd["enabled"])
	}
	if cimd["templateId"] != "tpl-1" {
		t.Errorf("expected templateId preserved, got %v", cimd["templateId"])
	}
}
