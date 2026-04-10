//go:build e2e

package e2e

import (
	"bytes"
	cryptorand "crypto/rand"
	"embed"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//go:embed fixtures
var fixturesFS embed.FS

// runSuffix is generated once per test process so sequential runs never collide
// on server-unique fields (API path, name) when a prior cleanup silently failed.
var runSuffix = func() string {
	b := make([]byte, 4)
	if _, err := cryptorand.Read(b); err != nil {
		panic("cannot init run suffix: " + err.Error())
	}

	return hex.EncodeToString(b)
}()

// writeFixture extracts an embedded fixture to a temp file, replacing the shared
// "gio-e2e-test" token with a per-test unique value.
func writeFixture(t *testing.T, name string) string {
	t.Helper()

	data, err := fixturesFS.ReadFile("fixtures/" + name)
	if err != nil {
		t.Fatalf("failed to read embedded fixture %q: %v", name, err)
	}

	topLevel := strings.SplitN(t.Name(), "/", 2)[0]
	suffix := strings.ToLower(topLevel) + "-" + runSuffix
	data = bytes.ReplaceAll(data, []byte("gio-e2e-test"), []byte("gio-e2e-test-"+suffix))

	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("failed to write fixture to temp file: %v", err)
	}

	return path
}
