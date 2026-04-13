package user

import (
	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

// newTestUser wires the mock service into the factory for testing.
func newTestUser(tc *testutil.TestContext, mock *am.MockService, _ string) {
	tc.Factory.SetAMService(mock)
}
