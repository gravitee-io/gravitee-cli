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
	"strings"
	"testing"

	"gravitee.io/gctl/internal/client"
)

func TestUserLock(t *testing.T) {
	fake := &client.FakeClient{
		PutFunc: func(path string, body interface{}) ([]byte, error) {
			if !strings.Contains(path, "/users/user-1/status") {
				t.Errorf("unexpected path: %s", path)
			}
			return nil, nil
		},
	}

	f, out := newTestFactory(fake, false)
	domainID := "test-domain"

	cmd := newLockCmd(f, &domainID)
	cmd.SetArgs([]string{"user-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "locked") {
		t.Errorf("expected 'locked' in output, got: %s", out.String())
	}
}

func TestUserUnlock(t *testing.T) {
	fake := &client.FakeClient{
		PutFunc: func(path string, body interface{}) ([]byte, error) {
			if !strings.Contains(path, "/users/user-1/status") {
				t.Errorf("unexpected path: %s", path)
			}
			return nil, nil
		},
	}

	f, out := newTestFactory(fake, false)
	domainID := "test-domain"

	cmd := newUnlockCmd(f, &domainID)
	cmd.SetArgs([]string{"user-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "unlocked") {
		t.Errorf("expected 'unlocked' in output, got: %s", out.String())
	}
}

func TestUserResetPassword(t *testing.T) {
	fake := &client.FakeClient{
		PostFunc: func(path string, body interface{}) ([]byte, error) {
			if !strings.Contains(path, "/users/user-1/resetPassword") {
				t.Errorf("unexpected path: %s", path)
			}

			// Verify the body contains the password
			var m map[string]interface{}
			switch b := body.(type) {
			case []byte:
				_ = json.Unmarshal(b, &m)
			case json.RawMessage:
				_ = json.Unmarshal(b, &m)
			}

			if pwd, ok := m["password"].(string); !ok || pwd != "newSecret123" {
				t.Errorf("expected password 'newSecret123', got: %v", m["password"])
			}

			return nil, nil
		},
	}

	f, out := newTestFactory(fake, false)
	domainID := "test-domain"

	cmd := newResetPasswordCmd(f, &domainID)
	cmd.SetArgs([]string{"user-1", "--password", "newSecret123"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Password reset for user 'user-1'") {
		t.Errorf("expected reset message in output, got: %s", out.String())
	}
}

func TestUserDelete(t *testing.T) {
	fake := &client.FakeClient{
		DeleteFunc: func(path string) error {
			if !strings.Contains(path, "/users/user-1") {
				t.Errorf("unexpected path: %s", path)
			}
			return nil
		},
	}

	f, out := newTestFactory(fake, false)
	domainID := "test-domain"

	cmd := newDeleteCmd(f, &domainID)
	cmd.SetArgs([]string{"user-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "User 'user-1' deleted.") {
		t.Errorf("unexpected output: %s", out.String())
	}
}
