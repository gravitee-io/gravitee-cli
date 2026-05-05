package am

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

type loginOptions struct {
	factory     *factory.Factory
	url         string
	token       string
	username    string
	password    string
	contextName string
	org         string
	envID       string
}

func newLoginCmd(f *factory.Factory) *cobra.Command {
	opts := &loginOptions{factory: f}
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with a Gravitee AM instance",
		Example: `  gio am login --url https://am.company.com --username admin --password admin
  gio am login --url https://am.company.com --token eyJhbG...`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return opts.run()
		},
	}
	cmd.Flags().StringVar(&opts.url, "url", "", "URL of the AM management API (required)")
	cmd.Flags().StringVar(&opts.token, "token", "", "Bearer token (skip username/password login)")
	cmd.Flags().StringVar(&opts.username, "username", "", "Username for authentication")
	cmd.Flags().StringVar(&opts.password, "password", "", "Password for authentication")
	cmd.Flags().StringVar(&opts.contextName, "context", "", "Context name (default: derived from URL)")
	cmd.Flags().StringVar(&opts.org, "org", config.DefaultOrg, "Organization ID")
	cmd.Flags().StringVar(&opts.envID, "env-id", config.DefaultEnv, "Environment ID")
	_ = cmd.MarkFlagRequired("url")

	return cmd
}

func (o *loginOptions) run() error {
	cfg := o.factory.Config
	if cfg == nil {
		cfg = &config.Config{Contexts: make(map[string]*config.Context)}
		o.factory.Config = cfg
	}

	token := o.token
	if token == "" {
		if o.username == "" || o.password == "" {
			return fmt.Errorf("either --token or both --username and --password are required")
		}

		var err error

		token, err = o.authenticate()
		if err != nil {
			return err
		}
	}

	contextName := o.contextName
	if contextName == "" {
		contextName = deriveContextName(o.url)
	}

	ctx := cfg.EnsureContext(contextName)
	ctx.Type = "am"
	ctx.Org = o.org
	ctx.Env = o.envID
	ctx.AM = &config.ProductConfig{
		URL:   o.url,
		Token: token,
	}

	cfg.Current = contextName

	if err := cfg.SaveTo(o.factory.ConfigPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Fprintf(o.factory.IOStreams.Out, "Context '%s' saved and set as current.\n", contextName)

	return nil
}

func (o *loginOptions) authenticate() (string, error) {
	authURL := strings.TrimRight(o.url, "/") + "/management/auth/token"
	credentials := base64.StdEncoding.EncodeToString([]byte(o.username + ":" + o.password))
	body := url.Values{
		"grant_type": {"password"},
		"username":   {o.username},
		"password":   {o.password},
	}

	req, err := http.NewRequest(http.MethodPost, authURL, strings.NewReader(body.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Authorization", "Basic "+credentials)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := &http.Client{Timeout: 30 * time.Second}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("login request failed: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read login response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("invalid username or password")
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("login failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse login response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("login succeeded but no access token in response")
	}

	return tokenResp.AccessToken, nil
}

func deriveContextName(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		return "am"
	}

	host := parsed.Hostname()
	host = strings.ReplaceAll(host, ".", "-")

	return host + "-am"
}
