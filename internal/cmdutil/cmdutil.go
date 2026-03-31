package cmdutil

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

// RequireContext returns an error if no context is configured.
func RequireContext(f *factory.Factory) error {
	if f.Resolved == nil {
		return fmt.Errorf("no context configured\nHint: run 'gio login' to get started")
	}

	return nil
}

// CheckReadOnly returns an error if the current context is read-only.
func CheckReadOnly(f *factory.Factory, cmdName string) error {
	if f.Resolved != nil && f.Resolved.ReadOnly {
		return fmt.Errorf("'%s' is not available in read-only mode (context: %s)", cmdName, f.Resolved.Name)
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

// NewPrinter creates a Printer from the factory settings.
func NewPrinter(f *factory.Factory) *printer.Printer {
	return printer.New(f.OutputFormat, f.IOStreams.Out, f.Quiet)
}

// StringField extracts a string value from a map[string]interface{}.
func StringField(item interface{}, key string) string {
	m, ok := item.(map[string]interface{})
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
