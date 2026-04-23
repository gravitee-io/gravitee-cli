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

package user

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestResetPassword(t *testing.T) {
	t.Run("resets user password", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			ResetPasswordFunc: func(d, userID string, body json.RawMessage) error {
				if d != domainID {
					t.Errorf("expected domain %q, got %q", domainID, d)
				}
				if userID != "user-1" {
					t.Errorf("expected userID %q, got %q", "user-1", userID)
				}

				return nil
			},
		}
		tc := testutil.NewFactory(nil)
		newTestUser(tc, mock, domainID)

		err := testutil.Execute(newResetPasswordCmd(tc.Factory, &domainID), "user-1", "--password", "newSecret123")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Password reset for user 'user-1'.")
	})

	t.Run("returns error on API failure", func(t *testing.T) {
		domainID := "dom-1"
		mock := &am.MockService{
			ResetPasswordFunc: func(_, _ string, _ json.RawMessage) error {
				return &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
			},
		}
		tc := testutil.NewFactory(nil)
		newTestUser(tc, mock, domainID)

		err := testutil.Execute(newResetPasswordCmd(tc.Factory, &domainID), "user-1", "--password", "newSecret123")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("requires password flag", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil)

		err := testutil.Execute(newResetPasswordCmd(tc.Factory, &domainID), "user-1")

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires user ID argument", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil)

		err := testutil.Execute(newResetPasswordCmd(tc.Factory, &domainID))

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		domainID := "dom-1"
		tc := testutil.NewFactory(nil)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newResetPasswordCmd(tc.Factory, &domainID), "user-1", "--password", "newSecret123")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
