package am

import (
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/client"
)

// ListScopesParams holds parameters for listing scopes.
type ListScopesParams struct {
	Query   string
	Page    int
	PerPage int
}

// ScopeService defines scope-related operations.
type ScopeService interface {
	ListScopes(domainID string, params ListScopesParams) (*PaginatedResponse, error)
	GetScope(domainID, scopeID string) (json.RawMessage, error)
	CreateScope(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateScope(domainID, scopeID string, body json.RawMessage) (json.RawMessage, error)
	PatchScope(domainID, scopeID string, body json.RawMessage) (json.RawMessage, error)
	DeleteScope(domainID, scopeID string) error
}

func (s *service) ListScopes(domainID string, params ListScopesParams) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{
		"page": client.Itoa(params.Page),
		"size": client.Itoa(params.PerPage),
		"q":    params.Query,
	})

	data, err := s.client.Get(s.domainPath(domainID, "scopes?"+q))
	if err != nil {
		return nil, fmt.Errorf("scope list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) GetScope(domainID, scopeID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("scopes/%s", scopeID)))
	if err != nil {
		return nil, fmt.Errorf("scope get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateScope(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "scopes"), body)
	if err != nil {
		return nil, fmt.Errorf("scope create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateScope(domainID, scopeID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("scopes/%s", scopeID)), body)
	if err != nil {
		return nil, fmt.Errorf("scope update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) PatchScope(domainID, scopeID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Patch(s.domainPath(domainID, fmt.Sprintf("scopes/%s", scopeID)), body)
	if err != nil {
		return nil, fmt.Errorf("scope patch failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteScope(domainID, scopeID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("scopes/%s", scopeID)))
	if err != nil {
		return fmt.Errorf("scope delete failed: %w", err)
	}

	return nil
}
