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

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

// --- Get ---

func TestGetAnalytics(t *testing.T) {
	t.Run("returns analytics data", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetAnalyticsFunc: func(domainID string, params am.AnalyticsParams) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if params.Type != "count" {
					t.Errorf("expected type 'count', got %q", params.Type)
				}

				return json.Marshal(map[string]any{"value": 42})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAnalyticsCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "--type", "count")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "42")
	})

	t.Run("passes all parameters", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetAnalyticsFunc: func(_ string, params am.AnalyticsParams) (json.RawMessage, error) {
				if params.Field != "status" {
					t.Errorf("expected field 'status', got %q", params.Field)
				}

				if params.From != "2024-01-01" {
					t.Errorf("expected from '2024-01-01', got %q", params.From)
				}

				if params.To != "2024-12-31" {
					t.Errorf("expected to '2024-12-31', got %q", params.To)
				}

				if params.Interval != "86400000" {
					t.Errorf("expected interval '86400000', got %q", params.Interval)
				}

				if params.Size != 10 {
					t.Errorf("expected size 10, got %d", params.Size)
				}

				return json.Marshal(map[string]any{"value": 1})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAnalyticsCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get",
			"--type", "count", "--field", "status",
			"--from", "2024-01-01", "--to", "2024-12-31",
			"--interval", "86400000", "--size", "10")

		testutil.AssertNoError(t, err)
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewAnalyticsCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get")

		testutil.AssertErrorContains(t, err, "no context configured")
	})

	t.Run("requires domain flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAnalyticsCmd(tc.Factory)
		err := testutil.Execute(cmd, "get")

		testutil.AssertErrorContains(t, err, "required")
	})
}
