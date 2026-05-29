// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apim

import (
	"encoding/json"
	"fmt"

	"gravitee.io/gctl/internal/client"
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
	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/subscriptions/%s/api-keys/_renew", apiID, subID)), json.RawMessage(`{}`))
	if err != nil {
		return nil, fmt.Errorf("API key renew failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) RevokeAPIKey(apiID, subID, keyID string) error {
	if _, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/subscriptions/%s/api-keys/%s/_revoke", apiID, subID, keyID)), json.RawMessage(`{}`)); err != nil {
		return fmt.Errorf("API key revoke failed: %w", err)
	}

	return nil
}

func (s *service) ReactivateAPIKey(apiID, subID, keyID string) (json.RawMessage, error) {
	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/subscriptions/%s/api-keys/%s/_reactivate", apiID, subID, keyID)), nil)
	if err != nil {
		return nil, fmt.Errorf("API key reactivate failed: %w", err)
	}

	return raw(data), nil
}
