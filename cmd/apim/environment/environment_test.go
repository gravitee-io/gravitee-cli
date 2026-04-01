package environment

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestListEnvironments(t *testing.T) {
	t.Run("returns environments from the organization", func(t *testing.T) {
		fake := environmentList(
			map[string]string{"id": "dev-1111", "name": "Development", "description": "Development environment"},
			map[string]string{"id": "prod-2222", "name": "Production", "description": "Production environment"},
		)
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Development")
		testutil.AssertOutputContains(t, tc.Out, "Production")
	})

	t.Run("calls the correct API path", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/management/organizations/DEFAULT/environments")

				data, _ := json.Marshal([]map[string]string{})

				return data, nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertNoError(t, err)
	})

	t.Run("rejects invalid token with hint", func(t *testing.T) {
		fake := testutil.APIFailingWith(401, "authentication failed (HTTP 401)")
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertErrorContains(t, err, "authentication failed")
	})
}

func TestGetEnvironment(t *testing.T) {
	t.Run("returns environment details", func(t *testing.T) {
		fake := testutil.APIReturningItem(map[string]any{
			"id": "prod-2222", "name": "Production", "description": "Production environment",
		})
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newGetCmd(tc.Factory), "prod-2222")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Production")
		testutil.AssertOutputContains(t, tc.Out, "prod-2222")
	})

	t.Run("calls the correct API path", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/management/organizations/DEFAULT/environments/prod-2222")

				resp, _ := json.Marshal(map[string]string{"id": "prod-2222", "name": "Production"})

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newGetCmd(tc.Factory), "prod-2222")

		testutil.AssertNoError(t, err)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		fake := testutil.APIFailingWith(404, "resource not found (HTTP 404)")
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newGetCmd(tc.Factory), "env-999")

		testutil.AssertErrorContains(t, err, "not found")
	})
}
