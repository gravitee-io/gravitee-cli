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

package scope

import (
	"encoding/json"
	"strings"
	"testing"

	"gravitee.io/gctl/internal/client"
)

func TestScopeList(t *testing.T) {
	resp := map[string]interface{}{
		"data": []map[string]interface{}{
			{"id": "scope-1", "key": "openid", "name": "OpenID", "description": "OpenID Connect scope"},
		},
		"currentPage": 0,
		"totalCount":  1,
	}

	data, _ := json.Marshal(resp)

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/scopes?") {
				t.Errorf("unexpected path: %s", path)
			}
			return data, nil
		},
	}

	f, out := newTestFactory(fake, false)
	domainID := "test-domain"

	cmd := newListCmd(f, &domainID)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "openid") {
		t.Errorf("expected 'openid' in output, got: %s", out.String())
	}
}

func TestScopeCreateWithFlags(t *testing.T) {
	fake := &client.FakeClient{
		PostFunc: func(path string, body interface{}) ([]byte, error) {
			if !strings.Contains(path, "/scopes") {
				t.Errorf("unexpected path: %s", path)
			}

			var m map[string]interface{}
			switch b := body.(type) {
			case []byte:
				_ = json.Unmarshal(b, &m)
			case json.RawMessage:
				_ = json.Unmarshal(b, &m)
			}

			if key, ok := m["key"].(string); !ok || key != "profile" {
				t.Errorf("expected key 'profile', got: %v", m["key"])
			}
			if name, ok := m["name"].(string); !ok || name != "Profile" {
				t.Errorf("expected name 'Profile', got: %v", m["name"])
			}

			resp := map[string]interface{}{"id": "scope-new", "key": "profile"}
			data, _ := json.Marshal(resp)
			return data, nil
		},
	}

	f, out := newTestFactory(fake, false)
	domainID := "test-domain"

	cmd := newCreateCmd(f, &domainID)
	cmd.SetArgs([]string{"--key", "profile", "--name", "Profile", "--description", "Profile scope"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "profile") {
		t.Errorf("expected 'profile' in output, got: %s", out.String())
	}
}

func TestScopeGet(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/scopes/scope-1") {
				t.Errorf("unexpected path: %s", path)
			}
			return []byte(`{"id":"scope-1","key":"openid","name":"OpenID"}`), nil
		},
	}
	f, out := newTestFactory(fake, false)
	domainID := "test-domain"
	cmd := newGetCmd(f, &domainID)
	cmd.SetArgs([]string{"scope-1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "scope-1") {
		t.Errorf("expected 'scope-1' in output, got: %s", out.String())
	}
}

func TestScopeDelete(t *testing.T) {
	deleted := false
	fake := &client.FakeClient{
		DeleteFunc: func(path string) error {
			if !strings.Contains(path, "/scopes/scope-1") {
				t.Errorf("unexpected path: %s", path)
			}
			deleted = true
			return nil
		},
	}
	f, _ := newTestFactory(fake, false)
	domainID := "test-domain"
	cmd := newDeleteCmd(f, &domainID)
	cmd.SetArgs([]string{"scope-1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !deleted {
		t.Error("expected Delete to be called")
	}
}

func TestScopeCreateReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)
	domainID := "test-domain"
	cmd := newCreateCmd(f, &domainID)
	cmd.SetArgs([]string{"--key", "openid", "--name", "OpenID"})
	if err := cmd.Execute(); err == nil {
		t.Error("expected read-only error")
	}
}
