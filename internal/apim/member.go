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

package apim

import (
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/client"
)

// MemberService defines member-related operations.
type MemberService interface {
	ListMembers(apiID string, page, perPage int) (*PaginatedResponse, error)
	AddMember(apiID, userID, role string) (json.RawMessage, error)
	RemoveMember(apiID, memberID string) error
}

func (s *service) ListMembers(apiID string, page, perPage int) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{"page": client.Itoa(page), "perPage": client.Itoa(perPage)})

	data, err := s.client.Get(s.v2(fmt.Sprintf("apis/%s/members?%s", apiID, q)))
	if err != nil {
		return nil, fmt.Errorf("member list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) AddMember(apiID, userID, role string) (json.RawMessage, error) {
	body := map[string]string{"userId": userID, "roleName": role}

	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/members", apiID)), body)
	if err != nil {
		return nil, fmt.Errorf("member add failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) RemoveMember(apiID, memberID string) error {
	if err := s.client.Delete(s.v2(fmt.Sprintf("apis/%s/members/%s", apiID, memberID))); err != nil {
		return fmt.Errorf("member removal failed: %w", err)
	}

	return nil
}
