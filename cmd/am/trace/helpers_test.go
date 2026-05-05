package trace

import (
	"bytes"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(c client.GraviteeClient) (*factory.Factory, *bytes.Buffer) {
	out := &bytes.Buffer{}
	cfg := &config.Config{
		Contexts: map[string]*config.Context{"test": {Org: "DEFAULT", Env: "DEFAULT", AM: &config.ProductConfig{URL: "http://am", Token: "tok"}}},
		Current:  "test",
	}
	f := &factory.Factory{
		Config:   cfg,
		Resolved: &config.ResolvedContext{Name: "test", URL: "http://am", Token: "tok", Org: "DEFAULT", Env: "DEFAULT", Domain: "dom1", Type: "am"},
		Client:   c,
		IOStreams: factory.IOStreams{Out: out, Err: &bytes.Buffer{}},
	}
	return f, out
}
