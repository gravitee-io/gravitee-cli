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

package app

import (
	"encoding/json"
	"testing"

	"gravitee.io/gctl/internal/am"
	"gravitee.io/gctl/internal/testutil"
)

func TestListApplications(t *testing.T) {
	t.Run("returns applications", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			ListApplicationsFunc: func(domainID string, p am.ListApplicationsParams) (*am.PaginatedResponse, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return &am.PaginatedResponse{
					Data: []json.RawMessage{
						json.RawMessage(`{"id":"app-1","name":"My App","type":"web","enabled":true}`),
						json.RawMessage(`{"id":"app-2","name":"Other","type":"native","enabled":false}`),
					},
					TotalCount:  2,
					CurrentPage: 0,
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My App")
		testutil.AssertOutputContains(t, tc.Out, "Other")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			ListApplicationsFunc: func(_ string, _ am.ListApplicationsParams) (*am.PaginatedResponse, error) {
				return &am.PaginatedResponse{
					Data:        []json.RawMessage{json.RawMessage(`{"id":"app-1","name":"Test"}`)},
					TotalCount:  1,
					CurrentPage: 0,
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"totalCount"`)
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertErrorContains(t, err, "no context configured")
	})

	t.Run("requires domain flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "list")

		testutil.AssertErrorContains(t, err, "required")
	})
}
