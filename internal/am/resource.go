package am

import (
	"encoding/json"
	"fmt"
)

// ResourceService defines resource-related operations.
type ResourceService interface {
	ListResources(domainID string) ([]json.RawMessage, error)
	GetResource(domainID, resourceID string) (json.RawMessage, error)
	CreateResource(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateResource(domainID, resourceID string, body json.RawMessage) (json.RawMessage, error)
	DeleteResource(domainID, resourceID string) error
}

func (s *service) ListResources(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "resources"))
	if err != nil {
		return nil, fmt.Errorf("resource list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse resource list: %w", err)
	}

	return items, nil
}

func (s *service) GetResource(domainID, resourceID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("resources/%s", resourceID)))
	if err != nil {
		return nil, fmt.Errorf("resource get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateResource(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "resources"), body)
	if err != nil {
		return nil, fmt.Errorf("resource create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateResource(domainID, resourceID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("resources/%s", resourceID)), body)
	if err != nil {
		return nil, fmt.Errorf("resource update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteResource(domainID, resourceID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("resources/%s", resourceID)))
	if err != nil {
		return fmt.Errorf("resource delete failed: %w", err)
	}

	return nil
}
