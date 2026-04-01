package api

import (
	"testing"

	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestGetAPI(t *testing.T) {
	t.Run("returns API details", func(t *testing.T) {
		fake := testutil.APIReturningItem(map[string]any{
			"id": "api-1", "name": "Weather API", "state": "STARTED",
		})
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newGetCmd(tc.Factory), "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Weather API")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		fake := testutil.APIReturningItem(map[string]any{
			"id": "api-1", "state": "STARTED",
		})
		tc := testutil.NewFactory(fake, false)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newGetCmd(tc.Factory), "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"state"`)
	})

	t.Run("returns not found with hint", func(t *testing.T) {
		fake := testutil.APIFailingWith(404, "resource not found")
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newGetCmd(tc.Factory), "api-999")

		testutil.AssertErrorContains(t, err, "not found")
	})
}
