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
