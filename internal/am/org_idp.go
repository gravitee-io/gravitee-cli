package am

import (
	"encoding/json"
	"fmt"
)

// OrgIdentityProviderService defines organization-level identity provider operations.
type OrgIdentityProviderService interface {
	ListOrgIdentityProviders() ([]json.RawMessage, error)
	GetOrgIdentityProvider(idpID string) (json.RawMessage, error)
	CreateOrgIdentityProvider(body json.RawMessage) (json.RawMessage, error)
	UpdateOrgIdentityProvider(idpID string, body json.RawMessage) (json.RawMessage, error)
	DeleteOrgIdentityProvider(idpID string) error
}

func (s *service) ListOrgIdentityProviders() ([]json.RawMessage, error) {
	data, err := s.client.Get(s.orgPath("identities"))
	if err != nil {
		return nil, fmt.Errorf("org identity provider list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse org identity provider list: %w", err)
	}

	return items, nil
}

func (s *service) GetOrgIdentityProvider(idpID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.orgPath(fmt.Sprintf("identities/%s", idpID)))
	if err != nil {
		return nil, fmt.Errorf("org identity provider get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateOrgIdentityProvider(body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.orgPath("identities"), body)
	if err != nil {
		return nil, fmt.Errorf("org identity provider create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateOrgIdentityProvider(idpID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.orgPath(fmt.Sprintf("identities/%s", idpID)), body)
	if err != nil {
		return nil, fmt.Errorf("org identity provider update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteOrgIdentityProvider(idpID string) error {
	err := s.client.Delete(s.orgPath(fmt.Sprintf("identities/%s", idpID)))
	if err != nil {
		return fmt.Errorf("org identity provider delete failed: %w", err)
	}

	return nil
}
