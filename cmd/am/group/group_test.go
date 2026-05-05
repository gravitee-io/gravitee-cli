package group

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
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
	cmd := newListCmd(f)
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
	cmd := newGetCmd(f)
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
	cmd := newCreateCmd(f)
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
	cmd := newDeleteCmd(f)
	cmd.SetArgs([]string{"group-1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !deleted {
		t.Error("expected Delete to be called")
	}
}
