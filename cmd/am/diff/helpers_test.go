package diff

import (
	"bytes"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
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
