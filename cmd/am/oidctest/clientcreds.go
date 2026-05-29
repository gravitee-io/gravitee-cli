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

package oidctest

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
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
