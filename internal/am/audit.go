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
