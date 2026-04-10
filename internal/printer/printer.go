package printer

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	FormatTable = "table"
	FormatJSON  = "json"
	FormatYAML  = "yaml"
	FormatID    = "id"
)

// IsStructured reports whether the format preserves the full server response
// (json/yaml). Table and id are per-item flat views and should route through PrintList.
func IsStructured(format string) bool {
	return format == FormatJSON || format == FormatYAML
}

// Column defines a column for table output.
type Column struct {
	Value func(item any) string
	Name  string
	Width int
}

// Printer handles output formatting in table, JSON, or YAML.
type Printer struct {
	Out       io.Writer
	Err       io.Writer
	Format    string
	NoHeaders bool
	Quiet     bool
}

// New creates a Printer with the given format and writer.
func New(format string, out, errOut io.Writer, quiet bool, noHeaders bool) *Printer {
	return &Printer{
		Format:    format,
		Out:       out,
		Err:       errOut,
		Quiet:     quiet,
		NoHeaders: noHeaders,
	}
}

// PrintList outputs a list of items in the configured format.
func (p *Printer) PrintList(items any, columns []Column) error {
	if p.Quiet {
		return nil
	}

	switch p.Format {
	case FormatJSON:
		return p.printJSON(items)
	case FormatYAML:
		return p.printYAML(items)
	case FormatID:
		return p.printIDList(items)
	default:
		return p.printTable(items, columns)
	}
}

// PrintDetail outputs a single item in the configured format.
func (p *Printer) PrintDetail(item any) error {
	if p.Quiet {
		return nil
	}

	switch p.Format {
	case FormatJSON:
		return p.printJSON(item)
	case FormatYAML:
		return p.printYAML(item)
	case FormatID:
		return p.printID(item)
	default:
		return p.printJSON(item)
	}
}

// PrintMessage outputs a plain text message (used for action confirmations).
func (p *Printer) PrintMessage(format string, args ...any) {
	if p.Quiet {
		return
	}

	fmt.Fprintf(p.Out, format+"\n", args...)
}

// PrintHint outputs a message to stderr (used for pagination hints and other metadata that should not pollute piped output).
func (p *Printer) PrintHint(format string, args ...any) {
	if p.Quiet {
		return
	}

	fmt.Fprintf(p.Err, format+"\n", args...)
}

func (p *Printer) printID(item any) error {
	m, err := toMap(item)
	if err != nil {
		return err
	}

	if id := idFromMap(m); id != "" {
		fmt.Fprintln(p.Out, id)
	}

	return nil
}

func (p *Printer) printIDList(items any) error {
	raw, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}

	var list []map[string]any
	if err := json.Unmarshal(raw, &list); err != nil {
		return fmt.Errorf("failed to unmarshal items: %w", err)
	}

	for _, m := range list {
		if id := idFromMap(m); id != "" {
			fmt.Fprintln(p.Out, id)
		}
	}

	return nil
}

func toMap(item any) (map[string]any, error) {
	raw, err := json.Marshal(item)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal item: %w", err)
	}

	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return m, nil
}

func idFromMap(m map[string]any) string {
	if id, ok := m["id"].(string); ok {
		return id
	}

	if key, ok := m["key"].(string); ok {
		return key
	}

	return ""
}

func (p *Printer) printJSON(data any) error {
	encoder := json.NewEncoder(p.Out)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

func (p *Printer) printYAML(data any) error {
	// Round-trip through JSON for consistent field naming, then parse as YAML
	// (which is a superset of JSON) to preserve int vs float types - avoids
	// epoch-millis ints being rendered as scientific-notation floats.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	var generic any
	if err := yaml.Unmarshal(jsonData, &generic); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	encoder := yaml.NewEncoder(p.Out)
	encoder.SetIndent(2)

	if err := encoder.Encode(generic); err != nil {
		return fmt.Errorf("failed to encode YAML: %w", err)
	}

	return encoder.Close()
}

func (p *Printer) printTable(data any, columns []Column) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	var items []any
	if err := json.Unmarshal(jsonData, &items); err != nil {
		return fmt.Errorf("failed to unmarshal data as array: %w", err)
	}

	if len(columns) == 0 {
		return nil
	}

	if len(items) == 0 {
		fmt.Fprintln(p.Out, "No results found.")

		return nil
	}

	widths := make([]int, len(columns))
	for i, col := range columns {
		widths[i] = len(col.Name)
		if col.Width > widths[i] {
			widths[i] = col.Width
		}
	}

	for _, item := range items {
		for i, col := range columns {
			val := col.Value(item)
			if len(val) > widths[i] {
				widths[i] = len(val)
			}
		}
	}

	if !p.NoHeaders {
		headers := make([]string, len(columns))
		for i, col := range columns {
			headers[i] = padRight(strings.ToUpper(col.Name), widths[i])
		}

		fmt.Fprintln(p.Out, strings.Join(headers, "  "))
	}

	for _, item := range items {
		values := make([]string, len(columns))
		for i, col := range columns {
			values[i] = padRight(col.Value(item), widths[i])
		}

		fmt.Fprintln(p.Out, strings.Join(values, "  "))
	}

	return nil
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}

	return s + strings.Repeat(" ", width-len(s))
}
