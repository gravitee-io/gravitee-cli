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
	"strings"
	"testing"

	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/config"
)

func newResolverService(fake *client.FakeClient) *service {
	return &service{
		client:   fake,
		resolved: &config.ResolvedContext{Org: "DEFAULT", Env: "DEFAULT"},
	}
}

func paginatedAPIsClient(items []map[string]any) *client.FakeClient {
	respond := func() []byte {
		resp := map[string]any{
			"data":       items,
			"pagination": map[string]int{"page": 1, "perPage": 100, "pageCount": 1, "totalCount": len(items), "pageItemsCount": len(items)},
		}

		data, _ := json.Marshal(resp)

		return data
	}

	return &client.FakeClient{
		GetFunc:  func(_ string) ([]byte, error) { return respond(), nil },
		PostFunc: func(_ string, _ any) ([]byte, error) { return respond(), nil },
	}
}

func TestResolveAPI(t *testing.T) {
	t.Run("passes non-path value through without API call", func(t *testing.T) {
		called := false
		fake := &client.FakeClient{
			GetFunc: func(_ string) ([]byte, error) {
				called = true

				return nil, nil
			},
		}
		s := newResolverService(fake)

		id, err := s.ResolveAPI("9e1884ee-ff53-4ecb-9884-eeff53fecbfe")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if id != "9e1884ee-ff53-4ecb-9884-eeff53fecbfe" {
			t.Fatalf("expected passthrough, got %q", id)
		}

		if called {
			t.Fatal("expected no API call for non-path input")
		}
	})

	t.Run("resolves V2 contextPath to id", func(t *testing.T) {
		fake := paginatedAPIsClient([]map[string]any{
			{"id": "api-123", "contextPath": "/my/api", "name": "My API"},
			{"id": "api-456", "contextPath": "/other/api", "name": "Other"},
		})
		s := newResolverService(fake)

		id, err := s.ResolveAPI("/my/api")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if id != "api-123" {
			t.Fatalf("expected api-123, got %q", id)
		}
	})

	t.Run("resolves V4 listener path to id", func(t *testing.T) {
		fake := paginatedAPIsClient([]map[string]any{
			{
				"id":   "api-789",
				"name": "V4 API",
				"listeners": []any{
					map[string]any{
						"paths": []any{
							map[string]any{"path": "/v4/api"},
						},
					},
				},
			},
		})
		s := newResolverService(fake)

		id, err := s.ResolveAPI("/v4/api")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if id != "api-789" {
			t.Fatalf("expected api-789, got %q", id)
		}
	})

	t.Run("errors when path not found", func(t *testing.T) {
		fake := paginatedAPIsClient([]map[string]any{
			{"id": "api-1", "contextPath": "/other"},
		})
		s := newResolverService(fake)

		_, err := s.ResolveAPI("/missing")
		if err == nil {
			t.Fatal("expected error for missing path")
		}

		if !strings.Contains(err.Error(), "no API found") {
			t.Fatalf("expected 'no API found' in error, got: %v", err)
		}
	})

	t.Run("errors when multiple APIs share the path", func(t *testing.T) {
		fake := paginatedAPIsClient([]map[string]any{
			{"id": "api-a", "contextPath": "/dup"},
			{"id": "api-b", "contextPath": "/dup"},
		})
		s := newResolverService(fake)

		_, err := s.ResolveAPI("/dup")
		if err == nil {
			t.Fatal("expected error for duplicate path")
		}

		if !strings.Contains(err.Error(), "multiple APIs") {
			t.Fatalf("expected 'multiple APIs' in error, got: %v", err)
		}
	})

	t.Run("ignores partial path matches returned by server search", func(t *testing.T) {
		// Server search on "/foo" may also return "/foobar" - resolver must match exactly.
		fake := paginatedAPIsClient([]map[string]any{
			{"id": "api-prefix", "contextPath": "/foobar"},
			{"id": "api-exact", "contextPath": "/foo"},
		})
		s := newResolverService(fake)

		id, err := s.ResolveAPI("/foo")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if id != "api-exact" {
			t.Fatalf("expected api-exact, got %q", id)
		}
	})
}
