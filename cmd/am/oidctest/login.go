package oidctest

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

type loginOpts struct {
	gateway, app, secret, username, password, scope string
	passwordStdin, secretStdin                      bool
}

func newLoginCmd(f *factory.Factory) *cobra.Command {
	o := &loginOpts{}

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Test Resource Owner Password Credentials (ROPC) flow",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runLogin(cmd.Context(), f, o)
		},
	}
	cmd.Flags().StringVar(&o.gateway, "gateway", "", "Gateway URL")
	cmd.Flags().StringVar(&o.app, "app", "", "Application client ID (required)")
	cmd.Flags().StringVar(&o.secret, "secret", "", "Application client secret (DEPRECATED: visible in process listings — prefer --secret-stdin)")
	cmd.Flags().BoolVar(&o.secretStdin, "secret-stdin", false, "Read application client secret from stdin")
	cmd.Flags().StringVar(&o.username, "username", "", "Username (required)")
	cmd.Flags().StringVar(&o.password, "password", "", "Password (DEPRECATED: visible in process listings — prefer --password-stdin)")
	cmd.Flags().BoolVar(&o.passwordStdin, "password-stdin", false, "Read password from stdin")
	cmd.Flags().StringVar(&o.scope, "scope", "", "Scopes to request (e.g. openid profile)")
	_ = cmd.MarkFlagRequired("app")
	_ = cmd.MarkFlagRequired("username")
	return cmd
}

func runLogin(ctx context.Context, f *factory.Factory, o *loginOpts) error {
	if err := cmdutil.RequireAMDomain(f); err != nil {
		return err
	}
	pw, err := cmdutil.ResolvePassword(o.password, o.passwordStdin, "Password: ", f.IOStreams.In, f.IOStreams.Err)
	if err != nil {
		return err
	}
	secret, err := resolveSecret(f, o)
	if err != nil {
		return err
	}
	tokenEndpoint, gw, discovery, err := resolveTokenEndpoint(ctx, f, o.gateway)
	if err != nil {
		return err
	}
	_ = gw
	tokenResp, err := requestRopcToken(ctx, tokenEndpoint, o.app, o.username, pw, secret, o.scope)
	if err != nil {
		return err
	}
	printTokenResult(f, tokenResp, discovery, o.app)
	return nil
}

func resolveSecret(f *factory.Factory, o *loginOpts) (string, error) {
	if !o.secretStdin {
		return o.secret, nil
	}
	return cmdutil.ResolvePassword(o.secret, true, "Client secret: ", f.IOStreams.In, f.IOStreams.Err)
}

func resolveTokenEndpoint(ctx context.Context, f *factory.Factory, gatewayFlag string) (tokenEndpoint, gw string, discovery map[string]interface{}, err error) {
	gw = gatewayURL(gatewayFlag, os.Getenv("AM_GATEWAY"), f.Resolved.URL)
	if err = validateGatewayURL(gw); err != nil {
		return "", "", nil, err
	}
	domainPath, err := fetchDomainPath(f)
	if err != nil {
		return "", "", nil, err
	}
	discovery, err = fetchDiscovery(ctx, gw, domainPath)
	if err != nil {
		return "", "", nil, err
	}
	tokenEndpoint, _ = discovery["token_endpoint"].(string)
	if tokenEndpoint == "" {
		return "", "", nil, fmt.Errorf("no token_endpoint in discovery")
	}
	if err = validateTokenEndpoint(tokenEndpoint, gw); err != nil {
		return "", "", nil, err
	}
	return tokenEndpoint, gw, discovery, nil
}

func requestRopcToken(ctx context.Context, tokenEndpoint, app, username, password, secret, scope string) (map[string]interface{}, error) {
	params := url.Values{
		"grant_type": {"password"},
		"username":   {username},
		"password":   {password},
	}
	if scope != "" {
		params.Set("scope", scope)
	}
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	if secret != "" {
		creds := base64.StdEncoding.EncodeToString([]byte(app + ":" + secret))
		headers["Authorization"] = "Basic " + creds
	} else {
		params.Set("client_id", app)
	}
	return httpPost(ctx, tokenEndpoint, params.Encode(), headers)
}

func fetchDiscovery(ctx context.Context, gw, domainPath string) (map[string]interface{}, error) {
	discoveryURL := fmt.Sprintf("%s%s/oidc/.well-known/openid-configuration", gw, domainPath)
	data, err := httpGet(ctx, discoveryURL, "")
	if err != nil {
		return nil, fmt.Errorf("OIDC discovery failed: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func httpPost(ctx context.Context, endpoint, body string, headers map[string]string) (map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	c := &http.Client{Timeout: 30 * time.Second}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	if resp.StatusCode >= 400 {
		var errResp map[string]interface{}
		_ = json.Unmarshal(respBody, &errResp)
		errStr, _ := errResp["error"].(string)
		desc, _ := errResp["error_description"].(string)
		return nil, fmt.Errorf("token request failed (HTTP %d): %s %s", resp.StatusCode, errStr, desc)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func printTokenResult(f *factory.Factory, tokenResp, discovery map[string]interface{}, clientID string) {
	out := f.IOStreams.Out
	printAccessTokenInfo(out, tokenResp)
	if idToken, ok := tokenResp["id_token"].(string); ok {
		printIDTokenInfo(out, idToken, discovery, clientID)
	}
}

func printAccessTokenInfo(out io.Writer, tokenResp map[string]interface{}) {
	accessToken, _ := tokenResp["access_token"].(string)
	fmt.Fprintf(out, "Access Token:  %s\n", truncateToken(accessToken, 40))
	fmt.Fprintf(out, "Token Type:    %v\n", tokenResp["token_type"])
	fmt.Fprintf(out, "Expires In:    %vs\n", tokenResp["expires_in"])
	fmt.Fprintf(out, "Scopes:        %v\n", tokenResp["scope"])
}

func printIDTokenInfo(out io.Writer, idToken string, discovery map[string]interface{}, clientID string) {
	header, payload, err := decodeJWT(idToken)
	if err != nil {
		fmt.Fprintf(out, "\nCould not decode ID token: %v\n", err)
		return
	}
	fmt.Fprintln(out, "\nID Token (decoded):")
	fmt.Fprintln(out, "  Header:")
	for k, v := range header {
		fmt.Fprintf(out, "    %s: %v\n", k, v)
	}
	fmt.Fprintln(out, "  Payload:")
	for k, v := range payload {
		fmt.Fprintf(out, "    %s: %v\n", k, v)
	}
	fmt.Fprintln(out, "  Claim sanity check (signature NOT verified — use jwks_uri to validate cryptographically):")
	validateIssuer(out, payload, discovery)
	validateAudience(out, payload, clientID)
}

func validateIssuer(out io.Writer, payload, discovery map[string]interface{}) {
	iss, ok := payload["iss"].(string)
	if !ok {
		return
	}
	discoveryIss, ok := discovery["issuer"].(string)
	if !ok {
		return
	}
	if iss == discoveryIss {
		fmt.Fprintln(out, "    [OK]   Issuer matches discovery")
	} else {
		fmt.Fprintf(out, "    [FAIL] Issuer mismatch: %s vs %s\n", iss, discoveryIss)
	}
}

func validateAudience(out io.Writer, payload map[string]interface{}, clientID string) {
	aud := payload["aud"]
	audMatches := false
	switch a := aud.(type) {
	case string:
		audMatches = a == clientID
	case []interface{}:
		for _, item := range a {
			if s, ok := item.(string); ok && s == clientID {
				audMatches = true
				break
			}
		}
	}
	if audMatches {
		fmt.Fprintln(out, "    [OK]   Audience matches client_id")
	} else {
		fmt.Fprintf(out, "    [FAIL] Audience mismatch: %v vs %s\n", aud, clientID)
	}
}
