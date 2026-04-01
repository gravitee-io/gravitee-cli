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
