package application

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func appJSON() map[string]any {
	return map[string]any{
		"id":           "app-1",
		"name":         "My Mobile App",
		"description":  "Mobile client for the Weather API",
		"type":         "SIMPLE",
		"status":       "ACTIVE",
		"owner":        map[string]any{"id": "user-1234", "displayName": "john.doe"},
		"api_key_mode": "UNSPECIFIED",
		"domain":       "https://my-app.com",
		"created_at":   "2026-03-15T10:00:00Z",
		"updated_at":   "2026-03-25T14:30:00Z",
	}
}

func paginatedApps(items ...map[string]any) *client.FakeClient {
	return &client.FakeClient{
		GetFunc: func(_ string) ([]byte, error) {
			resp := map[string]any{
				"data":       items,
				"pagination": map[string]int{"page": 1, "perPage": 10, "pageCount": 1, "totalCount": len(items), "pageItemsCount": len(items)},
			}

			data, _ := json.Marshal(resp)

			return data, nil
		},
	}
}

func emptyPaginatedResponse() []byte {
	resp := map[string]any{
		"data":       []map[string]any{},
		"pagination": map[string]int{"page": 1, "perPage": 10, "pageCount": 1, "totalCount": 0, "pageItemsCount": 0},
	}

	data, _ := json.Marshal(resp)

	return data
}

func writeTempJSON(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	file := filepath.Join(dir, "input.json")
	_ = os.WriteFile(file, []byte(content), 0600)

	return file
}
