package api

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func TestAnalyticsSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]interface{}{"type": "COUNT", "count": 4523})

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/analytics") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newAnalyticsCmd(f)
	cmd.SetArgs([]string{"api-1", "--type", "COUNT"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "4523") {
		t.Errorf("expected count in output, got: %s", out.String())
	}
}

func TestAnalyticsError(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(_ string) ([]byte, error) {
			return nil, &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
		},
	}

	f, _ := newTestFactory(fake, false)

	cmd := newAnalyticsCmd(f)
	cmd.SetArgs([]string{"api-999"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestHealthSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]interface{}{
		"availability": map[string]float64{"https://backend.example.com:443": 99.8},
	})

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/health/availability") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newHealthCmd(f)
	cmd.SetArgs([]string{"api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "99.8") {
		t.Errorf("expected availability in output, got: %s", out.String())
	}
}

func TestLogsSuccess(t *testing.T) {
	resp := map[string]interface{}{
		"data": []map[string]interface{}{
			{"requestId": "req-1", "method": "GET", "status": "200", "path": "/test", "timestamp": "2026-03-27"},
		},
		"pagination": map[string]int{"page": 1, "perPage": 10, "pageCount": 1, "totalCount": 1, "pageItemsCount": 1},
	}

	data, _ := json.Marshal(resp)

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/apis/api-1/logs") {
				t.Errorf("unexpected path: %s", path)
			}

			return data, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newLogsCmd(f)
	cmd.SetArgs([]string{"api-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "req-1") {
		t.Errorf("expected 'req-1' in output, got: %s", out.String())
	}
}

func TestLogSuccess(t *testing.T) {
	resp, _ := json.Marshal(map[string]interface{}{
		"requestId": "req-1", "method": "GET", "path": "/test", "status": 200, "responseTime": 45,
	})

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/logs/req-1") {
				t.Errorf("unexpected path: %s", path)
			}

			return resp, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newLogCmd(f)
	cmd.SetArgs([]string{"api-1", "req-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "req-1") {
		t.Errorf("expected 'req-1' in output, got: %s", out.String())
	}
}

func TestLogNotFound(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(_ string) ([]byte, error) {
			return nil, &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
		},
	}

	f, _ := newTestFactory(fake, false)

	cmd := newLogCmd(f)
	cmd.SetArgs([]string{"api-1", "req-999"})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected not found error, got: %v", err)
	}
}
