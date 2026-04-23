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

//go:build e2e

package e2e

import (
	"testing"
)

// TestAPIMPageLifecycle covers documentation page CRUD + publish/unpublish.
func TestAPIMPageLifecycle(t *testing.T) {
	apiFixture := writeFixture(t, "api.json")
	pageFixture := writeFixture(t, "page.json")
	pageUpdatedFixture := writeFixture(t, "page-updated.json")

	apiOut := runCLIExpectSuccess(t, "apim", "api", "create", "-f", apiFixture, "-o", "json")

	apiID := extractID(t, apiOut)
	if apiID == "" {
		t.Fatalf("api create returned no id: %s", apiOut)
	}

	t.Cleanup(func() {
		_, _ = runCLI("apim", "api", "stop", apiID)
		_, _ = runCLI("apim", "api", "delete", apiID, "--close-plans")
	})

	pageOut := runCLIExpectSuccess(t, "apim", "page", "create",
		"--api", apiID,
		"-f", pageFixture,
		"-o", "json")

	pageID := extractID(t, pageOut)
	if pageID == "" {
		t.Fatalf("page create returned no id: %s", pageOut)
	}

	runCLIExpectSuccess(t, "apim", "page", "get", pageID, "--api", apiID, "-o", "json")
	runCLIExpectSuccess(t, "apim", "page", "list", "--api", apiID, "-o", "json")
	runCLIExpectSuccess(t, "apim", "page", "update", pageID, "--api", apiID, "-f", pageUpdatedFixture)
	runCLIExpectSuccess(t, "apim", "page", "publish", pageID, "--api", apiID)
	runCLIExpectSuccess(t, "apim", "page", "unpublish", pageID, "--api", apiID)
	runCLIExpectSuccess(t, "apim", "page", "delete", pageID, "--api", apiID)
}
