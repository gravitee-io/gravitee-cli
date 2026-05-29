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

package subscription

import (
	"encoding/json"
	"fmt"

	"gravitee.io/gctl/internal/printer"
)

func extractID(m map[string]any, key string) string {
	if nested, ok := m[key].(map[string]any); ok {
		s, _ := nested["id"].(string)

		return s
	}

	s, _ := m[key+"Id"].(string)

	return s
}

func printSubDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	printSubField(p, m, "ID", "id")

	if plan := extractID(m, "plan"); plan != "" {
		p.PrintMessage("%-16s%s", "Plan:", plan)
	}

	if app := extractID(m, "application"); app != "" {
		p.PrintMessage("%-16s%s", "Application:", app)
	}

	for _, field := range []struct{ label, key string }{
		{"Status", "status"},
		{"Created", "createdAt"},
		{"Processed", "processedAt"},
		{"Starting", "startingAt"},
		{"Ending", "endingAt"},
		{"Closed", "closedAt"},
		{"Paused", "pausedAt"},
	} {
		printSubField(p, m, field.label, field.key)
	}

	return nil
}

func printSubCreateDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	printSubField(p, m, "ID", "id")

	if plan := extractID(m, "plan"); plan != "" {
		p.PrintMessage("%-16s%s", "Plan:", plan)
	}

	if app := extractID(m, "application"); app != "" {
		p.PrintMessage("%-16s%s", "Application:", app)
	}

	for _, field := range []struct{ label, key string }{
		{"Status", "status"},
		{"Created", "createdAt"},
	} {
		printSubField(p, m, field.label, field.key)
	}

	return nil
}

func printSubField(p *printer.Printer, m map[string]any, label, key string) {
	if v, ok := m[key]; ok && v != nil {
		p.PrintMessage("%-16s%v", label+":", v)
	}
}
