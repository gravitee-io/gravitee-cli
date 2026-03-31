package subscription

import (
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/printer"
)

func extractID(m map[string]interface{}, key string) string {
	if nested, ok := m[key].(map[string]interface{}); ok {
		s, _ := nested["id"].(string)

		return s
	}

	s, _ := m[key+"Id"].(string)

	return s
}

func printSubDetail(p *printer.Printer, data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	printSubField(p, m, "ID", "id")

	if plan := extractID(m, "plan"); plan != "" {
		p.PrintMessage("%-16s%s", "Plan:", plan)
	}

	if app := extractID(m, "application"); app != "" {
		p.PrintMessage("%-16s%s", "Application:", app)
	}

	for _, field := range []struct{ label, key string }{
		{"Status", "status"},
		{"Created", "createdAt"},
		{"Processed", "processedAt"},
		{"Starting", "startingAt"},
		{"Ending", "endingAt"},
		{"Closed", "closedAt"},
		{"Paused", "pausedAt"},
	} {
		printSubField(p, m, field.label, field.key)
	}

	return nil
}

func printSubCreateDetail(p *printer.Printer, data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	printSubField(p, m, "ID", "id")

	if plan := extractID(m, "plan"); plan != "" {
		p.PrintMessage("%-16s%s", "Plan:", plan)
	}

	if app := extractID(m, "application"); app != "" {
		p.PrintMessage("%-16s%s", "Application:", app)
	}

	for _, field := range []struct{ label, key string }{
		{"Status", "status"},
		{"Created", "createdAt"},
	} {
		printSubField(p, m, field.label, field.key)
	}

	return nil
}

func printSubField(p *printer.Printer, m map[string]interface{}, label, key string) {
	if v, ok := m[key]; ok && v != nil {
		p.PrintMessage("%-16s%v", label+":", v)
	}
}
