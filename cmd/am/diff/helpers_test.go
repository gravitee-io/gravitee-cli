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

package diff

import (
	"bytes"

	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/config"
	"gravitee.io/gctl/internal/factory"
)

func newTestFactory(c client.GraviteeClient) (*factory.Factory, *bytes.Buffer) {
	out := &bytes.Buffer{}
	cfg := &config.Config{
		Contexts: map[string]*config.Context{
			"ctx-a": {Org: "DEFAULT", Env: "DEFAULT", AM: &config.ProductConfig{URL: "http://am-a", Token: "tok-a"}},
			"ctx-b": {Org: "DEFAULT", Env: "DEFAULT", AM: &config.ProductConfig{URL: "http://am-b", Token: "tok-b"}},
		},
		Current: "ctx-a",
	}
	f := &factory.Factory{
		Config:    cfg,
		Resolved:  &config.ResolvedContext{Name: "ctx-a", URL: "http://am-a", Token: "tok-a", Org: "DEFAULT", Env: "DEFAULT", Domain: "dom1", Type: "am"},
		Client:    c,
		IOStreams: factory.IOStreams{Out: out},
	}
	return f, out
}
