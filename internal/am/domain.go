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

	"github.com/gravitee-io/gio-cli/internal/client"
)

// ListDomainsParams holds parameters for listing domains.
type ListDomainsParams struct {
	Query   string
	Page    int
	PerPage int
}

// DomainService defines domain-related operations.
type DomainService interface {
	ListDomains(params ListDomainsParams) (*PaginatedResponse, error)
	GetDomain(domainID string) (json.RawMessage, error)
	GetDomainByHRID(hrid string) (json.RawMessage, error)
	CreateDomain(body json.RawMessage) (json.RawMessage, error)
	UpdateDomain(domainID string, body json.RawMessage) (json.RawMessage, error)
	PatchDomain(domainID string, body json.RawMessage) (json.RawMessage, error)
	DeleteDomain(domainID string) error
	UpdateDomainCertificateSettings(domainID string, body json.RawMessage) (json.RawMessage, error)
	ListDataPlanes() (json.RawMessage, error)
}

func (s *service) ListDomains(params ListDomainsParams) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{
		"page": client.Itoa(params.Page),
		"size": client.Itoa(params.PerPage),
		"q":    params.Query,
	})

	data, err := s.client.Get(s.basePath("domains?" + q))
	if err != nil {
		return nil, fmt.Errorf("domain list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) GetDomain(domainID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.basePath(fmt.Sprintf("domains/%s", domainID)))
	if err != nil {
		return nil, err
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateDomain(body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.basePath("domains"), body)
	if err != nil {
		return nil, fmt.Errorf("domain create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateDomain(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.basePath(fmt.Sprintf("domains/%s", domainID)), body)
	if err != nil {
		return nil, fmt.Errorf("domain update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) PatchDomain(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Patch(s.basePath(fmt.Sprintf("domains/%s", domainID)), body)
	if err != nil {
		return nil, fmt.Errorf("domain patch failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteDomain(domainID string) error {
	err := s.client.Delete(s.basePath(fmt.Sprintf("domains/%s", domainID)))
	if err != nil {
		return fmt.Errorf("domain delete failed: %w", err)
	}

	return nil
}

func (s *service) GetDomainByHRID(hrid string) (json.RawMessage, error) {
	data, err := s.client.Get(s.basePath(fmt.Sprintf("domains/_hrid/%s", hrid)))
	if err != nil {
		return nil, fmt.Errorf("domain get by HRID failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateDomainCertificateSettings(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, "certificate-settings"), body)
	if err != nil {
		return nil, fmt.Errorf("domain certificate settings update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) ListDataPlanes() (json.RawMessage, error) {
	data, err := s.client.Get(s.basePath("data-planes"))
	if err != nil {
		return nil, fmt.Errorf("data plane list failed: %w", err)
	}

	return json.RawMessage(data), nil
}
