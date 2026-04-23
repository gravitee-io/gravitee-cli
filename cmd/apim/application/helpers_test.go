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
