package diff

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/config"
)

func TestDiffObjects(t *testing.T) {
	from := map[string]interface{}{"name": "foo", "description": "old desc"}
	to := map[string]interface{}{"name": "foo", "description": "new desc"}
	changes := diffObjects(from, to, []string{"description"})
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Field != "description" {
		t.Errorf("expected field 'description', got %q", changes[0].Field)
	}
}

func TestCompareResourcesAdded(t *testing.T) {
	fromItems := []map[string]interface{}{
		{"name": "scope-a"},
	}
	toItems := []map[string]interface{}{
		{"name": "scope-a"},
		{"name": "scope-b"},
	}
	result := compareResources(fromItems, toItems, "name", []string{"name"})
	if result.Added != 1 {
		t.Errorf("expected 1 added, got %d", result.Added)
	}
}

func TestCompareResourcesRemoved(t *testing.T) {
	fromItems := []map[string]interface{}{
		{"name": "scope-a"},
		{"name": "scope-b"},
	}
	toItems := []map[string]interface{}{
		{"name": "scope-a"},
	}
	result := compareResources(fromItems, toItems, "name", []string{"name"})
	if result.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", result.Removed)
	}
}

func TestCompareResourcesChanged(t *testing.T) {
	fromItems := []map[string]interface{}{
		{"name": "role-a", "description": "old"},
	}
	toItems := []map[string]interface{}{
		{"name": "role-a", "description": "new"},
	}
	result := compareResources(fromItems, toItems, "name", []string{"description"})
	if result.Changed != 1 {
		t.Errorf("expected 1 changed, got %d", result.Changed)
	}
	if result.Added != 0 || result.Removed != 0 {
		t.Errorf("expected no added/removed, got added=%d removed=%d", result.Added, result.Removed)
	}
	if len(result.Lines) == 0 {
		t.Error("expected diff lines for changed resource")
	}
}

func TestFormatDiffLine(t *testing.T) {
	line := formatDiffLine("+", "scope", "scope-b", nil)
	if !strings.Contains(line, "+") {
		t.Error("expected + prefix")
	}
	if !strings.Contains(line, "scope-b") {
		t.Error("expected resource name")
	}
}

func TestDiffCmd(t *testing.T) {
	// from: has scope-a only
	fromHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		if strings.Contains(r.URL.Path, "scopes") {
			_ = enc.Encode(map[string]interface{}{
				"data": []map[string]interface{}{{"key": "scope-a", "name": "Scope A"}},
			})
			return
		}
		if strings.Contains(r.URL.Path, "roles") ||
			strings.Contains(r.URL.Path, "groups") ||
			strings.Contains(r.URL.Path, "applications") {
			_ = enc.Encode(map[string]interface{}{"data": []interface{}{}})
			return
		}
		_ = enc.Encode([]interface{}{})
	})
	// to: has scope-a + scope-b (one extra)
	toHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		if strings.Contains(r.URL.Path, "scopes") {
			_ = enc.Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"key": "scope-a", "name": "Scope A"},
					{"key": "scope-b", "name": "Scope B"},
				},
			})
			return
		}
		if strings.Contains(r.URL.Path, "roles") ||
			strings.Contains(r.URL.Path, "groups") ||
			strings.Contains(r.URL.Path, "applications") {
			_ = enc.Encode(map[string]interface{}{"data": []interface{}{}})
			return
		}
		_ = enc.Encode([]interface{}{})
	})

	fromServer := httptest.NewServer(fromHandler)
	defer fromServer.Close()
	toServer := httptest.NewServer(toHandler)
	defer toServer.Close()

	f, out := newTestFactory(nil)
	// Point ctx-a and ctx-b at the httptest servers
	f.Config.Contexts["ctx-a"] = &config.Context{
		Org: "DEFAULT", Env: "DEFAULT", Domain: "dom-a", Type: "am",
		AM: &config.ProductConfig{URL: fromServer.URL, Token: "tok-a"},
	}
	f.Config.Contexts["ctx-b"] = &config.Context{
		Org: "DEFAULT", Env: "DEFAULT", Domain: "dom-b", Type: "am",
		AM: &config.ProductConfig{URL: toServer.URL, Token: "tok-b"},
	}

	cmd := NewDiffCmd(f)
	cmd.SetArgs([]string{"--from", "ctx-a", "--to", "ctx-b"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "scopes") {
		t.Errorf("expected 'scopes' in output, got: %s", output)
	}
	if !strings.Contains(output, "+1") {
		t.Errorf("expected '+1' added scope in diff output, got: %s", output)
	}
}
