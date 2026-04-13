package am

import (
	"encoding/json"
	"fmt"
)

// BotDetectionService defines bot detection-related operations.
type BotDetectionService interface {
	ListBotDetections(domainID string) ([]json.RawMessage, error)
	GetBotDetection(domainID, botDetectionID string) (json.RawMessage, error)
	CreateBotDetection(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateBotDetection(domainID, botDetectionID string, body json.RawMessage) (json.RawMessage, error)
	DeleteBotDetection(domainID, botDetectionID string) error
}

func (s *service) ListBotDetections(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "bot-detections"))
	if err != nil {
		return nil, fmt.Errorf("bot detection list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse bot detection list: %w", err)
	}

	return items, nil
}

func (s *service) GetBotDetection(domainID, botDetectionID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("bot-detections/%s", botDetectionID)))
	if err != nil {
		return nil, fmt.Errorf("bot detection get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateBotDetection(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "bot-detections"), body)
	if err != nil {
		return nil, fmt.Errorf("bot detection create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateBotDetection(domainID, botDetectionID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("bot-detections/%s", botDetectionID)), body)
	if err != nil {
		return nil, fmt.Errorf("bot detection update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteBotDetection(domainID, botDetectionID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("bot-detections/%s", botDetectionID)))
	if err != nil {
		return fmt.Errorf("bot detection delete failed: %w", err)
	}

	return nil
}
