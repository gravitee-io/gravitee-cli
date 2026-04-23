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

package page

import (
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/printer"
)

func boolField(item any, key string) string {
	m, ok := item.(map[string]any)
	if !ok {
		return ""
	}

	v, ok := m[key]
	if !ok {
		return ""
	}

	b, ok := v.(bool)
	if !ok {
		return ""
	}

	return fmt.Sprintf("%t", b)
}

func printPageDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, field := range []struct{ label, key string }{
		{"Name", "name"},
		{"ID", "id"},
		{"API", "apiId"},
		{"Type", "type"},
		{"Visibility", "visibility"},
		{"Published", "published"},
		{"Parent", "parentId"},
		{"Created", "createdAt"},
		{"Updated", "updatedAt"},
	} {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	return nil
}
