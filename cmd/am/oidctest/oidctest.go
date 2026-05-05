package oidctest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewTestCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test-oidc",
		Short: "OIDC testing utilities",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(newDiscoverCmd(f))
	cmd.AddCommand(newLoginCmd(f))
	cmd.AddCommand(newClientCredsCmd(f))
	return cmd
}

func deriveGatewayURL(mgmtURL string) string {
	parsed, err := url.Parse(mgmtURL)
	if err != nil {
		return "http://localhost:8092"
	}
	parsed.Host = parsed.Hostname() + ":8092"
	return strings.TrimRight(parsed.String(), "/")
}

func gatewayURL(flag, envVar, mgmtURL string) string {
	if flag != "" {
		return strings.TrimRight(flag, "/")
	}
	if envVar != "" {
		return strings.TrimRight(envVar, "/")
	}
	return deriveGatewayURL(mgmtURL)
}

func decodeJWT(token string) (header, payload map[string]interface{}, err error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, nil, fmt.Errorf("invalid JWT: expected 3 parts, got %d", len(parts))
	}
	decode := func(s string) (map[string]interface{}, error) {
		b, decErr := base64.RawURLEncoding.DecodeString(s)
		if decErr != nil {
			return nil, decErr
		}
		var m map[string]interface{}
		if decErr := json.Unmarshal(b, &m); decErr != nil {
			return nil, decErr
		}
		return m, nil
	}
	header, err = decode(parts[0])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode header: %w", err)
	}
	payload, err = decode(parts[1])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode payload: %w", err)
	}
	return header, payload, nil
}

func truncateToken(token string, maxLen int) string {
	if len(token) <= maxLen {
		return token
	}
	return token[:maxLen] + "...(truncated)"
}
