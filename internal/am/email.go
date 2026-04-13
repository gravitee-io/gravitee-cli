package am

import (
	"encoding/json"
	"fmt"
)

// EmailService defines email-related operations.
type EmailService interface {
	ListEmails(domainID string) ([]json.RawMessage, error)
	GetEmail(domainID, emailID string) (json.RawMessage, error)
	CreateEmail(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateEmail(domainID, emailID string, body json.RawMessage) (json.RawMessage, error)
	DeleteEmail(domainID, emailID string) error
}

func (s *service) ListEmails(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "emails"))
	if err != nil {
		return nil, fmt.Errorf("email list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse email list: %w", err)
	}

	return items, nil
}

func (s *service) GetEmail(domainID, emailID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("emails/%s", emailID)))
	if err != nil {
		return nil, fmt.Errorf("email get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateEmail(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "emails"), body)
	if err != nil {
		return nil, fmt.Errorf("email create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateEmail(domainID, emailID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("emails/%s", emailID)), body)
	if err != nil {
		return nil, fmt.Errorf("email update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteEmail(domainID, emailID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("emails/%s", emailID)))
	if err != nil {
		return fmt.Errorf("email delete failed: %w", err)
	}

	return nil
}
