package am

import (
	"encoding/json"
	"fmt"
)

// MemberService defines member-related operations.
type MemberService interface {
	ListMembers(domainID string) (json.RawMessage, error)
	AddMember(domainID string, body json.RawMessage) (json.RawMessage, error)
	RemoveMember(domainID, memberID string) error
	GetMemberPermissions(domainID string) (json.RawMessage, error)
}

func (s *service) ListMembers(domainID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "members"))
	if err != nil {
		return nil, fmt.Errorf("member list failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) AddMember(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "members"), body)
	if err != nil {
		return nil, fmt.Errorf("member add failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) RemoveMember(domainID, memberID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("members/%s", memberID)))
	if err != nil {
		return fmt.Errorf("member remove failed: %w", err)
	}

	return nil
}

func (s *service) GetMemberPermissions(domainID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "members/permissions"))
	if err != nil {
		return nil, fmt.Errorf("member permissions get failed: %w", err)
	}

	return json.RawMessage(data), nil
}
