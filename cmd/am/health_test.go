package am

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newAMTestFactory(fc *client.FakeClient) (*factory.Factory, *bytes.Buffer) {
	return newAMTestFactoryWithConfig(fc, &config.Config{
		CurrentContext: "am-test",
		Contexts: map[string]config.Context{
			"am-test": {URL: "https://am.example.com", Token: "tok", Type: "am", Org: "DEFAULT", Env: "DEFAULT"},
		},
	})
}

func newAMTestFactoryWithConfig(fc *client.FakeClient, cfg *config.Config) (*factory.Factory, *bytes.Buffer) {
	out := &bytes.Buffer{}
	var resolved *config.ResolvedContext
	if cfg != nil && cfg.CurrentContext != "" {
		if ctx, ok := cfg.Contexts[cfg.CurrentContext]; ok {
			resolved = &config.ResolvedContext{
				Name:  cfg.CurrentContext,
				URL:   ctx.URL,
				Token: ctx.Token,
				Org:   ctx.Org,
				Env:   ctx.Env,
				Type:  ctx.Type,
			}
		}
	}
	return &factory.Factory{
		Config:       cfg,
		Resolved:     resolved,
		Client:       fc,
		IOStreams:    factory.IOStreams{Out: out, Err: &bytes.Buffer{}},
		OutputFormat: "table",
	}, out
}

func TestHealth(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if path != "/management/health" {
				t.Errorf("unexpected path: %s", path)
			}
			return []byte(`{"status":"UP"}`), nil
		},
	}
	f, out := newAMTestFactory(fake)
	cmd := newHealthCmd(f)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "healthy") {
		t.Errorf("expected 'healthy' in output, got: %s", out.String())
	}
}

func TestWhoami(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if path != "/management/user" {
				t.Errorf("unexpected path: %s", path)
			}
			return []byte(`{"username":"admin","email":"admin@example.com"}`), nil
		},
	}
	f, out := newAMTestFactory(fake)
	cmd := newWhoamiCmd(f)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "admin") {
		t.Errorf("expected 'admin' in output, got: %s", out.String())
	}
}

func TestLogout(t *testing.T) {
	cfg := &config.Config{
		Contexts:       map[string]config.Context{"ctx1": {Token: "tok"}},
		CurrentContext: "ctx1",
	}
	f, out := newAMTestFactoryWithConfig(nil, cfg)
	cmd := newLogoutCmd(f)
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Config.Contexts["ctx1"].Token != "" {
		t.Error("expected token to be cleared")
	}
	if !strings.Contains(out.String(), "Logged out") {
		t.Errorf("expected success message, got: %s", out.String())
	}
}

func TestLogoutPersists(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.json")
	cfg := &config.Config{
		Contexts:       map[string]config.Context{"ctx1": {Token: "tok"}},
		CurrentContext: "ctx1",
	}
	f, _ := newAMTestFactoryWithConfig(nil, cfg)
	f.ConfigPath = cfgPath

	cmd := newLogoutCmd(f)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("config not written: %v", err)
	}
	var saved config.Config
	if err := json.Unmarshal(data, &saved); err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if saved.Contexts["ctx1"].Token != "" {
		t.Error("token not cleared in saved config")
	}
}

func TestLogoutAll(t *testing.T) {
	cfg := &config.Config{
		Contexts: map[string]config.Context{
			"ctx1": {Token: "tok1"},
			"ctx2": {Token: "tok2"},
		},
		CurrentContext: "ctx1",
	}
	f, out := newAMTestFactoryWithConfig(nil, cfg)
	cmd := newLogoutCmd(f)
	cmd.SetArgs([]string{"--all"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for name, ctx := range f.Config.Contexts {
		if ctx.Token != "" {
			t.Errorf("expected token cleared for %s", name)
		}
	}
	if !strings.Contains(out.String(), "2") {
		t.Errorf("expected count in message, got: %s", out.String())
	}
}

func TestStatus(t *testing.T) {
	cfg := &config.Config{
		Contexts: map[string]config.Context{
			"myws": {URL: "https://am.example.com", Token: "tok", Org: "DEFAULT", Env: "DEFAULT"},
		},
		CurrentContext: "myws",
	}
	f, out := newAMTestFactoryWithConfig(nil, cfg)
	cmd := newStatusCmd(f)
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "myws") {
		t.Errorf("expected workspace name, got: %s", out.String())
	}
	if !strings.Contains(out.String(), "https://am.example.com") {
		t.Errorf("expected URL, got: %s", out.String())
	}
	if !strings.Contains(out.String(), "authenticated: yes") {
		t.Errorf("expected authenticated:yes, got: %s", out.String())
	}
}

func TestStatusNoContext(t *testing.T) {
	cfg := &config.Config{Contexts: map[string]config.Context{}}
	f, out := newAMTestFactoryWithConfig(nil, cfg)
	cmd := newStatusCmd(f)
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "(not set)") {
		t.Errorf("expected (not set), got: %s", out.String())
	}
}

func TestDoctor(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if strings.Contains(path, "/management/user") {
				return []byte(`{"id":"u1","username":"admin"}`), nil
			}
			return nil, fmt.Errorf("unexpected path: %s", path)
		},
	}
	f, out := newAMTestFactory(fake)
	cmd := newDoctorCmd(f)
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "OK") {
		t.Errorf("expected OK checks, got: %s", output)
	}
	if strings.Contains(output, "FAIL") {
		t.Errorf("unexpected FAIL in output: %s", output)
	}
	if !strings.Contains(output, "connect") {
		t.Errorf("expected connect check in output, got: %s", output)
	}
}

func TestDoctorNoConfig(t *testing.T) {
	cfg := &config.Config{Contexts: map[string]config.Context{}}
	f, out := newAMTestFactoryWithConfig(nil, cfg)
	cmd := newDoctorCmd(f)
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "FAIL") {
		t.Errorf("expected FAIL for no config, got: %s", out.String())
	}
}

func TestDoctorNoCurrentContext(t *testing.T) {
	cfg := &config.Config{
		CurrentContext: "",
		Contexts: map[string]config.Context{
			"am-test": {URL: "https://am.example.com", Token: "tok", Type: "am"},
		},
	}
	f, out := newAMTestFactoryWithConfig(nil, cfg)
	cmd := newDoctorCmd(f)
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "WARN") {
		t.Errorf("expected WARN for no current context, got: %s", output)
	}
	if !strings.Contains(output, "context") {
		t.Errorf("expected context check in output, got: %s", output)
	}
}

func TestDoctorEmptyToken(t *testing.T) {
	cfg := &config.Config{
		CurrentContext: "am-test",
		Contexts: map[string]config.Context{
			"am-test": {URL: "https://am.example.com", Token: "", Type: "am"},
		},
	}
	f, out := newAMTestFactoryWithConfig(nil, cfg)
	cmd := newDoctorCmd(f)
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "FAIL") {
		t.Errorf("expected FAIL for empty token, got: %s", output)
	}
	if !strings.Contains(output, "auth") {
		t.Errorf("expected auth check in output, got: %s", output)
	}
}

func TestDoctorNoDomain(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			return []byte(`{"id":"u1","username":"admin"}`), nil
		},
	}
	cfg := &config.Config{
		CurrentContext: "am-test",
		Contexts: map[string]config.Context{
			"am-test": {URL: "https://am.example.com", Token: "tok", Type: "am", Org: "DEFAULT", Env: "DEFAULT"},
		},
	}
	f, out := newAMTestFactoryWithConfig(fake, cfg)
	// Resolved has no Domain set (default empty string)
	cmd := newDoctorCmd(f)
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "WARN") {
		t.Errorf("expected WARN for no domain, got: %s", output)
	}
	if !strings.Contains(output, "domain") {
		t.Errorf("expected domain check in output, got: %s", output)
	}
}

func TestDoctorConnectFail(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			return nil, fmt.Errorf("connection refused")
		},
	}
	f, out := newAMTestFactory(fake)
	cmd := newDoctorCmd(f)
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "FAIL") {
		t.Errorf("expected FAIL for connect error, got: %s", output)
	}
	if !strings.Contains(output, "connect") {
		t.Errorf("expected connect check in output, got: %s", output)
	}
}

func TestDoctorConnectSkippedWhenNotAMContext(t *testing.T) {
	cfg := &config.Config{
		CurrentContext: "apim-test",
		Contexts: map[string]config.Context{
			"apim-test": {URL: "https://apim.example.com", Token: "tok", Type: "apim"},
		},
	}
	f, out := newAMTestFactoryWithConfig(nil, cfg)
	cmd := newDoctorCmd(f)
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "connect") {
		t.Errorf("expected connect check in output, got: %s", output)
	}
	if !strings.Contains(output, "WARN") {
		t.Errorf("expected WARN for skipped connect, got: %s", output)
	}
	if strings.Contains(output, "FAIL") {
		t.Errorf("unexpected FAIL in output: %s", output)
	}
}
