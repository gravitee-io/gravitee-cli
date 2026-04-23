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

// ThemeService defines theme-related operations.
type ThemeService interface {
	ListThemes(domainID string) ([]json.RawMessage, error)
	GetTheme(domainID, themeID string) (json.RawMessage, error)
	CreateTheme(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateTheme(domainID, themeID string, body json.RawMessage) (json.RawMessage, error)
	DeleteTheme(domainID, themeID string) error
}

func (s *service) ListThemes(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "themes"))
	if err != nil {
		return nil, fmt.Errorf("theme list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse theme list: %w", err)
	}

	return items, nil
}

func (s *service) GetTheme(domainID, themeID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("themes/%s", themeID)))
	if err != nil {
		return nil, fmt.Errorf("theme get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateTheme(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "themes"), body)
	if err != nil {
		return nil, fmt.Errorf("theme create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateTheme(domainID, themeID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("themes/%s", themeID)), body)
	if err != nil {
		return nil, fmt.Errorf("theme update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteTheme(domainID, themeID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("themes/%s", themeID)))
	if err != nil {
		return fmt.Errorf("theme delete failed: %w", err)
	}

	return nil
}
