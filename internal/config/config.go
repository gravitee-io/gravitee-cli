package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	configDir  = ".gio"
	configFile = "config.yaml"

	DefaultOrg = "DEFAULT"
	DefaultEnv = "DEFAULT"
)

// ProductConfig holds product-specific connection details.
type ProductConfig struct {
	URL   string `yaml:"url"   json:"url"`
	Token string `yaml:"token" json:"token"`
}

// Context holds shared fields plus per-product config blocks.
type Context struct {
	Org  string         `yaml:"org,omitempty"  json:"org,omitempty"`
	Env  string         `yaml:"env,omitempty"  json:"env,omitempty"`
	APIM *ProductConfig `yaml:"apim,omitempty" json:"apim,omitempty"`
	AM   *ProductConfig `yaml:"am,omitempty"   json:"am,omitempty"`
}

// Config holds the unified CLI configuration.
type Config struct {
	Current  string              `yaml:"current"  json:"current"`
	Contexts map[string]*Context `yaml:"contexts" json:"contexts"`
}

// ResolvedContext holds the fully resolved context after applying overrides.
type ResolvedContext struct {
	Name  string
	URL   string
	Token string
	Org   string
	Env   string
}

// Overrides holds flag-based overrides applied on top of the config context.
type Overrides struct {
	Context string
	Org     string
	EnvID   string
}

// NormalizeContextName lowercases the name and replaces spaces with hyphens
// so context names are safe to display in tables and easy to type back.
// "Local Master" -> "local-master".
func NormalizeContextName(name string) string {
	return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(name)), " ", "-")
}

// Path returns the full path to the config file (~/.gio/config.yaml).
func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine home directory: %w", err)
	}

	return filepath.Join(home, configDir, configFile), nil
}

// Load reads the config file from the default path.
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
			return &Config{Contexts: make(map[string]*Context)}, nil
		}

		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if cfg.Contexts == nil {
		cfg.Contexts = make(map[string]*Context)
	}

	return &cfg, nil
}

// SaveTo writes the config to the given path, creating the directory if needed.
func (c *Config) SaveTo(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Resolve returns the fully resolved context for the given product after applying overrides.
func (c *Config) Resolve(overrides Overrides, product string) (*ResolvedContext, error) {
	contextName := c.Current
	if overrides.Context != "" {
		contextName = NormalizeContextName(overrides.Context)
	}

	if contextName == "" {
		return nil, fmt.Errorf("no context configured")
	}

	ctx, ok := c.Contexts[contextName]
	if !ok {
		available := c.ContextNames()

		return nil, fmt.Errorf("context '%s' not found\nHint: available contexts: %s", contextName, strings.Join(available, ", "))
	}

	pc := ctx.productConfig(product)
	if pc == nil {
		return nil, fmt.Errorf("%s not configured for context '%s'\nHint: run 'gio login %s' to configure", strings.ToUpper(product), contextName, product)
	}

	resolved := &ResolvedContext{
		Name:  contextName,
		URL:   pc.URL,
		Token: pc.Token,
		Org:   withDefault(ctx.Org, DefaultOrg),
		Env:   withDefault(ctx.Env, DefaultEnv),
	}

	if overrides.Org != "" {
		resolved.Org = overrides.Org
	}

	if overrides.EnvID != "" {
		resolved.Env = overrides.EnvID
	}

	return resolved, nil
}

// EnsureContext returns the context with the given name, creating it if it doesn't exist.
func (c *Config) EnsureContext(name string) *Context {
	if ctx, ok := c.Contexts[name]; ok {
		return ctx
	}

	ctx := &Context{}
	c.Contexts[name] = ctx

	return ctx
}

// DeleteContext removes a context. Clears Current if deleting the active context.
func (c *Config) DeleteContext(name string) error {
	if _, ok := c.Contexts[name]; !ok {
		return fmt.Errorf("context '%s' not found", name)
	}

	delete(c.Contexts, name)

	if c.Current == name {
		c.Current = ""
	}

	return nil
}

// ContextNames returns sorted context names.
func (c *Config) ContextNames() []string {
	names := make([]string, 0, len(c.Contexts))
	for name := range c.Contexts {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

func (ctx *Context) productConfig(product string) *ProductConfig {
	switch product {
	case "apim":
		return ctx.APIM
	case "am":
		return ctx.AM
	default:
		return nil
	}
}

// SetProductConfig sets the product-specific config block on the context.
func (ctx *Context) SetProductConfig(product string, pc *ProductConfig) {
	switch product {
	case "apim":
		ctx.APIM = pc
	case "am":
		ctx.AM = pc
	}
}

func withDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}
