package plan

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestListPlans(t *testing.T) {
	t.Run("returns plans for the API", func(t *testing.T) {
		fake := paginatedPlans(
			map[string]any{
				"id": "plan-1", "name": "Gold Plan", "status": "PUBLISHED",
				"security": map[string]string{"type": "API_KEY"}, "validation": "AUTO",
				"updatedAt": "2026-03-25T14:30:00Z",
			},
		)
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newListCmd(tc.Factory), "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Gold Plan")
		testutil.AssertOutputContains(t, tc.Out, "API_KEY")
	})

	t.Run("rejects invalid token with hint", func(t *testing.T) {
		fake := testutil.APIFailingWith(401, "authentication failed")
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newListCmd(tc.Factory), "--api", "api-1")

		testutil.AssertErrorContains(t, err, "authentication failed")
	})
}

func TestGetPlan(t *testing.T) {
	t.Run("returns the plan details", func(t *testing.T) {
		resp, _ := json.Marshal(map[string]any{
			"id": "plan-1", "name": "Gold Plan", "apiId": "api-1",
			"status": "PUBLISHED", "security": map[string]string{"type": "API_KEY"},
			"validation": "AUTO", "mode": "STANDARD",
			"createdAt": "2026-03-15T10:00:00Z", "updatedAt": "2026-03-25T14:30:00Z",
		})
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/plans/plan-1")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newGetCmd(tc.Factory), "plan-1", "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Gold Plan")
		testutil.AssertOutputContains(t, tc.Out, "API_KEY")
		testutil.AssertOutputContains(t, tc.Out, "PUBLISHED")
	})

	t.Run("reports not found", func(t *testing.T) {
		fake := testutil.APIFailingWith(404, "resource not found (HTTP 404)")
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newGetCmd(tc.Factory), "plan-999", "--api", "api-1")

		testutil.AssertErrorContains(t, err, "not found")
	})
}

func TestCreatePlan(t *testing.T) {
	t.Run("creates the plan from a JSON file", func(t *testing.T) {
		file := writeTempJSON(t, `{"name":"Gold Plan","security":{"type":"API_KEY"}}`)
		resp, _ := json.Marshal(map[string]any{
			"id": "plan-1", "name": "Gold Plan", "apiId": "api-1",
			"status": "PUBLISHED", "security": map[string]string{"type": "API_KEY"},
			"validation": "AUTO", "mode": "STANDARD",
			"createdAt": "2026-03-15T10:00:00Z", "updatedAt": "2026-03-25T14:30:00Z",
		})
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/plans")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newCreateCmd(tc.Factory), "--api", "api-1", "-f", file)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Gold Plan")
	})

}

func TestUpdatePlan(t *testing.T) {
	t.Run("updates the plan from a JSON file", func(t *testing.T) {
		file := writeTempJSON(t, `{"name":"Gold Plan v2"}`)
		resp, _ := json.Marshal(map[string]any{
			"id": "plan-1", "name": "Gold Plan v2", "apiId": "api-1",
			"status": "PUBLISHED", "security": map[string]string{"type": "API_KEY"},
		})
		fake := &client.FakeClient{
			PutFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/plans/plan-1")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newUpdateCmd(tc.Factory), "plan-1", "--api", "api-1", "-f", file)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Gold Plan v2")
	})

}

func TestDeletePlan(t *testing.T) {
	t.Run("deletes the plan", func(t *testing.T) {
		fake := &client.FakeClient{
			DeleteFunc: func(path string) error {
				testutil.AssertPathCalled(t, path, "/apis/api-1/plans/plan-1")

				return nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newDeleteCmd(tc.Factory), "plan-1", "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "deleted")
	})

}

func TestPublishPlan(t *testing.T) {
	t.Run("publishes the plan", func(t *testing.T) {
		resp, _ := json.Marshal(map[string]any{
			"id": "plan-1", "name": "Gold Plan", "apiId": "api-1",
			"status": "PUBLISHED", "security": map[string]string{"type": "API_KEY"},
			"validation": "AUTO", "mode": "STANDARD",
		})
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/plans/plan-1/_publish")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newPublishCmd(tc.Factory), "plan-1", "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "PUBLISHED")
		testutil.AssertOutputContains(t, tc.Out, "Gold Plan")
	})

	t.Run("reports API error when already published", func(t *testing.T) {
		fake := testutil.PostFailingWith(400, "invalid request (HTTP 400): plan is already published")
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newPublishCmd(tc.Factory), "plan-1", "--api", "api-1")

		testutil.AssertErrorContains(t, err, "already published")
	})
}

func TestDeprecatePlan(t *testing.T) {
	t.Run("deprecates the plan", func(t *testing.T) {
		resp, _ := json.Marshal(map[string]any{
			"id": "plan-1", "name": "Gold Plan", "apiId": "api-1",
			"status": "DEPRECATED", "security": map[string]string{"type": "API_KEY"},
			"validation": "AUTO", "mode": "STANDARD",
		})
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/plans/plan-1/_deprecate")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newDeprecateCmd(tc.Factory), "plan-1", "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "DEPRECATED")
	})

}

func TestClosePlan(t *testing.T) {
	t.Run("closes the plan", func(t *testing.T) {
		resp, _ := json.Marshal(map[string]any{
			"id": "plan-1", "name": "Gold Plan", "apiId": "api-1",
			"status": "CLOSED", "security": map[string]string{"type": "API_KEY"},
			"validation": "AUTO", "mode": "STANDARD",
			"closedAt": "2026-03-27T15:00:00Z",
		})
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/plans/plan-1/_close")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newCloseCmd(tc.Factory), "plan-1", "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "CLOSED")
		testutil.AssertOutputContains(t, tc.Out, "2026-03-27T15:00:00Z")
	})

}
