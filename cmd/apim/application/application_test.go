package application

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestListApplications(t *testing.T) {
	t.Run("returns applications from the environment", func(t *testing.T) {
		fake := paginatedApps(
			map[string]any{
				"id": "app-1", "name": "My Mobile App", "type": "SIMPLE",
				"status": "ACTIVE", "owner": map[string]any{"displayName": "john.doe"},
				"updated_at": "2026-03-25T14:30:00Z",
			},
		)
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Mobile App")
		testutil.AssertOutputContains(t, tc.Out, "john.doe")
	})

	t.Run("calls the V1 paged endpoint", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/applications/_paged?")
				testutil.AssertPathCalled(t, path, "organizations/DEFAULT/environments/DEFAULT")

				return emptyPaginatedResponse(), nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertNoError(t, err)
	})
}

func TestGetApplication(t *testing.T) {
	t.Run("returns application details", func(t *testing.T) {
		fake := testutil.APIReturningItem(appJSON())
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newGetCmd(tc.Factory), "app-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Mobile App")
		testutil.AssertOutputContains(t, tc.Out, "john.doe")
		testutil.AssertOutputContains(t, tc.Out, "ACTIVE")
	})

	t.Run("calls the V1 endpoint with org/env", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/applications/app-1")
				testutil.AssertPathCalled(t, path, "organizations/DEFAULT/environments/DEFAULT")

				resp, _ := json.Marshal(appJSON())

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newGetCmd(tc.Factory), "app-1")

		testutil.AssertNoError(t, err)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		fake := testutil.APIFailingWith(404, "resource not found (HTTP 404)")
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newGetCmd(tc.Factory), "app-999")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("rejects empty appId before calling the API", func(t *testing.T) {
		called := false
		fake := &client.FakeClient{
			GetFunc: func(_ string) ([]byte, error) {
				called = true

				return nil, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newGetCmd(tc.Factory), "")

		testutil.AssertErrorContains(t, err, "appId cannot be empty")

		if called {
			t.Fatal("expected no API call for empty appId")
		}
	})
}

func TestCreateApplication(t *testing.T) {
	t.Run("creates the application from a JSON file", func(t *testing.T) {
		file := writeTempJSON(t, `{"name":"My Mobile App","description":"Mobile client"}`)
		resp, _ := json.Marshal(appJSON())
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/applications")
				testutil.AssertPathCalled(t, path, "organizations/DEFAULT/environments/DEFAULT")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newCreateCmd(tc.Factory), "-f", file)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Mobile App")
	})
}

func TestDeleteApplication(t *testing.T) {
	t.Run("deletes the application", func(t *testing.T) {
		fake := &client.FakeClient{
			DeleteFunc: func(path string) error {
				testutil.AssertPathCalled(t, path, "/applications/app-1")
				testutil.AssertPathCalled(t, path, "organizations/DEFAULT/environments/DEFAULT")

				return nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newDeleteCmd(tc.Factory), "app-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "deleted")
	})
}
