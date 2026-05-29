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

package log

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/testutil"
)

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

func TestListLogs(t *testing.T) {
	t.Run("returns API request logs", func(t *testing.T) {
		fake := pagedLogs(map[int][]map[string]any{
			1: {{"requestId": "req-1", "method": "GET", "status": "200"}},
		})
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newListCmd(tc.Factory), "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "req-1")
	})

	t.Run("--all aggregates all pages in table output", func(t *testing.T) {
		fake := pagedLogs(map[int][]map[string]any{
			1: {{"requestId": "req-1"}},
			2: {{"requestId": "req-2"}},
		})
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newListCmd(tc.Factory), "--api", "api-1", "--all")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "req-1")
		testutil.AssertOutputContains(t, tc.Out, "req-2")
	})

	t.Run("--all aggregates all pages in json output", func(t *testing.T) {
		fake := pagedLogs(map[int][]map[string]any{
			1: {{"requestId": "req-1"}},
			2: {{"requestId": "req-2"}},
		})
		tc := testutil.NewFactory(fake)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newListCmd(tc.Factory), "--api", "api-1", "--all")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "req-1")
		testutil.AssertOutputContains(t, tc.Out, "req-2")
	})
}

func TestGetLog(t *testing.T) {
	t.Run("calls analytics endpoint with correct path", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/analytics/req-1")

				return json.Marshal(map[string]any{"requestId": "req-1", "method": "GET", "status": 200})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newGetCmd(tc.Factory), "req-1", "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "req-1")
	})

	t.Run("returns not found for unknown request", func(t *testing.T) {
		fake := testutil.APIFailingWith(404, "resource not found")
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newGetCmd(tc.Factory), "req-999", "--api", "api-1")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("returns not found for unknown request", func(t *testing.T) {
		fake := testutil.APIFailingWith(404, "resource not found")
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newGetCmd(tc.Factory), "req-999", "--api", "api-1")

		testutil.AssertErrorContains(t, err, "not found")
	})
}
