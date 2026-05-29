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

package role

import (
	"encoding/json"
	"strings"
	"testing"

	"gravitee.io/gctl/internal/client"
)

func TestRoleList(t *testing.T) {
	resp := map[string]interface{}{
		"data": []map[string]interface{}{
			{"id": "role-1", "name": "My Role", "assignableType": "DOMAIN", "description": "Test role"},
		},
		"currentPage": 0,
		"totalCount":  1,
	}

	data, _ := json.Marshal(resp)

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/roles?") {
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

	if !strings.Contains(out.String(), "My Role") {
		t.Errorf("expected 'My Role' in output, got: %s", out.String())
	}
}

func TestRoleCreateWithFlags(t *testing.T) {
	fake := &client.FakeClient{
		PostFunc: func(path string, body interface{}) ([]byte, error) {
			if !strings.Contains(path, "/roles") {
				t.Errorf("unexpected path: %s", path)
			}

			var m map[string]interface{}
			switch b := body.(type) {
			case []byte:
				_ = json.Unmarshal(b, &m)
			case json.RawMessage:
				_ = json.Unmarshal(b, &m)
			}

			if name, ok := m["name"].(string); !ok || name != "Admin Role" {
				t.Errorf("expected name 'Admin Role', got: %v", m["name"])
			}

			resp := map[string]interface{}{"id": "role-new", "name": "Admin Role"}
			data, _ := json.Marshal(resp)
			return data, nil
		},
	}

	f, out := newTestFactory(fake, false)
	domainID := "test-domain"

	cmd := newCreateCmd(f, &domainID)
	cmd.SetArgs([]string{"--name", "Admin Role", "--description", "Admin desc"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Admin Role") {
		t.Errorf("expected 'Admin Role' in output, got: %s", out.String())
	}
}

func TestRoleGet(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/roles/role-1") {
				t.Errorf("unexpected path: %s", path)
			}
			return []byte(`{"id":"role-1","name":"My Role","assignableType":"DOMAIN"}`), nil
		},
	}
	f, out := newTestFactory(fake, false)
	domainID := "test-domain"
	cmd := newGetCmd(f, &domainID)
	cmd.SetArgs([]string{"role-1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "role-1") {
		t.Errorf("expected 'role-1' in output, got: %s", out.String())
	}
}

func TestRoleDelete(t *testing.T) {
	deleted := false
	fake := &client.FakeClient{
		DeleteFunc: func(path string) error {
			if !strings.Contains(path, "/roles/role-1") {
				t.Errorf("unexpected path: %s", path)
			}
			deleted = true
			return nil
		},
	}
	f, _ := newTestFactory(fake, false)
	domainID := "test-domain"
	cmd := newDeleteCmd(f, &domainID)
	cmd.SetArgs([]string{"role-1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !deleted {
		t.Error("expected Delete to be called")
	}
}
