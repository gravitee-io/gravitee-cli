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

package health

import (
	"encoding/json"
	"testing"

	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/testutil"
)

func TestGetHealth(t *testing.T) {
	t.Run("returns health availability", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/health/availability")

				return json.Marshal(map[string]any{
					"availability": map[string]float64{"https://backend.example.com:443": 99.8},
				})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(NewHealthCmd(tc.Factory), "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "99.8")
	})
}
