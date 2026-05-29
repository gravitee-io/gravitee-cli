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
	"bytes"

	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/config"
	"gravitee.io/gctl/internal/factory"
)

// newTestFactory creates a factory for tests that use FakeClient directly.
func newTestFactory(fc *client.FakeClient, _ bool) (*factory.Factory, *bytes.Buffer) {
	out := &bytes.Buffer{}
	return &factory.Factory{
		Config: &config.Config{
			Current: "am-test",
			Contexts: map[string]*config.Context{
				"am-test": {
					Org: "DEFAULT", Env: "DEFAULT",
					Type: "am", Domain: "test-domain",
					AM: &config.ProductConfig{URL: "https://am-test.com", Token: "tok"},
				},
			},
		},
		Resolved: &config.ResolvedContext{
			Name: "am-test", URL: "https://am-test.com", Token: "tok",
			Org: "DEFAULT", Env: "DEFAULT",
			Type: "am", Domain: "test-domain",
		},
		Client:       fc,
		IOStreams:    factory.IOStreams{Out: out, Err: &bytes.Buffer{}},
		OutputFormat: "table",
	}, out
}
