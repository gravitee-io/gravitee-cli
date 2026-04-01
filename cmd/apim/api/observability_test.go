package api

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestGetAnalytics(t *testing.T) {
	t.Run("returns analytics data", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/analytics")

				return json.Marshal(map[string]any{"type": "COUNT", "count": 4523})
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newAnalyticsCmd(tc.Factory), "api-1", "--type", "COUNT")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "4523")
	})

	t.Run("returns not found for unknown API", func(t *testing.T) {
		fake := testutil.APIFailingWith(404, "resource not found")
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newAnalyticsCmd(tc.Factory), "api-999")

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
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newHealthCmd(tc.Factory), "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "99.8")
	})
}

func TestListLogs(t *testing.T) {
	t.Run("returns API request logs", func(t *testing.T) {
		fake := paginatedAPIs(
			map[string]any{"requestId": "req-1", "method": "GET", "status": "200"},
		)
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newLogsCmd(tc.Factory), "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "req-1")
	})
}

func TestGetLog(t *testing.T) {
	t.Run("returns a single log entry", func(t *testing.T) {
		fake := testutil.APIReturningItem(map[string]any{
			"requestId": "req-1", "method": "GET", "status": 200,
		})
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newLogCmd(tc.Factory), "api-1", "req-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "req-1")
	})

	t.Run("returns not found for unknown request", func(t *testing.T) {
		fake := testutil.APIFailingWith(404, "resource not found")
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newLogCmd(tc.Factory), "api-1", "req-999")

		testutil.AssertErrorContains(t, err, "not found")
	})
}
