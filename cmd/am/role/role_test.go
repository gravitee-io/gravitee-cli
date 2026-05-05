package role

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
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

	cmd := newListCmd(f)
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
			if assignableType, ok := m["assignableType"].(string); !ok || assignableType != "DOMAIN" {
				t.Errorf("expected assignableType 'DOMAIN', got: %v", m["assignableType"])
			}

			resp := map[string]interface{}{"id": "role-new", "name": "Admin Role"}
			data, _ := json.Marshal(resp)
			return data, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newCreateCmd(f)
	cmd.SetArgs([]string{"--name", "Admin Role", "--description", "Admin desc", "--type", "DOMAIN"})

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
	cmd := newGetCmd(f)
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
	cmd := newDeleteCmd(f)
	cmd.SetArgs([]string{"role-1", "--force"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !deleted {
		t.Error("expected Delete to be called")
	}
}

func TestRoleCreateReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)
	cmd := newCreateCmd(f)
	cmd.SetArgs([]string{"--name", "Admin", "--type", "DOMAIN"})
	if err := cmd.Execute(); err == nil {
		t.Error("expected read-only error")
	}
}
