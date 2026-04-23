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

// AuthorizationEngineService defines authorization engine-related operations.
// Note: This resource only supports List, Get, and Update (no Create or Delete).
type AuthorizationEngineService interface {
	ListAuthorizationEngines(domainID string) ([]json.RawMessage, error)
	GetAuthorizationEngine(domainID, engineID string) (json.RawMessage, error)
	UpdateAuthorizationEngine(domainID, engineID string, body json.RawMessage) (json.RawMessage, error)
}

func (s *service) ListAuthorizationEngines(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "authorization-engines"))
	if err != nil {
		return nil, fmt.Errorf("authorization engine list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse authorization engine list: %w", err)
	}

	return items, nil
}

func (s *service) GetAuthorizationEngine(domainID, engineID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("authorization-engines/%s", engineID)))
	if err != nil {
		return nil, fmt.Errorf("authorization engine get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateAuthorizationEngine(domainID, engineID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("authorization-engines/%s", engineID)), body)
	if err != nil {
		return nil, fmt.Errorf("authorization engine update failed: %w", err)
	}

	return json.RawMessage(data), nil
}
