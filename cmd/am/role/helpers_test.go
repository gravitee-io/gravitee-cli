package role

import (
	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

// newTestRole wires the mock service into the factory for testing.
func newTestRole(tc *testutil.TestContext, mock *am.MockService, _ string) {
	tc.Factory.SetAMService(mock)
}
