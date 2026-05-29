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
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"gravitee.io/gctl/internal/client"
)

func TestCheckUserStatusEnabled(t *testing.T) {
	user := map[string]interface{}{"enabled": true, "accountNonLocked": true}
	step := checkUserStatus(user)
	if step.Status != "ok" {
		t.Errorf("expected ok, got %q: %s", step.Status, step.Detail)
	}
}

func TestCheckUserStatusDisabled(t *testing.T) {
	user := map[string]interface{}{"enabled": false}
	step := checkUserStatus(user)
	if step.Status != "block" {
		t.Errorf("expected block, got %q", step.Status)
	}
}

func TestCheckUserStatusLocked(t *testing.T) {
	user := map[string]interface{}{"enabled": true, "accountNonLocked": false}
	step := checkUserStatus(user)
	if step.Status != "block" {
		t.Errorf("expected block, got %q", step.Status)
	}
}

func TestCheckIdpMatchSuccess(t *testing.T) {
	user := map[string]interface{}{"source": "idp-1"}
	app := map[string]interface{}{
		"identityProviders": []interface{}{
			map[string]interface{}{"identity": "idp-1"},
		},
	}
	domainIdps := []map[string]interface{}{{"id": "idp-1", "name": "GitHub"}}
	step := checkIdpMatch(user, app, domainIdps)
	if step.Status != "ok" {
		t.Errorf("expected ok, got %q: %s", step.Status, step.Detail)
	}
}

func TestCheckIdpMatchFail(t *testing.T) {
	user := map[string]interface{}{"source": "idp-other"}
	app := map[string]interface{}{
		"identityProviders": []interface{}{
			map[string]interface{}{"identity": "idp-1"},
		},
	}
	domainIdps := []map[string]interface{}{{"id": "idp-1", "name": "GitHub"}}
	step := checkIdpMatch(user, app, domainIdps)
	if step.Status != "block" {
		t.Errorf("expected block, got %q", step.Status)
	}
}

func TestBuildVerdict_AllOk(t *testing.T) {
	steps := []TraceStep{
		{Status: "ok"}, {Status: "ok"},
	}
	v := buildVerdict(steps)
	if !v.CanAuthenticate {
		t.Error("expected can authenticate")
	}
	if v.Reason != "All checks passed" {
		t.Errorf("unexpected reason: %s", v.Reason)
	}
}

func TestBuildVerdict_Blocked(t *testing.T) {
	steps := []TraceStep{
		{Status: "ok"},
		{Status: "block", Detail: "User disabled"},
	}
	v := buildVerdict(steps)
	if v.CanAuthenticate {
		t.Error("expected cannot authenticate")
	}
}

func TestCheckGrantTypes_UserFacing(t *testing.T) {
	app := map[string]interface{}{
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"grantTypes": []interface{}{"authorization_code", "client_credentials"},
			},
		},
	}
	step := checkGrantTypes(app)
	if step.Status != "ok" {
		t.Errorf("expected ok, got %q: %s", step.Status, step.Detail)
	}
	if !strings.Contains(step.Detail, "authorization_code") {
		t.Errorf("expected 'authorization_code' in detail, got %q", step.Detail)
	}
}

func TestCheckGrantTypes_NoUserFacing(t *testing.T) {
	app := map[string]interface{}{
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"grantTypes": []interface{}{"client_credentials"},
			},
		},
	}
	step := checkGrantTypes(app)
	if step.Status != "warn" {
		t.Errorf("expected warn, got %q", step.Status)
	}
}

func TestCheckMfa_NoFactors(t *testing.T) {
	step := checkMfa(map[string]interface{}{}, nil)
	if step.Status != "ok" {
		t.Errorf("expected ok when no domain factors, got %q", step.Status)
	}
}

func TestCheckMfa_Enrolled(t *testing.T) {
	user := map[string]interface{}{
		"factors": []interface{}{"factor-1"},
	}
	domainFactors := []map[string]interface{}{{"id": "factor-1", "name": "TOTP"}}
	step := checkMfa(user, domainFactors)
	if step.Status != "ok" {
		t.Errorf("expected ok when factor enrolled, got %q: %s", step.Status, step.Detail)
	}
}

func TestCheckMfa_RequiredNotEnrolled(t *testing.T) {
	user := map[string]interface{}{}
	domainFactors := []map[string]interface{}{{"id": "factor-1", "name": "TOTP"}}
	step := checkMfa(user, domainFactors)
	if step.Status != "warn" {
		t.Errorf("expected warn when factors required but not enrolled, got %q", step.Status)
	}
	if !strings.Contains(step.Detail, "TOTP") {
		t.Errorf("expected factor name 'TOTP' in detail, got %q", step.Detail)
	}
}

func TestCheckFlows_WithPolicies(t *testing.T) {
	flows := []map[string]interface{}{
		{
			"type": "LOGIN",
			"pre": []interface{}{
				map[string]interface{}{"name": "IP-Filter"},
			},
		},
	}
	step := checkFlows(flows)
	if step.Status != "ok" {
		t.Errorf("expected ok, got %q", step.Status)
	}
	if !strings.Contains(step.Detail, "IP-Filter") {
		t.Errorf("expected policy name in detail, got %q", step.Detail)
	}
}

func TestCheckFlows_NoPolicies(t *testing.T) {
	step := checkFlows(nil)
	if step.Status != "ok" {
		t.Errorf("expected ok, got %q", step.Status)
	}
	if !strings.Contains(step.Detail, "No pre-login") {
		t.Errorf("expected 'No pre-login' in detail, got %q", step.Detail)
	}
}

func TestCheckConsent_Skipped(t *testing.T) {
	app := map[string]interface{}{
		"settings": map[string]interface{}{
			"advanced": map[string]interface{}{"skipConsent": true},
		},
	}
	step := checkConsent(app)
	if !strings.Contains(step.Detail, "skipped") {
		t.Errorf("expected 'skipped' in detail, got %q", step.Detail)
	}
}

func TestCheckTokenConfig(t *testing.T) {
	app := map[string]interface{}{
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"accessTokenValiditySeconds":  float64(3600),
				"refreshTokenValiditySeconds": float64(86400),
				"idTokenValiditySeconds":      nil,
			},
		},
	}
	step := checkTokenConfig(app)
	if step.Status != "ok" {
		t.Errorf("expected ok, got %q", step.Status)
	}
	if !strings.Contains(step.Detail, "3600") {
		t.Errorf("expected access token value in detail, got %q", step.Detail)
	}
	if !strings.Contains(step.Detail, "default") {
		t.Errorf("expected 'default' for nil id token, got %q", step.Detail)
	}
}

func TestRunTrace(t *testing.T) {
	user := map[string]interface{}{
		"id": "user-1", "username": "john", "email": "john@example.com",
		"enabled": true, "accountNonLocked": true,
		"source": "idp-1",
	}
	app := map[string]interface{}{
		"id": "app-1", "name": "MyApp",
		"identityProviders": []interface{}{
			map[string]interface{}{"identity": "idp-1"},
		},
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"grantTypes": []interface{}{"authorization_code"},
			},
			"advanced": map[string]interface{}{"skipConsent": false},
		},
	}
	userBytes, _ := json.Marshal(user)
	appBytes, _ := json.Marshal(app)
	empty, _ := json.Marshal([]interface{}{})

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			switch {
			case strings.Contains(path, "/users/user-1"):
				return userBytes, nil
			case strings.Contains(path, "/applications/app-1"):
				return appBytes, nil
			case strings.Contains(path, "/identities"):
				return empty, nil
			case strings.Contains(path, "/factors"):
				return empty, nil
			case strings.Contains(path, "/flows"):
				return empty, nil
			}
			return nil, fmt.Errorf("unexpected path: %s", path)
		},
	}
	f, out := newTestFactory(fake)
	cmd := NewTraceCmd(f)
	cmd.SetArgs([]string{"--user", "user-1", "--app", "app-1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "john@example.com") {
		t.Errorf("expected user email in output, got: %s", output)
	}
	if !strings.Contains(output, "All checks passed") {
		t.Errorf("expected 'All checks passed' verdict, got: %s", output)
	}
}
