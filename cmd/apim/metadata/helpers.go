package metadata

import (
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/printer"
)

func printMetadataDetail(p *printer.Printer, data []byte, apiID string) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, field := range []struct{ label, key string }{
		{"Key", "key"},
		{"Name", "name"},
		{"Value", "value"},
		{"Format", "format"},
	} {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	p.PrintMessage("%-16s%s", "API:", apiID)

	return nil
}
