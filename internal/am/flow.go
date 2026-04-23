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

// FlowService defines flow-related operations.
type FlowService interface {
	ListFlows(domainID string) ([]json.RawMessage, error)
	GetFlow(domainID, flowID string) (json.RawMessage, error)
	UpdateFlows(domainID string, body json.RawMessage) (json.RawMessage, error)
}

func (s *service) ListFlows(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "flows"))
	if err != nil {
		return nil, fmt.Errorf("flow list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse flow list: %w", err)
	}

	return items, nil
}

func (s *service) GetFlow(domainID, flowID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("flows/%s", flowID)))
	if err != nil {
		return nil, fmt.Errorf("flow get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateFlows(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, "flows"), body)
	if err != nil {
		return nil, fmt.Errorf("flow update failed: %w", err)
	}

	return json.RawMessage(data), nil
}
