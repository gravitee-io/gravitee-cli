package am

import (
	"encoding/json"
	"fmt"
)

// OrgFormService defines organization-level form operations.
type OrgFormService interface {
	GetOrgForm(template string) (json.RawMessage, error)
	CreateOrgForm(body json.RawMessage) (json.RawMessage, error)
	UpdateOrgForm(formID string, body json.RawMessage) (json.RawMessage, error)
	DeleteOrgForm(formID string) error
}

func (s *service) GetOrgForm(template string) (json.RawMessage, error) {
	data, err := s.client.Get(s.orgPath(fmt.Sprintf("forms?template=%s", template)))
	if err != nil {
		return nil, fmt.Errorf("org form get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateOrgForm(body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.orgPath("forms"), body)
	if err != nil {
		return nil, fmt.Errorf("org form create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateOrgForm(formID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.orgPath(fmt.Sprintf("forms/%s", formID)), body)
	if err != nil {
		return nil, fmt.Errorf("org form update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteOrgForm(formID string) error {
	err := s.client.Delete(s.orgPath(fmt.Sprintf("forms/%s", formID)))
	if err != nil {
		return fmt.Errorf("org form delete failed: %w", err)
	}

	return nil
}
