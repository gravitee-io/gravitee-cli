package testutil

import (
	"bytes"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// TestContext holds a Factory and captured output buffers for testing.
type TestContext struct {
	Factory *factory.Factory
	Out     *bytes.Buffer
	Err     *bytes.Buffer
}

// NewTestFactory creates a Factory configured for testing with the given client and read-only setting.
func NewTestFactory(c client.GraviteeClient, readOnly bool) *TestContext {
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}

	return &TestContext{
		Factory: &factory.Factory{
			Config: &config.Config{
				CurrentContext: "test",
				Contexts: map[string]config.Context{
					"test": {
						URL:      "https://apim-test.company.com",
						Token:    "test-token",
						Org:      "DEFAULT",
						Env:      "DEFAULT",
						ReadOnly: readOnly,
					},
				},
			},
			Client:    c,
			IOStreams: factory.IOStreams{Out: out, Err: errOut},
		},
		Out: out,
		Err: errOut,
	}
}
