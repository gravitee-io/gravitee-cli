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

// ReporterService defines reporter-related operations.
type ReporterService interface {
	ListReporters(domainID string) ([]json.RawMessage, error)
	GetReporter(domainID, reporterID string) (json.RawMessage, error)
	CreateReporter(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateReporter(domainID, reporterID string, body json.RawMessage) (json.RawMessage, error)
	DeleteReporter(domainID, reporterID string) error
}

func (s *service) ListReporters(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "reporters"))
	if err != nil {
		return nil, fmt.Errorf("reporter list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse reporter list: %w", err)
	}

	return items, nil
}

func (s *service) GetReporter(domainID, reporterID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("reporters/%s", reporterID)))
	if err != nil {
		return nil, fmt.Errorf("reporter get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateReporter(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "reporters"), body)
	if err != nil {
		return nil, fmt.Errorf("reporter create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateReporter(domainID, reporterID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("reporters/%s", reporterID)), body)
	if err != nil {
		return nil, fmt.Errorf("reporter update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteReporter(domainID, reporterID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("reporters/%s", reporterID)))
	if err != nil {
		return fmt.Errorf("reporter delete failed: %w", err)
	}

	return nil
}
