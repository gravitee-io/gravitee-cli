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
