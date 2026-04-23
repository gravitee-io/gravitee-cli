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

// DeviceIdentifierService defines device identifier-related operations.
type DeviceIdentifierService interface {
	ListDeviceIdentifiers(domainID string) ([]json.RawMessage, error)
	GetDeviceIdentifier(domainID, deviceIdentifierID string) (json.RawMessage, error)
	CreateDeviceIdentifier(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateDeviceIdentifier(domainID, deviceIdentifierID string, body json.RawMessage) (json.RawMessage, error)
	DeleteDeviceIdentifier(domainID, deviceIdentifierID string) error
}

func (s *service) ListDeviceIdentifiers(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "device-identifiers"))
	if err != nil {
		return nil, fmt.Errorf("device identifier list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse device identifier list: %w", err)
	}

	return items, nil
}

func (s *service) GetDeviceIdentifier(domainID, deviceIdentifierID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("device-identifiers/%s", deviceIdentifierID)))
	if err != nil {
		return nil, fmt.Errorf("device identifier get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateDeviceIdentifier(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "device-identifiers"), body)
	if err != nil {
		return nil, fmt.Errorf("device identifier create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateDeviceIdentifier(domainID, deviceIdentifierID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("device-identifiers/%s", deviceIdentifierID)), body)
	if err != nil {
		return nil, fmt.Errorf("device identifier update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteDeviceIdentifier(domainID, deviceIdentifierID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("device-identifiers/%s", deviceIdentifierID)))
	if err != nil {
		return fmt.Errorf("device identifier delete failed: %w", err)
	}

	return nil
}
