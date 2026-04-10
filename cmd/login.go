package cmd

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

const loginLongFmt = `Configure credentials for a Gravitee %s instance.

Three modes:
  Interactive (for humans): run without arguments and paste the curl from
      Gravitee's UI. URL, token, org, env are parsed automatically.
  Flags (for agents/scripts): --url, --token, --context, --org, --env.
  Env vars (for CI): GIO_%s_URL, GIO_%s_TOKEN, GIO_ORG, GIO_ENV. These
      bypass the config file entirely, no 'gio login' needed.`

type loginProductOptions struct {
	factory     *factory.Factory
	product     string
	url         string
	token       string
	contextName string
	org         string
	envID       string
}

func newLoginCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [apim|am]",
		Short: "Configure credentials for a Gravitee product",
		Long: `Configure credentials for a Gravitee product (APIM or AM).

Three modes:
  Interactive (for humans): 'gio login apim' or 'gio login am' without args.
      Paste the curl from Gravitee's UI, URL/token/org/env are parsed automatically.
  Flags (for agents/scripts): --url, --token, --context, --org, --env.
  Env vars (for CI): GIO_APIM_URL, GIO_APIM_TOKEN, GIO_AM_URL, GIO_AM_TOKEN,
      GIO_ORG, GIO_ENV. These bypass the config file entirely.`,
		Example: `  gio login apim    # interactive for APIM
  gio login am      # interactive for AM`,
		Args: cobra.NoArgs,
	}

	cmd.AddCommand(newLoginProductCmd(f, "apim"))
	cmd.AddCommand(newLoginProductCmd(f, "am"))

	return cmd
}

func newLoginProductCmd(f *factory.Factory, product string) *cobra.Command {
	opts := &loginProductOptions{factory: f, product: product}
	productUpper := strings.ToUpper(product)

	cmd := &cobra.Command{
		Use:   product,
		Short: fmt.Sprintf("Configure credentials for a Gravitee %s instance", productUpper),
		Long:  fmt.Sprintf(loginLongFmt, productUpper, productUpper, productUpper),
		Example: fmt.Sprintf(`  gio login %s
      Interactive: paste the curl command from Gravitee's UI.

  gio login %s --url https://%s.example.com --token gioat_abc123 --context prod
      Non-interactive with flags (for CI / scripts).`,
			product, product, product),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return opts.run(cmd)
		},
	}

	cmd.Flags().StringVar(&opts.url, "url", "", "URL of the Gravitee control plane")
	cmd.Flags().StringVar(&opts.token, "token", "", "Personal Access Token")
	cmd.Flags().StringVar(&opts.contextName, "context", "default", "Context name")
	cmd.Flags().StringVar(&opts.org, "org", config.DefaultOrg, "Organization ID")
	cmd.Flags().StringVar(&opts.envID, "env", config.DefaultEnv, "Environment ID")

	return cmd
}

func (o *loginProductOptions) run(cmd *cobra.Command) error {
	o.contextName = config.NormalizeContextName(o.contextName)

	if o.url == "" && o.token == "" {
		return runInteractiveLogin(o.factory, o.product)
	}

	if o.url == "" {
		return fmt.Errorf("--url is required\nHint: run 'gio login %s' without flags for interactive mode", o.product)
	}

	if o.token == "" {
		return fmt.Errorf("--token is required\nHint: run 'gio login %s' without flags for interactive mode", o.product)
	}

	baseURL, org, env, err := cmdutil.ParseLoginURL(o.url)
	if err != nil {
		return err
	}

	o.url = baseURL

	orgFromURL := org != ""
	envFromURL := env != ""

	if orgFromURL && !cmd.Flags().Changed("org") {
		o.org = org
	}

	if envFromURL && !cmd.Flags().Changed("env") {
		o.envID = env
	}

	return o.save(
		orgFromURL || cmd.Flags().Changed("org"),
		envFromURL || cmd.Flags().Changed("env"),
	)
}

func (o *loginProductOptions) save(setOrg, setEnv bool) error {
	if err := cmdutil.SetupConfig(o.factory); err != nil {
		return err
	}

	cfg := o.factory.Config
	ctx := cfg.EnsureContext(o.contextName)

	if setOrg {
		ctx.Org = o.org
	}

	if setEnv {
		ctx.Env = o.envID
	}

	ctx.SetProductConfig(o.product, &config.ProductConfig{
		URL:   o.url,
		Token: o.token,
	})

	cfg.Current = o.contextName

	if err := cfg.SaveTo(o.factory.ConfigPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Fprintf(o.factory.IOStreams.Out, "Context '%s' saved and set as current (%s configured).\n", o.contextName, strings.ToUpper(o.product))

	return nil
}

// runInteractiveLogin drives the interactive flow for a given product:
// context name (default = current) -> URL/curl -> token -> org -> env.
// Prompts show the existing context's org/env as defaults so Enter preserves them.
// When a curl is pasted, token/org/env are extracted and the remaining prompts skipped.
func runInteractiveLogin(f *factory.Factory, product string) error {
	if err := cmdutil.SetupConfig(f); err != nil {
		return err
	}

	p := newPrompter(f.IOStreams.In, f.IOStreams.Out)
	cfg := f.Config

	contextName := p.promptContext(cfg)
	ctx := cfg.EnsureContext(contextName)

	baseURL, token, org, env, err := p.promptURLOrCurl()
	if err != nil {
		return err
	}

	if token == "" {
		token, err = p.promptToken()
		if err != nil {
			return err
		}
	}

	if org == "" {
		org = p.promptOrg(defaultOr(ctx.Org, config.DefaultOrg))
	}

	if env == "" {
		env = p.promptEnv(defaultOr(ctx.Env, config.DefaultEnv))
	}

	ctx.Org = org
	ctx.Env = env
	ctx.SetProductConfig(product, &config.ProductConfig{URL: baseURL, Token: token})
	cfg.Current = contextName

	if err := cfg.SaveTo(f.ConfigPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Fprintf(p.out, "Context '%s' saved and set as current (%s configured).\n", contextName, strings.ToUpper(product))

	return nil
}

func defaultOr(value, fallback string) string {
	if value != "" {
		return value
	}

	return fallback
}

// prompter wraps stdin in a single bufio.Reader so all prompts read through
// the same buffer. Multiple readers on the same underlying stream would cause
// the second reader to miss bytes already buffered by the first.
type prompter struct {
	out    io.Writer
	reader *bufio.Reader
}

func newPrompter(in io.Reader, out io.Writer) *prompter {
	return &prompter{out: out, reader: bufio.NewReader(in)}
}

// readLine reads one line from stdin, trimmed. Returns "" on blank input,
// EOF, or read error.
func (p *prompter) readLine() string {
	line, _ := p.reader.ReadString('\n')

	return strings.TrimSpace(line)
}

// promptURLOrCurl asks for either a bare URL or a full "curl ..." command.
// When a curl is detected, token / org / env are extracted from it; otherwise
// the returned token is empty and the caller falls back to promptToken.
func (p *prompter) promptURLOrCurl() (baseURL, token, org, env string, err error) {
	fmt.Fprint(p.out, "URL (or paste full curl command): ")

	input := p.readLine()

	if strings.HasPrefix(input, "curl ") || strings.HasPrefix(input, "curl\t") {
		rawURL, tok, perr := cmdutil.ParseCurl(input)
		if perr != nil {
			return "", "", "", "", perr
		}

		base, o, e, perr := cmdutil.ParseLoginURL(rawURL)
		if perr != nil {
			return "", "", "", "", perr
		}

		return base, tok, o, e, nil
	}

	base, o, e, perr := cmdutil.ParseLoginURL(input)
	if perr != nil {
		return "", "", "", "", perr
	}

	return base, "", o, e, nil
}

func (p *prompter) promptToken() (string, error) {
	fmt.Fprint(p.out, "Token: ")

	token := p.readLine()
	if token == "" {
		return "", fmt.Errorf("token is required")
	}

	return token, nil
}

// promptOrg asks for the organization ID, showing the given default. Enter
// returns the default (which is the existing context's value when reusing it,
// or DefaultOrg for a new context) so the prompt is honest about what applies.
func (p *prompter) promptOrg(defaultValue string) string {
	fmt.Fprintf(p.out, "Organization ID (%s): ", defaultValue)

	if s := p.readLine(); s != "" {
		return s
	}

	return defaultValue
}

// promptEnv asks for the environment ID, showing the given default. Enter
// returns the default (existing context's value when reusing, or DefaultEnv
// for a new context).
func (p *prompter) promptEnv(defaultValue string) string {
	fmt.Fprintf(p.out, "Environment ID (%s): ", defaultValue)

	if s := p.readLine(); s != "" {
		return s
	}

	return defaultValue
}

func (p *prompter) promptContext(cfg *config.Config) string {
	names := cfg.ContextNames()

	defaultName := cfg.Current
	if defaultName == "" {
		defaultName = "default"
	}

	if len(names) > 0 {
		fmt.Fprintln(p.out, "Existing contexts (reuse a name to update, or type a new name to create):")

		for _, name := range names {
			marker := "  "
			if name == cfg.Current {
				marker = "* "
			}

			fmt.Fprintf(p.out, "  %s%s\n", marker, name)
		}
	}

	fmt.Fprintf(p.out, "Context name (%s): ", defaultName)

	if name := p.readLine(); name != "" {
		return config.NormalizeContextName(name)
	}

	return defaultName
}
