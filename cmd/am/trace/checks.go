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

package trace

import (
	"fmt"
	"strings"
)

type TraceStep struct {
	Phase  string
	Status string // "ok", "warn", "block"
	Label  string
	Detail string
}

type TraceVerdict struct {
	CanAuthenticate bool
	Reason          string
}

func checkUserStatus(user map[string]interface{}) TraceStep {
	enabled, _ := user["enabled"].(bool)
	locked, hasLocked := user["accountNonLocked"].(bool)

	if !enabled {
		return TraceStep{Phase: "user-status", Status: "block", Label: "User status", Detail: "User account is disabled"}
	}
	if hasLocked && !locked {
		return TraceStep{Phase: "user-status", Status: "block", Label: "User status", Detail: "User account is locked"}
	}
	credExp, hasCredExp := user["credentialNonExpired"].(bool)
	if hasCredExp && !credExp {
		return TraceStep{Phase: "user-status", Status: "warn", Label: "User status", Detail: "User credentials are expired"}
	}
	return TraceStep{Phase: "user-status", Status: "ok", Label: "User status", Detail: "enabled, account not locked"}
}

func checkIdpMatch(user, app map[string]interface{}, domainIdps []map[string]interface{}) TraceStep {
	appIdps := extractIdpIds(app)
	if len(appIdps) == 0 {
		return TraceStep{Phase: "idp-match", Status: "block", Label: "Identity source", Detail: "Application has no identity providers assigned"}
	}
	userSource, _ := user["source"].(string)
	if userSource == "" {
		return TraceStep{Phase: "idp-match", Status: "warn", Label: "Identity source", Detail: "User has no source identity provider set"}
	}
	for _, idpID := range appIdps {
		if idpID == userSource {
			sourceName := lookupIdpName(domainIdps, userSource)
			return TraceStep{Phase: "idp-match", Status: "ok", Label: "Identity source", Detail: fmt.Sprintf("%s matches app IdP", sourceName)}
		}
	}
	sourceName := lookupIdpName(domainIdps, userSource)
	appNames := make([]string, 0, len(appIdps))
	for _, id := range appIdps {
		appNames = append(appNames, lookupIdpName(domainIdps, id))
	}
	return TraceStep{Phase: "idp-match", Status: "block", Label: "Identity source",
		Detail: fmt.Sprintf("User's IdP '%s' not in app (app has: %s)", sourceName, strings.Join(appNames, ", "))}
}

func checkGrantTypes(app map[string]interface{}) TraceStep {
	oauth, _ := app["settings"].(map[string]interface{})
	oauthMap, _ := oauth["oauth"].(map[string]interface{})
	grantTypes, _ := oauthMap["grantTypes"].([]interface{})
	var userFacing []string
	for _, g := range grantTypes {
		if s, ok := g.(string); ok && (s == "password" || s == "authorization_code") {
			userFacing = append(userFacing, s)
		}
	}
	if len(userFacing) > 0 {
		return TraceStep{Phase: "grant-type", Status: "ok", Label: "Grant types", Detail: strings.Join(userFacing, ", ") + " available"}
	}
	allGrants := make([]string, 0, len(grantTypes))
	for _, g := range grantTypes {
		if s, ok := g.(string); ok {
			allGrants = append(allGrants, s)
		}
	}
	return TraceStep{Phase: "grant-type", Status: "warn", Label: "Grant types",
		Detail: "No user-facing grant type (only: " + strings.Join(allGrants, ", ") + ")"}
}

func checkMfa(user map[string]interface{}, domainFactors []map[string]interface{}) TraceStep {
	if len(domainFactors) == 0 {
		return TraceStep{Phase: "mfa", Status: "ok", Label: "MFA", Detail: "MFA not required"}
	}
	userFactors, _ := user["factors"].([]interface{})
	if len(userFactors) > 0 {
		return TraceStep{Phase: "mfa", Status: "ok", Label: "MFA", Detail: fmt.Sprintf("MFA factor enrolled (%d)", len(userFactors))}
	}
	available := make([]string, 0, len(domainFactors))
	for _, f := range domainFactors {
		if name, ok := f["name"].(string); ok {
			available = append(available, name)
		}
	}
	return TraceStep{Phase: "mfa", Status: "warn", Label: "MFA",
		Detail: "MFA required but user has no enrolled factor. Available: " + strings.Join(available, ", ")}
}

func checkFlows(flows []map[string]interface{}) TraceStep {
	var policies []string
	for _, flow := range flows {
		flowType, _ := flow["type"].(string)
		if flowType != "ROOT" && flowType != "LOGIN" {
			continue
		}
		pre, _ := flow["pre"].([]interface{})
		for _, step := range pre {
			if sm, ok := step.(map[string]interface{}); ok {
				if name, ok := sm["name"].(string); ok && name != "" {
					policies = append(policies, name)
				}
			}
		}
	}
	if len(policies) > 0 {
		return TraceStep{Phase: "pre-login", Status: "ok", Label: "Pre-login flows",
			Detail: fmt.Sprintf("%d policies (%s)", len(policies), strings.Join(policies, ", "))}
	}
	return TraceStep{Phase: "pre-login", Status: "ok", Label: "Pre-login flows", Detail: "No pre-login policies"}
}

func checkConsent(app map[string]interface{}) TraceStep {
	advanced, _ := app["settings"].(map[string]interface{})
	adv, _ := advanced["advanced"].(map[string]interface{})
	skipConsent, _ := adv["skipConsent"].(bool)
	detail := "will be requested"
	if skipConsent {
		detail = "will be skipped"
	}
	return TraceStep{Phase: "consent", Status: "ok", Label: "Consent", Detail: detail}
}

func checkTokenConfig(app map[string]interface{}) TraceStep {
	settings, _ := app["settings"].(map[string]interface{})
	oauth, _ := settings["oauth"].(map[string]interface{})
	access := fmtTokenVal(oauth["accessTokenValiditySeconds"])
	refresh := fmtTokenVal(oauth["refreshTokenValiditySeconds"])
	id := fmtTokenVal(oauth["idTokenValiditySeconds"])
	return TraceStep{Phase: "token", Status: "ok", Label: "Token config",
		Detail: fmt.Sprintf("access=%ss, refresh=%ss, id=%ss", access, refresh, id)}
}

func buildVerdict(steps []TraceStep) TraceVerdict {
	for _, s := range steps {
		if s.Status == "block" {
			return TraceVerdict{CanAuthenticate: false, Reason: s.Detail}
		}
	}
	for _, s := range steps {
		if s.Status == "warn" {
			return TraceVerdict{CanAuthenticate: true, Reason: "Likely yes, but warnings present"}
		}
	}
	return TraceVerdict{CanAuthenticate: true, Reason: "All checks passed"}
}

func extractIdpIds(app map[string]interface{}) []string {
	var ids []string
	if idps, ok := app["identityProviders"].([]interface{}); ok {
		for _, idp := range idps {
			if m, ok := idp.(map[string]interface{}); ok {
				if id, ok := m["identity"].(string); ok {
					ids = append(ids, id)
				}
			}
		}
	}
	if idps, ok := app["identities"].([]interface{}); ok {
		for _, idp := range idps {
			if id, ok := idp.(string); ok {
				ids = append(ids, id)
			}
		}
	}
	return ids
}

func lookupIdpName(idps []map[string]interface{}, id string) string {
	for _, idp := range idps {
		if idp["id"] == id {
			if name, ok := idp["name"].(string); ok {
				return name
			}
		}
	}
	return id
}

func fmtTokenVal(v interface{}) string {
	if v == nil {
		return "default"
	}
	return fmt.Sprintf("%v", v)
}
