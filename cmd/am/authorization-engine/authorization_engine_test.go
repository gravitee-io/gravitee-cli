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

package authorizationengine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

// --- List ---

func TestListAuthorizationEngines(t *testing.T) {
	t.Run("returns authorization engines", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			ListAuthorizationEnginesFunc: func(domainID string) ([]json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return []json.RawMessage{
					json.RawMessage(`{"id":"ae-1","name":"My Engine","type":"uma"}`),
					json.RawMessage(`{"id":"ae-2","name":"Other","type":"opa"}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAuthorizationEngineCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Engine")
		testutil.AssertOutputContains(t, tc.Out, "Other")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			ListAuthorizationEnginesFunc: func(_ string) ([]json.RawMessage, error) {
				return []json.RawMessage{
					json.RawMessage(`{"id":"ae-1","name":"Test"}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAuthorizationEngineCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewAuthorizationEngineCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertErrorContains(t, err, "no context configured")
	})

	t.Run("requires domain flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAuthorizationEngineCmd(tc.Factory)
		err := testutil.Execute(cmd, "list")

		testutil.AssertErrorContains(t, err, "required")
	})
}

// --- Get ---

func TestGetAuthorizationEngine(t *testing.T) {
	t.Run("returns authorization engine details", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetAuthorizationEngineFunc: func(domainID, engineID string) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if engineID != "ae-1" {
					t.Errorf("expected engineID 'ae-1', got %q", engineID)
				}

				return json.Marshal(map[string]any{
					"id": "ae-1", "name": "My Engine", "type": "uma",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAuthorizationEngineCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "ae-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Engine")
		testutil.AssertOutputContains(t, tc.Out, "ae-1")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetAuthorizationEngineFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "ae-1", "name": "Test"})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAuthorizationEngineCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "get", "ae-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires engine ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAuthorizationEngineCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewAuthorizationEngineCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "ae-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Update ---

func TestUpdateAuthorizationEngine(t *testing.T) {
	t.Run("updates an authorization engine from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			UpdateAuthorizationEngineFunc: func(domainID, engineID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if engineID != "ae-1" {
					t.Errorf("expected engineID 'ae-1', got %q", engineID)
				}

				return json.Marshal(map[string]any{
					"id": "ae-1", "name": "Updated", "type": "uma",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"name":"Updated"}`)

		cmd := NewAuthorizationEngineCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "ae-1", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Updated")
	})

	t.Run("requires file flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAuthorizationEngineCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "ae-1")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires engine ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewAuthorizationEngineCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"name":"Test"}`)

		cmd := NewAuthorizationEngineCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "ae-1", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Helper ---

func writeTempJSON(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.json")

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	return path
}
