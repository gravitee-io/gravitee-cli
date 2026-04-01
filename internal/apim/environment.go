package apim

import (
	"encoding/json"
	"fmt"
)

// EnvironmentService defines environment-related operations.
type EnvironmentService interface {
	ListEnvironments() (json.RawMessage, error)
	GetEnvironment(envID string) (json.RawMessage, error)
}

func (s *service) ListEnvironments() (json.RawMessage, error) {
	path := fmt.Sprintf("/management/organizations/%s/environments", s.resolved.Org)

	data, err := s.client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("environment list failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) GetEnvironment(envID string) (json.RawMessage, error) {
	path := fmt.Sprintf("/management/organizations/%s/environments/%s", s.resolved.Org, envID)

	data, err := s.client.Get(path)
	if err != nil {
		return nil, err
	}

	return raw(data), nil
}
