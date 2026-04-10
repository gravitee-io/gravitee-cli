package api

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func paginatedAPIs(items ...map[string]any) *client.FakeClient {
	respond := func() []byte {
		resp := map[string]any{
			"data":       items,
			"pagination": map[string]int{"page": 1, "perPage": 10, "pageCount": 1, "totalCount": len(items), "pageItemsCount": len(items)},
		}

		data, _ := json.Marshal(resp)

		return data
	}

	return &client.FakeClient{
		GetFunc:  func(_ string) ([]byte, error) { return respond(), nil },
		PostFunc: func(_ string, _ any) ([]byte, error) { return respond(), nil },
	}
}

// pagedLogs returns a FakeClient that serves different items per page,
// keyed by the "page" query parameter on the request URL.
func pagedLogs(pages map[int][]map[string]any) *client.FakeClient {
	pageCount := 0
	total := 0

	for p, items := range pages {
		if p > pageCount {
			pageCount = p
		}

		total += len(items)
	}

	return &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			page := 1

			if q := strings.Index(path, "?"); q >= 0 {
				if parsed, err := url.ParseQuery(path[q+1:]); err == nil {
					if n, err := strconv.Atoi(parsed.Get("page")); err == nil {
						page = n
					}
				}
			}

			items := pages[page]
			resp := map[string]any{
				"data": items,
				"pagination": map[string]int{
					"page":           page,
					"perPage":        10,
					"pageCount":      pageCount,
					"totalCount":     total,
					"pageItemsCount": len(items),
				},
			}

			return json.Marshal(resp)
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
