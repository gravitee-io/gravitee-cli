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

// OrgFormService defines organization-level form operations.
type OrgFormService interface {
	GetOrgForm(template string) (json.RawMessage, error)
	CreateOrgForm(body json.RawMessage) (json.RawMessage, error)
	UpdateOrgForm(formID string, body json.RawMessage) (json.RawMessage, error)
	DeleteOrgForm(formID string) error
}

func (s *service) GetOrgForm(template string) (json.RawMessage, error) {
	data, err := s.client.Get(s.orgPath(fmt.Sprintf("forms?template=%s", template)))
	if err != nil {
		return nil, fmt.Errorf("org form get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateOrgForm(body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.orgPath("forms"), body)
	if err != nil {
		return nil, fmt.Errorf("org form create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateOrgForm(formID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.orgPath(fmt.Sprintf("forms/%s", formID)), body)
	if err != nil {
		return nil, fmt.Errorf("org form update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteOrgForm(formID string) error {
	err := s.client.Delete(s.orgPath(fmt.Sprintf("forms/%s", formID)))
	if err != nil {
		return fmt.Errorf("org form delete failed: %w", err)
	}

	return nil
}
