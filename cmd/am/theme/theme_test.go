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

package theme

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"gravitee.io/gctl/internal/am"
	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/testutil"
)

// --- List ---

func TestListThemes(t *testing.T) {
	t.Run("returns themes", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			ListThemesFunc: func(domainID string) ([]json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return []json.RawMessage{
					json.RawMessage(`{"id":"theme-1","name":"Default","primaryColor":"#0000FF"}`),
					json.RawMessage(`{"id":"theme-2","name":"Dark","primaryColor":"#333333"}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewThemeCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Default")
		testutil.AssertOutputContains(t, tc.Out, "Dark")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			ListThemesFunc: func(_ string) ([]json.RawMessage, error) {
				return []json.RawMessage{
					json.RawMessage(`{"id":"theme-1","name":"Default"}`),
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewThemeCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewThemeCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertErrorContains(t, err, "no context configured")
	})

	t.Run("requires domain flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewThemeCmd(tc.Factory)
		err := testutil.Execute(cmd, "list")

		testutil.AssertErrorContains(t, err, "required")
	})
}

// --- Get ---

func TestGetTheme(t *testing.T) {
	t.Run("returns theme details", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetThemeFunc: func(domainID, themeID string) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if themeID != "theme-1" {
					t.Errorf("expected themeID 'theme-1', got %q", themeID)
				}

				return json.Marshal(map[string]any{
					"id": "theme-1", "name": "Default", "primaryColor": "#0000FF",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewThemeCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "theme-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Default")
		testutil.AssertOutputContains(t, tc.Out, "theme-1")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetThemeFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "theme-1", "name": "Default"})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewThemeCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "get", "theme-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires theme ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewThemeCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewThemeCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "theme-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Create ---

func TestCreateTheme(t *testing.T) {
	t.Run("creates a theme from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			CreateThemeFunc: func(domainID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return json.Marshal(map[string]any{
					"id": "new-theme", "name": "Custom", "primaryColor": "#FF0000",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"name":"Custom","primaryColor":"#FF0000"}`)

		cmd := NewThemeCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Custom")
		testutil.AssertOutputContains(t, tc.Out, "new-theme")
	})

	t.Run("requires json input", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewThemeCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create")

		testutil.AssertErrorContains(t, err, "input")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"name":"Custom"}`)

		cmd := NewThemeCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "create", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Update ---

func TestUpdateTheme(t *testing.T) {
	t.Run("updates a theme from file", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			UpdateThemeFunc: func(domainID, themeID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if themeID != "theme-1" {
					t.Errorf("expected themeID 'theme-1', got %q", themeID)
				}

				return json.Marshal(map[string]any{
					"id": "theme-1", "name": "Updated", "primaryColor": "#00FF00",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		tmpFile := writeTempJSON(t, `{"name":"Updated"}`)

		cmd := NewThemeCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "theme-1", "--file", tmpFile)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Updated")
	})

	t.Run("requires json input", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewThemeCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "theme-1")

		testutil.AssertErrorContains(t, err, "input")
	})

	t.Run("requires theme ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		tmpFile := writeTempJSON(t, `{"name":"Updated"}`)

		cmd := NewThemeCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		tmpFile := writeTempJSON(t, `{"name":"Updated"}`)

		cmd := NewThemeCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "theme-1", "--file", tmpFile)

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

// --- Delete ---

func TestDeleteTheme(t *testing.T) {
	t.Run("deletes a theme", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			DeleteThemeFunc: func(domainID, themeID string) error {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if themeID != "theme-1" {
					t.Errorf("expected themeID 'theme-1', got %q", themeID)
				}

				return nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewThemeCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "theme-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Theme 'theme-1' deleted.")
	})

	t.Run("returns error on API failure", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			DeleteThemeFunc: func(_, _ string) error {
				return &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewThemeCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "theme-1")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("requires theme ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewThemeCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewThemeCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "theme-1")

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
