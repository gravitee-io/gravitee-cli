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

package metadata

import (
	"testing"

	"gravitee.io/gctl/internal/testutil"
)

func TestListMetadata(t *testing.T) {
	t.Run("returns metadata from the API", func(t *testing.T) {
		fake := paginatedMetadata(map[string]any{
			"key": "team-email", "name": "Team Email",
			"value": "platform-team@company.com", "format": "MAIL",
			"updatedAt": "2026-03-25T14:30:00Z",
		})
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newListCmd(tc.Factory), "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Team Email")
		testutil.AssertOutputContains(t, tc.Out, "MAIL")
	})
}
