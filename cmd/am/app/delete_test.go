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
	"testing"

	"gravitee.io/gctl/internal/am"
	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/testutil"
)

func TestDeleteApplication(t *testing.T) {
	t.Run("deletes an application", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			DeleteApplicationFunc: func(domainID, appID string) error {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if appID != "app-1" {
					t.Errorf("expected appID 'app-1', got %q", appID)
				}

				return nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "app-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Application 'app-1' deleted.")
	})

	t.Run("returns error on API failure", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			DeleteApplicationFunc: func(_, _ string) error {
				return &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "app-1")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("requires app ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "delete", "app-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
