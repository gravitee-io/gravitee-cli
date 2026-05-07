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
