package member

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func TestListSuccess(t *testing.T) {
	resp := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"id":          "aaaa1111-2222-3333-4444-555566667777",
				"displayName": "Alice Martin",
				"roles":       []map[string]interface{}{{"name": "PRIMARY_OWNER", "scope": "API"}},
				"type":        "USER",
			},
			{
				"id":          "bbbb1111-2222-3333-4444-555566667777",
				"displayName": "Bob Dupont",
				"roles":       []map[string]interface{}{{"name": "OWNER", "scope": "API"}},
				"type":        "USER",
			},
		},
		"pagination": map[string]int{
			"page": 1, "perPage": 10, "pageCount": 1, "totalCount": 2, "pageItemsCount": 2,
		},
	}

	data, _ := json.Marshal(resp)

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/members?") {
				t.Errorf("unexpected path: %s", path)
			}

			return data, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newListCmd(f)
	cmd.SetArgs([]string{"--api", "api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()

	if !strings.Contains(output, "Alice Martin") {
		t.Errorf("expected 'Alice Martin' in output, got: %s", output)
	}

	if !strings.Contains(output, "PRIMARY_OWNER") {
		t.Errorf("expected 'PRIMARY_OWNER' in output, got: %s", output)
	}

	if !strings.Contains(output, "Bob Dupont") {
		t.Errorf("expected 'Bob Dupont' in output, got: %s", output)
	}
}

func TestListAPIError(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(_ string) ([]byte, error) {
			return nil, &client.APIError{Status: 401, Message: "authentication failed (HTTP 401)"}
		},
	}

	f, _ := newTestFactory(fake, false)

	cmd := newListCmd(f)
	cmd.SetArgs([]string{"--api", "api-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "authentication failed") {
		t.Errorf("expected auth error, got: %v", err)
	}
}

func TestAddSuccess(t *testing.T) {
	resp := map[string]interface{}{
		"id":          "bbbb1111-2222-3333-4444-555566667777",
		"displayName": "Bob Dupont",
		"roles":       []map[string]interface{}{{"name": "OWNER", "scope": "API"}},
		"type":        "USER",
	}

	data, _ := json.Marshal(resp)

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.HasSuffix(path, "/apis/api-1/members") {
				t.Errorf("unexpected path: %s", path)
			}

			return data, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newAddCmd(f)
	cmd.SetArgs([]string{"--api", "api-1", "--user", "bbbb1111-2222-3333-4444-555566667777", "--role", "OWNER"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()

	if !strings.Contains(output, "Bob Dupont") {
		t.Errorf("expected 'Bob Dupont' in output, got: %s", output)
	}

	if !strings.Contains(output, "OWNER") {
		t.Errorf("expected 'OWNER' in output, got: %s", output)
	}
}

func TestAddReadOnly(t *testing.T) {
	fake := &client.FakeClient{}
	f, _ := newTestFactory(fake, true)

	cmd := newAddCmd(f)
	cmd.SetArgs([]string{"--api", "api-1", "--user", "user-1", "--role", "OWNER"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}

func TestRemoveSuccess(t *testing.T) {
	fake := &client.FakeClient{
		DeleteFunc: func(path string) error {
			if !strings.Contains(path, "/apis/api-1/members/member-1") {
				t.Errorf("unexpected path: %s", path)
			}

			return nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newRemoveCmd(f)
	cmd.SetArgs([]string{"member-1", "--api", "api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Member 'member-1' removed.") {
		t.Errorf("unexpected output: %s", out.String())
	}
}

func TestRemoveReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)

	cmd := newRemoveCmd(f)
	cmd.SetArgs([]string{"member-1", "--api", "api-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}
