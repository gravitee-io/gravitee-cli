package am

import (
	"encoding/json"
	"fmt"
)

// FormService defines form-related operations.
type FormService interface {
	ListForms(domainID string) ([]json.RawMessage, error)
	GetForm(domainID, formID string) (json.RawMessage, error)
	CreateForm(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateForm(domainID, formID string, body json.RawMessage) (json.RawMessage, error)
	DeleteForm(domainID, formID string) error
}

func (s *service) ListForms(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "forms"))
	if err != nil {
		return nil, fmt.Errorf("form list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse form list: %w", err)
	}

	return items, nil
}

func (s *service) GetForm(domainID, formID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("forms/%s", formID)))
	if err != nil {
		return nil, fmt.Errorf("form get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateForm(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "forms"), body)
	if err != nil {
		return nil, fmt.Errorf("form create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateForm(domainID, formID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("forms/%s", formID)), body)
	if err != nil {
		return nil, fmt.Errorf("form update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteForm(domainID, formID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("forms/%s", formID)))
	if err != nil {
		return fmt.Errorf("form delete failed: %w", err)
	}

	return nil
}
