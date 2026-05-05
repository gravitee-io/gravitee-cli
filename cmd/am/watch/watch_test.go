package watch

import (
	"strings"
	"testing"
)

func TestBuildDashboardData(t *testing.T) {
	rawEvents := []map[string]interface{}{
		{
			"id": "e1", "type": "USER_LOGIN", "timestamp": float64(1700000000000),
			"outcome": map[string]interface{}{"status": "SUCCESS"},
			"actor":   map[string]interface{}{"displayName": "admin"},
		},
		{
			"id": "e2", "type": "USER_LOGIN", "timestamp": float64(1700000001000),
			"outcome": map[string]interface{}{"status": "FAILURE"},
			"actor":   map[string]interface{}{"displayName": "bob"},
		},
	}
	data := buildDashboardData(rawEvents, "my-domain", "test-ws")
	if data.Stats.Total != 2 {
		t.Errorf("expected 2 total, got %d", data.Stats.Total)
	}
	if data.Stats.Successes != 1 {
		t.Errorf("expected 1 success, got %d", data.Stats.Successes)
	}
	if data.Stats.Failures != 1 {
		t.Errorf("expected 1 failure, got %d", data.Stats.Failures)
	}
	if len(data.Stats.TopTypes) == 0 {
		t.Error("expected top types")
	}
}

func TestRender(t *testing.T) {
	data := DashboardData{
		DomainName:  "my-domain",
		Workspace:   "test-ws",
		RefreshedAt: "2023-11-14 22:13:20",
		Events: []AuditEvent{
			{ID: "e1", EventType: "USER_LOGIN", Status: "SUCCESS", Actor: "admin", Timestamp: "2023-11-14 22:13:20"},
		},
		Stats: DashboardStats{Total: 1, Successes: 1, Failures: 0},
	}
	out := render(data, 5)
	if !strings.Contains(out, "my-domain") {
		t.Error("expected domain name in render")
	}
	if !strings.Contains(out, "USER_LOGIN") {
		t.Error("expected event type in render")
	}
}
