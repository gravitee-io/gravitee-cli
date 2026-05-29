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

package group

import (
	"encoding/json"
	"strings"
	"testing"

	"gravitee.io/gctl/internal/client"
)

func TestGroupList(t *testing.T) {
	resp := map[string]interface{}{
		"data":        []map[string]interface{}{{"id": "group-1", "name": "Admins", "description": "Admin group"}},
		"currentPage": 0,
		"totalCount":  1,
	}
	data, _ := json.Marshal(resp)
	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/groups?") {
				t.Errorf("unexpected path: %s", path)
			}
			return data, nil
		},
	}
	f, out := newTestFactory(fake, false)
	domainID := "test-domain"
	cmd := newListCmd(f, &domainID)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "Admins") {
		t.Errorf("expected 'Admins' in output, got: %s", out.String())
	}
}

func TestGroupGet(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/groups/group-1") {
				t.Errorf("unexpected path: %s", path)
			}
			return []byte(`{"id":"group-1","name":"Admins","description":"Admin group"}`), nil
		},
	}
	f, out := newTestFactory(fake, false)
	domainID := "test-domain"
	cmd := newGetCmd(f, &domainID)
	cmd.SetArgs([]string{"group-1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "group-1") {
		t.Errorf("expected 'group-1' in output, got: %s", out.String())
	}
}

func TestGroupCreate(t *testing.T) {
	fake := &client.FakeClient{
		PostFunc: func(path string, body interface{}) ([]byte, error) {
			if !strings.Contains(path, "/groups") {
				t.Errorf("unexpected path: %s", path)
			}
			return []byte(`{"id":"group-new","name":"DevTeam"}`), nil
		},
	}
	f, out := newTestFactory(fake, false)
	domainID := "test-domain"
	cmd := newCreateCmd(f, &domainID)
	cmd.SetArgs([]string{"--name", "DevTeam", "--description", "Developers"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "DevTeam") {
		t.Errorf("expected 'DevTeam' in output, got: %s", out.String())
	}
}

func TestGroupDelete(t *testing.T) {
	deleted := false
	fake := &client.FakeClient{
		DeleteFunc: func(path string) error {
			if !strings.Contains(path, "/groups/group-1") {
				t.Errorf("unexpected path: %s", path)
			}
			deleted = true
			return nil
		},
	}
	f, _ := newTestFactory(fake, false)
	domainID := "test-domain"
	cmd := newDeleteCmd(f, &domainID)
	cmd.SetArgs([]string{"group-1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !deleted {
		t.Error("expected Delete to be called")
	}
}
