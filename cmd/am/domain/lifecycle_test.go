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

package domain

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func TestDomainEnable(t *testing.T) {
	fake := &client.FakeClient{
		PatchFunc: func(path string, body interface{}) ([]byte, error) {
			if !strings.Contains(path, "/domains/domain-1") {
				t.Errorf("unexpected path: %s", path)
			}

			var m map[string]interface{}
			switch b := body.(type) {
			case []byte:
				_ = json.Unmarshal(b, &m)
			case json.RawMessage:
				_ = json.Unmarshal(b, &m)
			}

			if enabled, ok := m["enabled"].(bool); !ok || !enabled {
				t.Errorf("expected enabled=true, got: %v", m["enabled"])
			}

			return nil, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newEnableCmd(f)
	cmd.SetArgs([]string{"domain-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "enabled") {
		t.Errorf("expected 'enabled' in output, got: %s", out.String())
	}
}

func TestDomainDisable(t *testing.T) {
	fake := &client.FakeClient{
		PatchFunc: func(path string, body interface{}) ([]byte, error) {
			if !strings.Contains(path, "/domains/domain-1") {
				t.Errorf("unexpected path: %s", path)
			}

			var m map[string]interface{}
			switch b := body.(type) {
			case []byte:
				_ = json.Unmarshal(b, &m)
			case json.RawMessage:
				_ = json.Unmarshal(b, &m)
			}

			if enabled, ok := m["enabled"].(bool); !ok || enabled {
				t.Errorf("expected enabled=false, got: %v", m["enabled"])
			}

			return nil, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newDisableCmd(f)
	cmd.SetArgs([]string{"domain-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "disabled") {
		t.Errorf("expected 'disabled' in output, got: %s", out.String())
	}
}

func TestDomainDelete(t *testing.T) {
	fake := &client.FakeClient{
		DeleteFunc: func(path string) error {
			if !strings.Contains(path, "/domains/domain-1") {
				t.Errorf("unexpected path: %s", path)
			}
			return nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newDeleteCmd(f)
	cmd.SetArgs([]string{"domain-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Domain 'domain-1' deleted.") {
		t.Errorf("unexpected output: %s", out.String())
	}
}
