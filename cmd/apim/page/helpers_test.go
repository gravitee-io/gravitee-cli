package page

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func fakePagesResponse(items ...map[string]any) *client.FakeClient {
	return &client.FakeClient{
		GetFunc: func(_ string) ([]byte, error) {
			resp := map[string]any{
				"pages": items,
			}

			data, _ := json.Marshal(resp)

			return data, nil
		},
	}
}

func writeTempJSON(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	file := filepath.Join(dir, "input.json")
	_ = os.WriteFile(file, []byte(content), 0600)

	return file
}
