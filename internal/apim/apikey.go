package apim

import (
	"encoding/json"
	"fmt"
	"github.com/gravitee-io/gio-cli/internal/client"
)

// APIKeyService defines API key operations.
type APIKeyService interface {
	ListAPIKeys(apiID, subID string, page, perPage int) (*PaginatedResponse, error)
	RenewAPIKey(apiID, subID string) (json.RawMessage, error)
	RevokeAPIKey(apiID, subID, keyID string) error
	ReactivateAPIKey(apiID, subID, keyID string) (json.RawMessage, error)
}

func (s *service) ListAPIKeys(apiID, subID string, page, perPage int) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{"page": client.Itoa(page), "perPage": client.Itoa(perPage)})

	data, err := s.client.Get(s.v2(fmt.Sprintf("apis/%s/subscriptions/%s/api-keys?%s", apiID, subID, q)))
	if err != nil {
		return nil, fmt.Errorf("API key list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) RenewAPIKey(apiID, subID string) (json.RawMessage, error) {
	if err := s.requireWrite(); err != nil {
		return nil, err
	}

	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/subscriptions/%s/api-keys/_renew", apiID, subID)), nil)
	if err != nil {
		return nil, fmt.Errorf("API key renew failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) RevokeAPIKey(apiID, subID, keyID string) error {
	if err := s.requireWrite(); err != nil {
		return err
	}

	if _, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/subscriptions/%s/api-keys/%s/_revoke", apiID, subID, keyID)), nil); err != nil {
		return fmt.Errorf("API key revoke failed: %w", err)
	}

	return nil
}

func (s *service) ReactivateAPIKey(apiID, subID, keyID string) (json.RawMessage, error) {
	if err := s.requireWrite(); err != nil {
		return nil, err
	}

	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/subscriptions/%s/api-keys/%s/_reactivate", apiID, subID, keyID)), nil)
	if err != nil {
		return nil, fmt.Errorf("API key reactivate failed: %w", err)
	}

	return raw(data), nil
}
