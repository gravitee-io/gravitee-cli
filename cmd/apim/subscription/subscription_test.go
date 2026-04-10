package subscription

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestListSubscriptions(t *testing.T) {
	t.Run("returns subscriptions for the API", func(t *testing.T) {
		fake := paginatedSubscriptions(
			map[string]any{
				"id": "sub-1", "planId": "plan-1", "applicationId": "app-1",
				"status": "ACCEPTED", "createdAt": "2026-03-20T10:30:00Z",
			},
		)
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newListCmd(tc.Factory), "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "sub-1")
		testutil.AssertOutputContains(t, tc.Out, "ACCEPTED")
	})
}

func TestGetSubscription(t *testing.T) {
	t.Run("returns the subscription details", func(t *testing.T) {
		resp, _ := json.Marshal(map[string]string{
			"id": "sub-1", "planId": "plan-1", "applicationId": "app-1",
			"status": "ACCEPTED", "createdAt": "2026-03-20T10:30:00Z",
			"processedAt": "2026-03-20T10:35:00Z",
		})
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/subscriptions/sub-1")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newGetCmd(tc.Factory), "sub-1", "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "sub-1")
		testutil.AssertOutputContains(t, tc.Out, "ACCEPTED")
	})

	t.Run("reports not found", func(t *testing.T) {
		fake := testutil.APIFailingWith(404, "resource not found (HTTP 404)")
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newGetCmd(tc.Factory), "sub-999", "--api", "api-1")

		testutil.AssertErrorContains(t, err, "not found")
	})
}

func TestCreateSubscription(t *testing.T) {
	t.Run("creates the subscription", func(t *testing.T) {
		resp, _ := json.Marshal(map[string]string{
			"id": "sub-new", "planId": "plan-1", "applicationId": "app-1",
			"status": "PENDING", "createdAt": "2026-03-27T09:00:00Z",
		})
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/subscriptions")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newCreateCmd(tc.Factory), "--api", "api-1", "--plan", "plan-1", "--app", "app-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "sub-new")
		testutil.AssertOutputContains(t, tc.Out, "PENDING")
	})
}

func TestAcceptSubscription(t *testing.T) {
	t.Run("accepts the subscription", func(t *testing.T) {
		resp, _ := json.Marshal(map[string]string{
			"id": "sub-1", "planId": "plan-1", "applicationId": "app-1",
			"status": "ACCEPTED", "createdAt": "2026-03-27T09:00:00Z",
			"processedAt": "2026-03-27T09:10:00Z",
		})
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/subscriptions/sub-1/_accept")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newAcceptCmd(tc.Factory), "sub-1", "--api", "api-1", "--reason", "Approved")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "ACCEPTED")
	})
}

func TestRejectSubscription(t *testing.T) {
	t.Run("rejects the subscription", func(t *testing.T) {
		resp, _ := json.Marshal(map[string]string{
			"id": "sub-1", "planId": "plan-1", "applicationId": "app-1",
			"status": "REJECTED", "createdAt": "2026-03-27T09:00:00Z",
			"processedAt": "2026-03-27T09:10:00Z",
		})
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/subscriptions/sub-1/_reject")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newRejectCmd(tc.Factory), "sub-1", "--api", "api-1", "--reason", "Denied")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "REJECTED")
	})
}

func TestPauseSubscription(t *testing.T) {
	t.Run("pauses the subscription", func(t *testing.T) {
		resp, _ := json.Marshal(map[string]string{
			"id": "sub-1", "planId": "plan-1", "applicationId": "app-1",
			"status": "PAUSED", "createdAt": "2026-03-20T10:30:00Z",
			"pausedAt": "2026-03-27T11:00:00Z",
		})
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/subscriptions/sub-1/_pause")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newPauseCmd(tc.Factory), "sub-1", "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "PAUSED")
	})
}

func TestResumeSubscription(t *testing.T) {
	t.Run("resumes the subscription", func(t *testing.T) {
		resp, _ := json.Marshal(map[string]string{
			"id": "sub-1", "planId": "plan-1", "applicationId": "app-1",
			"status": "ACCEPTED", "createdAt": "2026-03-20T10:30:00Z",
			"processedAt": "2026-03-20T10:35:00Z",
		})
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/subscriptions/sub-1/_resume")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newResumeCmd(tc.Factory), "sub-1", "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "ACCEPTED")
	})
}

func TestCloseSubscription(t *testing.T) {
	t.Run("closes the subscription", func(t *testing.T) {
		resp, _ := json.Marshal(map[string]string{
			"id": "sub-1", "planId": "plan-1", "applicationId": "app-1",
			"status": "CLOSED", "createdAt": "2026-03-20T10:30:00Z",
			"closedAt": "2026-03-27T15:00:00Z",
		})
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/subscriptions/sub-1/_close")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newCloseCmd(tc.Factory), "sub-1", "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "CLOSED")
	})
}

func TestTransferSubscription(t *testing.T) {
	t.Run("transfers the subscription to another plan", func(t *testing.T) {
		resp, _ := json.Marshal(map[string]string{
			"id": "sub-1", "planId": "plan-2", "applicationId": "app-1",
			"status": "ACCEPTED", "createdAt": "2026-03-20T10:30:00Z",
			"processedAt": "2026-03-20T10:35:00Z",
		})
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/subscriptions/sub-1/_transfer")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newTransferCmd(tc.Factory), "sub-1", "--api", "api-1", "--plan", "plan-2")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "plan-2")
	})
}
