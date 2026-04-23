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
			Pagination: Pagination{Page: 1, PerPage: 10, PageCount: 1, TotalCount: 2, PageItemsCount: 2},
		}, nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 items, got %d", len(result))
	}
}

func TestFetchAllPages_CapsAtMaxItems(t *testing.T) {
	callCount := 0
	itemsPerPage := 100

	result, err := FetchAllPages(func(page int) (*PaginatedResponse, error) {
		callCount++

		items := make([]json.RawMessage, itemsPerPage)
		for i := range items {
			items[i] = json.RawMessage(fmt.Sprintf(`{"id":"%d-%d"}`, page, i))
		}

		return &PaginatedResponse{
			Data:       items,
			Pagination: Pagination{Page: page, PerPage: itemsPerPage, PageCount: 200, TotalCount: 20000, PageItemsCount: itemsPerPage},
		}, nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) > maxItems+itemsPerPage {
		t.Errorf("expected at most %d items (with one page overshoot), got %d", maxItems+itemsPerPage, len(result))
	}

	if callCount > 101 {
		t.Errorf("expected at most 101 pages fetched, got %d", callCount)
	}
}

func TestFetchAllPages_ErrorPropagation(t *testing.T) {
	_, err := FetchAllPages(func(page int) (*PaginatedResponse, error) {
		return nil, fmt.Errorf("API error")
	})

	if err == nil || err.Error() != "API error" {
		t.Errorf("expected 'API error', got: %v", err)
	}
}
