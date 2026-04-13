package am

import (
	"encoding/json"
	"fmt"
)

// OrgTagService defines organization-level sharding tag operations.
type OrgTagService interface {
	ListOrgTags() ([]json.RawMessage, error)
	GetOrgTag(tagID string) (json.RawMessage, error)
	CreateOrgTag(body json.RawMessage) (json.RawMessage, error)
	UpdateOrgTag(tagID string, body json.RawMessage) (json.RawMessage, error)
	DeleteOrgTag(tagID string) error
}

func (s *service) ListOrgTags() ([]json.RawMessage, error) {
	data, err := s.client.Get(s.orgPath("tags"))
	if err != nil {
		return nil, fmt.Errorf("org tag list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse org tag list: %w", err)
	}

	return items, nil
}

func (s *service) GetOrgTag(tagID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.orgPath(fmt.Sprintf("tags/%s", tagID)))
	if err != nil {
		return nil, fmt.Errorf("org tag get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateOrgTag(body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.orgPath("tags"), body)
	if err != nil {
		return nil, fmt.Errorf("org tag create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateOrgTag(tagID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.orgPath(fmt.Sprintf("tags/%s", tagID)), body)
	if err != nil {
		return nil, fmt.Errorf("org tag update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteOrgTag(tagID string) error {
	err := s.client.Delete(s.orgPath(fmt.Sprintf("tags/%s", tagID)))
	if err != nil {
		return fmt.Errorf("org tag delete failed: %w", err)
	}

	return nil
}
