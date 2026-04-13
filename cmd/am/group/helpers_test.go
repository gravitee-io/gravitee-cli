package group

import (
	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

// newTestGroup wires the mock service into the factory for testing.
func newTestGroup(tc *testutil.TestContext, mock *am.MockService, _ string) {
	tc.Factory.SetAMService(mock)
}
