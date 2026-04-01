package printer

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintListTable(t *testing.T) {
	var buf bytes.Buffer

	p := New(FormatTable, &buf, false, false)

	items := []map[string]string{
		{"name": "Weather API", "status": "STARTED"},
		{"name": "Petstore", "status": "STOPPED"},
	}

	columns := []Column{
		{
			Name: "Name",
			Value: func(item any) string {
				m, ok := item.(map[string]any)
				if !ok {
					return ""
				}

				s, _ := m["name"].(string)

				return s
			},
		},
		{
			Name: "Status",
			Value: func(item any) string {
				m, ok := item.(map[string]any)
				if !ok {
					return ""
				}

				s, _ := m["status"].(string)

				return s
			},
		},
	}

	if err := p.PrintList(items, columns); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "NAME") {
		t.Error("expected table header NAME")
	}

	if !strings.Contains(output, "STATUS") {
		t.Error("expected table header STATUS")
	}

	if !strings.Contains(output, "Weather API") {
		t.Error("expected Weather API in output")
	}

	if !strings.Contains(output, "STOPPED") {
		t.Error("expected STOPPED in output")
	}
}

func TestPrintListJSON(t *testing.T) {
	var buf bytes.Buffer

	p := New(FormatJSON, &buf, false, false)

	items := []map[string]string{
		{"id": "123", "name": "Test API"},
	}

	if err := p.PrintList(items, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, `"id"`) {
		t.Error("expected JSON id field")
	}

	if !strings.Contains(output, `"Test API"`) {
		t.Error("expected Test API in JSON output")
	}
}

func TestPrintListYAML(t *testing.T) {
	var buf bytes.Buffer

	p := New(FormatYAML, &buf, false, false)

	items := []map[string]string{
		{"id": "123", "name": "Test API"},
	}

	if err := p.PrintList(items, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "name: Test API") {
		t.Errorf("expected YAML name field, got: %s", output)
	}
}

func TestPrintDetailJSON(t *testing.T) {
	var buf bytes.Buffer

	p := New(FormatJSON, &buf, false, false)

	item := map[string]string{"id": "abc", "status": "STARTED"}

	if err := p.PrintDetail(item); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, `"status": "STARTED"`) {
		t.Error("expected status field in JSON output")
	}
}

func TestPrintDetailYAML(t *testing.T) {
	var buf bytes.Buffer

	p := New(FormatYAML, &buf, false, false)

	item := map[string]string{"id": "abc", "status": "STARTED"}

	if err := p.PrintDetail(item); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "status: STARTED") {
		t.Errorf("expected YAML status field, got: %s", buf.String())
	}
}

func TestQuietSuppressesOutput(t *testing.T) {
	var buf bytes.Buffer

	p := New(FormatTable, &buf, true, false)

	items := []map[string]string{{"name": "test"}}
	columns := []Column{{Name: "Name", Value: func(_ any) string { return "test" }}}

	if err := p.PrintList(items, columns); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf("quiet mode should produce no output, got: %s", buf.String())
	}
}

func TestPrintMessage(t *testing.T) {
	var buf bytes.Buffer

	p := New(FormatTable, &buf, false, false)
	p.PrintMessage("Plan '%s' published.", "plan-123")

	if !strings.Contains(buf.String(), "Plan 'plan-123' published.") {
		t.Errorf("unexpected message: %s", buf.String())
	}
}

func TestPrintMessageQuiet(t *testing.T) {
	var buf bytes.Buffer

	p := New(FormatTable, &buf, true, false)
	p.PrintMessage("should not appear")

	if buf.Len() != 0 {
		t.Error("quiet mode should suppress PrintMessage")
	}
}

func TestPrintListTable_NoHeaders(t *testing.T) {
	var buf bytes.Buffer

	p := New(FormatTable, &buf, false, true)

	items := []map[string]string{
		{"name": "Weather API"},
		{"name": "Petstore"},
	}

	columns := []Column{
		{
			Name: "Name",
			Value: func(item any) string {
				m, ok := item.(map[string]any)
				if !ok {
					return ""
				}

				s, _ := m["name"].(string)

				return s
			},
		},
	}

	if err := p.PrintList(items, columns); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	if strings.Contains(output, "NAME") {
		t.Error("expected no header row, but found 'NAME' in output")
	}

	if !strings.Contains(output, "Weather API") {
		t.Error("expected data row 'Weather API'")
	}

	if !strings.Contains(output, "Petstore") {
		t.Error("expected data row 'Petstore'")
	}
}
