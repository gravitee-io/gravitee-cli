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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

var (
	cliBinary string
	amURL     = "http://localhost:8093"
	apimURL   = "http://localhost:18183"
)

func TestMain(m *testing.M) {
	if u := os.Getenv("AM_URL"); u != "" {
		amURL = u
	}

	if u := os.Getenv("APIM_URL"); u != "" {
		apimURL = u
	}

	binary, err := buildCLI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build CLI: %v\n", err)
		os.Exit(1)
	}

	cliBinary = binary

	defer os.Remove(binary)

	if err := waitForAM(amURL, 3*time.Minute); err != nil {
		fmt.Fprintf(os.Stderr, "AM not ready: %v\n", err)
		os.Exit(1)
	}

	if err := waitForAPIM(apimURL, 3*time.Minute); err != nil {
		fmt.Fprintf(os.Stderr, "APIM not ready: %v\n", err)
		os.Exit(1)
	}

	amToken, err := fetchAMToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch AM token: %v\n", err)
		os.Exit(1)
	}

	apimToken, err := fetchAPIMToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch APIM token: %v\n", err)
		os.Exit(1)
	}

	// Env vars bypass the config file entirely (see cmdutil.ResolveProductContext).
	// The login/config persistence flow has its own dedicated test in login_config_test.go.
	os.Setenv("GCTL_AM_URL", amURL)
	os.Setenv("GCTL_AM_TOKEN", amToken)
	os.Setenv("GCTL_APIM_URL", apimURL)
	os.Setenv("GCTL_APIM_TOKEN", apimToken)

	os.Exit(m.Run())
}

func buildCLI() (string, error) {
	binary := "../dist/gctl-e2e-test"

	cmd := exec.Command("go", "build", "-o", binary, "..")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return binary, cmd.Run()
}

func waitForAM(baseURL string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		resp, err := http.Get(baseURL + "/management/organizations/DEFAULT/environments/DEFAULT/domains")
		if err == nil {
			resp.Body.Close()
			// AM is up - it returns 401 for unauthenticated requests.
			if resp.StatusCode == 401 {
				return nil
			}
		}

		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("AM did not become ready within %s", timeout)
}

func fetchAMToken() (string, error) {
	req, _ := http.NewRequest("POST", amURL+"/management/auth/token", nil)
	req.SetBasicAuth("admin", "adminadmin")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("empty access token in response")
	}

	return tokenResp.AccessToken, nil
}

// waitForAPIM polls the APIM Management API until it accepts requests.
func waitForAPIM(baseURL string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		req, _ := http.NewRequest("GET", baseURL+"/management/organizations/DEFAULT/environments", nil)
		req.SetBasicAuth("admin", "admin")

		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}

		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("APIM did not become ready within %s", timeout)
}

// fetchAPIMToken obtains a Bearer Personal Access Token for the admin user
// on a default APIM Management API instance (basic auth admin:admin).
//
// Two-step flow (mirrors the init script in gravitee-ai-assistant):
//  1. POST /management/organizations/DEFAULT/user/login → JWT, sub claim = user ID
//  2. POST /management/organizations/DEFAULT/users/{id}/tokens → bearer PAT
//
// The token name includes a unique suffix so repeated runs against the same
// APIM instance don't collide.
func fetchAPIMToken() (string, error) {
	userID, err := apimUserLogin()
	if err != nil {
		return "", fmt.Errorf("login: %w", err)
	}

	tokenName := fmt.Sprintf("gctl-e2e-%d", time.Now().UnixNano())

	req, _ := http.NewRequest("POST",
		apimURL+"/management/organizations/DEFAULT/users/"+userID+"/tokens",
		strings.NewReader(fmt.Sprintf(`{"name":%q}`, tokenName)))
	req.SetBasicAuth("admin", "admin")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("create token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("create token: HTTP %d: %s", resp.StatusCode, body)
	}

	var tokenResp struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("decode token response: %w", err)
	}

	if tokenResp.Token == "" {
		return "", fmt.Errorf("empty token in response")
	}

	return tokenResp.Token, nil
}

func apimUserLogin() (string, error) {
	req, _ := http.NewRequest("POST", apimURL+"/management/organizations/DEFAULT/user/login", nil)
	req.SetBasicAuth("admin", "admin")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, body)
	}

	var loginResp struct {
		ID    string `json:"id"`
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", fmt.Errorf("decode login response: %w", err)
	}

	if loginResp.ID != "" {
		return loginResp.ID, nil
	}

	if loginResp.Token == "" {
		return "", fmt.Errorf("no id or token in login response")
	}

	// Fall back to extracting the sub claim from the JWT payload.
	parts := strings.Split(loginResp.Token, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("malformed JWT in login response")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		// Some encoders use standard base64 with padding.
		padded := parts[1]
		if pad := len(padded) % 4; pad != 0 {
			padded += strings.Repeat("=", 4-pad)
		}

		payload, err = base64.URLEncoding.DecodeString(padded)
		if err != nil {
			return "", fmt.Errorf("decode JWT payload: %w", err)
		}
	}

	var claims struct {
		Sub string `json:"sub"`
	}

	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", fmt.Errorf("parse JWT claims: %w", err)
	}

	if claims.Sub == "" {
		return "", fmt.Errorf("empty sub claim in JWT")
	}

	return claims.Sub, nil
}

func runCLI(args ...string) (string, error) {
	cmd := exec.Command(cliBinary, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	return strings.TrimSpace(output), err
}

func runCLIExpectSuccess(t *testing.T, args ...string) string {
	t.Helper()

	out, err := runCLI(args...)
	if err != nil {
		t.Fatalf("CLI command failed: %v\nOutput: %s\nArgs: %v", err, out, args)
	}

	return out
}

func runCLIExpectError(t *testing.T, args ...string) string {
	t.Helper()

	out, err := runCLI(args...)
	if err == nil {
		t.Fatalf("expected CLI command to fail, but it succeeded\nOutput: %s\nArgs: %v", out, args)
	}

	return out
}

func extractID(t *testing.T, jsonOutput string) string {
	t.Helper()

	var obj map[string]any
	if err := json.Unmarshal([]byte(jsonOutput), &obj); err != nil {
		t.Fatalf("failed to parse JSON output: %v\nOutput: %s", err, jsonOutput)
	}

	id, ok := obj["id"].(string)
	if !ok {
		t.Fatalf("no 'id' field in JSON output: %s", jsonOutput)
	}

	return id
}

// getDefaultDomainID lists domains and returns the first domain ID.
// If no domains exist, it creates one for testing.
func getDefaultDomainID(t *testing.T) string {
	t.Helper()

	out := runCLIExpectSuccess(t, "am", "domain", "list", "-o", "json")

	var resp struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("failed to parse domain list: %v", err)
	}

	if len(resp.Data) > 0 {
		return resp.Data[0].ID
	}

	// No domains - create one for the test suite.
	out = runCLIExpectSuccess(t, "am", "domain", "create",
		"--name", "e2e-shared-domain",
		"--description", "Shared E2E test domain",
		"-o", "json")

	return extractID(t, out)
}
