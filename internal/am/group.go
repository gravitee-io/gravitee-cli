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

	"gravitee.io/gctl/internal/client"
)

// ListGroupsParams holds parameters for listing groups.
type ListGroupsParams struct {
	Query   string
	Page    int
	PerPage int
}

// GroupService defines group-related operations.
type GroupService interface {
	ListGroups(domainID string, params ListGroupsParams) (*PaginatedResponse, error)
	GetGroup(domainID, groupID string) (json.RawMessage, error)
	CreateGroup(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateGroup(domainID, groupID string, body json.RawMessage) (json.RawMessage, error)
	DeleteGroup(domainID, groupID string) error

	// Group members
	ListGroupMembers(domainID, groupID string) (json.RawMessage, error)
	AddGroupMember(domainID, groupID, memberID string) error
	RemoveGroupMember(domainID, groupID, memberID string) error

	// Group roles
	ListGroupRoles(domainID, groupID string) (json.RawMessage, error)
	AssignGroupRoles(domainID, groupID string, body json.RawMessage) (json.RawMessage, error)
	RevokeGroupRole(domainID, groupID, roleID string) error
}

func (s *service) ListGroups(domainID string, params ListGroupsParams) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{
		"page": client.Itoa(params.Page),
		"size": client.Itoa(params.PerPage),
		"q":    params.Query,
	})

	data, err := s.client.Get(s.domainPath(domainID, "groups?"+q))
	if err != nil {
		return nil, fmt.Errorf("group list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) GetGroup(domainID, groupID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("groups/%s", groupID)))
	if err != nil {
		return nil, fmt.Errorf("group get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateGroup(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "groups"), body)
	if err != nil {
		return nil, fmt.Errorf("group create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateGroup(domainID, groupID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("groups/%s", groupID)), body)
	if err != nil {
		return nil, fmt.Errorf("group update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteGroup(domainID, groupID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("groups/%s", groupID)))
	if err != nil {
		return fmt.Errorf("group delete failed: %w", err)
	}

	return nil
}

// Group members

func (s *service) ListGroupMembers(domainID, groupID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("groups/%s/members", groupID)))
	if err != nil {
		return nil, fmt.Errorf("group member list failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) AddGroupMember(domainID, groupID, memberID string) error {
	_, err := s.client.Post(s.domainPath(domainID, fmt.Sprintf("groups/%s/members/%s", groupID, memberID)), nil)
	if err != nil {
		return fmt.Errorf("group member add failed: %w", err)
	}

	return nil
}

func (s *service) RemoveGroupMember(domainID, groupID, memberID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("groups/%s/members/%s", groupID, memberID)))
	if err != nil {
		return fmt.Errorf("group member remove failed: %w", err)
	}

	return nil
}

// Group roles

func (s *service) ListGroupRoles(domainID, groupID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("groups/%s/roles", groupID)))
	if err != nil {
		return nil, fmt.Errorf("group role list failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) AssignGroupRoles(domainID, groupID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, fmt.Sprintf("groups/%s/roles", groupID)), body)
	if err != nil {
		return nil, fmt.Errorf("group role assign failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) RevokeGroupRole(domainID, groupID, roleID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("groups/%s/roles/%s", groupID, roleID)))
	if err != nil {
		return fmt.Errorf("group role revoke failed: %w", err)
	}

	return nil
}
