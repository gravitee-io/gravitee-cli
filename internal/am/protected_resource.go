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
)

// ProtectedResourceService defines protected resource-related operations.
type ProtectedResourceService interface {
	ListProtectedResources(domainID string) ([]json.RawMessage, error)
	GetProtectedResource(domainID, protectedResourceID string) (json.RawMessage, error)
	CreateProtectedResource(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateProtectedResource(domainID, protectedResourceID string, body json.RawMessage) (json.RawMessage, error)
	DeleteProtectedResource(domainID, protectedResourceID string) error

	// Protected resource sub-resources
	ListProtectedResourceMembers(domainID, prID string) (json.RawMessage, error)
	RemoveProtectedResourceMember(domainID, prID, memberID string) error
	ListProtectedResourceSecrets(domainID, prID string) (json.RawMessage, error)
}

func (s *service) ListProtectedResources(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "protected-resources"))
	if err != nil {
		return nil, fmt.Errorf("protected resource list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse protected resource list: %w", err)
	}

	return items, nil
}

func (s *service) GetProtectedResource(domainID, protectedResourceID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("protected-resources/%s", protectedResourceID)))
	if err != nil {
		return nil, fmt.Errorf("protected resource get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateProtectedResource(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "protected-resources"), body)
	if err != nil {
		return nil, fmt.Errorf("protected resource create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateProtectedResource(domainID, protectedResourceID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("protected-resources/%s", protectedResourceID)), body)
	if err != nil {
		return nil, fmt.Errorf("protected resource update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteProtectedResource(domainID, protectedResourceID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("protected-resources/%s", protectedResourceID)))
	if err != nil {
		return fmt.Errorf("protected resource delete failed: %w", err)
	}

	return nil
}

// Protected resource sub-resources

func (s *service) ListProtectedResourceMembers(domainID, prID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("protected-resources/%s/members", prID)))
	if err != nil {
		return nil, fmt.Errorf("protected resource member list failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) RemoveProtectedResourceMember(domainID, prID, memberID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("protected-resources/%s/members/%s", prID, memberID)))
	if err != nil {
		return fmt.Errorf("protected resource member remove failed: %w", err)
	}

	return nil
}

func (s *service) ListProtectedResourceSecrets(domainID, prID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("protected-resources/%s/secrets", prID)))
	if err != nil {
		return nil, fmt.Errorf("protected resource secret list failed: %w", err)
	}

	return json.RawMessage(data), nil
}
