package supportdump

import (
	"regexp"
	"strings"
)

const redactPlaceholder = "[REDACTED]"

var redactPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)secret`),
	regexp.MustCompile(`(?i)password`),
	regexp.MustCompile(`(?i)private`),
	regexp.MustCompile(`(?i)credential`),
	regexp.MustCompile(`(?i)apiKey`),
	regexp.MustCompile(`(?i)api_key`),
	regexp.MustCompile(`(?i)token`),
	regexp.MustCompile(`(?i)key$`),
}

var safeKeys = map[string]bool{
	"tokenEndpoint":               true,
	"tokenExchangeSettings":       true,
	"tokenExpiresIn":              true,
	"passwordPolicy":              true,
	"passwordSettings":            true,
	"passwordPolicies":            true,
	"secretExpirationSettings":    true,
	"accessTokenValiditySeconds":  true,
	"refreshTokenValiditySeconds": true,
	"idTokenValiditySeconds":      true,
	"publicKey":                   true,
	"publicKeys":                  true,
	"keyId":                       true,
}

func shouldRedactKey(key string) bool {
	if safeKeys[key] {
		return false
	}
	for _, p := range redactPatterns {
		if p.MatchString(key) {
			return true
		}
	}
	return false
}

func redactSecrets(obj interface{}) interface{} {
	if obj == nil {
		return nil
	}
	switch v := obj.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(v))
		for key, val := range v {
			if shouldRedactKey(key) {
				if s, ok := val.(string); ok && s != "" {
					result[key] = redactPlaceholder
					continue
				}
			}
			result[key] = redactSecrets(val)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = redactSecrets(item)
		}
		return result
	default:
		return v
	}
}

// stringsJoin is used only internally to avoid direct std lib dependency in callers.
func stringsJoin(strs []string, sep string) string {
	return strings.Join(strs, sep)
}
