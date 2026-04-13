package am

import (
	"encoding/json"
	"fmt"
)

// FactorService defines factor-related operations.
type FactorService interface {
	ListFactors(domainID string) ([]json.RawMessage, error)
	GetFactor(domainID, factorID string) (json.RawMessage, error)
	CreateFactor(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateFactor(domainID, factorID string, body json.RawMessage) (json.RawMessage, error)
	DeleteFactor(domainID, factorID string) error
}

func (s *service) ListFactors(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "factors"))
	if err != nil {
		return nil, fmt.Errorf("factor list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse factor list: %w", err)
	}

	return items, nil
}

func (s *service) GetFactor(domainID, factorID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("factors/%s", factorID)))
	if err != nil {
		return nil, fmt.Errorf("factor get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateFactor(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "factors"), body)
	if err != nil {
		return nil, fmt.Errorf("factor create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateFactor(domainID, factorID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("factors/%s", factorID)), body)
	if err != nil {
		return nil, fmt.Errorf("factor update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteFactor(domainID, factorID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("factors/%s", factorID)))
	if err != nil {
		return fmt.Errorf("factor delete failed: %w", err)
	}

	return nil
}
