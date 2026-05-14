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
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestEntrypointsGet(t *testing.T) {
	t.Run("prints vhost mode summary", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "domains/dom-1")
				return json.Marshal(map[string]any{
					"id":        "dom-1",
					"vhostMode": true,
					"vhosts": []map[string]any{
						{"host": "auth.example.com", "path": "/", "overrideEntrypoint": true},
					},
				})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newEntrypointsGetCmd(tc.Factory), "dom-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "vhost")
		testutil.AssertOutputContains(t, tc.Out, "auth.example.com")
		testutil.AssertOutputContains(t, tc.Out, "(override)")
	})

	t.Run("prints context-path mode summary", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(_ string) ([]byte, error) {
				return json.Marshal(map[string]any{
					"id":        "dom-1",
					"vhostMode": false,
					"path":      "/auth",
				})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newEntrypointsGetCmd(tc.Factory), "dom-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "context-path")
		testutil.AssertOutputContains(t, tc.Out, "/auth")
	})
}

func TestEntrypointsSetPath(t *testing.T) {
	var captured map[string]any
	fake := &client.FakeClient{
		PatchFunc: func(path string, body any) ([]byte, error) {
			testutil.AssertPathCalled(t, path, "domains/dom-1")
			raw, _ := json.Marshal(body)
			_ = json.Unmarshal(raw, &captured)
			return json.Marshal(map[string]any{
				"id": "dom-1", "vhostMode": false, "path": "/auth",
			})
		},
	}
	tc := testutil.NewFactory(fake)

	err := testutil.Execute(newEntrypointsSetPathCmd(tc.Factory), "dom-1", "/auth")

	testutil.AssertNoError(t, err)
	if captured["vhostMode"] != false {
		t.Errorf("expected vhostMode=false, got %v", captured["vhostMode"])
	}
	if captured["path"] != "/auth" {
		t.Errorf("expected path=/auth, got %v", captured["path"])
	}
}

func TestEntrypointsAddVhost(t *testing.T) {
	t.Run("appends to existing vhosts", func(t *testing.T) {
		var captured map[string]any
		fake := &client.FakeClient{
			GetFunc: func(_ string) ([]byte, error) {
				return json.Marshal(map[string]any{
					"id":        "dom-1",
					"vhostMode": true,
					"vhosts": []map[string]any{
						{"host": "first.example.com", "path": "/", "overrideEntrypoint": true},
					},
				})
			},
			PatchFunc: func(_ string, body any) ([]byte, error) {
				raw, _ := json.Marshal(body)
				_ = json.Unmarshal(raw, &captured)
				return raw, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newEntrypointsAddVhostCmd(tc.Factory),
			"dom-1", "second.example.com", "--path", "/auth")

		testutil.AssertNoError(t, err)
		if captured["vhostMode"] != true {
			t.Errorf("expected vhostMode=true, got %v", captured["vhostMode"])
		}

		hosts, ok := captured["vhosts"].([]any)
		if !ok || len(hosts) != 2 {
			t.Fatalf("expected 2 vhosts, got %v", captured["vhosts"])
		}
	})

	t.Run("rejects duplicate host/path", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(_ string) ([]byte, error) {
				return json.Marshal(map[string]any{
					"vhostMode": true,
					"vhosts": []map[string]any{
						{"host": "dup.example.com", "path": "/"},
					},
				})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newEntrypointsAddVhostCmd(tc.Factory),
			"dom-1", "dup.example.com", "--path", "/")

		testutil.AssertErrorContains(t, err, "already exists")
	})
}

func TestEntrypointsRemoveVhost(t *testing.T) {
	t.Run("drops matching vhost by host", func(t *testing.T) {
		var captured map[string]any
		fake := &client.FakeClient{
			GetFunc: func(_ string) ([]byte, error) {
				return json.Marshal(map[string]any{
					"vhostMode": true,
					"vhosts": []map[string]any{
						{"host": "a.example.com", "path": "/"},
						{"host": "b.example.com", "path": "/"},
					},
				})
			},
			PatchFunc: func(_ string, body any) ([]byte, error) {
				raw, _ := json.Marshal(body)
				_ = json.Unmarshal(raw, &captured)
				return raw, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newEntrypointsRemoveVhostCmd(tc.Factory),
			"dom-1", "a.example.com")

		testutil.AssertNoError(t, err)
		hosts := captured["vhosts"].([]any)
		if len(hosts) != 1 {
			t.Fatalf("expected 1 remaining vhost, got %d", len(hosts))
		}
	})

	t.Run("errors when no vhost matches", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(_ string) ([]byte, error) {
				return json.Marshal(map[string]any{
					"vhostMode": true,
					"vhosts":    []map[string]any{{"host": "a.example.com", "path": "/"}},
				})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newEntrypointsRemoveVhostCmd(tc.Factory),
			"dom-1", "missing.example.com")

		testutil.AssertErrorContains(t, err, "no vhost matched")
	})
}

func TestEntrypointsClearVhosts(t *testing.T) {
	var captured map[string]any
	fake := &client.FakeClient{
		PatchFunc: func(_ string, body any) ([]byte, error) {
			raw, _ := json.Marshal(body)
			_ = json.Unmarshal(raw, &captured)
			return raw, nil
		},
	}
	tc := testutil.NewFactory(fake)

	err := testutil.Execute(newEntrypointsClearVhostsCmd(tc.Factory), "dom-1")

	testutil.AssertNoError(t, err)
	if captured["vhostMode"] != false {
		t.Errorf("expected vhostMode=false, got %v", captured["vhostMode"])
	}
	if hosts, _ := captured["vhosts"].([]any); len(hosts) != 0 {
		t.Errorf("expected empty vhosts, got %v", hosts)
	}
}
