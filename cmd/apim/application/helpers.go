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

package application

import (
	"encoding/json"
	"fmt"
	"time"

	"gravitee.io/gctl/internal/printer"
)

func ownerDisplayName(item any) string {
	m, ok := item.(map[string]any)
	if !ok {
		return ""
	}

	owner, ok := m["owner"].(map[string]any)
	if !ok {
		return ""
	}

	s, _ := owner["displayName"].(string)

	return s
}

func printAppDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, f := range []struct{ label, key string }{
		{"Name", "name"},
		{"ID", "id"},
		{"Description", "description"},
		{"Type", "type"},
		{"Status", "status"},
	} {
		printField(p, m, f.label, f.key)
	}

	printOwner(p, m)

	for _, f := range []struct{ label, key string }{
		{"API Key Mode", "api_key_mode"},
		{"Domain", "domain"},
		{"Created", "created_at"},
		{"Updated", "updated_at"},
	} {
		printField(p, m, f.label, f.key)
	}

	return nil
}

func printField(p *printer.Printer, m map[string]any, label, key string) {
	v, ok := m[key]
	if !ok || v == nil {
		return
	}

	if f, isNum := v.(float64); isNum && (key == "created_at" || key == "updated_at") {
		p.PrintMessage("%-16s%s", label+":", time.UnixMilli(int64(f)).UTC().Format(time.RFC3339))

		return
	}

	p.PrintMessage("%-16s%v", label+":", v)
}

func printOwner(p *printer.Printer, m map[string]any) {
	owner, ok := m["owner"].(map[string]any)
	if !ok {
		return
	}

	if dn, ok := owner["displayName"].(string); ok && dn != "" {
		p.PrintMessage("%-16s%s", "Owner:", dn)
	}
}
