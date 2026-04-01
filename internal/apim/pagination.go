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
