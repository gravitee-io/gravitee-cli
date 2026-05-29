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

package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
	"gravitee.io/gctl/internal/printer"
)

func newSettingsCmd(f *factory.Factory) *cobra.Command {
	var (
		grantTypes           string
		responseTypes        string
		redirectURIs         string
		postLogoutURIs       string
		tokenLifetime        int
		refreshTokenLifetime int
		idTokenLifetime      int
	)

	cmd := &cobra.Command{
		Use:   "settings <appId>",
		Short: "View or update OAuth2 settings for an application",
		Example: `  gctl am app settings my-app-id
  gctl am app settings my-app-id --grant-types "authorization_code,refresh_token"
  gctl am app settings my-app-id --redirect-uris "https://myapp.com/callback"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}

			isUpdate := grantTypes != "" || responseTypes != "" || redirectURIs != "" ||
				postLogoutURIs != "" ||
				cmd.Flags().Changed("token-lifetime") ||
				cmd.Flags().Changed("refresh-token-lifetime") ||
				cmd.Flags().Changed("id-token-lifetime")

			if isUpdate {
				return runSettingsUpdate(f, args[0], grantTypes, responseTypes, redirectURIs, postLogoutURIs, tokenLifetime, refreshTokenLifetime, idTokenLifetime, cmd)
			}

			return runSettingsView(f, args[0])
		},
	}

	cmd.Flags().StringVar(&grantTypes, "grant-types", "", "Comma-separated list of grant types")
	cmd.Flags().StringVar(&responseTypes, "response-types", "", "Comma-separated list of response types")
	cmd.Flags().StringVar(&redirectURIs, "redirect-uris", "", "Comma-separated list of redirect URIs")
	cmd.Flags().StringVar(&postLogoutURIs, "post-logout-uris", "", "Comma-separated list of post logout URIs")
	cmd.Flags().IntVar(&tokenLifetime, "token-lifetime", 0, "Access token lifetime in seconds")
	cmd.Flags().IntVar(&refreshTokenLifetime, "refresh-token-lifetime", 0, "Refresh token lifetime in seconds")
	cmd.Flags().IntVar(&idTokenLifetime, "id-token-lifetime", 0, "ID token lifetime in seconds")

	return cmd
}

func runSettingsView(f *factory.Factory, appID string) error {
	path := cmdutil.AMDomainPath(f, fmt.Sprintf("applications/%s", appID))

	data, err := f.Client.Get(path)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	settings, _ := m["settings"].(map[string]interface{})
	if settings == nil {
		p.PrintMessage("No settings found for application '%s'.", appID)
		return nil
	}

	oauth, _ := settings["oauth"].(map[string]interface{})
	if oauth == nil {
		p.PrintMessage("No OAuth2 settings found for application '%s'.", appID)
		return nil
	}

	printOAuthField(p, "Client ID", oauth, "clientId")
	printOAuthField(p, "Client Secret", oauth, "clientSecret")
	printOAuthSliceField(p, "Grant Types", oauth, "grantTypes")
	printOAuthSliceField(p, "Response Types", oauth, "responseTypes")
	printOAuthSliceField(p, "Redirect URIs", oauth, "redirectUris")
	printOAuthSliceField(p, "Post Logout URIs", oauth, "postLogoutRedirectUris")
	printOAuthIntField(p, "Token Lifetime", oauth, "accessTokenValiditySeconds")
	printOAuthIntField(p, "Refresh Token", oauth, "refreshTokenValiditySeconds")
	printOAuthIntField(p, "ID Token", oauth, "idTokenValiditySeconds")

	return nil
}

func printOAuthField(p *printer.Printer, label string, oauth map[string]interface{}, key string) {
	if v, ok := oauth[key].(string); ok && v != "" {
		p.PrintMessage("%-20s%s", label+":", v)
	}
}

func printOAuthSliceField(p *printer.Printer, label string, oauth map[string]interface{}, key string) {
	if v, ok := oauth[key].([]interface{}); ok && len(v) > 0 {
		parts := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				parts = append(parts, s)
			}
		}
		p.PrintMessage("%-20s%s", label+":", strings.Join(parts, ", "))
	}
}

func printOAuthIntField(p *printer.Printer, label string, oauth map[string]interface{}, key string) {
	if v, ok := oauth[key].(float64); ok {
		p.PrintMessage("%-20s%s", label+":", strconv.Itoa(int(v)))
	}
}

func runSettingsUpdate(f *factory.Factory, appID, grantTypes, responseTypes, redirectURIs, postLogoutURIs string, tokenLifetime, refreshTokenLifetime, idTokenLifetime int, cmd *cobra.Command) error {
	if grantTypes == "" && responseTypes == "" && redirectURIs == "" && postLogoutURIs == "" &&
		!cmd.Flags().Changed("token-lifetime") &&
		!cmd.Flags().Changed("refresh-token-lifetime") &&
		!cmd.Flags().Changed("id-token-lifetime") {
		return errors.New("no settings specified")
	}

	oauth := map[string]interface{}{}

	if grantTypes != "" {
		oauth["grantTypes"] = splitTrimmed(grantTypes)
	}
	if responseTypes != "" {
		oauth["responseTypes"] = splitTrimmed(responseTypes)
	}
	if redirectURIs != "" {
		oauth["redirectUris"] = splitTrimmed(redirectURIs)
	}
	if postLogoutURIs != "" {
		oauth["postLogoutRedirectUris"] = splitTrimmed(postLogoutURIs)
	}
	if cmd.Flags().Changed("token-lifetime") {
		oauth["accessTokenValiditySeconds"] = tokenLifetime
	}
	if cmd.Flags().Changed("refresh-token-lifetime") {
		oauth["refreshTokenValiditySeconds"] = refreshTokenLifetime
	}
	if cmd.Flags().Changed("id-token-lifetime") {
		oauth["idTokenValiditySeconds"] = idTokenLifetime
	}

	payload := map[string]interface{}{
		"settings": map[string]interface{}{
			"oauth": oauth,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to build request body: %w", err)
	}

	path := cmdutil.AMDomainPath(f, fmt.Sprintf("applications/%s", appID))

	data, err := f.Client.Patch(path, body)
	if err != nil {
		return fmt.Errorf("application settings update failed: %w", err)
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}

	return printAppDetail(p, data)
}

func splitTrimmed(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			result = append(result, t)
		}
	}
	return result
}
