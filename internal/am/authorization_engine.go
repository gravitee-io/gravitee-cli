package am

import (
	"encoding/json"
	"fmt"
)

// AuthorizationEngineService defines authorization engine-related operations.
// Note: This resource only supports List, Get, and Update (no Create or Delete).
type AuthorizationEngineService interface {
	ListAuthorizationEngines(domainID string) ([]json.RawMessage, error)
	GetAuthorizationEngine(domainID, engineID string) (json.RawMessage, error)
	UpdateAuthorizationEngine(domainID, engineID string, body json.RawMessage) (json.RawMessage, error)
}

func (s *service) ListAuthorizationEngines(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "authorization-engines"))
	if err != nil {
		return nil, fmt.Errorf("authorization engine list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse authorization engine list: %w", err)
	}

	return items, nil
}

func (s *service) GetAuthorizationEngine(domainID, engineID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("authorization-engines/%s", engineID)))
	if err != nil {
		return nil, fmt.Errorf("authorization engine get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateAuthorizationEngine(domainID, engineID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("authorization-engines/%s", engineID)), body)
	if err != nil {
		return nil, fmt.Errorf("authorization engine update failed: %w", err)
	}

	return json.RawMessage(data), nil
}
