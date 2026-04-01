package api

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestDeleteAPI(t *testing.T) {
	t.Run("deletes the API", func(t *testing.T) {
		fake := &client.FakeClient{
			DeleteFunc: func(path string) error {
				testutil.AssertPathCalled(t, path, "/apis/api-1")

				return nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newDeleteCmd(tc.Factory), "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "deleted")
	})

	t.Run("passes --close-plans to the API", func(t *testing.T) {
		fake := &client.FakeClient{
			DeleteFunc: func(path string) error {
				if !strings.Contains(path, "closePlans=true") {
					t.Errorf("expected closePlans param, got: %s", path)
				}

				return nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newDeleteCmd(tc.Factory), "api-1", "--close-plans")

		testutil.AssertNoError(t, err)
	})

}

func TestStartAPI(t *testing.T) {
	t.Run("starts the API", func(t *testing.T) {
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/_start")

				return nil, nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newStartCmd(tc.Factory), "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "started")
	})

}

func TestStopAPI(t *testing.T) {
	t.Run("stops the API", func(t *testing.T) {
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/_stop")

				return nil, nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newStopCmd(tc.Factory), "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "stopped")
	})
}

func TestDeployAPI(t *testing.T) {
	t.Run("requests a deployment", func(t *testing.T) {
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/deployments")

				return nil, nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newDeployCmd(tc.Factory), "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "deployment requested")
	})

	t.Run("rejects labels exceeding 32 characters", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		err := testutil.Execute(newDeployCmd(tc.Factory), "api-1", "--label", strings.Repeat("x", 33))

		testutil.AssertErrorContains(t, err, "exceeds 32 characters")
	})
}

func TestUpdateAPI(t *testing.T) {
	t.Run("updates the API from a JSON file", func(t *testing.T) {
		file := writeTempJSON(t, `{"name":"Updated"}`)
		fake := &client.FakeClient{
			PutFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1")

				return json.Marshal(map[string]string{"id": "api-1", "name": "Updated"})
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newUpdateCmd(tc.Factory), "api-1", "-f", file)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Updated")
	})
}

func TestImportAPI(t *testing.T) {
	t.Run("imports an API definition", func(t *testing.T) {
		file := writeTempJSON(t, `{"api":{"name":"Imported"}}`)
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/_import/definition")

				return json.Marshal(map[string]string{"id": "new-id", "name": "Imported"})
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newImportCmd(tc.Factory), "-f", file)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Imported")
	})

}

func TestExportAPI(t *testing.T) {
	t.Run("exports the API definition", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/_export/definition")

				return []byte(`{"api":{"name":"Weather API"}}`), nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newExportCmd(tc.Factory), "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Weather API")
	})

	t.Run("passes exclude params to the API", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				if !strings.Contains(path, "excludeAdditionalData=members,pages") {
					t.Errorf("expected exclude params, got: %s", path)
				}

				return []byte(`{}`), nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newExportCmd(tc.Factory), "api-1", "--exclude", "members", "--exclude", "pages")

		testutil.AssertNoError(t, err)
	})
}

func TestRollbackAPI(t *testing.T) {
	t.Run("rolls back to a specific event", func(t *testing.T) {
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/_rollback")

				return nil, nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newRollbackCmd(tc.Factory), "api-1", "--event-id", "evt-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "rolled back")
	})

}
