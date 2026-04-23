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

// PageService defines page-related operations.
type PageService interface {
	ListPages(apiID, parentID string) ([]json.RawMessage, error)
	GetPage(apiID, pageID string) (json.RawMessage, error)
	CreatePage(apiID string, body json.RawMessage) (json.RawMessage, error)
	UpdatePage(apiID, pageID string, body json.RawMessage) (json.RawMessage, error)
	DeletePage(apiID, pageID string) error
	PublishPage(apiID, pageID string) (json.RawMessage, error)
	UnpublishPage(apiID, pageID string) (json.RawMessage, error)
}

func (s *service) ListPages(apiID, parentID string) ([]json.RawMessage, error) {
	q := client.BuildQuery(map[string]string{"parentId": parentID})

	path := fmt.Sprintf("apis/%s/pages", apiID)
	if parentID != "" {
		path += "?" + q
	}

	data, err := s.client.Get(s.v2(path))
	if err != nil {
		return nil, fmt.Errorf("page list failed: %w", err)
	}

	var resp struct {
		Pages []json.RawMessage `json:"pages"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse pages response: %w", err)
	}

	return resp.Pages, nil
}

func (s *service) GetPage(apiID, pageID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.v2(fmt.Sprintf("apis/%s/pages/%s", apiID, pageID)))
	if err != nil {
		return nil, err
	}

	return raw(data), nil
}

func (s *service) CreatePage(apiID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/pages", apiID)), body)
	if err != nil {
		return nil, fmt.Errorf("page creation failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) UpdatePage(apiID, pageID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.v2(fmt.Sprintf("apis/%s/pages/%s", apiID, pageID)), body)
	if err != nil {
		return nil, fmt.Errorf("page update failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) DeletePage(apiID, pageID string) error {
	if err := s.client.Delete(s.v2(fmt.Sprintf("apis/%s/pages/%s", apiID, pageID))); err != nil {
		return fmt.Errorf("page deletion failed: %w", err)
	}

	return nil
}

func (s *service) PublishPage(apiID, pageID string) (json.RawMessage, error) {
	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/pages/%s/_publish", apiID, pageID)), nil)
	if err != nil {
		return nil, fmt.Errorf("page publish failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) UnpublishPage(apiID, pageID string) (json.RawMessage, error) {
	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/pages/%s/_unpublish", apiID, pageID)), nil)
	if err != nil {
		return nil, fmt.Errorf("page unpublish failed: %w", err)
	}

	return raw(data), nil
}
