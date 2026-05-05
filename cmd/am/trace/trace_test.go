package trace

import (
	"testing"
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
