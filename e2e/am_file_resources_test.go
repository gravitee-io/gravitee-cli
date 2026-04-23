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
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestFileBasedCreateErrors(t *testing.T) {
	domainID := getDefaultDomainID(t)

	t.Run("nonexistent file returns error", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "idp", "create", "--domain", domainID, "--file", "/nonexistent/path.json")
		if !strings.Contains(out, "no such file") {
			t.Errorf("expected 'no such file' error, got: %s", out)
		}
	})

	t.Run("invalid JSON content returns error", func(t *testing.T) {
		tmp, err := os.CreateTemp("", "e2e-invalid-*.json")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmp.Name())
		tmp.Close()

		if err := os.WriteFile(tmp.Name(), []byte("not json"), 0o644); err != nil {
			t.Fatal(err)
		}

		out := runCLIExpectError(t, "am", "idp", "create", "--domain", domainID, "--file", tmp.Name())
		if !strings.Contains(out, "invalid JSON") {
			t.Errorf("expected 'invalid JSON' error, got: %s", out)
		}
	})

	t.Run("empty file returns error", func(t *testing.T) {
		tmp, err := os.CreateTemp("", "e2e-empty-*.json")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmp.Name())
		tmp.Close()

		out := runCLIExpectError(t, "am", "idp", "create", "--domain", domainID, "--file", tmp.Name())
		if !strings.Contains(out, "invalid JSON") {
			t.Errorf("expected 'invalid JSON' error, got: %s", out)
		}
	})

	t.Run("missing file flag returns error", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "idp", "create", "--domain", domainID)
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})
}

func TestFileBasedResourceListGet(t *testing.T) {
	domainID := getDefaultDomainID(t)

	resources := []struct {
		name      string
		hasData   bool // true if fresh AM has default data for this resource
		checkJSON bool
	}{
		{name: "idp", hasData: true, checkJSON: true},
		{name: "certificate", hasData: false},
		{name: "factor", hasData: false},
		{name: "flow", hasData: true, checkJSON: true},
		// Skipped: form (requires formTemplate), email (requires emailTemplate),
		// protected-resource (500), password-policy (empty JSON), theme (empty JSON).
		{name: "reporter", hasData: false},
		{name: "bot-detection", hasData: false},
		{name: "device-identifier", hasData: false},
		{name: "auth-device-notifier", hasData: false},
		{name: "authorization-engine", hasData: false},
		{name: "extension-grant", hasData: false},
		{name: "resource", hasData: false},
		{name: "dictionary", hasData: false},
	}

	for _, r := range resources {
		r := r
		t.Run("list "+r.name, func(t *testing.T) {
			runCLIExpectSuccess(t, "am", r.name, "list", "--domain", domainID)
		})

		if r.checkJSON {
			t.Run("list "+r.name+" json", func(t *testing.T) {
				out := runCLIExpectSuccess(t, "am", r.name, "list", "--domain", domainID, "-o", "json")

				var arr json.RawMessage
				if err := json.Unmarshal([]byte(out), &arr); err != nil {
					t.Fatalf("expected valid JSON array for %s, got error: %v\nOutput: %s", r.name, err, out)
				}
			})
		}
	}
}

func TestFileBasedResourceOutputFormats(t *testing.T) {
	domainID := getDefaultDomainID(t)

	t.Run("idp list yaml", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "idp", "list", "--domain", domainID, "-o", "yaml")
		if !strings.Contains(out, "-") && !strings.Contains(out, "name:") {
			t.Errorf("expected YAML output with list items or 'name:' field, got: %s", out)
		}
	})

	t.Run("idp list quiet", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "idp", "list", "--domain", domainID, "-q")
		if out != "" {
			t.Errorf("expected empty output in quiet mode, got: %s", out)
		}
	})

	t.Run("flow list yaml", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "flow", "list", "--domain", domainID, "-o", "yaml")
	})
}

func TestAlertOperations(t *testing.T) {
	domainID := getDefaultDomainID(t)

	t.Run("alert notifier list", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "alert", "notifier", "list", "--domain", domainID)
	})

	t.Run("alert trigger get", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "alert", "trigger", "get", "--domain", domainID)
	})
}
