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

// MemberService defines member-related operations.
type MemberService interface {
	ListMembers(domainID string) (json.RawMessage, error)
	AddMember(domainID string, body json.RawMessage) (json.RawMessage, error)
	RemoveMember(domainID, memberID string) error
	GetMemberPermissions(domainID string) (json.RawMessage, error)
}

func (s *service) ListMembers(domainID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "members"))
	if err != nil {
		return nil, fmt.Errorf("member list failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) AddMember(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "members"), body)
	if err != nil {
		return nil, fmt.Errorf("member add failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) RemoveMember(domainID, memberID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("members/%s", memberID)))
	if err != nil {
		return fmt.Errorf("member remove failed: %w", err)
	}

	return nil
}

func (s *service) GetMemberPermissions(domainID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "members/permissions"))
	if err != nil {
		return nil, fmt.Errorf("member permissions get failed: %w", err)
	}

	return json.RawMessage(data), nil
}
