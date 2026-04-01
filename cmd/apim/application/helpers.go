package application

import (
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/printer"
)

func ownerDisplayName(item any) string {
	m, ok := item.(map[string]any)
	if !ok {
		return ""
	}

	owner, ok := m["owner"].(map[string]any)
	if !ok {
		return ""
	}

	s, _ := owner["displayName"].(string)

	return s
}

func printAppDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, f := range []struct{ label, key string }{
		{"Name", "name"},
		{"ID", "id"},
		{"Description", "description"},
		{"Type", "type"},
		{"Status", "status"},
	} {
		printField(p, m, f.label, f.key)
	}

	printOwner(p, m)

	for _, f := range []struct{ label, key string }{
		{"API Key Mode", "api_key_mode"},
		{"Domain", "domain"},
		{"Created", "created_at"},
		{"Updated", "updated_at"},
	} {
		printField(p, m, f.label, f.key)
	}

	return nil
}

func printField(p *printer.Printer, m map[string]any, label, key string) {
	if v, ok := m[key]; ok && v != nil {
		p.PrintMessage("%-16s%v", label+":", v)
	}
}

func printOwner(p *printer.Printer, m map[string]any) {
	owner, ok := m["owner"].(map[string]any)
	if !ok {
		return
	}

	if dn, ok := owner["displayName"].(string); ok && dn != "" {
		p.PrintMessage("%-16s%s", "Owner:", dn)
	}
}
