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
)

// Column defines a column for table output.
type Column struct {
	Value func(item any) string
	Name  string
	Width int
}

// Printer handles output formatting in table, JSON, or YAML.
type Printer struct {
	Out       io.Writer
	Format    string
	NoHeaders bool
	Quiet     bool
}

// New creates a Printer with the given format and writer.
func New(format string, out io.Writer, quiet bool, noHeaders bool) *Printer {
	return &Printer{
		Format:    format,
		Out:       out,
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

func (p *Printer) printJSON(data any) error {
	encoder := json.NewEncoder(p.Out)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

func (p *Printer) printYAML(data any) error {
	// Round-trip through JSON to ensure consistent field naming.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	var generic any
	if err := json.Unmarshal(jsonData, &generic); err != nil {
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
	// Convert data to a slice of interfaces via JSON round-trip.
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

	// Calculate column widths.
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

	// Print header.
	if !p.NoHeaders {
		headers := make([]string, len(columns))
		for i, col := range columns {
			headers[i] = padRight(strings.ToUpper(col.Name), widths[i])
		}

		fmt.Fprintln(p.Out, strings.Join(headers, "  "))
	}

	// Print rows.
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
