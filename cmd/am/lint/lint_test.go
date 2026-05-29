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

package lint

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"gravitee.io/gctl/internal/client"
)

func TestRuleImplicitGrant(t *testing.T) {
	app := map[string]interface{}{
		"name": "My App",
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"grantTypes": []interface{}{"implicit"},
			},
		},
	}
	findings := ruleImplicitGrant([]map[string]interface{}{app})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != "critical" {
		t.Errorf("expected critical, got %q", findings[0].Severity)
	}
}

func TestRuleWildcardRedirect(t *testing.T) {
	app := map[string]interface{}{
		"name": "My App",
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"redirectUris": []interface{}{"https://*.example.com/callback"},
			},
		},
	}
	findings := ruleWildcardRedirect([]map[string]interface{}{app})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestRuleCertExpiry(t *testing.T) {
	soon := time.Now().Add(10 * 24 * time.Hour).UnixMilli()
	cert := map[string]interface{}{
		"name":      "My Cert",
		"expiresAt": float64(soon),
	}
	findings := ruleCertExpiry([]map[string]interface{}{cert})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != "critical" {
		t.Errorf("expected critical, got %q", findings[0].Severity)
	}
}

func TestRuleEmptyDomain(t *testing.T) {
	findings := ruleEmptyDomain([]map[string]interface{}{})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for empty domain, got %d", len(findings))
	}
}

func TestCalculateScore(t *testing.T) {
	findings := []LintFinding{
		{Severity: "critical"},
		{Severity: "critical"},
		{Severity: "warning"},
	}
	score := calculateScore(findings)
	// 10 - 2*2 - 1*1 = 5
	if score != 5 {
		t.Errorf("expected score 5, got %d", score)
	}
}

func TestScoreFloorZero(t *testing.T) {
	var findings []LintFinding
	for i := 0; i < 10; i++ {
		findings = append(findings, LintFinding{Severity: "critical"})
	}
	score := calculateScore(findings)
	if score != 0 {
		t.Errorf("expected score 0, got %d", score)
	}
}

func TestRuleNoIdp(t *testing.T) {
	app := map[string]interface{}{"name": "NoIdpApp"}
	findings := ruleNoIdp([]map[string]interface{}{app})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != "critical" {
		t.Errorf("expected critical, got %q", findings[0].Severity)
	}
}

func TestRuleAppDisabled(t *testing.T) {
	app := map[string]interface{}{"name": "DisabledApp", "enabled": false}
	findings := ruleAppDisabled([]map[string]interface{}{app})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != "warning" {
		t.Errorf("expected warning, got %q", findings[0].Severity)
	}
}

func TestRuleLocalhostRedirect(t *testing.T) {
	app := map[string]interface{}{
		"name": "LocalApp",
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"redirectUris": []interface{}{"http://localhost:3000/callback"},
			},
		},
	}
	findings := ruleLocalhostRedirect([]map[string]interface{}{app})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestRuleNoPkce(t *testing.T) {
	app := map[string]interface{}{
		"name": "AuthCodeApp",
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"grantTypes": []interface{}{"authorization_code"},
				// no forcePKCE or requirePKCE
			},
		},
	}
	findings := ruleNoPkce([]map[string]interface{}{app})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != "warning" {
		t.Errorf("expected warning, got %q", findings[0].Severity)
	}
}

func TestRuleLongTokenLifetime(t *testing.T) {
	app := map[string]interface{}{
		"name": "LongTokenApp",
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"accessTokenValiditySeconds": float64(7200),
			},
		},
	}
	findings := ruleLongTokenLifetime([]map[string]interface{}{app})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	// boundary: exactly 3600 should NOT fire
	exact := map[string]interface{}{
		"name": "ExactApp",
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"accessTokenValiditySeconds": float64(3600),
			},
		},
	}
	if f := ruleLongTokenLifetime([]map[string]interface{}{exact}); len(f) != 0 {
		t.Errorf("expected no findings for 3600s, got %d", len(f))
	}
}

func TestRuleHttpRedirect(t *testing.T) {
	app := map[string]interface{}{
		"name": "HttpApp",
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"redirectUris": []interface{}{"http://evil.example.com/callback"},
			},
		},
	}
	findings := ruleHttpRedirect([]map[string]interface{}{app})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for http redirect, got %d", len(findings))
	}
	// localhost http should NOT fire (covered by localhost-redirect rule)
	localApp := map[string]interface{}{
		"name": "LocalApp",
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"redirectUris": []interface{}{"http://localhost:3000/callback"},
			},
		},
	}
	if f := ruleHttpRedirect([]map[string]interface{}{localApp}); len(f) != 0 {
		t.Errorf("expected no findings for localhost http, got %d", len(f))
	}
}

func TestLintCmd(t *testing.T) {
	// App with implicit grant triggers a critical finding
	apps := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"name": "bad-app",
				"settings": map[string]interface{}{
					"oauth": map[string]interface{}{
						"grantTypes": []interface{}{"implicit"},
					},
				},
			},
		},
		"totalCount": 1,
	}
	appsBytes, _ := json.Marshal(apps)
	emptyArr, _ := json.Marshal([]interface{}{})
	emptyList, _ := json.Marshal(map[string]interface{}{"data": []interface{}{}, "totalCount": 0})

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			switch {
			case strings.Contains(path, "applications"):
				return appsBytes, nil
			case strings.Contains(path, "scopes"):
				return emptyList, nil
			default:
				return emptyArr, nil
			}
		},
	}
	f, out := newTestFactory(fake, false)
	cmd := NewLintCmd(f)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "critical") {
		t.Errorf("expected 'critical' finding in output, got: %s", output)
	}
	if !strings.Contains(output, "Score:") {
		t.Errorf("expected 'Score:' in output, got: %s", output)
	}
}

func TestRuleUnusedScope(t *testing.T) {
	scopes := []map[string]interface{}{
		{"key": "read:data"},
		{"key": "admin"},
	}
	apps := []map[string]interface{}{
		{"name": "App", "settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"scopeSettings": []interface{}{
					map[string]interface{}{"scope": "read:data"},
				},
			},
		}},
	}
	findings := ruleUnusedScope(apps, scopes)
	if len(findings) != 1 {
		t.Fatalf("expected 1 unused scope, got %d", len(findings))
	}
	if !strings.Contains(findings[0].Resource, "admin") {
		t.Errorf("expected 'admin' to be unused, got: %s", findings[0].Resource)
	}
}
