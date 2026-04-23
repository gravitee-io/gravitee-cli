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

// AuthDeviceNotifierService defines auth device notifier-related operations.
type AuthDeviceNotifierService interface {
	ListAuthDeviceNotifiers(domainID string) ([]json.RawMessage, error)
	GetAuthDeviceNotifier(domainID, authDeviceNotifierID string) (json.RawMessage, error)
	CreateAuthDeviceNotifier(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateAuthDeviceNotifier(domainID, authDeviceNotifierID string, body json.RawMessage) (json.RawMessage, error)
	DeleteAuthDeviceNotifier(domainID, authDeviceNotifierID string) error
}

func (s *service) ListAuthDeviceNotifiers(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "auth-device-notifiers"))
	if err != nil {
		return nil, fmt.Errorf("auth device notifier list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse auth device notifier list: %w", err)
	}

	return items, nil
}

func (s *service) GetAuthDeviceNotifier(domainID, authDeviceNotifierID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("auth-device-notifiers/%s", authDeviceNotifierID)))
	if err != nil {
		return nil, fmt.Errorf("auth device notifier get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateAuthDeviceNotifier(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "auth-device-notifiers"), body)
	if err != nil {
		return nil, fmt.Errorf("auth device notifier create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateAuthDeviceNotifier(domainID, authDeviceNotifierID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("auth-device-notifiers/%s", authDeviceNotifierID)), body)
	if err != nil {
		return nil, fmt.Errorf("auth device notifier update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteAuthDeviceNotifier(domainID, authDeviceNotifierID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("auth-device-notifiers/%s", authDeviceNotifierID)))
	if err != nil {
		return fmt.Errorf("auth device notifier delete failed: %w", err)
	}

	return nil
}
