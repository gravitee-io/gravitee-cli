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

// EmailService defines email-related operations.
type EmailService interface {
	ListEmails(domainID string) ([]json.RawMessage, error)
	GetEmail(domainID, emailID string) (json.RawMessage, error)
	CreateEmail(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateEmail(domainID, emailID string, body json.RawMessage) (json.RawMessage, error)
	DeleteEmail(domainID, emailID string) error
}

func (s *service) ListEmails(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "emails"))
	if err != nil {
		return nil, fmt.Errorf("email list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse email list: %w", err)
	}

	return items, nil
}

func (s *service) GetEmail(domainID, emailID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("emails/%s", emailID)))
	if err != nil {
		return nil, fmt.Errorf("email get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateEmail(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "emails"), body)
	if err != nil {
		return nil, fmt.Errorf("email create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateEmail(domainID, emailID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("emails/%s", emailID)), body)
	if err != nil {
		return nil, fmt.Errorf("email update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteEmail(domainID, emailID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("emails/%s", emailID)))
	if err != nil {
		return fmt.Errorf("email delete failed: %w", err)
	}

	return nil
}
