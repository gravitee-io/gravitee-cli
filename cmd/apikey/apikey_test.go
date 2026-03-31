package apikey

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func TestListSuccess(t *testing.T) {
	resp := map[string]interface{}{
		"data": []map[string]interface{}{
			{"key": "key-1", "revoked": false, "expired": false, "createdAt": "2026-03-20T10:00:00Z"},
			{"key": "key-2", "revoked": true, "expired": false, "createdAt": "2026-03-15T08:00:00Z"},
		},
		"pagination": map[string]int{
			"page": 1, "perPage": 10, "pageCount": 1, "totalCount": 2, "pageItemsCount": 2,
		},
	}

	data, _ := json.Marshal(resp)

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/subscriptions/sub-1/api-keys?") {
				t.Errorf("unexpected path: %s", path)
			}

			return data, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newListCmd(f)
	cmd.SetArgs([]string{"--api", "api-1", "--subscription", "sub-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "key-1") {
		t.Errorf("expected 'key-1' in output, got: %s", output)
	}

	if !strings.Contains(output, "key-2") {
		t.Errorf("expected 'key-2' in output, got: %s", output)
	}
}

func TestRenewSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]interface{}{
		"key":          "new-key-1",
		"subscription": "sub-1",
		"api":          "api-1",
		"revoked":      false,
		"expired":      false,
		"createdAt":    "2026-03-27T10:00:00Z",
	})

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/subscriptions/sub-1/api-keys/_renew") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newRenewCmd(f)
	cmd.SetArgs([]string{"--api", "api-1", "--subscription", "sub-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "new-key-1") {
		t.Errorf("expected 'new-key-1' in output, got: %s", output)
	}

	if !strings.Contains(output, "sub-1") {
		t.Errorf("expected subscription ID in output, got: %s", output)
	}
}

func TestRenewReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)

	cmd := newRenewCmd(f)
	cmd.SetArgs([]string{"--api", "api-1", "--subscription", "sub-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}

func TestRevokeSuccess(t *testing.T) {
	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/subscriptions/sub-1/api-keys/key-1/_revoke") {
				t.Errorf("unexpected path: %s", path)
			}

			return nil, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newRevokeCmd(f)
	cmd.SetArgs([]string{"key-1", "--api", "api-1", "--subscription", "sub-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "API key 'key-1' revoked.") {
		t.Errorf("unexpected output: %s", out.String())
	}
}

func TestRevokeReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)

	cmd := newRevokeCmd(f)
	cmd.SetArgs([]string{"key-1", "--api", "api-1", "--subscription", "sub-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}

func TestReactivateSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]interface{}{
		"key":          "key-1",
		"subscription": "sub-1",
		"api":          "api-1",
		"revoked":      false,
		"expired":      false,
		"createdAt":    "2026-03-20T10:00:00Z",
	})

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/subscriptions/sub-1/api-keys/key-1/_reactivate") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newReactivateCmd(f)
	cmd.SetArgs([]string{"key-1", "--api", "api-1", "--subscription", "sub-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "key-1") {
		t.Errorf("expected 'key-1' in output, got: %s", output)
	}

	if !strings.Contains(output, "sub-1") {
		t.Errorf("expected subscription ID in output, got: %s", output)
	}
}

func TestReactivateAPIError(t *testing.T) {
	fake := &client.FakeClient{
		PostFunc: func(_ string, _ interface{}) ([]byte, error) {
			return nil, fmt.Errorf("resource not found (HTTP 404)")
		},
	}

	f, _ := newTestFactory(fake, false)

	cmd := newReactivateCmd(f)
	cmd.SetArgs([]string{"key-999", "--api", "api-1", "--subscription", "sub-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected not found error, got: %v", err)
	}
}
