package am

import (
	"encoding/json"
	"fmt"
)

// OrgReporterService defines organization-level reporter operations.
type OrgReporterService interface {
	ListOrgReporters() ([]json.RawMessage, error)
	GetOrgReporter(reporterID string) (json.RawMessage, error)
	CreateOrgReporter(body json.RawMessage) (json.RawMessage, error)
	UpdateOrgReporter(reporterID string, body json.RawMessage) (json.RawMessage, error)
	DeleteOrgReporter(reporterID string) error
}

func (s *service) ListOrgReporters() ([]json.RawMessage, error) {
	data, err := s.client.Get(s.orgPath("reporters"))
	if err != nil {
		return nil, fmt.Errorf("org reporter list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse org reporter list: %w", err)
	}

	return items, nil
}

func (s *service) GetOrgReporter(reporterID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.orgPath(fmt.Sprintf("reporters/%s", reporterID)))
	if err != nil {
		return nil, fmt.Errorf("org reporter get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateOrgReporter(body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.orgPath("reporters"), body)
	if err != nil {
		return nil, fmt.Errorf("org reporter create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateOrgReporter(reporterID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.orgPath(fmt.Sprintf("reporters/%s", reporterID)), body)
	if err != nil {
		return nil, fmt.Errorf("org reporter update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteOrgReporter(reporterID string) error {
	err := s.client.Delete(s.orgPath(fmt.Sprintf("reporters/%s", reporterID)))
	if err != nil {
		return fmt.Errorf("org reporter delete failed: %w", err)
	}

	return nil
}
