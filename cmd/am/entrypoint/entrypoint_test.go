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

package entrypoint

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

// --- Get ---

func TestGetEntrypoints(t *testing.T) {
	t.Run("returns entrypoint data", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetEntrypointsFunc: func(domainID string) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return json.Marshal(map[string]any{
					"url": "https://example.com",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewEntrypointCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "example.com")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetEntrypointsFunc: func(_ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"url": "https://example.com"})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewEntrypointCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "get")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"url"`)
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewEntrypointCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get")

		testutil.AssertErrorContains(t, err, "no context configured")
	})

	t.Run("requires domain flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewEntrypointCmd(tc.Factory)
		err := testutil.Execute(cmd, "get")

		testutil.AssertErrorContains(t, err, "required")
	})
}
