package metadata

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func TestListSuccess(t *testing.T) {
	resp := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"key": "team-email", "name": "Team Email",
				"value": "platform-team@company.com", "format": "MAIL",
				"updatedAt": "2026-03-25T14:30:00Z",
			},
		},
		"pagination": map[string]int{
			"page": 1, "perPage": 10, "pageCount": 1, "totalCount": 1, "pageItemsCount": 1,
		},
	}

	data, _ := json.Marshal(resp)

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/metadata?") {
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

	if !strings.Contains(out.String(), "Team Email") {
		t.Errorf("expected 'Team Email' in output, got: %s", out.String())
	}

	if !strings.Contains(out.String(), "MAIL") {
		t.Errorf("expected 'MAIL' in output, got: %s", out.String())
	}
}

func TestCreateSuccess(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "metadata.json")
	_ = os.WriteFile(file, []byte(`{"name":"Team Email","value":"platform-team@company.com","format":"MAIL"}`), 0600)

	resp, _ := json.Marshal(map[string]interface{}{
		"key": "team-email", "name": "Team Email",
		"value": "platform-team@company.com", "format": "MAIL",
	})

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.HasSuffix(path, "/apis/api-1/metadata") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newCreateCmd(f)
	cmd.SetArgs([]string{"--api", "api-1", "-f", file})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Team Email") {
		t.Errorf("expected 'Team Email' in output, got: %s", out.String())
	}

	if !strings.Contains(out.String(), "team-email") {
		t.Errorf("expected 'team-email' in output, got: %s", out.String())
	}
}

func TestCreateReadOnly(t *testing.T) {
	fake := &client.FakeClient{}
	f, _ := newTestFactory(fake, true)

	cmd := newCreateCmd(f)
	cmd.SetArgs([]string{"--api", "api-1", "-f", "test.json"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}

func TestUpdateSuccess(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "metadata.json")
	_ = os.WriteFile(file, []byte(`{"name":"Team Email","value":"new-team@company.com","format":"MAIL"}`), 0600)

	resp, _ := json.Marshal(map[string]interface{}{
		"key": "team-email", "name": "Team Email",
		"value": "new-team@company.com", "format": "MAIL",
	})

	fake := &client.FakeClient{
		PutFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/metadata/team-email") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newUpdateCmd(f)
	cmd.SetArgs([]string{"team-email", "--api", "api-1", "-f", file})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "new-team@company.com") {
		t.Errorf("expected 'new-team@company.com' in output, got: %s", out.String())
	}
}

func TestDeleteSuccess(t *testing.T) {
	fake := &client.FakeClient{
		DeleteFunc: func(path string) error {
			if !strings.Contains(path, "/apis/api-1/metadata/team-email") {
				t.Errorf("unexpected path: %s", path)
			}

			return nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newDeleteCmd(f)
	cmd.SetArgs([]string{"team-email", "--api", "api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Metadata 'team-email' deleted.") {
		t.Errorf("unexpected output: %s", out.String())
	}
}

func TestDeleteReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)

	cmd := newDeleteCmd(f)
	cmd.SetArgs([]string{"team-email", "--api", "api-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}
