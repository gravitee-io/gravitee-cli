package am

import (
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/client"
)

// ListDomainsParams holds parameters for listing domains.
type ListDomainsParams struct {
	Query   string
	Page    int
	PerPage int
}

// DomainService defines domain-related operations.
type DomainService interface {
	ListDomains(params ListDomainsParams) (*PaginatedResponse, error)
	GetDomain(domainID string) (json.RawMessage, error)
}

func (s *service) ListDomains(params ListDomainsParams) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{
		"page": client.Itoa(params.Page),
		"size": client.Itoa(params.PerPage),
		"q":    params.Query,
	})

	data, err := s.client.Get(s.basePath("domains?" + q))
	if err != nil {
		return nil, fmt.Errorf("domain list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) GetDomain(domainID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.basePath(fmt.Sprintf("domains/%s", domainID)))
	if err != nil {
		return nil, err
	}

	return json.RawMessage(data), nil
}
