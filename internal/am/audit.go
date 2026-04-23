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

package am

import (
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/client"
)

// ListAuditsParams holds parameters for listing audits.
type ListAuditsParams struct {
	Type    string
	Status  string
	From    string
	To      string
	Page    int
	PerPage int
}

// AuditService defines audit-related operations.
type AuditService interface {
	ListAudits(domainID string, params ListAuditsParams) (*PaginatedResponse, error)
	GetAudit(domainID, auditID string) (json.RawMessage, error)
}

func (s *service) ListAudits(domainID string, params ListAuditsParams) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{
		"page":   client.Itoa(params.Page),
		"size":   client.Itoa(params.PerPage),
		"type":   params.Type,
		"status": params.Status,
		"from":   params.From,
		"to":     params.To,
	})

	data, err := s.client.Get(s.domainPath(domainID, "audits?"+q))
	if err != nil {
		return nil, fmt.Errorf("audit list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) GetAudit(domainID, auditID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("audits/%s", auditID)))
	if err != nil {
		return nil, fmt.Errorf("audit get failed: %w", err)
	}

	return json.RawMessage(data), nil
}
