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

// EntrypointService defines entrypoint-related operations.
type EntrypointService interface {
	GetEntrypoints(domainID string) (json.RawMessage, error)
	CreateEntrypoint(body json.RawMessage) (json.RawMessage, error)
	UpdateEntrypoint(entrypointID string, body json.RawMessage) (json.RawMessage, error)
	DeleteEntrypoint(entrypointID string) error
}

func (s *service) GetEntrypoints(domainID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "entrypoints"))
	if err != nil {
		return nil, fmt.Errorf("entrypoint get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateEntrypoint(body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.orgPath("entrypoints"), body)
	if err != nil {
		return nil, fmt.Errorf("entrypoint create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateEntrypoint(entrypointID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.orgPath(fmt.Sprintf("entrypoints/%s", entrypointID)), body)
	if err != nil {
		return nil, fmt.Errorf("entrypoint update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteEntrypoint(entrypointID string) error {
	err := s.client.Delete(s.orgPath(fmt.Sprintf("entrypoints/%s", entrypointID)))
	if err != nil {
		return fmt.Errorf("entrypoint delete failed: %w", err)
	}

	return nil
}
