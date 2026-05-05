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

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
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
