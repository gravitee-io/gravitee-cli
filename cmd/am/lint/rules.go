package lint

import (
	"fmt"
	"strings"
	"time"
)

type LintFinding struct {
	Rule     string
	Severity string // "critical" or "warning"
	Resource string
	Message  string
}

func ruleImplicitGrant(apps []map[string]interface{}) []LintFinding {
	var out []LintFinding
	for _, app := range apps {
		for _, g := range oauthGrantTypes(app) {
			if g == "implicit" {
				out = append(out, LintFinding{
					Rule: "implicit-grant", Severity: "critical",
					Resource: appName(app),
					Message:  "Uses implicit grant type (deprecated, insecure)",
				})
				break
			}
		}
	}
	return out
}

func ruleNoPkce(apps []map[string]interface{}) []LintFinding {
	var out []LintFinding
	for _, app := range apps {
		hasAuthCode := false
		for _, g := range oauthGrantTypes(app) {
			if g == "authorization_code" {
				hasAuthCode = true
			}
		}
		if !hasAuthCode {
			continue
		}
		oauth := oauthSettings(app)
		forcePkce, _ := oauth["forcePKCE"].(bool)
		requirePkce, _ := oauth["requirePKCE"].(bool)
		if !forcePkce && !requirePkce {
			out = append(out, LintFinding{
				Rule: "no-pkce", Severity: "warning",
				Resource: appName(app),
				Message:  "Authorization code flow without PKCE enforcement",
			})
		}
	}
	return out
}

func ruleLongTokenLifetime(apps []map[string]interface{}) []LintFinding {
	var out []LintFinding
	for _, app := range apps {
		oauth := oauthSettings(app)
		if v, ok := oauth["accessTokenValiditySeconds"].(float64); ok && v > 3600 {
			out = append(out, LintFinding{
				Rule: "long-token-lifetime", Severity: "warning",
				Resource: appName(app),
				Message:  fmt.Sprintf("Access token lifetime %.0fs exceeds 1 hour", v),
			})
		}
	}
	return out
}

func ruleLongRefreshLifetime(apps []map[string]interface{}) []LintFinding {
	var out []LintFinding
	for _, app := range apps {
		oauth := oauthSettings(app)
		if v, ok := oauth["refreshTokenValiditySeconds"].(float64); ok && v > 2592000 {
			out = append(out, LintFinding{
				Rule: "long-refresh-lifetime", Severity: "warning",
				Resource: appName(app),
				Message:  fmt.Sprintf("Refresh token lifetime %.0fs exceeds 30 days", v),
			})
		}
	}
	return out
}

func ruleNoIdp(apps []map[string]interface{}) []LintFinding {
	var out []LintFinding
	for _, app := range apps {
		idps, _ := app["identityProviders"].([]interface{})
		if len(idps) == 0 {
			idps2, _ := app["identities"].([]interface{})
			if len(idps2) == 0 {
				out = append(out, LintFinding{
					Rule: "no-idp", Severity: "critical",
					Resource: appName(app),
					Message:  "No identity providers assigned",
				})
			}
		}
	}
	return out
}

func ruleLocalhostRedirect(apps []map[string]interface{}) []LintFinding {
	var out []LintFinding
	for _, app := range apps {
		for _, uri := range redirectUris(app) {
			if strings.Contains(uri, "localhost") || strings.Contains(uri, "127.0.0.1") {
				out = append(out, LintFinding{
					Rule: "localhost-redirect", Severity: "warning",
					Resource: appName(app),
					Message:  fmt.Sprintf("Redirect URI contains localhost: %s", uri),
				})
				break
			}
		}
	}
	return out
}

func ruleHttpRedirect(apps []map[string]interface{}) []LintFinding {
	var out []LintFinding
	for _, app := range apps {
		for _, uri := range redirectUris(app) {
			if strings.HasPrefix(uri, "http://") && !strings.Contains(uri, "localhost") {
				out = append(out, LintFinding{
					Rule: "http-redirect", Severity: "warning",
					Resource: appName(app),
					Message:  fmt.Sprintf("Non-secure redirect URI: %s", uri),
				})
				break
			}
		}
	}
	return out
}

func ruleWildcardRedirect(apps []map[string]interface{}) []LintFinding {
	var out []LintFinding
	for _, app := range apps {
		for _, uri := range redirectUris(app) {
			if strings.Contains(uri, "*") {
				out = append(out, LintFinding{
					Rule: "wildcard-redirect", Severity: "critical",
					Resource: appName(app),
					Message:  fmt.Sprintf("Wildcard redirect URI: %s", uri),
				})
				break
			}
		}
	}
	return out
}

func ruleAppDisabled(apps []map[string]interface{}) []LintFinding {
	var out []LintFinding
	for _, app := range apps {
		enabled, _ := app["enabled"].(bool)
		if !enabled {
			out = append(out, LintFinding{
				Rule: "app-disabled", Severity: "warning",
				Resource: appName(app),
				Message:  "Application is disabled",
			})
		}
	}
	return out
}

func ruleNoFactors(apps []map[string]interface{}, factors []map[string]interface{}) []LintFinding {
	if len(factors) == 0 {
		return nil
	}
	var out []LintFinding
	for _, app := range apps {
		appFactors, _ := app["factors"].([]interface{})
		if len(appFactors) == 0 {
			out = append(out, LintFinding{
				Rule: "no-factors", Severity: "warning",
				Resource: appName(app),
				Message:  "Domain has MFA factors but app has none assigned",
			})
		}
	}
	return out
}

func ruleCertExpiry(certs []map[string]interface{}) []LintFinding {
	var out []LintFinding
	threshold := time.Now().Add(30 * 24 * time.Hour).UnixMilli()
	for _, cert := range certs {
		expiresAt, ok := cert["expiresAt"].(float64)
		if !ok {
			continue
		}
		if int64(expiresAt) <= threshold {
			name, _ := cert["name"].(string)
			out = append(out, LintFinding{
				Rule: "cert-expiry", Severity: "critical",
				Resource: name,
				Message:  fmt.Sprintf("Certificate expires within 30 days (expiresAt: %d)", int64(expiresAt)),
			})
		}
	}
	return out
}

func ruleUnusedScope(apps []map[string]interface{}, scopes []map[string]interface{}) []LintFinding {
	used := make(map[string]bool)
	for _, app := range apps {
		oauth := oauthSettings(app)
		if scopeSettings, ok := oauth["scopeSettings"].([]interface{}); ok {
			for _, s := range scopeSettings {
				if sm, ok := s.(map[string]interface{}); ok {
					if key, ok := sm["scope"].(string); ok {
						used[key] = true
					}
				}
			}
		}
	}
	var out []LintFinding
	for _, scope := range scopes {
		key, _ := scope["key"].(string)
		if key != "" && !used[key] {
			out = append(out, LintFinding{
				Rule: "unused-scope", Severity: "warning",
				Resource: key,
				Message:  fmt.Sprintf("Scope %q is not used by any application", key),
			})
		}
	}
	return out
}

func rulePasswordGrantNoMfa(apps []map[string]interface{}, factors []map[string]interface{}) []LintFinding {
	if len(factors) == 0 {
		return nil
	}
	var out []LintFinding
	for _, app := range apps {
		hasPassword := false
		for _, g := range oauthGrantTypes(app) {
			if g == "password" {
				hasPassword = true
			}
		}
		if !hasPassword {
			continue
		}
		appFactors, _ := app["factors"].([]interface{})
		if len(appFactors) == 0 {
			out = append(out, LintFinding{
				Rule: "password-grant-no-mfa", Severity: "warning",
				Resource: appName(app),
				Message:  "Password grant without MFA factor",
			})
		}
	}
	return out
}

func ruleEmptyDomain(apps []map[string]interface{}) []LintFinding {
	if len(apps) == 0 {
		return []LintFinding{{
			Rule: "empty-domain", Severity: "warning",
			Resource: "domain",
			Message:  "Domain has no applications",
		}}
	}
	return nil
}

func calculateScore(findings []LintFinding) int {
	criticals, warnings := 0, 0
	for _, f := range findings {
		switch f.Severity {
		case "critical":
			criticals++
		case "warning":
			warnings++
		}
	}
	score := 10 - 2*criticals - warnings
	if score < 0 {
		return 0
	}
	return score
}

// Helpers

func appName(app map[string]interface{}) string {
	if name, ok := app["name"].(string); ok {
		return name
	}
	return "unknown"
}

func oauthSettings(app map[string]interface{}) map[string]interface{} {
	settings, _ := app["settings"].(map[string]interface{})
	oauth, _ := settings["oauth"].(map[string]interface{})
	return oauth
}

func oauthGrantTypes(app map[string]interface{}) []string {
	oauth := oauthSettings(app)
	raw, _ := oauth["grantTypes"].([]interface{})
	result := make([]string, 0, len(raw))
	for _, g := range raw {
		if s, ok := g.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

func redirectUris(app map[string]interface{}) []string {
	oauth := oauthSettings(app)
	raw, _ := oauth["redirectUris"].([]interface{})
	result := make([]string, 0, len(raw))
	for _, u := range raw {
		if s, ok := u.(string); ok {
			result = append(result, s)
		}
	}
	return result
}
