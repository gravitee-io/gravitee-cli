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

package printer

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintListTable(t *testing.T) {
	var buf bytes.Buffer

	p := New(FormatTable, &buf, &buf, false, false)

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

	p := New(FormatJSON, &buf, &buf, false, false)

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

	p := New(FormatYAML, &buf, &buf, false, false)

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

	p := New(FormatJSON, &buf, &buf, false, false)

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

	p := New(FormatYAML, &buf, &buf, false, false)

	item := map[string]any{
		"id":        "abc",
		"status":    "STARTED",
		"updatedAt": int64(1776784037000),
		"ratio":     1.5,
		"active":    true,
		"meta":      map[string]any{"version": 3},
	}

	if err := p.PrintDetail(item); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	for _, want := range []string{
		"id: abc",
		"status: STARTED",
		"updatedAt: 1776784037000",
		"ratio: 1.5",
		"active: true",
		"version: 3",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("expected %q in YAML output, got: %s", want, output)
		}
	}

	if strings.ContainsAny(output, "eE") && strings.Contains(output, "+") {
		if strings.Contains(output, "e+") || strings.Contains(output, "E+") {
			t.Errorf("unexpected scientific notation in YAML output: %s", output)
		}
	}
}

func TestQuietSuppressesOutput(t *testing.T) {
	var buf bytes.Buffer

	p := New(FormatTable, &buf, &buf, true, false)

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

	p := New(FormatTable, &buf, &buf, false, false)
	p.PrintMessage("Plan '%s' published.", "plan-123")

	if !strings.Contains(buf.String(), "Plan 'plan-123' published.") {
		t.Errorf("unexpected message: %s", buf.String())
	}
}

func TestPrintMessageQuiet(t *testing.T) {
	var buf bytes.Buffer

	p := New(FormatTable, &buf, &buf, true, false)
	p.PrintMessage("should not appear")

	if buf.Len() != 0 {
		t.Error("quiet mode should suppress PrintMessage")
	}
}

func TestPrintDetailID(t *testing.T) {
	var buf bytes.Buffer

	p := New(FormatID, &buf, &buf, false, false)

	item := map[string]string{"id": "abc-123", "name": "Test API", "status": "STARTED"}

	if err := p.PrintDetail(item); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := buf.String(); got != "abc-123\n" {
		t.Errorf("expected %q, got %q", "abc-123\n", got)
	}
}

func TestPrintDetailID_KeyFallback(t *testing.T) {
	var buf bytes.Buffer

	p := New(FormatID, &buf, &buf, false, false)

	// APIKeys expose their identifier as `key`, not `id`.
	item := map[string]string{"key": "apikey-value", "applicationName": "X"}

	if err := p.PrintDetail(item); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := buf.String(); got != "apikey-value\n" {
		t.Errorf("expected %q, got %q", "apikey-value\n", got)
	}
}

func TestPrintDetailID_MissingID(t *testing.T) {
	var buf bytes.Buffer

	p := New(FormatID, &buf, &buf, false, false)

	item := map[string]string{"name": "no-id-here"}

	if err := p.PrintDetail(item); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf("expected empty output when id missing, got %q", buf.String())
	}
}

func TestPrintListID(t *testing.T) {
	var buf bytes.Buffer

	p := New(FormatID, &buf, &buf, false, false)

	items := []map[string]string{
		{"id": "id-1", "name": "A"},
		{"id": "id-2", "name": "B"},
		{"key": "key-3", "name": "C"}, // fallback to key
	}

	if err := p.PrintList(items, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := buf.String(); got != "id-1\nid-2\nkey-3\n" {
		t.Errorf("expected %q, got %q", "id-1\nid-2\nkey-3\n", got)
	}
}

func TestIsStructured(t *testing.T) {
	cases := []struct {
		format string
		want   bool
	}{
		{FormatJSON, true},
		{FormatYAML, true},
		{FormatTable, false},
		{FormatID, false},
		{"", false},
	}

	for _, c := range cases {
		if got := IsStructured(c.format); got != c.want {
			t.Errorf("IsStructured(%q) = %v, want %v", c.format, got, c.want)
		}
	}
}

func TestPrintListTable_NoHeaders(t *testing.T) {
	var buf bytes.Buffer

	p := New(FormatTable, &buf, &buf, false, true)

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
