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

// AlertService defines alert-related operations (notifiers and triggers).
type AlertService interface {
	ListAlertNotifiers(domainID string) ([]json.RawMessage, error)
	GetAlertNotifier(domainID, notifierID string) (json.RawMessage, error)
	CreateAlertNotifier(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateAlertNotifier(domainID, notifierID string, body json.RawMessage) (json.RawMessage, error)
	DeleteAlertNotifier(domainID, notifierID string) error
	GetAlertTriggers(domainID string) (json.RawMessage, error)
	UpdateAlertTriggers(domainID string, body json.RawMessage) (json.RawMessage, error)
}

func (s *service) ListAlertNotifiers(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "alerts/notifiers"))
	if err != nil {
		return nil, fmt.Errorf("alert notifier list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse alert notifier list: %w", err)
	}

	return items, nil
}

func (s *service) GetAlertNotifier(domainID, notifierID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("alerts/notifiers/%s", notifierID)))
	if err != nil {
		return nil, fmt.Errorf("alert notifier get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateAlertNotifier(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "alerts/notifiers"), body)
	if err != nil {
		return nil, fmt.Errorf("alert notifier create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateAlertNotifier(domainID, notifierID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("alerts/notifiers/%s", notifierID)), body)
	if err != nil {
		return nil, fmt.Errorf("alert notifier update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteAlertNotifier(domainID, notifierID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("alerts/notifiers/%s", notifierID)))
	if err != nil {
		return fmt.Errorf("alert notifier delete failed: %w", err)
	}

	return nil
}

func (s *service) GetAlertTriggers(domainID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "alerts/triggers"))
	if err != nil {
		return nil, fmt.Errorf("alert trigger get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateAlertTriggers(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, "alerts/triggers"), body)
	if err != nil {
		return nil, fmt.Errorf("alert trigger update failed: %w", err)
	}

	return json.RawMessage(data), nil
}
