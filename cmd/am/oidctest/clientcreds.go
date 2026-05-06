package oidctest

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newClientCredsCmd(f *factory.Factory) *cobra.Command {
	var gatewayFlag, app, secret, scope string
	var secretStdin bool

	cmd := &cobra.Command{
		Use:   "client-credentials",
		Short: "Test client_credentials grant flow",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			s, err := cmdutil.ResolvePassword(secret, secretStdin, "Client secret: ", f.IOStreams.In, f.IOStreams.Err)
			if err != nil {
				return err
			}
			gw := gatewayURL(gatewayFlag, os.Getenv("AM_GATEWAY"), f.Resolved.URL)
			if vErr := validateGatewayURL(gw); vErr != nil {
				return vErr
			}
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
			if vErr := validateTokenEndpoint(tokenEndpoint, gw); vErr != nil {
				return vErr
			}

			params := url.Values{"grant_type": {"client_credentials"}}
			if scope != "" {
				params.Set("scope", scope)
			}
			creds := base64.StdEncoding.EncodeToString([]byte(app + ":" + s))
			headers := map[string]string{
				"Content-Type":  "application/x-www-form-urlencoded",
				"Authorization": "Basic " + creds,
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
	cmd.Flags().StringVar(&secret, "secret", "", "Application client secret (DEPRECATED: visible in process listings — prefer --secret-stdin)")
	cmd.Flags().BoolVar(&secretStdin, "secret-stdin", false, "Read application client secret from stdin")
	cmd.Flags().StringVar(&scope, "scope", "", "Scopes to request")
	_ = cmd.MarkFlagRequired("app")
	return cmd
}
