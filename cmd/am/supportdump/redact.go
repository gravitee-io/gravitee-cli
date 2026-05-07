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

package supportdump

import (
	"regexp"
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
