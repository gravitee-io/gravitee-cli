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

// CertificateService defines certificate-related operations.
type CertificateService interface {
	ListCertificates(domainID string) ([]json.RawMessage, error)
	GetCertificate(domainID, certID string) (json.RawMessage, error)
	CreateCertificate(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateCertificate(domainID, certID string, body json.RawMessage) (json.RawMessage, error)
	DeleteCertificate(domainID, certID string) error
	GetCertificateKey(domainID, certID string) (json.RawMessage, error)
	GetCertificateKeys(domainID, certID string) (json.RawMessage, error)
	RotateCertificates(domainID string) (json.RawMessage, error)
}

func (s *service) ListCertificates(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "certificates"))
	if err != nil {
		return nil, fmt.Errorf("certificate list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse certificate list: %w", err)
	}

	return items, nil
}

func (s *service) GetCertificate(domainID, certID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("certificates/%s", certID)))
	if err != nil {
		return nil, fmt.Errorf("certificate get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateCertificate(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "certificates"), body)
	if err != nil {
		return nil, fmt.Errorf("certificate create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateCertificate(domainID, certID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("certificates/%s", certID)), body)
	if err != nil {
		return nil, fmt.Errorf("certificate update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteCertificate(domainID, certID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("certificates/%s", certID)))
	if err != nil {
		return fmt.Errorf("certificate delete failed: %w", err)
	}

	return nil
}

func (s *service) GetCertificateKey(domainID, certID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("certificates/%s/key", certID)))
	if err != nil {
		return nil, fmt.Errorf("certificate key get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) GetCertificateKeys(domainID, certID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("certificates/%s/keys", certID)))
	if err != nil {
		return nil, fmt.Errorf("certificate keys get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) RotateCertificates(domainID string) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "certificates/rotate"), nil)
	if err != nil {
		return nil, fmt.Errorf("certificate rotate failed: %w", err)
	}

	return json.RawMessage(data), nil
}
