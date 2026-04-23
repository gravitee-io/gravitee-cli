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

// FormService defines form-related operations.
type FormService interface {
	ListForms(domainID string) ([]json.RawMessage, error)
	GetForm(domainID, formID string) (json.RawMessage, error)
	CreateForm(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateForm(domainID, formID string, body json.RawMessage) (json.RawMessage, error)
	DeleteForm(domainID, formID string) error
}

func (s *service) ListForms(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "forms"))
	if err != nil {
		return nil, fmt.Errorf("form list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse form list: %w", err)
	}

	return items, nil
}

func (s *service) GetForm(domainID, formID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("forms/%s", formID)))
	if err != nil {
		return nil, fmt.Errorf("form get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateForm(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "forms"), body)
	if err != nil {
		return nil, fmt.Errorf("form create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateForm(domainID, formID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("forms/%s", formID)), body)
	if err != nil {
		return nil, fmt.Errorf("form update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteForm(domainID, formID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("forms/%s", formID)))
	if err != nil {
		return fmt.Errorf("form delete failed: %w", err)
	}

	return nil
}
