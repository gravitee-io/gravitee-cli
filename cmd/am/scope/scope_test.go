package scope

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
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

	cmd := newListCmd(f)
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

	cmd := newCreateCmd(f)
	cmd.SetArgs([]string{"--key", "profile", "--name", "Profile", "--description", "Profile scope"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "profile") {
		t.Errorf("expected 'profile' in output, got: %s", out.String())
	}
}
