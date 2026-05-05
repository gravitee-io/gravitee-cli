package app

import (
	"bytes"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

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
