package am

import (
	"testing"

	"github.com/gravitee-io/gio-cli/internal/factory"
)

// TestAMCommandsRegistered guards against subcommands being implemented but
// never wired into NewAMCmd. Add new expected commands here when they ship.
func TestAMCommandsRegistered(t *testing.T) {
	expected := []string{
		"alert",
		"analytics",
		"app",
		"audit",
		"auth-device-notifier",
		"authorization-engine",
		"bot-detection",
		"certificate",
		"data-plane",
		"device-identifier",
		"dictionary",
		"diff",
		"doctor",
		"domain",
		"email",
		"entrypoint",
		"extension-grant",
		"factor",
		"flow",
		"form",
		"group",
		"health",
		"idp",
		"lint",
		"logout",
		"member",
		"org",
		"password-policy",
		"plugin",
		"protected-resource",
		"reporter",
		"resource",
		"role",
		"scope",
		"set",
		"shell",
		"status",
		"support-dump",
		"theme",
		"token",
		"trace",
		"test-oidc",
		"user",
		"watch",
		"whoami",
	}

	cmd := NewAMCmd(&factory.Factory{})

	registered := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		registered[sub.Name()] = true
	}

	for _, name := range expected {
		if !registered[name] {
			t.Errorf("expected subcommand %q to be registered on `am`, but it isn't — wire it up in NewAMCmd", name)
		}
	}
}
