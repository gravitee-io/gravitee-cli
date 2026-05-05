package lint

import (
	"strings"
	"testing"
	"time"
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
