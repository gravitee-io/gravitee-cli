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

package am

import (
	"encoding/json"
	"fmt"
)

// PaginatedResponse holds a page of results from the AM API.
type PaginatedResponse struct {
	Data        []json.RawMessage `json:"data"`
	TotalCount  int               `json:"totalCount"`
	CurrentPage int               `json:"currentPage"`
}

const maxItems = 10_000

// FetchAllPages fetches all pages using the given single-page fetcher (0-indexed).
func FetchAllPages(fetch func(page int) (*PaginatedResponse, error), perPage int) ([]json.RawMessage, error) {
	var all []json.RawMessage

	for page := 0; ; page++ {
		resp, err := fetch(page)
		if err != nil {
			return nil, err
		}

		all = append(all, resp.Data...)

		if len(all) >= maxItems || resp.TotalCount <= 0 || len(all) >= resp.TotalCount || page > 1000 {
			break
		}

		if len(resp.Data) < perPage {
			break
		}
	}

	return all, nil
}

func parsePaginatedResponse(data []byte) (*PaginatedResponse, error) {
	var resp PaginatedResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse paginated response: %w", err)
	}

	return &resp, nil
}
