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

package am

import (
	"testing"

	"github.com/gravitee-io/gio-cli/internal/factory"
)

// TestAMCmdROCommandsRegistered guards against RO subcommands drifting from NewAMCmdRO.
// Add new expected commands here when they are wired into the RO bundle.
func TestAMCmdROCommandsRegistered(t *testing.T) {
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
		"member",
		"password-policy",
		"plugin",
		"protected-resource",
		"reporter",
		"resource",
		"role",
		"scope",
		"status",
		"theme",
		"token",
		"user",
		"whoami",
	}

	cmd := NewAMCmdRO(&factory.Factory{})

	registered := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		registered[sub.Name()] = true
	}

	for _, name := range expected {
		if !registered[name] {
			t.Errorf("expected subcommand %q to be registered on `am` (RO), but it isn't — wire it up in NewAMCmdRO", name)
		}
	}
}

// TestAMCommandsRegistered guards against subcommands being implemented but
// never wired into NewAMCmd. Add new expected commands here when they ship.
func TestAMCommandsRegistered(t *testing.T) {
	expected := []string{
		"alert",
		"analytics",
		"app",
		"audit",
		"auth",
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
