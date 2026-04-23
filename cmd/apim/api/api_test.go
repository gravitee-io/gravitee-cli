// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestListAPIs(t *testing.T) {
	t.Run("returns APIs from the environment", func(t *testing.T) {
		fake := paginatedAPIs(
			map[string]any{"name": "Weather API", "state": "STARTED"},
			map[string]any{"name": "Petstore", "state": "STOPPED"},
		)
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Weather API")
		testutil.AssertOutputContains(t, tc.Out, "Petstore")
	})

	t.Run("rejects invalid token with hint", func(t *testing.T) {
		fake := testutil.APIFailingWith(401, "authentication failed")
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertErrorContains(t, err, "authentication failed")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}

func TestGetAPI(t *testing.T) {
	t.Run("returns API details", func(t *testing.T) {
		fake := testutil.APIReturningItem(map[string]any{
			"id": "api-1", "name": "Weather API", "state": "STARTED",
		})
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newGetCmd(tc.Factory), "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Weather API")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		fake := testutil.APIReturningItem(map[string]any{
			"id": "api-1", "state": "STARTED",
		})
		tc := testutil.NewFactory(fake)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newGetCmd(tc.Factory), "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"state"`)
	})

	t.Run("returns not found with hint", func(t *testing.T) {
		fake := testutil.APIFailingWith(404, "resource not found")
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newGetCmd(tc.Factory), "api-999")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("rejects empty apiId before calling the API", func(t *testing.T) {
		called := false
		fake := &client.FakeClient{
			GetFunc: func(_ string) ([]byte, error) {
				called = true

				return nil, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newGetCmd(tc.Factory), "")

		testutil.AssertErrorContains(t, err, "apiId cannot be empty")

		if called {
			t.Fatal("expected no API call for empty apiId")
		}
	})
}

func TestCreateAPI(t *testing.T) {
	t.Run("creates an API from a JSON file", func(t *testing.T) {
		file := writeTempJSON(t, `{"name":"Test API"}`)
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis")

				return json.Marshal(map[string]string{"id": "new-id", "name": "Test API"})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newCreateCmd(tc.Factory), "-f", file)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Test API")
	})

	t.Run("fails when input file is missing", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		err := testutil.Execute(newCreateCmd(tc.Factory), "-f", "/nonexistent/api.json")

		testutil.AssertErrorContains(t, err, "failed to read")
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
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newUpdateCmd(tc.Factory), "api-1", "-f", file)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Updated")
	})
}

func TestDeleteAPI(t *testing.T) {
	t.Run("deletes the API", func(t *testing.T) {
		fake := &client.FakeClient{
			DeleteFunc: func(path string) error {
				testutil.AssertPathCalled(t, path, "/apis/api-1")

				if strings.Contains(path, "closePlans") {
					t.Errorf("expected no closePlans without --force, got: %s", path)
				}

				return nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newDeleteCmd(tc.Factory), "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "deleted")
	})

	t.Run("--force closes plans and deletes when API is not running", func(t *testing.T) {
		var deletePath string

		fake := &client.FakeClient{
			GetFunc: func(_ string) ([]byte, error) {
				data, _ := json.Marshal(map[string]any{"id": "api-1", "state": "STOPPED"})

				return data, nil
			},
			PostFunc: func(path string, _ any) ([]byte, error) {
				t.Errorf("expected no stop call when API already stopped, got POST %s", path)

				return nil, nil
			},
			DeleteFunc: func(path string) error {
				deletePath = path

				return nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newDeleteCmd(tc.Factory), "api-1", "--force")

		testutil.AssertNoError(t, err)

		if !strings.Contains(deletePath, "closePlans=true") {
			t.Errorf("expected closePlans=true with --force, got: %s", deletePath)
		}

		testutil.AssertOutputContains(t, tc.Out, "plans closed")
	})

	t.Run("--force stops the API when running, then deletes", func(t *testing.T) {
		var stopped bool

		fake := &client.FakeClient{
			GetFunc: func(_ string) ([]byte, error) {
				data, _ := json.Marshal(map[string]any{"id": "api-1", "state": "STARTED"})

				return data, nil
			},
			PostFunc: func(path string, _ any) ([]byte, error) {
				if !strings.Contains(path, "/_stop") {
					t.Errorf("expected /_stop call, got: %s", path)
				}

				stopped = true

				return nil, nil
			},
			DeleteFunc: func(_ string) error { return nil },
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newDeleteCmd(tc.Factory), "api-1", "--force")

		testutil.AssertNoError(t, err)

		if !stopped {
			t.Fatal("expected stop call before delete")
		}

		testutil.AssertOutputContains(t, tc.Out, "stopped and deleted")
	})

	t.Run("hints --force when server rejects with 400", func(t *testing.T) {
		fake := &client.FakeClient{
			DeleteFunc: func(_ string) error {
				return &client.APIError{Status: 400, Message: "API is running"}
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newDeleteCmd(tc.Factory), "api-1")

		testutil.AssertErrorContains(t, err, "--force")
	})

	t.Run("emits {id, status} envelope in json even with --force", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(_ string) ([]byte, error) {
				data, _ := json.Marshal(map[string]any{"id": "api-1", "state": "STOPPED"})

				return data, nil
			},
			DeleteFunc: func(_ string) error { return nil },
		}
		tc := testutil.NewFactory(fake)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newDeleteCmd(tc.Factory), "api-1", "--force")

		testutil.AssertNoError(t, err)

		var got map[string]string
		if jsonErr := json.Unmarshal(tc.Out.Bytes(), &got); jsonErr != nil {
			t.Fatalf("expected valid JSON, got: %s", tc.Out.String())
		}

		if got["id"] != "api-1" || got["status"] != "deleted" {
			t.Errorf("unexpected json: %+v", got)
		}
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
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newStartCmd(tc.Factory), "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "started")
	})

	t.Run("emits {id, status} envelope in json", func(t *testing.T) {
		fake := &client.FakeClient{
			PostFunc: func(_ string, _ any) ([]byte, error) {
				return nil, nil
			},
		}
		tc := testutil.NewFactory(fake)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newStartCmd(tc.Factory), "api-1")

		testutil.AssertNoError(t, err)

		var got map[string]string
		if jsonErr := json.Unmarshal(tc.Out.Bytes(), &got); jsonErr != nil {
			t.Fatalf("expected valid JSON, got: %s", tc.Out.String())
		}

		if got["id"] != "api-1" || got["status"] != "started" {
			t.Errorf("unexpected json: %+v", got)
		}
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
		tc := testutil.NewFactory(fake)

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
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newDeployCmd(tc.Factory), "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "deployment requested")
	})

	t.Run("rejects labels exceeding 32 characters", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		err := testutil.Execute(newDeployCmd(tc.Factory), "api-1", "--label", strings.Repeat("x", 33))

		testutil.AssertErrorContains(t, err, "exceeds 32 characters")
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
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newRollbackCmd(tc.Factory), "api-1", "--event-id", "evt-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "rolled back")
	})

	t.Run("emits {id, eventId, status} envelope in json", func(t *testing.T) {
		fake := &client.FakeClient{
			PostFunc: func(_ string, _ any) ([]byte, error) {
				return nil, nil
			},
		}
		tc := testutil.NewFactory(fake)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newRollbackCmd(tc.Factory), "api-1", "--event-id", "evt-1")

		testutil.AssertNoError(t, err)

		var got map[string]string
		if jsonErr := json.Unmarshal(tc.Out.Bytes(), &got); jsonErr != nil {
			t.Fatalf("expected valid JSON, got: %s", tc.Out.String())
		}

		if got["id"] != "api-1" || got["eventId"] != "evt-1" || got["status"] != "rolled-back" {
			t.Errorf("unexpected json: %+v", got)
		}
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
		tc := testutil.NewFactory(fake)

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
		tc := testutil.NewFactory(fake)

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
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newExportCmd(tc.Factory), "api-1", "--exclude", "members", "--exclude", "pages")

		testutil.AssertNoError(t, err)
	})
}

func TestListLogs(t *testing.T) {
	t.Run("returns API request logs", func(t *testing.T) {
		fake := paginatedAPIs(
			map[string]any{"requestId": "req-1", "method": "GET", "status": "200"},
		)
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newLogsCmd(tc.Factory), "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "req-1")
	})

	t.Run("--all aggregates all pages in table output", func(t *testing.T) {
		fake := pagedLogs(map[int][]map[string]any{
			1: {{"requestId": "req-1"}},
			2: {{"requestId": "req-2"}},
		})
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newLogsCmd(tc.Factory), "api-1", "--all")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "req-1")
		testutil.AssertOutputContains(t, tc.Out, "req-2")
	})

	t.Run("--all aggregates all pages in json output", func(t *testing.T) {
		fake := pagedLogs(map[int][]map[string]any{
			1: {{"requestId": "req-1"}},
			2: {{"requestId": "req-2"}},
		})
		tc := testutil.NewFactory(fake)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newLogsCmd(tc.Factory), "api-1", "--all")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "req-1")
		testutil.AssertOutputContains(t, tc.Out, "req-2")
	})
}

func TestGetLog(t *testing.T) {
	t.Run("returns a single log entry", func(t *testing.T) {
		fake := testutil.APIReturningItem(map[string]any{
			"requestId": "req-1", "method": "GET", "status": 200,
		})
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newLogCmd(tc.Factory), "api-1", "req-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "req-1")
	})

	t.Run("returns not found for unknown request", func(t *testing.T) {
		fake := testutil.APIFailingWith(404, "resource not found")
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newLogCmd(tc.Factory), "api-1", "req-999")

		testutil.AssertErrorContains(t, err, "not found")
	})
}

func TestGetAnalytics(t *testing.T) {
	t.Run("returns analytics data", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/analytics")

				return json.Marshal(map[string]any{"type": "COUNT", "count": 4523})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newAnalyticsCmd(tc.Factory), "api-1",
			"--type", "COUNT", "--from", "1700000000000", "--to", "1700000001000")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "4523")
	})

	t.Run("returns not found for unknown API", func(t *testing.T) {
		fake := testutil.APIFailingWith(404, "resource not found")
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newAnalyticsCmd(tc.Factory), "api-999",
			"--from", "1700000000000", "--to", "1700000001000")

		testutil.AssertErrorContains(t, err, "not found")
	})
}

func TestGetHealth(t *testing.T) {
	t.Run("returns health availability", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/health/availability")

				return json.Marshal(map[string]any{
					"availability": map[string]float64{"https://backend.example.com:443": 99.8},
				})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newHealthCmd(tc.Factory), "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "99.8")
	})
}
