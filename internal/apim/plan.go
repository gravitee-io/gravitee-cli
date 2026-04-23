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

// ListPlansParams holds parameters for listing plans.
type ListPlansParams struct {
	Status   string
	Security string
	Page     int
	PerPage  int
}

// PlanService defines plan-related operations.
type PlanService interface {
	ListPlans(apiID string, params ListPlansParams) (*PaginatedResponse, error)
	GetPlan(apiID, planID string) (json.RawMessage, error)
	CreatePlan(apiID string, body json.RawMessage) (json.RawMessage, error)
	UpdatePlan(apiID, planID string, body json.RawMessage) (json.RawMessage, error)
	DeletePlan(apiID, planID string) error
	PublishPlan(apiID, planID string) (json.RawMessage, error)
	DeprecatePlan(apiID, planID string) (json.RawMessage, error)
	ClosePlan(apiID, planID string) (json.RawMessage, error)
}

func (s *service) ListPlans(apiID string, params ListPlansParams) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{
		"page": client.Itoa(params.Page), "perPage": client.Itoa(params.PerPage),
		"statuses": params.Status, "securities": params.Security,
	})

	data, err := s.client.Get(s.v2(fmt.Sprintf("apis/%s/plans?%s", apiID, q)))
	if err != nil {
		return nil, fmt.Errorf("plan list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) GetPlan(apiID, planID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.v2(fmt.Sprintf("apis/%s/plans/%s", apiID, planID)))
	if err != nil {
		return nil, err
	}

	return raw(data), nil
}

func (s *service) CreatePlan(apiID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/plans", apiID)), body)
	if err != nil {
		return nil, fmt.Errorf("plan creation failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) UpdatePlan(apiID, planID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.v2(fmt.Sprintf("apis/%s/plans/%s", apiID, planID)), body)
	if err != nil {
		return nil, fmt.Errorf("plan update failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) DeletePlan(apiID, planID string) error {
	if err := s.client.Delete(s.v2(fmt.Sprintf("apis/%s/plans/%s", apiID, planID))); err != nil {
		return fmt.Errorf("plan deletion failed: %w", err)
	}

	return nil
}

func (s *service) PublishPlan(apiID, planID string) (json.RawMessage, error) {
	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/plans/%s/_publish", apiID, planID)), nil)
	if err != nil {
		return nil, fmt.Errorf("plan publish failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) DeprecatePlan(apiID, planID string) (json.RawMessage, error) {
	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/plans/%s/_deprecate", apiID, planID)), nil)
	if err != nil {
		return nil, fmt.Errorf("plan deprecate failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) ClosePlan(apiID, planID string) (json.RawMessage, error) {
	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/plans/%s/_close", apiID, planID)), nil)
	if err != nil {
		return nil, fmt.Errorf("plan close failed: %w", err)
	}

	return raw(data), nil
}
