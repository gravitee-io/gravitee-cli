package am

import (
	"encoding/json"
	"fmt"
)

// OrgEntrypointService defines organization-level entrypoint operations.
type OrgEntrypointService interface {
	ListOrgEntrypoints() ([]json.RawMessage, error)
	GetOrgEntrypoint(entrypointID string) (json.RawMessage, error)
	CreateOrgEntrypoint(body json.RawMessage) (json.RawMessage, error)
	UpdateOrgEntrypoint(entrypointID string, body json.RawMessage) (json.RawMessage, error)
	DeleteOrgEntrypoint(entrypointID string) error
}

func (s *service) ListOrgEntrypoints() ([]json.RawMessage, error) {
	data, err := s.client.Get(s.orgPath("entrypoints"))
	if err != nil {
		return nil, fmt.Errorf("org entrypoint list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse org entrypoint list: %w", err)
	}

	return items, nil
}

func (s *service) GetOrgEntrypoint(entrypointID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.orgPath(fmt.Sprintf("entrypoints/%s", entrypointID)))
	if err != nil {
		return nil, fmt.Errorf("org entrypoint get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateOrgEntrypoint(body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.orgPath("entrypoints"), body)
	if err != nil {
		return nil, fmt.Errorf("org entrypoint create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateOrgEntrypoint(entrypointID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.orgPath(fmt.Sprintf("entrypoints/%s", entrypointID)), body)
	if err != nil {
		return nil, fmt.Errorf("org entrypoint update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteOrgEntrypoint(entrypointID string) error {
	err := s.client.Delete(s.orgPath(fmt.Sprintf("entrypoints/%s", entrypointID)))
	if err != nil {
		return fmt.Errorf("org entrypoint delete failed: %w", err)
	}

	return nil
}
