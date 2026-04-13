package am

import (
	"encoding/json"
	"fmt"
)

// OrgUserTokenService defines organization-level user token (account access token) operations.
type OrgUserTokenService interface {
	ListOrgUserTokens(userID string) (json.RawMessage, error)
	CreateOrgUserToken(userID string, body json.RawMessage) (json.RawMessage, error)
	RevokeOrgUserToken(userID, tokenID string) error
}

func (s *service) ListOrgUserTokens(userID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.orgPath(fmt.Sprintf("users/%s/tokens", userID)))
	if err != nil {
		return nil, fmt.Errorf("org user token list failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateOrgUserToken(userID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.orgPath(fmt.Sprintf("users/%s/tokens", userID)), body)
	if err != nil {
		return nil, fmt.Errorf("org user token create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) RevokeOrgUserToken(userID, tokenID string) error {
	err := s.client.Delete(s.orgPath(fmt.Sprintf("users/%s/tokens/%s", userID, tokenID)))
	if err != nil {
		return fmt.Errorf("org user token revoke failed: %w", err)
	}

	return nil
}
