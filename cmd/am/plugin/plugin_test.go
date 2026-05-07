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

package plugin

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func TestPluginList(t *testing.T) {
	plugins := []map[string]interface{}{
		{"id": "github-am-idp", "name": "GitHub Identity Provider", "version": "2.4.0"},
	}
	data, _ := json.Marshal(plugins)
	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/platform/plugins/identities") {
				t.Errorf("unexpected path: %s", path)
			}
			return data, nil
		},
	}
	f, out := newTestFactory(fake, false)
	cmd := NewPluginCmd(f)
	cmd.SetArgs([]string{"list", "idp"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "GitHub Identity Provider") {
		t.Errorf("expected plugin name, got: %s", out.String())
	}
}

func TestPluginSchema(t *testing.T) {
	schema := map[string]interface{}{
		"properties": map[string]interface{}{
			"clientId":     map[string]interface{}{"type": "string", "title": "Client ID"},
			"clientSecret": map[string]interface{}{"type": "string", "title": "Client Secret"},
		},
	}
	data, _ := json.Marshal(schema)
	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/platform/plugins/identities/github-am-idp/schema") {
				t.Errorf("unexpected path: %s", path)
			}
			return data, nil
		},
	}
	f, out := newTestFactory(fake, false)
	cmd := NewPluginCmd(f)
	cmd.SetArgs([]string{"schema", "idp", "github-am-idp"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "clientId") {
		t.Errorf("expected field in output, got: %s", out.String())
	}
}

func TestPluginCreateWithFile(t *testing.T) {
	var posted map[string]interface{}
	fake := &client.FakeClient{
		PostFunc: func(path string, body interface{}) ([]byte, error) {
			if b, ok := body.([]byte); ok {
				_ = json.Unmarshal(b, &posted)
			}
			return []byte(`{"id":"idp-new","name":"My GitHub"}`), nil
		},
	}
	f, out := newTestFactory(fake, false)

	tmp, err := os.CreateTemp("", "*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	_, err = tmp.WriteString(`{"clientId":"abc","clientSecret":"xyz"}`)
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	cmd := NewPluginCmd(f)
	cmd.SetArgs([]string{"create", "idp", "github-am-idp", "--name", "My GitHub", "--config-file", tmp.Name()})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "My GitHub") {
		t.Errorf("expected name in output, got: %s", out.String())
	}
}
