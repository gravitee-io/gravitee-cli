package apikey

import (
	"bytes"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(fc *client.FakeClient, readOnly bool) (*factory.Factory, *bytes.Buffer) {
	out := &bytes.Buffer{}

	return &factory.Factory{
		Config: &config.Config{
			CurrentContext: "test",
			Contexts: map[string]config.Context{
				"test": {URL: "https://test.com", Token: "tok", Org: "DEFAULT", Env: "DEFAULT", ReadOnly: readOnly},
			},
		},
		Resolved: &config.ResolvedContext{
			Name: "test", URL: "https://test.com", Token: "tok",
			Org: "DEFAULT", Env: "DEFAULT", ReadOnly: readOnly,
		},
		Client:       fc,
		IOStreams:    factory.IOStreams{Out: out, Err: &bytes.Buffer{}},
		OutputFormat: "table",
	}, out
}
