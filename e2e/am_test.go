//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
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
)

func TestMain(m *testing.M) {
	if u := os.Getenv("AM_URL"); u != "" {
		amURL = u
	}

	// Build the CLI binary.
	binary, err := buildCLI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build CLI: %v\n", err)
		os.Exit(1)
	}

	cliBinary = binary

	defer os.Remove(binary)

	// Wait for AM to be ready.
	if err := waitForAM(amURL, 3*time.Minute); err != nil {
		fmt.Fprintf(os.Stderr, "AM not ready: %v\n", err)
		os.Exit(1)
	}

	// Login.
	if err := loginToAM(); err != nil {
		fmt.Fprintf(os.Stderr, "login failed: %v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func buildCLI() (string, error) {
	binary := "../dist/gio-e2e-test"

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
			// AM is up — it returns 401 for unauthenticated requests.
			if resp.StatusCode == 401 {
				return nil
			}
		}

		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("AM did not become ready within %s", timeout)
}

func loginToAM() error {
	// Get token via basic auth.
	req, _ := http.NewRequest("POST", amURL+"/management/auth/token", nil)
	req.SetBasicAuth("admin", "adminadmin")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token: %w", err)
	}

	// Login via CLI.
	out, err := runCLI("login", "am", "--url", amURL, "--token", tokenResp.AccessToken)
	if err != nil {
		return fmt.Errorf("CLI login failed: %s: %w", out, err)
	}

	return nil
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

	// No domains — create one for the test suite.
	out = runCLIExpectSuccess(t, "am", "domain", "create",
		"--name", "e2e-shared-domain",
		"--description", "Shared E2E test domain",
		"-o", "json")

	return extractID(t, out)
}
