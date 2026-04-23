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

// ResourceService defines resource-related operations.
type ResourceService interface {
	ListResources(domainID string) ([]json.RawMessage, error)
	GetResource(domainID, resourceID string) (json.RawMessage, error)
	CreateResource(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateResource(domainID, resourceID string, body json.RawMessage) (json.RawMessage, error)
	DeleteResource(domainID, resourceID string) error
}

func (s *service) ListResources(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "resources"))
	if err != nil {
		return nil, fmt.Errorf("resource list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse resource list: %w", err)
	}

	return items, nil
}

func (s *service) GetResource(domainID, resourceID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("resources/%s", resourceID)))
	if err != nil {
		return nil, fmt.Errorf("resource get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateResource(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "resources"), body)
	if err != nil {
		return nil, fmt.Errorf("resource create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateResource(domainID, resourceID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("resources/%s", resourceID)), body)
	if err != nil {
		return nil, fmt.Errorf("resource update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteResource(domainID, resourceID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("resources/%s", resourceID)))
	if err != nil {
		return fmt.Errorf("resource delete failed: %w", err)
	}

	return nil
}
