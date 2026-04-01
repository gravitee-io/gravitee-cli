package cmdutil

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

// SetupConfig loads the unified config file. Idempotent — skips if already loaded.
// Called from product PersistentPreRunE and root-level commands (login, context).
func SetupConfig(f *factory.Factory) error {
	if f.Config != nil {
		return nil
	}

	p, err := config.Path()
	if err != nil {
		return err
	}

	f.ConfigPath = p

	cfg, err := config.LoadFrom(p)
	if err != nil {
		return err
	}

	f.Config = cfg

	return nil
}

// ResolveProductContext resolves the product-specific block from the unified config
// and creates the HTTP client. Called from each product's PersistentPreRunE.
func ResolveProductContext(f *factory.Factory, product string) error {
	f.Product = product

	resolved, err := f.Config.Resolve(f.Overrides, product)
	if err != nil {
		// Allow commands that don't need a context (login, context).
		f.ContextResolveErr = err

		return nil
	}

	f.Resolved = resolved
	f.Client = client.NewHTTPClient(client.HTTPClientConfig{
		BaseURL:  resolved.URL,
		Token:    resolved.Token,
		Debug:    f.Debug,
		DebugOut: f.IOStreams.Err,
	})

	return nil
}

// ContextHeader returns a formatted string showing the current context info for help display.
func ContextHeader(f *factory.Factory, product string) string {
	if f.Resolved == nil {
		return ""
	}

	r := f.Resolved
	mode := "read-write"
	if r.ReadOnly {
		mode = "read-only"
	}

	return fmt.Sprintf("\n  Context:   %s\n  URL:       %s\n  Org:       %s\n  Env:       %s\n  Mode:      %s\n",
		r.Name, r.URL, r.Org, r.Env, mode)
}

// AddOutputFlags adds -o/--output, -q/--quiet, and --no-headers as persistent flags on a command.
// Call this on resource parent commands (api, plan, domain...) so that only
// commands producing output expose these flags.
func AddOutputFlags(cmd *cobra.Command, f *factory.Factory) {
	cmd.PersistentFlags().StringVarP(&f.OutputFormat, "output", "o", "table", "Output format: table, json, yaml")
	cmd.PersistentFlags().BoolVarP(&f.Quiet, "quiet", "q", false, "Suppress output except errors")
	cmd.PersistentFlags().BoolVar(&f.NoHeaders, "no-headers", false, "Hide table headers (only applies to table output)")
}

// RequireContext returns an error if no context is configured.
func RequireContext(f *factory.Factory) error {
	if f.ContextResolveErr != nil {
		return f.ContextResolveErr
	}

	if f.Resolved == nil {
		product := f.Product
		if product == "" {
			product = "apim"
		}

		return fmt.Errorf("no context configured\nHint: run 'gio login %s' to get started", product)
	}

	return nil
}

// V2EnvPath builds a V2 environment-scoped API path.
func V2EnvPath(f *factory.Factory, path string) string {
	return client.V2Path(f.Resolved.Env, path)
}

// V1EnvPath builds a V1 org+env-scoped API path.
func V1EnvPath(f *factory.Factory, path string) string {
	return client.V1Path(f.Resolved.Org, f.Resolved.Env, path)
}

// ReadJSONFile reads a JSON file and returns its raw content.
func ReadJSONFile(path string) (json.RawMessage, error) {
	path = filepath.Clean(path)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read '%s': %w", path, err)
	}

	if !json.Valid(data) {
		return nil, fmt.Errorf("failed to read '%s': invalid JSON", path)
	}

	return data, nil
}

// NewPrinter creates a Printer from the factory settings, validating the output format.
func NewPrinter(f *factory.Factory) (*printer.Printer, error) {
	switch f.OutputFormat {
	case printer.FormatTable, printer.FormatJSON, printer.FormatYAML:
		return printer.New(f.OutputFormat, f.IOStreams.Out, f.Quiet, f.NoHeaders), nil
	default:
		return nil, fmt.Errorf("invalid output format %q\nHint: allowed values are table, json, yaml", f.OutputFormat)
	}
}

// MaskToken masks a token for display, showing only the last 3 characters.
func MaskToken(token string) string {
	if len(token) <= 3 {
		return "***"
	}

	return strings.Repeat("*", len(token)-3) + token[len(token)-3:]
}

// PrintPaginationHint prints the "Showing X-Y of Z (page M/N)" message for APIM paginated responses.
func PrintPaginationHint(p *printer.Printer, page, perPage, pageCount, totalCount, pageItemsCount int, isAll bool) {
	if totalCount == 0 {
		return
	}

	start := (page-1)*perPage + 1
	end := start + pageItemsCount - 1

	if pageCount > 1 {
		hint := " Use --all to fetch all results."
		if isAll || page == pageCount {
			hint = ""
		}

		p.PrintMessage("Showing %d-%d of %d (page %d/%d).%s",
			start, end, totalCount, page, pageCount, hint)
	} else {
		p.PrintMessage("Showing %d-%d of %d (page %d/%d).",
			start, end, totalCount, page, pageCount)
	}
}

// ValidateURL checks that a URL is well-formed with a scheme and host.
func ValidateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL is required")
	}

	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return fmt.Errorf("invalid URL %q\nHint: URL must start with http:// or https://", rawURL)
	}

	return nil
}

// StringField extracts a string value from a map[string]any.
func StringField(item any, key string) string {
	m, ok := item.(map[string]any)
	if !ok {
		return ""
	}

	s, _ := m[key].(string)

	return s
}

// ValidateEnum checks that a value is in the allowed set.
func ValidateEnum(value, flag string, allowed []string) error {
	for _, a := range allowed {
		if value == a {
			return nil
		}
	}

	return fmt.Errorf("invalid value '%s' for flag --%s\nHint: allowed values are %s", value, flag, strings.Join(allowed, ", "))
}
