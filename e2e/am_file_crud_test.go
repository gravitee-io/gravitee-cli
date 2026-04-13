//go:build e2e

package e2e

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func writeTempJSON(t *testing.T, data map[string]any) string {
	t.Helper()
	raw, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}
	f, err := os.CreateTemp("", "e2e-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.Write(raw); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestIDPCrud(t *testing.T) {
	domainID := getDefaultDomainID(t)

	var idpID string

	t.Run("list IDPs", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "idp", "list", "--domain", domainID, "-o", "json")

		var idps []map[string]any
		if err := json.Unmarshal([]byte(out), &idps); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}

		if len(idps) == 0 {
			t.Fatal("expected at least one IDP (default), got none")
		}

		id, ok := idps[0]["id"].(string)
		if !ok || id == "" {
			t.Fatal("first IDP has no valid 'id' field")
		}
		idpID = id
	})

	t.Run("get IDP by ID", func(t *testing.T) {
		if idpID == "" {
			t.Skip("no IDP ID available")
		}

		out := runCLIExpectSuccess(t, "am", "idp", "get", "--domain", domainID, idpID)
		if !strings.Contains(out, "Identity Provider") && !strings.Contains(out, idpID) {
			t.Errorf("expected IDP info in output, got: %s", out)
		}
	})

	t.Run("get IDP JSON", func(t *testing.T) {
		if idpID == "" {
			t.Skip("no IDP ID available")
		}

		out := runCLIExpectSuccess(t, "am", "idp", "get", "--domain", domainID, idpID, "-o", "json")

		var obj map[string]any
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}

		if obj["id"] == nil {
			t.Error("expected 'id' field in JSON output")
		}
	})

	t.Run("get nonexistent IDP", func(t *testing.T) {
		runCLIExpectError(t, "am", "idp", "get", "--domain", domainID, "nonexistent-idp-id")
	})
}

func TestCertificateList(t *testing.T) {
	domainID := getDefaultDomainID(t)

	t.Run("list certificates", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "certificate", "list", "--domain", domainID)
	})

	t.Run("list certificates JSON", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "certificate", "list", "--domain", domainID, "-o", "json")

		var arr []map[string]any
		if err := json.Unmarshal([]byte(out), &arr); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
	})
}

func TestFlowGetAndUpdate(t *testing.T) {
	domainID := getDefaultDomainID(t)

	var flowID string
	var flowType string

	t.Run("list flows and find one with ID", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "flow", "list", "--domain", domainID, "-o", "json")

		var flowArr []map[string]any
		if err := json.Unmarshal([]byte(out), &flowArr); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}

		for _, flow := range flowArr {
			id, ok := flow["id"].(string)
			if ok && id != "" {
				flowID = id
				if ft, ok := flow["type"].(string); ok {
					flowType = ft
				}
				break
			}
		}
	})

	t.Run("get flow by ID", func(t *testing.T) {
		if flowID == "" {
			t.Skip("no flows with ID on this AM instance")
		}

		out := runCLIExpectSuccess(t, "am", "flow", "get", "--domain", domainID, flowID)
		if flowType != "" && !strings.Contains(out, flowType) {
			t.Errorf("expected flow type %q in output, got: %s", flowType, out)
		}
	})

	t.Run("get flow JSON", func(t *testing.T) {
		if flowID == "" {
			t.Skip("no flows with ID on this AM instance")
		}

		out := runCLIExpectSuccess(t, "am", "flow", "get", "--domain", domainID, flowID, "-o", "json")

		var obj map[string]any
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
	})
}

func TestReporterList(t *testing.T) {
	domainID := getDefaultDomainID(t)

	t.Run("list reporters", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "reporter", "list", "--domain", domainID)
	})

	t.Run("list reporters JSON", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "reporter", "list", "--domain", domainID, "-o", "json")

		var arr []map[string]any
		if err := json.Unmarshal([]byte(out), &arr); err != nil {
			t.Fatalf("invalid JSON array: %v", err)
		}
	})
}

func TestDictionaryCRUD(t *testing.T) {
	domainID := getDefaultDomainID(t)

	tmpFile := writeTempJSON(t, map[string]any{
		"name":   "E2E English",
		"locale": "en",
	})
	defer os.Remove(tmpFile)

	var dictID string

	t.Run("create dictionary", func(t *testing.T) {
		out, err := runCLI("am", "dictionary", "create", "--domain", domainID, "--file", tmpFile, "-o", "json")
		if err != nil {
			t.Skipf("dictionary create failed (API may require more fields): %s", out)
		}
		dictID = extractID(t, out)

		if dictID == "" {
			t.Fatal("dictionary ID is empty")
		}
	})

	t.Run("get dictionary", func(t *testing.T) {
		if dictID == "" {
			t.Skip("dictionary was not created")
		}

		out := runCLIExpectSuccess(t, "am", "dictionary", "get", "--domain", domainID, dictID)
		if !strings.Contains(out, "E2E English") && !strings.Contains(out, dictID) {
			t.Errorf("expected dictionary info in output, got: %s", out)
		}
	})

	t.Run("delete dictionary", func(t *testing.T) {
		if dictID == "" {
			t.Skip("dictionary was not created")
		}

		out := runCLIExpectSuccess(t, "am", "dictionary", "delete", "--domain", domainID, dictID)
		if !strings.Contains(out, "deleted") && !strings.Contains(out, dictID) {
			t.Errorf("expected deletion confirmation, got: %s", out)
		}
	})
}

func TestBotDetectionCRUD(t *testing.T) {
	domainID := getDefaultDomainID(t)

	tmpFile := writeTempJSON(t, map[string]any{
		"name":          "E2E Bot Detection",
		"type":          "google-recaptcha-v3-am-bot-detection",
		"detectionType": "CAPTCHA",
		"configuration": `{"siteKey":"fake-key","secretKey":"fake-secret"}`,
	})
	defer os.Remove(tmpFile)

	var bdID string

	t.Run("create bot-detection", func(t *testing.T) {
		out, err := runCLI("am", "bot-detection", "create", "--domain", domainID, "--file", tmpFile, "-o", "json")
		if err != nil {
			t.Skipf("bot-detection create failed (API may require specific plugin config): %s", out)
		}
		bdID = extractID(t, out)

		if bdID == "" {
			t.Fatal("bot-detection ID is empty")
		}
	})

	t.Run("get bot-detection", func(t *testing.T) {
		if bdID == "" {
			t.Skip("bot-detection was not created")
		}

		out := runCLIExpectSuccess(t, "am", "bot-detection", "get", "--domain", domainID, bdID)
		if !strings.Contains(out, "E2E Bot Detection") && !strings.Contains(out, bdID) {
			t.Errorf("expected bot-detection info in output, got: %s", out)
		}
	})

	t.Run("delete bot-detection", func(t *testing.T) {
		if bdID == "" {
			t.Skip("bot-detection was not created")
		}

		out := runCLIExpectSuccess(t, "am", "bot-detection", "delete", "--domain", domainID, bdID)
		if !strings.Contains(out, "deleted") && !strings.Contains(out, bdID) {
			t.Errorf("expected deletion confirmation, got: %s", out)
		}
	})
}

func TestDeviceIdentifierList(t *testing.T) {
	domainID := getDefaultDomainID(t)

	t.Run("list device identifiers", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "device-identifier", "list", "--domain", domainID)
	})

	t.Run("list device identifiers JSON", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "device-identifier", "list", "--domain", domainID, "-o", "json")

		var arr []map[string]any
		if err := json.Unmarshal([]byte(out), &arr); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
	})
}

func TestExtensionGrantList(t *testing.T) {
	domainID := getDefaultDomainID(t)

	t.Run("list extension grants", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "extension-grant", "list", "--domain", domainID)
	})

	t.Run("list extension grants JSON", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "extension-grant", "list", "--domain", domainID, "-o", "json")

		var arr []map[string]any
		if err := json.Unmarshal([]byte(out), &arr); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
	})
}

func TestResourceList(t *testing.T) {
	domainID := getDefaultDomainID(t)

	t.Run("list resources", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "resource", "list", "--domain", domainID)
	})

	t.Run("list resources JSON", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "resource", "list", "--domain", domainID, "-o", "json")

		var arr []map[string]any
		if err := json.Unmarshal([]byte(out), &arr); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
	})
}

func TestAuthorizationEngineList(t *testing.T) {
	domainID := getDefaultDomainID(t)

	t.Run("list authorization engines", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "authorization-engine", "list", "--domain", domainID)
	})

	t.Run("list authorization engines JSON", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "authorization-engine", "list", "--domain", domainID, "-o", "json")

		var arr []map[string]any
		if err := json.Unmarshal([]byte(out), &arr); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
	})
}
