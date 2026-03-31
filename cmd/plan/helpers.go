package plan

import (
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/printer"
)

func securityType(item interface{}) string {
	m, ok := item.(map[string]interface{})
	if !ok {
		return ""
	}

	sec, ok := m["security"].(map[string]interface{})
	if !ok {
		return ""
	}

	s, _ := sec["type"].(string)

	return s
}

func printPlanDetail(p *printer.Printer, data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	fields := []struct {
		label string
		key   string
	}{
		{"Name", "name"},
		{"ID", "id"},
		{"API", "apiId"},
		{"Status", "status"},
		{"Validation", "validation"},
		{"Mode", "mode"},
		{"Created", "createdAt"},
		{"Updated", "updatedAt"},
		{"Published", "publishedAt"},
		{"Closed", "closedAt"},
	}

	securityStr := ""
	if sec, ok := m["security"].(map[string]interface{}); ok {
		if t, ok := sec["type"].(string); ok {
			securityStr = t
		}
	}

	for _, field := range fields {
		if field.key == "status" {
			if v, ok := m[field.key]; ok && v != nil {
				p.PrintMessage("%-16s%v", field.label+":", v)
			}

			if securityStr != "" {
				p.PrintMessage("%-16s%s", "Security:", securityStr)
			}

			continue
		}

		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	return nil
}
