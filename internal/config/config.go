package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	configDir  = ".gio"
	configFile = "config.json"

	DefaultOrg = "DEFAULT"
	DefaultEnv = "DEFAULT"
)

// Context holds the connection details for a Gravitee APIM instance.
type Context struct {
	URL      string `json:"url"`
	Token    string `json:"token"`
	Org      string `json:"org,omitempty"`
	Env      string `json:"env,omitempty"`
	ReadOnly bool   `json:"readOnly,omitempty"`
}

// Config holds the CLI configuration with named contexts.
type Config struct {
	Contexts       map[string]Context `json:"contexts"`
	CurrentContext string             `json:"currentContext"`
}

// ResolvedContext holds the fully resolved context after applying overrides.
type ResolvedContext struct {
	Name     string
	URL      string
	Token    string
	Org      string
	Env      string
	ReadOnly bool
}

// Overrides holds flag-based overrides applied on top of the config context.
type Overrides struct {
	Context string
	Org     string
	EnvID   string
}

// Path returns the full path to the config file.
func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine home directory: %w", err)
	}

	return filepath.Join(home, configDir, configFile), nil
}

// Load reads the config file from disk. Returns an empty config if the file does not exist.
func Load() (*Config, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}

	return LoadFrom(path)
}

// LoadFrom reads the config file from the given path. Returns an empty config if the file does not exist.
func LoadFrom(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Contexts: make(map[string]Context)}, nil
		}

		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if cfg.Contexts == nil {
		cfg.Contexts = make(map[string]Context)
	}

	return &cfg, nil
}

// Save writes the config to disk, creating the directory if needed.
func (c *Config) Save() error {
	path, err := Path()
	if err != nil {
		return err
	}

	return c.SaveTo(path)
}

// SaveTo writes the config to the given path, creating the directory if needed.
func (c *Config) SaveTo(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Resolve returns the fully resolved context after applying overrides and defaults.
func (c *Config) Resolve(overrides Overrides) (*ResolvedContext, error) {
	contextName := c.CurrentContext
	if overrides.Context != "" {
		contextName = overrides.Context
	}

	if contextName == "" {
		return nil, fmt.Errorf("no context configured\nHint: run 'gio login' to get started")
	}

	ctx, ok := c.Contexts[contextName]
	if !ok {
		available := make([]string, 0, len(c.Contexts))
		for name := range c.Contexts {
			available = append(available, name)
		}

		sort.Strings(available)

		return nil, fmt.Errorf("context '%s' not found\nHint: available contexts: %s. See 'gio config get-contexts'", contextName, strings.Join(available, ", "))
	}

	resolved := &ResolvedContext{
		Name:     contextName,
		URL:      ctx.URL,
		Token:    ctx.Token,
		Org:      withDefault(ctx.Org, DefaultOrg),
		Env:      withDefault(ctx.Env, DefaultEnv),
		ReadOnly: ctx.ReadOnly,
	}

	if overrides.Org != "" {
		resolved.Org = overrides.Org
	}

	if overrides.EnvID != "" {
		resolved.Env = overrides.EnvID
	}

	return resolved, nil
}

func withDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}
