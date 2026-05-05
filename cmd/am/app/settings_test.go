package app

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func TestAppSettingsView(t *testing.T) {
	resp := map[string]interface{}{
		"id": "app-1", "name": "My App", "type": "WEB",
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"clientId":     "client-123",
				"clientSecret": "secret-abc",
				"grantTypes":   []string{"authorization_code", "refresh_token"},
				"redirectUris": []string{"https://myapp.com/callback"},
			},
		},
	}

	data, _ := json.Marshal(resp)

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/applications/app-1") {
				t.Errorf("unexpected path: %s", path)
			}
			return data, nil
		},
	}

	f, out := newTestFactory(fake, false)

	cmd := newSettingsCmd(f)
	cmd.SetArgs([]string{"app-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "client-123") {
		t.Errorf("expected clientId in output, got: %s", out.String())
	}
}

func TestAppSettingsUpdate(t *testing.T) {
	appResp := map[string]interface{}{
		"id": "app-1", "name": "My App", "type": "WEB",
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"clientId":   "client-123",
				"grantTypes": []string{"authorization_code", "refresh_token"},
			},
		},
	}

	data, _ := json.Marshal(appResp)

	var capturedBody map[string]interface{}

	fake := &client.FakeClient{
		PatchFunc: func(path string, body interface{}) ([]byte, error) {
			if !strings.Contains(path, "/applications/app-1") {
				t.Errorf("unexpected path: %s", path)
			}

			switch b := body.(type) {
			case []byte:
				_ = json.Unmarshal(b, &capturedBody)
			case json.RawMessage:
				_ = json.Unmarshal(b, &capturedBody)
			}

			return data, nil
		},
	}

	f, _ := newTestFactory(fake, false)

	cmd := newSettingsCmd(f)
	cmd.SetArgs([]string{"app-1", "--grant-types", "authorization_code,refresh_token"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the PATCH body contains settings.oauth.grantTypes
	settings, _ := capturedBody["settings"].(map[string]interface{})
	if settings == nil {
		t.Fatalf("expected settings in request body")
	}
	oauth, _ := settings["oauth"].(map[string]interface{})
	if oauth == nil {
		t.Fatalf("expected settings.oauth in request body")
	}
	if _, ok := oauth["grantTypes"]; !ok {
		t.Errorf("expected grantTypes in settings.oauth, got: %v", oauth)
	}
}
