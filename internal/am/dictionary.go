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

// DictionaryService defines i18n dictionary-related operations.
type DictionaryService interface {
	ListDictionaries(domainID string) ([]json.RawMessage, error)
	GetDictionary(domainID, dictID string) (json.RawMessage, error)
	CreateDictionary(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateDictionary(domainID, dictID string, body json.RawMessage) (json.RawMessage, error)
	DeleteDictionary(domainID, dictID string) error
	ListDictionaryEntries(domainID, dictID string) (json.RawMessage, error)
	UpdateDictionaryEntries(domainID, dictID string, body json.RawMessage) (json.RawMessage, error)
}

func (s *service) ListDictionaries(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "i18n/dictionaries"))
	if err != nil {
		return nil, fmt.Errorf("dictionary list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse dictionary list: %w", err)
	}

	return items, nil
}

func (s *service) GetDictionary(domainID, dictID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("i18n/dictionaries/%s", dictID)))
	if err != nil {
		return nil, fmt.Errorf("dictionary get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateDictionary(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "i18n/dictionaries"), body)
	if err != nil {
		return nil, fmt.Errorf("dictionary create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateDictionary(domainID, dictID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("i18n/dictionaries/%s", dictID)), body)
	if err != nil {
		return nil, fmt.Errorf("dictionary update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteDictionary(domainID, dictID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("i18n/dictionaries/%s", dictID)))
	if err != nil {
		return fmt.Errorf("dictionary delete failed: %w", err)
	}

	return nil
}

func (s *service) ListDictionaryEntries(domainID, dictID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("i18n/dictionaries/%s/entries", dictID)))
	if err != nil {
		return nil, fmt.Errorf("dictionary entry list failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateDictionaryEntries(domainID, dictID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Patch(s.domainPath(domainID, fmt.Sprintf("i18n/dictionaries/%s/entries", dictID)), body)
	if err != nil {
		return nil, fmt.Errorf("dictionary entry update failed: %w", err)
	}

	return json.RawMessage(data), nil
}
