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

package domain

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func TestDomainExport(t *testing.T) {
	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			switch {
			case strings.HasSuffix(path, "/domains/dom-1"):
				return []byte(`{"id":"dom-1","name":"Test Domain"}`), nil
			case strings.Contains(path, "/applications"):
				return []byte(`{"data":[],"currentPage":0,"totalCount":0}`), nil
			case strings.Contains(path, "/identityProviders"):
				return []byte(`[]`), nil
			case strings.Contains(path, "/roles"):
				return []byte(`{"data":[],"currentPage":0,"totalCount":0}`), nil
			case strings.Contains(path, "/scopes"):
				return []byte(`{"data":[],"currentPage":0,"totalCount":0}`), nil
			case strings.Contains(path, "/factors"):
				return []byte(`[]`), nil
			case strings.Contains(path, "/groups"):
				return []byte(`{"data":[],"currentPage":0,"totalCount":0}`), nil
			case strings.Contains(path, "/flows"):
				return []byte(`[]`), nil
			case strings.Contains(path, "/certificates"):
				return []byte(`[]`), nil
			}
			t.Logf("unmatched path: %s", path)
			return []byte(`{}`), nil
		},
	}
	f, out := newTestFactory(fake, false)
	cmd := newExportCmd(f)
	cmd.SetArgs([]string{"dom-1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]json.RawMessage
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("export output is not valid JSON: %v\noutput: %s", err, out.String())
	}
	if _, ok := result["domain"]; !ok {
		t.Error("export JSON missing 'domain' key")
	}
}

func TestDomainImportCreatesNewDomain(t *testing.T) {
	exportData := map[string]interface{}{
		"domain":            map[string]interface{}{"name": "Source Domain", "description": "desc"},
		"applications":      []interface{}{},
		"identityProviders": []interface{}{},
		"roles":             []interface{}{},
		"scopes":            []interface{}{},
		"factors":           []interface{}{},
		"groups":            []interface{}{},
		"flows":             []interface{}{},
		"certificates":      []interface{}{},
	}
	raw, _ := json.Marshal(exportData)
	tmpFile := filepath.Join(t.TempDir(), "export.json")
	if err := os.WriteFile(tmpFile, raw, 0600); err != nil {
		t.Fatal(err)
	}

	domainCreated := false
	fake := &client.FakeClient{
		PostFunc: func(path string, body interface{}) ([]byte, error) {
			if strings.Contains(path, "/domains") && !strings.Contains(path, "/domains/") {
				domainCreated = true
				return []byte(`{"id":"new-domain-id","name":"Source Domain"}`), nil
			}
			return []byte(`{}`), nil
		},
	}
	f, out := newTestFactory(fake, false)
	cmd := newImportCmd(f)
	cmd.SetArgs([]string{tmpFile})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !domainCreated {
		t.Error("expected a new domain to be created via POST")
	}
	if !strings.Contains(out.String(), "Import complete") {
		t.Errorf("expected 'Import complete' in output, got: %s", out.String())
	}
}

func TestDomainImportReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)
	cmd := newImportCmd(f)
	cmd.SetArgs([]string{"nonexistent.json"})
	if err := cmd.Execute(); err == nil {
		t.Error("expected read-only error")
	}
}

func TestDomainCopyCreatesAndCopies(t *testing.T) {
	fake := &client.FakeClient{
		PostFunc: func(path string, body interface{}) ([]byte, error) {
			if strings.Contains(path, "/domains") && !strings.Contains(path, "/domains/") {
				return []byte(`{"id":"copy-id","name":"My Copy"}`), nil
			}
			return []byte(`{}`), nil
		},
		GetFunc: func(path string) ([]byte, error) {
			switch {
			case strings.HasSuffix(path, "/domains/src-id"):
				return []byte(`{"id":"src-id","name":"Source"}`), nil
			case strings.Contains(path, "/applications"):
				return []byte(`{"data":[],"currentPage":0,"totalCount":0}`), nil
			case strings.Contains(path, "/identityProviders"):
				return []byte(`[]`), nil
			case strings.Contains(path, "/roles"):
				return []byte(`{"data":[],"currentPage":0,"totalCount":0}`), nil
			case strings.Contains(path, "/scopes"):
				return []byte(`{"data":[],"currentPage":0,"totalCount":0}`), nil
			case strings.Contains(path, "/factors"):
				return []byte(`[]`), nil
			case strings.Contains(path, "/groups"):
				return []byte(`{"data":[],"currentPage":0,"totalCount":0}`), nil
			case strings.Contains(path, "/flows"):
				return []byte(`[]`), nil
			case strings.Contains(path, "/certificates"):
				return []byte(`[]`), nil
			}
			return []byte(`{}`), nil
		},
	}
	f, out := newTestFactory(fake, false)
	cmd := newCopyCmd(f)
	cmd.SetArgs([]string{"src-id", "--name", "My Copy"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "Copy complete") {
		t.Errorf("expected 'Copy complete' in output, got: %s", out.String())
	}
}

func TestDomainCopyReadOnly(t *testing.T) {
	f, _ := newTestFactory(&client.FakeClient{}, true)
	cmd := newCopyCmd(f)
	cmd.SetArgs([]string{"src-id", "--name", "My Copy"})
	if err := cmd.Execute(); err == nil {
		t.Error("expected read-only error")
	}
}
