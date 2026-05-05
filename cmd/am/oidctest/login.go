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

func newLoginCmd(f *factory.Factory) *cobra.Command {
	var gatewayFlag, app, secret, username, password, scope string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Test Resource Owner Password Credentials (ROPC) flow",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			gw := gatewayURL(gatewayFlag, os.Getenv("AM_GATEWAY"), f.Resolved.URL)
			domainPath, err := fetchDomainPath(f)
			if err != nil {
				return err
			}
			discovery, err := fetchDiscovery(cmd.Context(), gw, domainPath)
			if err != nil {
				return err
			}
			tokenEndpoint, _ := discovery["token_endpoint"].(string)
			if tokenEndpoint == "" {
				return fmt.Errorf("no token_endpoint in discovery")
			}

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

			tokenResp, err := httpPost(cmd.Context(), tokenEndpoint, params.Encode(), headers)
			if err != nil {
				return err
			}
			printTokenResult(f, tokenResp, discovery, app)
			return nil
		},
	}
	cmd.Flags().StringVar(&gatewayFlag, "gateway", "", "Gateway URL")
	cmd.Flags().StringVar(&app, "app", "", "Application client ID (required)")
	cmd.Flags().StringVar(&secret, "secret", "", "Application client secret (omit for public clients)")
	cmd.Flags().StringVar(&username, "username", "", "Username (required)")
	cmd.Flags().StringVar(&password, "password", "", "Password (required)")
	cmd.Flags().StringVar(&scope, "scope", "", "Scopes to request (e.g. openid profile)")
	_ = cmd.MarkFlagRequired("app")
	_ = cmd.MarkFlagRequired("username")
	_ = cmd.MarkFlagRequired("password")
	return cmd
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
	accessToken, _ := tokenResp["access_token"].(string)
	fmt.Fprintf(out, "Access Token:  %s\n", truncateToken(accessToken, 40))
	fmt.Fprintf(out, "Token Type:    %v\n", tokenResp["token_type"])
	fmt.Fprintf(out, "Expires In:    %vs\n", tokenResp["expires_in"])
	fmt.Fprintf(out, "Scopes:        %v\n", tokenResp["scope"])

	if idToken, ok := tokenResp["id_token"].(string); ok {
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
		fmt.Fprintln(out, "  Validation:")
		if iss, ok := payload["iss"].(string); ok {
			if discoveryIss, ok := discovery["issuer"].(string); ok {
				if iss == discoveryIss {
					fmt.Fprintln(out, "    ✓ Issuer matches discovery")
				} else {
					fmt.Fprintf(out, "    ✗ Issuer mismatch: %s vs %s\n", iss, discoveryIss)
				}
			}
		}
		aud := payload["aud"]
		audMatches := false
		switch a := aud.(type) {
		case string:
			audMatches = a == clientID
		case []interface{}:
			for _, item := range a {
				if s, ok := item.(string); ok && s == clientID {
					audMatches = true
				}
			}
		}
		if audMatches {
			fmt.Fprintln(out, "    ✓ Audience matches client_id")
		} else {
			fmt.Fprintf(out, "    ✗ Audience mismatch: %v vs %s\n", aud, clientID)
		}
	}
}
