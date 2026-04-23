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
	"fmt"
)

// PaginatedResponse holds a page of results from the APIM API.
type PaginatedResponse struct {
	Data       []json.RawMessage `json:"data"`
	Pagination Pagination        `json:"pagination"`
}

// Pagination holds the pagination metadata.
type Pagination struct {
	Page           int `json:"page"`
	PerPage        int `json:"perPage"`
	PageCount      int `json:"pageCount"`
	TotalCount     int `json:"totalCount"`
	PageItemsCount int `json:"pageItemsCount"`
}

const (
	maxPages = 1000
	maxItems = 10_000
)

// FetchAllPages fetches all pages using the given single-page fetcher.
func FetchAllPages(fetch func(page int) (*PaginatedResponse, error)) ([]json.RawMessage, error) {
	var all []json.RawMessage

	for page := 1; ; page++ {
		resp, err := fetch(page)
		if err != nil {
			return nil, err
		}

		all = append(all, resp.Data...)

		if len(all) >= maxItems || resp.Pagination.PageCount <= 0 || page >= resp.Pagination.PageCount || page > maxPages {
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
