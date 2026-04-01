package api

import (
	"testing"

	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestListAPIs(t *testing.T) {
	t.Run("returns APIs from the environment", func(t *testing.T) {
		fake := paginatedAPIs(
			map[string]any{"name": "Weather API", "state": "STARTED"},
			map[string]any{"name": "Petstore", "state": "STOPPED"},
		)
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Weather API")
		testutil.AssertOutputContains(t, tc.Out, "Petstore")
	})

	t.Run("rejects invalid token with hint", func(t *testing.T) {
		fake := testutil.APIFailingWith(401, "authentication failed")
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertErrorContains(t, err, "authentication failed")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
