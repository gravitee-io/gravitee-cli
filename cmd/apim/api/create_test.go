package api

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestCreateAPI(t *testing.T) {
	t.Run("creates an API from a JSON file", func(t *testing.T) {
		file := writeTempJSON(t, `{"name":"Test API"}`)
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis")

				return json.Marshal(map[string]string{"id": "new-id", "name": "Test API"})
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newCreateCmd(tc.Factory), "-f", file)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Test API")
	})

	t.Run("fails when input file is missing", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		err := testutil.Execute(newCreateCmd(tc.Factory), "-f", "/nonexistent/api.json")

		testutil.AssertErrorContains(t, err, "failed to read")
	})
}
