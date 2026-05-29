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

package apim

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/config"
)

func newAppService(fake *client.FakeClient) *service {
	return &service{
		client:   fake,
		resolved: &config.ResolvedContext{Org: "DEFAULT", Env: "DEFAULT"},
	}
}

// v1PagedAppsClient serves the V1 applications envelope, paged by the "page" query param.
func v1PagedAppsClient(pages map[int][]map[string]any) *client.FakeClient {
	totalPages := 0
	totalElements := 0

	for p, items := range pages {
		if p > totalPages {
			totalPages = p
		}

		totalElements += len(items)
	}

	return &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			page := 1
			size := 10

			if q := strings.Index(path, "?"); q >= 0 {
				if parsed, err := url.ParseQuery(path[q+1:]); err == nil {
					if n, err := strconv.Atoi(parsed.Get("page")); err == nil {
						page = n
					}

					if n, err := strconv.Atoi(parsed.Get("size")); err == nil {
						size = n
					}
				}
			}

			items := pages[page]
			resp := map[string]any{
				"data": items,
				"page": map[string]int{
					"current":        page,
					"size":           size,
					"per_page":       size,
					"total_pages":    totalPages,
					"total_elements": totalElements,
				},
			}

			return json.Marshal(resp)
		},
	}
}

func TestListApplications_TranslatesV1Envelope(t *testing.T) {
	fake := v1PagedAppsClient(map[int][]map[string]any{
		1: {{"id": "app-1"}, {"id": "app-2"}},
	})
	s := newAppService(fake)

	resp, err := s.ListApplications(ListApplicationsParams{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.Data))
	}

	if resp.Pagination.Page != 1 {
		t.Errorf("Page: expected 1, got %d", resp.Pagination.Page)
	}

	if resp.Pagination.PageCount != 1 {
		t.Errorf("PageCount: expected 1, got %d", resp.Pagination.PageCount)
	}

	if resp.Pagination.TotalCount != 2 {
		t.Errorf("TotalCount: expected 2, got %d", resp.Pagination.TotalCount)
	}

	if resp.Pagination.PageItemsCount != 2 {
		t.Errorf("PageItemsCount: expected 2, got %d", resp.Pagination.PageItemsCount)
	}
}

func TestListApplications_FetchAllPages_AggregatesAcrossPages(t *testing.T) {
	fake := v1PagedAppsClient(map[int][]map[string]any{
		1: {{"id": "app-1"}, {"id": "app-2"}},
		2: {{"id": "app-3"}},
	})
	s := newAppService(fake)

	all, err := FetchAllPages(func(page int) (*PaginatedResponse, error) {
		return s.ListApplications(ListApplicationsParams{Page: page, PerPage: 2})
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(all) != 3 {
		t.Fatalf("expected 3 aggregated items, got %d", len(all))
	}
}
