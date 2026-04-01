package am

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestFetchAllPages_SmallResultSet(t *testing.T) {
	items := []json.RawMessage{
		json.RawMessage(`{"id":"1"}`),
		json.RawMessage(`{"id":"2"}`),
	}

	result, err := FetchAllPages(func(page int) (*PaginatedResponse, error) {
		return &PaginatedResponse{
			Data:       items,
			TotalCount: 2,
		}, nil
	}, 10)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 items, got %d", len(result))
	}
}

func TestFetchAllPages_CapsAtMaxItems(t *testing.T) {
	callCount := 0
	perPage := 100

	result, err := FetchAllPages(func(page int) (*PaginatedResponse, error) {
		callCount++

		items := make([]json.RawMessage, perPage)
		for i := range items {
			items[i] = json.RawMessage(fmt.Sprintf(`{"id":"%d-%d"}`, page, i))
		}

		return &PaginatedResponse{
			Data:       items,
			TotalCount: 20000,
		}, nil
	}, perPage)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) > maxItems+perPage {
		t.Errorf("expected at most %d items (with one page overshoot), got %d", maxItems+perPage, len(result))
	}
}

func TestFetchAllPages_ErrorPropagation(t *testing.T) {
	_, err := FetchAllPages(func(page int) (*PaginatedResponse, error) {
		return nil, fmt.Errorf("API error")
	}, 10)

	if err == nil || err.Error() != "API error" {
		t.Errorf("expected 'API error', got: %v", err)
	}
}
