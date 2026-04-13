package scope

import (
	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

// newTestScope wires the mock service into the factory for testing.
func newTestScope(tc *testutil.TestContext, mock *am.MockService, _ string) {
	tc.Factory.SetAMService(mock)
}
