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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

func newDiscoverCmd(f *factory.Factory) *cobra.Command {
	var gatewayFlag string

	return &cobra.Command{
		Use:   "discover",
		Short: "Fetch and display the OIDC discovery document",
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
			discoveryURL := fmt.Sprintf("%s%s/oidc/.well-known/openid-configuration", gw, domainPath)
			data, err := httpGet(cmd.Context(), discoveryURL, "")
			if err != nil {
				return fmt.Errorf("OIDC discovery failed: %w", err)
			}
			var discovery map[string]interface{}
			if err := json.Unmarshal(data, &discovery); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "OIDC Discovery for %s%s\n\n", gw, domainPath)
			keys := make([]string, 0, len(discovery))
			for k := range discovery {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, key := range keys {
				val := discovery[key]
				switch v := val.(type) {
				case []interface{}:
					fmt.Fprintf(f.IOStreams.Out, "%s:\n", key)
					for _, item := range v {
						fmt.Fprintf(f.IOStreams.Out, "  - %v\n", item)
					}
				default:
					fmt.Fprintf(f.IOStreams.Out, "%s: %v\n", key, v)
				}
			}
			return nil
		},
	}
}

func fetchDomainPath(f *factory.Factory) (string, error) {
	data, err := f.Client.Get(cmdutil.AMEnvPath(f, "domains/"+f.Resolved.Domain))
	if err != nil {
		return "", err
	}
	var domain map[string]interface{}
	if err := json.Unmarshal(data, &domain); err != nil {
		return "", err
	}
	if path, ok := domain["path"].(string); ok && path != "" {
		return path, nil
	}
	if hrid, ok := domain["hrid"].(string); ok && hrid != "" {
		return "/" + hrid, nil
	}
	return "/" + f.Resolved.Domain, nil
}

func httpGet(ctx context.Context, url, bearerToken string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}
	c := &http.Client{Timeout: 30 * time.Second}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}
