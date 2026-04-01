package apim

import (
	"encoding/json"
	"fmt"
	"github.com/gravitee-io/gio-cli/internal/client"
)

// MemberService defines member-related operations.
type MemberService interface {
	ListMembers(apiID string, page, perPage int) (*PaginatedResponse, error)
	AddMember(apiID, userID, role string) (json.RawMessage, error)
	RemoveMember(apiID, memberID string) error
}

func (s *service) ListMembers(apiID string, page, perPage int) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{"page": client.Itoa(page), "perPage": client.Itoa(perPage)})

	data, err := s.client.Get(s.v2(fmt.Sprintf("apis/%s/members?%s", apiID, q)))
	if err != nil {
		return nil, fmt.Errorf("member list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) AddMember(apiID, userID, role string) (json.RawMessage, error) {
	if err := s.requireWrite(); err != nil {
		return nil, err
	}

	body := map[string]string{"userId": userID, "roleName": role}

	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/members", apiID)), body)
	if err != nil {
		return nil, fmt.Errorf("member add failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) RemoveMember(apiID, memberID string) error {
	if err := s.requireWrite(); err != nil {
		return err
	}

	if err := s.client.Delete(s.v2(fmt.Sprintf("apis/%s/members/%s", apiID, memberID))); err != nil {
		return fmt.Errorf("member removal failed: %w", err)
	}

	return nil
}
