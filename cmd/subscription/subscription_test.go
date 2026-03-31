package subscription

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func TestListSuccess(t *testing.T) {
	resp := map[string]interface{}{
		"data": []map[string]string{
			{
				"id":            "sub-1",
				"planId":        "plan-1",
				"applicationId": "app-1",
				"status":        "ACCEPTED",
				"createdAt":     "2026-03-20T10:30:00Z",
			},
		},
		"pagination": map[string]int{
			"page": 1, "perPage": 10, "pageCount": 1, "totalCount": 1, "pageItemsCount": 1,
		},
	}

	data, _ := json.Marshal(resp)

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/subscriptions?") {
				t.Errorf("unexpected path: %s", path)
			}

			return data, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newListCmd(f)
	cmd.SetArgs([]string{"--api", "api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "sub-1") {
		t.Errorf("expected 'sub-1' in output, got: %s", out.String())
	}

	if !strings.Contains(out.String(), "ACCEPTED") {
		t.Errorf("expected 'ACCEPTED' in output, got: %s", out.String())
	}
}

func TestGetSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]string{
		"id":            "sub-1",
		"planId":        "plan-1",
		"applicationId": "app-1",
		"status":        "ACCEPTED",
		"createdAt":     "2026-03-20T10:30:00Z",
		"processedAt":   "2026-03-20T10:35:00Z",
	})

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/subscriptions/sub-1") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newGetCmd(f)
	cmd.SetArgs([]string{"sub-1", "--api", "api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "sub-1") {
		t.Errorf("expected 'sub-1' in output, got: %s", out.String())
	}

	if !strings.Contains(out.String(), "ACCEPTED") {
		t.Errorf("expected 'ACCEPTED' in output, got: %s", out.String())
	}
}

func TestGetNotFound(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(_ string) ([]byte, error) {
			return nil, &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
		},
	}

	f, _ := newTestFactory(fake, false)

	cmd := newGetCmd(f)
	cmd.SetArgs([]string{"sub-999", "--api", "api-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestCreateSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]string{
		"id":            "sub-new",
		"planId":        "plan-1",
		"applicationId": "app-1",
		"status":        "PENDING",
		"createdAt":     "2026-03-27T09:00:00Z",
	})

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/subscriptions") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newCreateCmd(f)
	cmd.SetArgs([]string{"--api", "api-1", "--plan", "plan-1", "--app", "app-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "sub-new") {
		t.Errorf("expected 'sub-new' in output, got: %s", out.String())
	}

	if !strings.Contains(out.String(), "PENDING") {
		t.Errorf("expected 'PENDING' in output, got: %s", out.String())
	}
}

func TestCreateReadOnly(t *testing.T) {
	fake := &client.FakeClient{}
	f, _ := newTestFactory(fake, true)

	cmd := newCreateCmd(f)
	cmd.SetArgs([]string{"--api", "api-1", "--plan", "plan-1", "--app", "app-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}

func TestAcceptSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]string{
		"id":            "sub-1",
		"planId":        "plan-1",
		"applicationId": "app-1",
		"status":        "ACCEPTED",
		"createdAt":     "2026-03-27T09:00:00Z",
		"processedAt":   "2026-03-27T09:10:00Z",
	})

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/subscriptions/sub-1/_accept") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newAcceptCmd(f)
	cmd.SetArgs([]string{"sub-1", "--api", "api-1", "--reason", "Approved"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "ACCEPTED") {
		t.Errorf("expected 'ACCEPTED' in output, got: %s", out.String())
	}
}

func TestAcceptReadOnly(t *testing.T) {
	fake := &client.FakeClient{}
	f, _ := newTestFactory(fake, true)

	cmd := newAcceptCmd(f)
	cmd.SetArgs([]string{"sub-1", "--api", "api-1"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}

func TestRejectSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]string{
		"id":            "sub-1",
		"planId":        "plan-1",
		"applicationId": "app-1",
		"status":        "REJECTED",
		"createdAt":     "2026-03-27T09:00:00Z",
		"processedAt":   "2026-03-27T09:10:00Z",
	})

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/subscriptions/sub-1/_reject") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newRejectCmd(f)
	cmd.SetArgs([]string{"sub-1", "--api", "api-1", "--reason", "Denied"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "REJECTED") {
		t.Errorf("expected 'REJECTED' in output, got: %s", out.String())
	}
}

func TestPauseSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]string{
		"id":            "sub-1",
		"planId":        "plan-1",
		"applicationId": "app-1",
		"status":        "PAUSED",
		"createdAt":     "2026-03-20T10:30:00Z",
		"pausedAt":      "2026-03-27T11:00:00Z",
	})

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/subscriptions/sub-1/_pause") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newPauseCmd(f)
	cmd.SetArgs([]string{"sub-1", "--api", "api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "PAUSED") {
		t.Errorf("expected 'PAUSED' in output, got: %s", out.String())
	}
}

func TestResumeSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]string{
		"id":            "sub-1",
		"planId":        "plan-1",
		"applicationId": "app-1",
		"status":        "ACCEPTED",
		"createdAt":     "2026-03-20T10:30:00Z",
		"processedAt":   "2026-03-20T10:35:00Z",
	})

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/subscriptions/sub-1/_resume") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newResumeCmd(f)
	cmd.SetArgs([]string{"sub-1", "--api", "api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "ACCEPTED") {
		t.Errorf("expected 'ACCEPTED' in output, got: %s", out.String())
	}
}

func TestCloseSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]string{
		"id":            "sub-1",
		"planId":        "plan-1",
		"applicationId": "app-1",
		"status":        "CLOSED",
		"createdAt":     "2026-03-20T10:30:00Z",
		"closedAt":      "2026-03-27T15:00:00Z",
	})

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/subscriptions/sub-1/_close") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newCloseCmd(f)
	cmd.SetArgs([]string{"sub-1", "--api", "api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "CLOSED") {
		t.Errorf("expected 'CLOSED' in output, got: %s", out.String())
	}
}

func TestTransferSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]string{
		"id":            "sub-1",
		"planId":        "plan-2",
		"applicationId": "app-1",
		"status":        "ACCEPTED",
		"createdAt":     "2026-03-20T10:30:00Z",
		"processedAt":   "2026-03-20T10:35:00Z",
	})

	fake := &client.FakeClient{
		PostFunc: func(path string, _ interface{}) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/subscriptions/sub-1/_transfer") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newTransferCmd(f)
	cmd.SetArgs([]string{"sub-1", "--api", "api-1", "--plan", "plan-2"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "plan-2") {
		t.Errorf("expected 'plan-2' in output, got: %s", out.String())
	}
}

func TestTransferReadOnly(t *testing.T) {
	fake := &client.FakeClient{}
	f, _ := newTestFactory(fake, true)

	cmd := newTransferCmd(f)
	cmd.SetArgs([]string{"sub-1", "--api", "api-1", "--plan", "plan-2"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "read-only mode") {
		t.Errorf("expected read-only error, got: %v", err)
	}
}
