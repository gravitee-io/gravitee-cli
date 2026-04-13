package am

import (
	"encoding/json"
	"fmt"
)

// FlowService defines flow-related operations.
type FlowService interface {
	ListFlows(domainID string) ([]json.RawMessage, error)
	GetFlow(domainID, flowID string) (json.RawMessage, error)
	UpdateFlows(domainID string, body json.RawMessage) (json.RawMessage, error)
}

func (s *service) ListFlows(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "flows"))
	if err != nil {
		return nil, fmt.Errorf("flow list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse flow list: %w", err)
	}

	return items, nil
}

func (s *service) GetFlow(domainID, flowID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("flows/%s", flowID)))
	if err != nil {
		return nil, fmt.Errorf("flow get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateFlows(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, "flows"), body)
	if err != nil {
		return nil, fmt.Errorf("flow update failed: %w", err)
	}

	return json.RawMessage(data), nil
}
