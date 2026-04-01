package metadata

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestListMetadata(t *testing.T) {
	t.Run("returns metadata from the API", func(t *testing.T) {
		fake := paginatedMetadata(map[string]any{
			"key": "team-email", "name": "Team Email",
			"value": "platform-team@company.com", "format": "MAIL",
			"updatedAt": "2026-03-25T14:30:00Z",
		})
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newListCmd(tc.Factory), "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Team Email")
		testutil.AssertOutputContains(t, tc.Out, "MAIL")
	})
}

func TestCreateMetadata(t *testing.T) {
	t.Run("creates metadata from a JSON file", func(t *testing.T) {
		file := writeTempJSON(t, `{"name":"Team Email","value":"platform-team@company.com","format":"MAIL"}`)
		resp, _ := json.Marshal(map[string]any{
			"key": "team-email", "name": "Team Email",
			"value": "platform-team@company.com", "format": "MAIL",
		})
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/metadata")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newCreateCmd(tc.Factory), "--api", "api-1", "-f", file)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Team Email")
		testutil.AssertOutputContains(t, tc.Out, "team-email")
	})

}

func TestUpdateMetadata(t *testing.T) {
	t.Run("updates metadata from a JSON file", func(t *testing.T) {
		file := writeTempJSON(t, `{"name":"Team Email","value":"new-team@company.com","format":"MAIL"}`)
		resp, _ := json.Marshal(map[string]any{
			"key": "team-email", "name": "Team Email",
			"value": "new-team@company.com", "format": "MAIL",
		})
		fake := &client.FakeClient{
			PutFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/metadata/team-email")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newUpdateCmd(tc.Factory), "team-email", "--api", "api-1", "-f", file)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "new-team@company.com")
	})
}

func TestDeleteMetadata(t *testing.T) {
	t.Run("deletes the metadata entry", func(t *testing.T) {
		fake := &client.FakeClient{
			DeleteFunc: func(path string) error {
				testutil.AssertPathCalled(t, path, "/apis/api-1/metadata/team-email")

				return nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newDeleteCmd(tc.Factory), "team-email", "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Metadata 'team-email' deleted.")
	})

}
