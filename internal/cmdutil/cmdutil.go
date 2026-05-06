// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmdutil

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

// Environment variable names for CLI overrides.
// Per-product: GIO_APIM_URL, GIO_APIM_TOKEN, GIO_AM_URL, GIO_AM_TOKEN.
// Shared: GIO_ORG, GIO_ENV (apply to whichever product is used).
const (
	EnvOrg = "GIO_ORG"
	EnvEnv = "GIO_ENV"
)

func productEnvURL(product string) string {
	return "GIO_" + strings.ToUpper(product) + "_URL"
}

func productEnvToken(product string) string {
	return "GIO_" + strings.ToUpper(product) + "_TOKEN"
}

// SetupConfig loads the unified config file. Idempotent - skips if already loaded.
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
// Priority: env vars > CLI flags > config file > defaults.
func ResolveProductContext(f *factory.Factory, product string) error {
	f.Product = product

	// Env vars can bypass config entirely (useful for CI/CD).
	envURL := os.Getenv(productEnvURL(product))
	envToken := os.Getenv(productEnvToken(product))

	if envURL != "" && envToken != "" {
		f.Resolved = &config.ResolvedContext{
			Name:  "env",
			URL:   envURL,
			Token: envToken,
			Org:   envOrDefault(EnvOrg, f.Overrides.Org, config.DefaultOrg),
			Env:   envOrDefault(EnvEnv, f.Overrides.EnvID, config.DefaultEnv),
		}

		f.Client = client.NewHTTPClient(client.HTTPClientConfig{
			BaseURL:  f.Resolved.URL,
			Token:    f.Resolved.Token,
			Debug:    f.Debug,
			DebugOut: f.IOStreams.Err,
		})

		return nil
	}

	resolved, err := f.Config.Resolve(f.Overrides, product)
	if err != nil {
		// Allow commands that don't need a context (login, context).
		f.ContextResolveErr = err

		return nil
	}

	// Env vars override individual fields even when using config file.
	if envOrg := os.Getenv(EnvOrg); envOrg != "" {
		resolved.Org = envOrg
	}

	if envEnv := os.Getenv(EnvEnv); envEnv != "" {
		resolved.Env = envEnv
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

func envOrDefault(envKey, flagValue, defaultValue string) string {
	if v := os.Getenv(envKey); v != "" {
		return v
	}

	if flagValue != "" {
		return flagValue
	}

	return defaultValue
}

// ContextHeader returns a formatted string showing the current context info for help display.
// Shows the actual resolved values (env vars > flags > config).
func ContextHeader(f *factory.Factory, product string) string {
	if f.Resolved == nil {
		return ""
	}

	r := f.Resolved

	source := r.Name
	if r.Name == "env" {
		source = fmt.Sprintf("env (%s, %s)", productEnvURL(product), productEnvToken(product))
	}

	return fmt.Sprintf("\n  Context:   %s\n  URL:       %s\n  Org:       %s\n  Env:       %s\n",
		source, r.URL, r.Org, r.Env)
}

// AddOutputFlags adds -o/--output, -q/--quiet, and --no-headers as persistent flags on a command.
// Call this on resource parent commands (api, plan, domain...) so that only
// commands producing output expose these flags.
func AddOutputFlags(cmd *cobra.Command, f *factory.Factory) {
	cmd.PersistentFlags().StringVarP(&f.OutputFormat, "output", "o", "table", "Output format: table, json, yaml, id")
	cmd.PersistentFlags().BoolVarP(&f.Quiet, "quiet", "q", false, "Suppress output on success (errors still go to stderr)")
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

// AMEnvPath builds an AM environment-scoped API path using factory context.
func AMEnvPath(f *factory.Factory, path string) string {
	return client.AMEnvPath(f.Resolved.Org, f.Resolved.Env, path)
}

// AMDomainPath builds an AM domain-scoped API path using factory context.
func AMDomainPath(f *factory.Factory, path string) string {
	return client.AMDomainPath(f.Resolved.Org, f.Resolved.Env, f.Resolved.Domain, path)
}

// AMDomainPathFor builds a domain-scoped path using an explicit domainID.
func AMDomainPathFor(f *factory.Factory, domainID, path string) string {
	return client.AMDomainPath(f.Resolved.Org, f.Resolved.Env, domainID, path)
}

// RequireAMContext returns an error if no AM context is configured.
func RequireAMContext(f *factory.Factory) error {
	if f.Resolved == nil {
		return fmt.Errorf("no AM context configured\nHint: run 'gio login am' to get started")
	}

	if f.Resolved.Type != "am" {
		return fmt.Errorf("current context '%s' is not an AM context (type: %s)\nHint: switch to an AM context with 'gio config use-context <am-context>'", f.Resolved.Name, f.Resolved.Type)
	}

	return nil
}

// RequireAMDomain returns an error if no AM domain is set in the context.
func RequireAMDomain(f *factory.Factory) error {
	if err := RequireAMContext(f); err != nil {
		return err
	}

	if f.Resolved.Domain == "" {
		return fmt.Errorf("no domain selected\nHint: run 'gio am set domain <id>' to select a domain")
	}

	return nil
}

// ReadJSONFile reads a JSON file, expands ${VAR} references from the environment, and returns the raw content.
func ReadJSONFile(path string) (json.RawMessage, error) {
	path = filepath.Clean(path)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read '%s': %w", path, err)
	}

	data, err = expandEnvVars(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read '%s': %w", path, err)
	}

	if !json.Valid(data) {
		return nil, fmt.Errorf("failed to read '%s': invalid JSON", path)
	}

	return data, nil
}

func expandEnvVars(data []byte) ([]byte, error) {
	for _, m := range envVarRE.FindAllSubmatch(data, -1) {
		name := string(m[1])
		if _, ok := os.LookupEnv(name); !ok {
			return nil, fmt.Errorf("environment variable '%s' is not defined", name)
		}
	}

	return envVarRE.ReplaceAllFunc(data, func(match []byte) []byte {
		return []byte(os.Getenv(string(match[2 : len(match)-1])))
	}), nil
}

// NewPrinter creates a Printer from the factory settings, validating the output format.
func NewPrinter(f *factory.Factory) (*printer.Printer, error) {
	switch f.OutputFormat {
	case printer.FormatTable, printer.FormatJSON, printer.FormatYAML, printer.FormatID:
		return printer.New(f.OutputFormat, f.IOStreams.Out, f.IOStreams.Err, f.Quiet, f.NoHeaders), nil
	default:
		return nil, fmt.Errorf("invalid output format %q\nHint: allowed values are table, json, yaml, id", f.OutputFormat)
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

		p.PrintHint("Showing %d-%d of %d (page %d/%d).%s",
			start, end, totalCount, page, pageCount, hint)
	} else {
		p.PrintHint("Showing %d-%d of %d (page %d/%d).",
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

var loginURLPathRE = regexp.MustCompile(`/management/organizations/([^/]+)(?:/environments/([^/]+))?`)

var envVarRE = regexp.MustCompile(`\$\{([^}]+)\}`)

// ParseLoginURL takes a raw management URL - either a bare base (https://host) or one
// including a /management/organizations/<org>[/environments/<env>[/...]] path - and
// returns the base URL plus the extracted org and env. When the path does not match
// the management pattern, org and env are returned empty and the base is the full
// scheme+host+port (path is discarded). Only the URL scheme prefix is validated.
func ParseLoginURL(raw string) (baseURL, org, env string, err error) {
	if err := ValidateURL(raw); err != nil {
		return "", "", "", err
	}

	u, perr := url.Parse(raw)
	if perr != nil {
		return "", "", "", fmt.Errorf("invalid URL %q: %w", raw, perr)
	}

	base := u.Scheme + "://" + u.Host

	if m := loginURLPathRE.FindStringSubmatch(u.Path); m != nil {
		org = m[1]
		env = m[2]
	}

	return base, org, env, nil
}

// ParseCurl extracts the URL and bearer token from a curl command string.
// Accepts shell-style quoting (single or double quotes) and ignores curl flags
// other than -H / --header. The first http(s) token is taken as the URL; the
// first -H value whose content matches "Authorization: Bearer <token>" (case
// insensitive on the keywords) yields the token. Returns errors when the URL or
// the Authorization header cannot be found.
func ParseCurl(cmd string) (rawURL, token string, err error) {
	tokens, err := shellSplit(cmd)
	if err != nil {
		return "", "", fmt.Errorf("curl command: %w", err)
	}

	rawURL, token = scanCurlTokens(tokens)

	if rawURL == "" {
		return "", "", fmt.Errorf("curl command: missing URL\nHint: use devtools \"Copy as cURL\" so the request URL is included")
	}

	if token == "" {
		return "", "", fmt.Errorf("curl command: missing Authorization: Bearer header\nHint: the curl must include -H 'Authorization: Bearer <token>'")
	}

	return rawURL, token, nil
}

func scanCurlTokens(tokens []string) (rawURL, token string) {
	for i := 0; i < len(tokens); i++ {
		tk := tokens[i]

		if rawURL == "" && (strings.HasPrefix(tk, "http://") || strings.HasPrefix(tk, "https://")) {
			rawURL = tk
			continue
		}

		switch {
		case tk == "-H" || tk == "--header":
			if i+1 < len(tokens) {
				if t := extractBearer(tokens[i+1]); t != "" && token == "" {
					token = t
				}

				i++
			}
		case strings.HasPrefix(tk, "--header="):
			if t := extractBearer(strings.TrimPrefix(tk, "--header=")); t != "" && token == "" {
				token = t
			}
		}
	}

	return rawURL, token
}

// extractBearer returns the token from a header value like
// "Authorization: Bearer xxx" (case insensitive on the keywords), or "" if the
// value does not match.
func extractBearer(headerValue string) string {
	h := strings.TrimSpace(headerValue)

	const authPrefix = "Authorization:"
	if len(h) < len(authPrefix) || !strings.EqualFold(h[:len(authPrefix)], authPrefix) {
		return ""
	}

	rest := strings.TrimSpace(h[len(authPrefix):])

	const bearerPrefix = "Bearer "
	if len(rest) < len(bearerPrefix) || !strings.EqualFold(rest[:len(bearerPrefix)], bearerPrefix) {
		return ""
	}

	return strings.TrimSpace(rest[len(bearerPrefix):])
}

// shellSplit tokenizes a command line the way a POSIX shell would, for the
// subset of syntax used by curl commands pasted from browser devtools: single
// quotes (literal), double quotes (with \" and \\ escapes), backslash escapes
// outside quotes, and whitespace as token separators.
func shellSplit(s string) ([]string, error) {
	l := &shellLexer{s: s}

	for i := 0; i < len(s); i++ {
		skip, err := l.step(i)
		if err != nil {
			return nil, err
		}

		i += skip
	}

	if l.inSingle || l.inDouble {
		return nil, fmt.Errorf("unterminated quoted string")
	}

	l.flush()

	return l.tokens, nil
}

type shellLexer struct {
	s          string
	tokens     []string
	cur        strings.Builder
	hasContent bool
	inSingle   bool
	inDouble   bool
}

func (l *shellLexer) step(i int) (int, error) {
	switch {
	case l.inSingle:
		l.stepSingle(l.s[i])

		return 0, nil
	case l.inDouble:
		return l.stepDouble(i), nil
	default:
		return l.stepUnquoted(i)
	}
}

func (l *shellLexer) stepSingle(c byte) {
	if c == '\'' {
		l.inSingle = false
	} else {
		l.cur.WriteByte(c)
	}
}

func (l *shellLexer) stepDouble(i int) int {
	c := l.s[i]

	if c == '"' {
		l.inDouble = false

		return 0
	}

	if c == '\\' && i+1 < len(l.s) {
		next := l.s[i+1]
		if next == '"' || next == '\\' || next == '$' || next == '`' {
			l.cur.WriteByte(next)

			return 1
		}
	}

	l.cur.WriteByte(c)

	return 0
}

func (l *shellLexer) stepUnquoted(i int) (int, error) {
	c := l.s[i]

	switch c {
	case '\'':
		l.inSingle = true
		l.hasContent = true
	case '"':
		l.inDouble = true
		l.hasContent = true
	case '\\':
		if i+1 >= len(l.s) {
			return 0, fmt.Errorf("trailing backslash")
		}

		l.cur.WriteByte(l.s[i+1])
		l.hasContent = true

		return 1, nil
	case ' ', '\t', '\n', '\r':
		l.flush()
	default:
		l.cur.WriteByte(c)
		l.hasContent = true
	}

	return 0, nil
}

func (l *shellLexer) flush() {
	if l.hasContent {
		l.tokens = append(l.tokens, l.cur.String())
		l.cur.Reset()

		l.hasContent = false
	}
}

// StringField extracts a value from a map[string]any and formats it as a string.
func StringField(item any, key string) string {
	m, ok := item.(map[string]any)
	if !ok {
		return ""
	}

	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}

	if f, ok := v.(float64); ok {
		return strconv.FormatInt(int64(f), 10)
	}

	return fmt.Sprintf("%v", v)
}

// ValidatePagination checks that page and perPage are positive.
func ValidatePagination(page, perPage int) error {
	if page < 1 {
		return fmt.Errorf("--page must be >= 1, got %d", page)
	}

	if perPage < 1 {
		return fmt.Errorf("--per-page must be >= 1, got %d", perPage)
	}

	return nil
}

// TimestampField extracts a timestamp from a map[string]any and formats it as ISO 8601.
// Handles both string timestamps (returned as-is) and numeric epoch milliseconds.
func TimestampField(item any, key string) string {
	m, ok := item.(map[string]any)
	if !ok {
		return ""
	}

	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}

	switch ts := v.(type) {
	case string:
		return ts
	case float64:
		return time.UnixMilli(int64(ts)).UTC().Format(time.RFC3339)
	}

	return fmt.Sprintf("%v", v)
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

// RequireNonEmpty rejects empty or whitespace-only values for required
// positional args and flag values. cobra.ExactArgs and MarkFlagRequired
// check presence, not content, so "" passes through otherwise.
func RequireNonEmpty(name, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s cannot be empty", name)
	}

	return nil
}

// PrintActionResult renders the outcome of an action verb (delete, start,
// stop, deploy...) that the server acknowledges without returning a body.
// In structured formats it emits a stable {id, status} envelope so scripts
// can parse. In -o id it prints the id alone. Otherwise it falls back to
// the human message.
func PrintActionResult(p *printer.Printer, id, status, humanMsg string) error {
	if printer.IsStructured(p.Format) || p.Format == printer.FormatID {
		return p.PrintDetail(map[string]string{"id": id, "status": status})
	}

	p.PrintMessage("%s", humanMsg)

	return nil
}
