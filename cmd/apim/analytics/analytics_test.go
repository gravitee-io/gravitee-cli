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

package analytics

import (
	"encoding/json"
	"testing"

	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/testutil"
)

func TestGetAnalytics(t *testing.T) {
	t.Run("returns analytics data", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/analytics")

				return json.Marshal(map[string]any{"type": "COUNT", "count": 4523})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(NewAnalyticsCmd(tc.Factory),
			"--api", "api-1", "--type", "COUNT", "--from", "1700000000000", "--to", "1700000001000")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "4523")
	})

	t.Run("returns not found for unknown API", func(t *testing.T) {
		fake := testutil.APIFailingWith(404, "resource not found")
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(NewAnalyticsCmd(tc.Factory),
			"--api", "api-999", "--from", "1700000000000", "--to", "1700000001000")

		testutil.AssertErrorContains(t, err, "not found")
	})
}
