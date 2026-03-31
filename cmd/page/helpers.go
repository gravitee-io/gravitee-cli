package page

import (
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/printer"
)

func boolField(item interface{}, key string) string {
	m, ok := item.(map[string]interface{})
	if !ok {
		return ""
	}

	v, ok := m[key]
	if !ok {
		return ""
	}

	b, ok := v.(bool)
	if !ok {
		return ""
	}

	return fmt.Sprintf("%t", b)
}

func printPageDetail(p *printer.Printer, data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, field := range []struct{ label, key string }{
		{"Name", "name"},
		{"ID", "id"},
		{"API", "apiId"},
		{"Type", "type"},
		{"Visibility", "visibility"},
		{"Published", "published"},
		{"Parent", "parentId"},
		{"Created", "createdAt"},
		{"Updated", "updatedAt"},
	} {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	return nil
}
