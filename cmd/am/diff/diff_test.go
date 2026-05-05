package diff

import (
	"strings"
	"testing"
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
