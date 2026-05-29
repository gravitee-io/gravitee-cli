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

package plan

import (
	"encoding/json"
	"fmt"

	"gravitee.io/gctl/internal/printer"
)

func securityType(item any) string {
	m, ok := item.(map[string]any)
	if !ok {
		return ""
	}

	sec, ok := m["security"].(map[string]any)
	if !ok {
		return ""
	}

	s, _ := sec["type"].(string)

	return s
}

func printPlanDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	fields := []struct {
		label string
		key   string
	}{
		{"Name", "name"},
		{"ID", "id"},
		{"API", "apiId"},
		{"Status", "status"},
		{"Validation", "validation"},
		{"Mode", "mode"},
		{"Created", "createdAt"},
		{"Updated", "updatedAt"},
		{"Published", "publishedAt"},
		{"Closed", "closedAt"},
	}

	securityStr := ""
	if sec, ok := m["security"].(map[string]any); ok {
		if t, ok := sec["type"].(string); ok {
			securityStr = t
		}
	}

	for _, field := range fields {
		if field.key == "status" {
			if v, ok := m[field.key]; ok && v != nil {
				p.PrintMessage("%-16s%v", field.label+":", v)
			}

			if securityStr != "" {
				p.PrintMessage("%-16s%s", "Security:", securityStr)
			}

			continue
		}

		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	return nil
}
