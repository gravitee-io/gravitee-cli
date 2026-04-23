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

// FactorService defines factor-related operations.
type FactorService interface {
	ListFactors(domainID string) ([]json.RawMessage, error)
	GetFactor(domainID, factorID string) (json.RawMessage, error)
	CreateFactor(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateFactor(domainID, factorID string, body json.RawMessage) (json.RawMessage, error)
	DeleteFactor(domainID, factorID string) error
}

func (s *service) ListFactors(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "factors"))
	if err != nil {
		return nil, fmt.Errorf("factor list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse factor list: %w", err)
	}

	return items, nil
}

func (s *service) GetFactor(domainID, factorID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("factors/%s", factorID)))
	if err != nil {
		return nil, fmt.Errorf("factor get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateFactor(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "factors"), body)
	if err != nil {
		return nil, fmt.Errorf("factor create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateFactor(domainID, factorID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("factors/%s", factorID)), body)
	if err != nil {
		return nil, fmt.Errorf("factor update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteFactor(domainID, factorID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("factors/%s", factorID)))
	if err != nil {
		return fmt.Errorf("factor delete failed: %w", err)
	}

	return nil
}
