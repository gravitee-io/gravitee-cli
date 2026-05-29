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

//go:build e2e

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestLoginPersistsConfig covers the full login flow:
// `gctl login am` writes ~/.gctl/config.yaml, and subsequent commands read from it.
//
// Isolates HOME so the real user config is never touched, and temporarily
// unsets GCTL_AM_* env vars so they don't bypass the config file.
func TestLoginPersistsConfig(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("GCTL_AM_URL", "")
	t.Setenv("GCTL_AM_TOKEN", "")

	token, err := fetchAMToken()
	if err != nil {
		t.Fatalf("failed to fetch AM token: %v", err)
	}

	runInEnvCLIExpectSuccess(t, "login", "am", "--url", amURL, "--token", token)

	// The config file must exist at the isolated HOME path.
	cfgPath := filepath.Join(tmpHome, ".gctl", "config.yaml")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("config file not created at %s: %v", cfgPath, err)
	}

	var cfg struct {
		Current  string `yaml:"current"`
		Contexts map[string]struct {
			AM *struct {
				URL   string `yaml:"url"`
				Token string `yaml:"token"`
			} `yaml:"am"`
		} `yaml:"contexts"`
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("failed to parse config file: %v", err)
	}

	if cfg.Current != "default" {
		t.Errorf("expected current context 'default', got %q", cfg.Current)
	}

	ctx, ok := cfg.Contexts["default"]
	if !ok {
		t.Fatalf("expected 'default' context in config, got %+v", cfg.Contexts)
	}

	if ctx.AM == nil {
		t.Fatal("expected AM block in default context")
	}

	if ctx.AM.URL != amURL {
		t.Errorf("expected AM URL %q, got %q", amURL, ctx.AM.URL)
	}

	if ctx.AM.Token != token {
		t.Errorf("AM token in config does not match fetched token")
	}

	// A subsequent command must resolve its context from the written file,
	// not from env vars (which we cleared above).
	out := runInEnvCLIExpectSuccess(t, "am", "domain", "list", "-o", "json")
	if !strings.HasPrefix(out, "[") && !strings.HasPrefix(out, "{") {
		t.Errorf("expected JSON output from domain list, got: %s", out)
	}
}

// TestEnvVarsOverrideConfig verifies that GCTL_APIM_URL + GCTL_APIM_TOKEN bypass
// the config file entirely: even when the on-disk config points to a broken
// URL, the CLI uses the env vars and the command still succeeds.
func TestEnvVarsOverrideConfig(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Write a config pointing to a host that doesn't exist - if env vars didn't
	// override, the CLI would try this URL and fail.
	cfgDir := filepath.Join(tmpHome, ".gctl")
	if err := os.MkdirAll(cfgDir, 0o700); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	brokenCfg := `current: broken
contexts:
  broken:
    apim:
      url: http://127.0.0.1:1
      token: not-a-real-token
`

	if err := os.WriteFile(filepath.Join(cfgDir, "config.yaml"), []byte(brokenCfg), 0o600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// GCTL_APIM_URL + GCTL_APIM_TOKEN are already set by TestMain to the real
	// APIM - so the command must succeed, proving the env vars override the
	// broken config.
	runInEnvCLIExpectSuccess(t, "apim", "env", "list", "-o", "json")
}

// runInEnvCLIExpectSuccess runs the CLI with the test's current env (including
// the isolated HOME) and fails the test if the command errors.
func runInEnvCLIExpectSuccess(t *testing.T, args ...string) string {
	t.Helper()

	cmd := exec.Command(cliBinary, args...)
	cmd.Env = os.Environ()

	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))

	if err != nil {
		t.Fatalf("CLI command failed: %v\nOutput: %s\nArgs: %v", err, output, args)
	}

	return output
}
