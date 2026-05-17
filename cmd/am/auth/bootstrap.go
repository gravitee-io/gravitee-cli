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

// Package auth provides first-mile authentication helpers — notably
// `am auth bootstrap`, which mints a Personal Access Token from a username +
// password on a fresh AM stack where no PAT is available yet.
package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

const sessionCookieName = "Auth-Graviteeio-AM"

type bootstrapOptions struct {
	factory       *factory.Factory
	amURL         string
	username      string
	password      string
	passwordStdin bool
	tokenName     string
	org           string
	contextName   string
	save          bool
}

// bootstrapClient is the minimal interface bootstrap needs — defined so
// tests can stub HTTP without spinning up a server.
type bootstrapClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func newBootstrapCmd(f *factory.Factory) *cobra.Command {
	opts := &bootstrapOptions{factory: f}

	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Mint a Personal Access Token from username + password",
		Long: `Bootstrap a Personal Access Token (PAT) on a fresh AM stack where no
PAT is available yet. Performs a form login at /management/auth/login, captures
the session cookie, then POSTs to /organizations/{org}/users/{userId}/tokens to
mint a PAT.

Useful for local-stack first-time CLI setup (e.g. admin/adminadmin), avoiding
the UI click-through that's currently required before 'gio login am'.`,
		Example: `  gio am auth bootstrap --url http://localhost:8093 --user admin --password-stdin --save
  gio am auth bootstrap --url http://localhost:8093 --user admin --password adminadmin
  gio am auth bootstrap --user admin --password-stdin --token-name ci-token --save`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return opts.run(http.DefaultClient)
		},
	}

	cmd.Flags().StringVar(&opts.amURL, "url", "", "AM management URL (falls back to current context's AM URL)")
	cmd.Flags().StringVar(&opts.username, "user", "", "Username (required)")
	cmd.Flags().StringVar(&opts.password, "password", "", "Password (use --password-stdin to avoid leaks)")
	cmd.Flags().BoolVar(&opts.passwordStdin, "password-stdin", false, "Read password from stdin")
	cmd.Flags().StringVar(&opts.tokenName, "token-name", "gio-cli", "Name to give the minted PAT")
	cmd.Flags().StringVar(&opts.org, "org", config.DefaultOrg, "Organization ID")
	cmd.Flags().StringVar(&opts.contextName, "context", "", "Context name to update when --save is set (defaults to current context)")
	cmd.Flags().BoolVar(&opts.save, "save", false, "Write the minted PAT into ~/.gio/config.yaml for the chosen context")
	_ = cmd.MarkFlagRequired("user")

	return cmd
}

func (o *bootstrapOptions) run(httpClient bootstrapClient) error {
	if err := o.resolveURL(); err != nil {
		return err
	}

	pw, err := cmdutil.ResolvePassword(
		o.password,
		o.passwordStdin,
		fmt.Sprintf("Password for %s: ", o.username),
		o.factory.IOStreams.In,
		o.factory.IOStreams.Err,
	)
	if err != nil {
		return err
	}

	cookie, err := loginAndGetCookie(httpClient, o.amURL, o.username, pw)
	if err != nil {
		return err
	}

	userID, err := lookupCurrentUserID(httpClient, o.amURL, o.org, cookie)
	if err != nil {
		return err
	}

	tokenValue, tokenID, err := mintToken(httpClient, o.amURL, o.org, userID, cookie, o.tokenName)
	if err != nil {
		return err
	}

	out := o.factory.IOStreams.Out
	fmt.Fprintf(out, "Minted PAT %q (ID: %s).\n", o.tokenName, tokenID)
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Token value (store it now — it will not be shown again):")
	fmt.Fprintf(out, "  %s\n", tokenValue)

	if o.save {
		if err := o.saveTokenToConfig(tokenValue); err != nil {
			return err
		}
		fmt.Fprintln(out)
		fmt.Fprintf(out, "Saved to context %q in %s.\n", o.effectiveContext(), o.factory.ConfigPath)
	}

	return nil
}

func (o *bootstrapOptions) resolveURL() error {
	if o.amURL != "" {
		o.amURL = strings.TrimRight(o.amURL, "/")
		return nil
	}

	if o.factory.Resolved != nil && o.factory.Resolved.URL != "" {
		o.amURL = strings.TrimRight(o.factory.Resolved.URL, "/")
		return nil
	}

	return fmt.Errorf("no AM URL: pass --url or configure a context with 'gio login am'")
}

func loginAndGetCookie(httpClient bootstrapClient, amURL, username, password string) (string, error) {
	form := url.Values{}
	form.Set("username", username)
	form.Set("password", password)

	req, err := http.NewRequest(http.MethodPost, amURL+"/management/auth/login", strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("build login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return "", fmt.Errorf("login failed (status %d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	for _, c := range resp.Cookies() {
		if c.Name == sessionCookieName && c.Value != "" {
			return c.Value, nil
		}
	}

	return "", fmt.Errorf("login succeeded but no %s cookie returned", sessionCookieName)
}

func lookupCurrentUserID(httpClient bootstrapClient, amURL, org, cookie string) (string, error) {
	path := fmt.Sprintf("%s/management/organizations/%s/user", amURL, url.PathEscape(org))
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return "", fmt.Errorf("build current-user request: %w", err)
	}
	addSessionCookie(req, cookie)
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("current-user request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return "", fmt.Errorf("current-user lookup failed (status %d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var payload struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("parse current-user response: %w", err)
	}
	if payload.ID == "" {
		return "", fmt.Errorf("current-user response missing id")
	}
	return payload.ID, nil
}

func mintToken(httpClient bootstrapClient, amURL, org, userID, cookie, tokenName string) (value, id string, err error) {
	body, _ := json.Marshal(map[string]string{"name": tokenName})
	path := fmt.Sprintf("%s/management/organizations/%s/users/%s/tokens",
		amURL, url.PathEscape(org), url.PathEscape(userID))

	req, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return "", "", fmt.Errorf("build token request: %w", err)
	}
	addSessionCookie(req, cookie)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return "", "", fmt.Errorf("token mint failed (status %d): %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var payload struct {
		ID    string `json:"id"`
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", "", fmt.Errorf("parse token response: %w", err)
	}
	if payload.Token == "" {
		return "", "", fmt.Errorf("token response missing token value")
	}
	return payload.Token, payload.ID, nil
}

func addSessionCookie(req *http.Request, cookie string) {
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: cookie})
}

func (o *bootstrapOptions) effectiveContext() string {
	if o.contextName != "" {
		return o.contextName
	}
	if o.factory.Config != nil && o.factory.Config.Current != "" {
		return o.factory.Config.Current
	}
	return "default"
}

func (o *bootstrapOptions) saveTokenToConfig(token string) error {
	if o.factory.ConfigPath == "" {
		return fmt.Errorf("--save requires a config file path; none configured")
	}

	cfg := o.factory.Config
	if cfg == nil {
		cfg = &config.Config{Contexts: map[string]*config.Context{}}
		o.factory.Config = cfg
	}

	ctxName := o.effectiveContext()
	ctx := cfg.EnsureContext(ctxName)
	if cfg.Current == "" {
		cfg.Current = ctxName
	}
	if ctx.AM == nil {
		ctx.AM = &config.ProductConfig{}
	}
	ctx.AM.URL = o.amURL
	ctx.AM.Token = token
	if ctx.Org == "" {
		ctx.Org = o.org
	}

	return cfg.SaveTo(o.factory.ConfigPath)
}

