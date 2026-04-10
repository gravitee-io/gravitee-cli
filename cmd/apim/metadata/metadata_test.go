package metadata

import (
	"testing"

	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestListMetadata(t *testing.T) {
	t.Run("returns metadata from the API", func(t *testing.T) {
		fake := paginatedMetadata(map[string]any{
			"key": "team-email", "name": "Team Email",
			"value": "platform-team@company.com", "format": "MAIL",
			"updatedAt": "2026-03-25T14:30:00Z",
		})
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newListCmd(tc.Factory), "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Team Email")
		testutil.AssertOutputContains(t, tc.Out, "MAIL")
	})
}
