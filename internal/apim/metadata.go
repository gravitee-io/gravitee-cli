package apim

import (
	"encoding/json"
	"fmt"
	"github.com/gravitee-io/gio-cli/internal/client"
)

// MetadataService defines metadata-related operations.
type MetadataService interface {
	ListMetadata(apiID string, page, perPage int) (*PaginatedResponse, error)
	CreateMetadata(apiID string, body json.RawMessage) (json.RawMessage, error)
	UpdateMetadata(apiID, key string, body json.RawMessage) (json.RawMessage, error)
	DeleteMetadata(apiID, key string) error
}

func (s *service) ListMetadata(apiID string, page, perPage int) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{"page": client.Itoa(page), "perPage": client.Itoa(perPage)})

	data, err := s.client.Get(s.v2(fmt.Sprintf("apis/%s/metadata?%s", apiID, q)))
	if err != nil {
		return nil, fmt.Errorf("metadata list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) CreateMetadata(apiID string, body json.RawMessage) (json.RawMessage, error) {
	if err := s.requireWrite(); err != nil {
		return nil, err
	}

	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/metadata", apiID)), body)
	if err != nil {
		return nil, fmt.Errorf("metadata creation failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) UpdateMetadata(apiID, key string, body json.RawMessage) (json.RawMessage, error) {
	if err := s.requireWrite(); err != nil {
		return nil, err
	}

	data, err := s.client.Put(s.v2(fmt.Sprintf("apis/%s/metadata/%s", apiID, key)), body)
	if err != nil {
		return nil, fmt.Errorf("metadata update failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) DeleteMetadata(apiID, key string) error {
	if err := s.requireWrite(); err != nil {
		return err
	}

	if err := s.client.Delete(s.v2(fmt.Sprintf("apis/%s/metadata/%s", apiID, key))); err != nil {
		return fmt.Errorf("metadata deletion failed: %w", err)
	}

	return nil
}
