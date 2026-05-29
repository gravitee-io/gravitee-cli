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

package am

import (
	"encoding/json"
	"fmt"

	"gravitee.io/gctl/internal/client"
)

// IdentityProviderService defines identity provider-related operations.
type IdentityProviderService interface {
	ListIdentityProviders(domainID string, userProvider bool) ([]json.RawMessage, error)
	GetIdentityProvider(domainID, idpID string) (json.RawMessage, error)
	CreateIdentityProvider(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateIdentityProvider(domainID, idpID string, body json.RawMessage) (json.RawMessage, error)
	DeleteIdentityProvider(domainID, idpID string) error
	UpdateIDPPasswordPolicy(domainID, idpID string, body json.RawMessage) (json.RawMessage, error)
}

func (s *service) ListIdentityProviders(domainID string, userProvider bool) ([]json.RawMessage, error) {
	q := client.BuildQuery(map[string]string{
		"userProvider": fmt.Sprintf("%t", userProvider),
	})

	data, err := s.client.Get(s.domainPath(domainID, "identities?"+q))
	if err != nil {
		return nil, fmt.Errorf("identity provider list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse identity provider list: %w", err)
	}

	return items, nil
}

func (s *service) GetIdentityProvider(domainID, idpID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("identities/%s", idpID)))
	if err != nil {
		return nil, fmt.Errorf("identity provider get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateIdentityProvider(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "identities"), body)
	if err != nil {
		return nil, fmt.Errorf("identity provider create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateIdentityProvider(domainID, idpID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("identities/%s", idpID)), body)
	if err != nil {
		return nil, fmt.Errorf("identity provider update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteIdentityProvider(domainID, idpID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("identities/%s", idpID)))
	if err != nil {
		return fmt.Errorf("identity provider delete failed: %w", err)
	}

	return nil
}

func (s *service) UpdateIDPPasswordPolicy(domainID, idpID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("identities/%s/password-policy", idpID)), body)
	if err != nil {
		return nil, fmt.Errorf("idp password policy update failed: %w", err)
	}

	return json.RawMessage(data), nil
}
