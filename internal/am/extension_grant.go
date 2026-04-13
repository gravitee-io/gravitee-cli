package am

import (
	"encoding/json"
	"fmt"
)

// ExtensionGrantService defines extension grant-related operations.
type ExtensionGrantService interface {
	ListExtensionGrants(domainID string) ([]json.RawMessage, error)
	GetExtensionGrant(domainID, grantID string) (json.RawMessage, error)
	CreateExtensionGrant(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateExtensionGrant(domainID, grantID string, body json.RawMessage) (json.RawMessage, error)
	DeleteExtensionGrant(domainID, grantID string) error
}

func (s *service) ListExtensionGrants(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "extensionGrants"))
	if err != nil {
		return nil, fmt.Errorf("extension grant list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse extension grant list: %w", err)
	}

	return items, nil
}

func (s *service) GetExtensionGrant(domainID, grantID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("extensionGrants/%s", grantID)))
	if err != nil {
		return nil, fmt.Errorf("extension grant get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateExtensionGrant(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "extensionGrants"), body)
	if err != nil {
		return nil, fmt.Errorf("extension grant create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateExtensionGrant(domainID, grantID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("extensionGrants/%s", grantID)), body)
	if err != nil {
		return nil, fmt.Errorf("extension grant update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteExtensionGrant(domainID, grantID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("extensionGrants/%s", grantID)))
	if err != nil {
		return fmt.Errorf("extension grant delete failed: %w", err)
	}

	return nil
}
