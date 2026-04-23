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

// ListApplicationsParams holds parameters for listing applications.
type ListApplicationsParams struct {
	Query   string
	Status  string
	Page    int
	PerPage int
}

// ApplicationService defines application-related operations (V1 API).
type ApplicationService interface {
	ListApplications(params ListApplicationsParams) (*PaginatedResponse, error)
	GetApplication(appID string) (json.RawMessage, error)
	CreateApplication(body json.RawMessage) (json.RawMessage, error)
	UpdateApplication(appID string, body json.RawMessage) (json.RawMessage, error)
	DeleteApplication(appID string) error
}

func (s *service) ListApplications(params ListApplicationsParams) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{
		"page": client.Itoa(params.Page), "size": client.Itoa(params.PerPage),
		"query": params.Query, "status": params.Status,
	})

	data, err := s.client.Get(s.v1(fmt.Sprintf("applications/_paged?%s", q)))
	if err != nil {
		return nil, fmt.Errorf("application list failed: %w", err)
	}

	// V1 wraps pagination in a "page" object with snake_case fields;
	// translate to the V2 shape the rest of the code expects.
	var v1 struct {
		Data []json.RawMessage `json:"data"`
		Page struct {
			Current       int `json:"current"`
			PerPage       int `json:"per_page"`
			TotalPages    int `json:"total_pages"`
			TotalElements int `json:"total_elements"`
		} `json:"page"`
	}
	if err := json.Unmarshal(data, &v1); err != nil {
		return nil, fmt.Errorf("failed to parse paginated response: %w", err)
	}

	return &PaginatedResponse{
		Data: v1.Data,
		Pagination: Pagination{
			Page:           v1.Page.Current,
			PerPage:        v1.Page.PerPage,
			PageCount:      v1.Page.TotalPages,
			TotalCount:     v1.Page.TotalElements,
			PageItemsCount: len(v1.Data),
		},
	}, nil
}

func (s *service) GetApplication(appID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.v1(fmt.Sprintf("applications/%s", appID)))
	if err != nil {
		return nil, err
	}

	return raw(data), nil
}

func (s *service) CreateApplication(body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.v1("applications"), body)
	if err != nil {
		return nil, fmt.Errorf("application creation failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) UpdateApplication(appID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.v1(fmt.Sprintf("applications/%s", appID)), body)
	if err != nil {
		return nil, fmt.Errorf("application update failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) DeleteApplication(appID string) error {
	if err := s.client.Delete(s.v1(fmt.Sprintf("applications/%s", appID))); err != nil {
		return fmt.Errorf("application deletion failed: %w", err)
	}

	return nil
}
