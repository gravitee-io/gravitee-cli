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

package lint

import (
	"bytes"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(c client.GraviteeClient, _ bool) (*factory.Factory, *bytes.Buffer) {
	out := &bytes.Buffer{}
	cfg := &config.Config{
		Contexts: map[string]*config.Context{"test": {Org: "DEFAULT", Env: "DEFAULT", AM: &config.ProductConfig{URL: "http://am", Token: "tok"}}},
		Current:  "test",
	}
	f := &factory.Factory{
		Config:    cfg,
		Resolved:  &config.ResolvedContext{Name: "test", URL: "http://am", Token: "tok", Org: "DEFAULT", Env: "DEFAULT", Domain: "dom1", Type: "am"},
		Client:    c,
		IOStreams: factory.IOStreams{Out: out},
	}
	return f, out
}
