package apim

import (
	"encoding/json"
	"fmt"
	"github.com/gravitee-io/gio-cli/internal/client"
)

// PageService defines page-related operations.
type PageService interface {
	ListPages(apiID string, page, perPage int) (*PaginatedResponse, error)
	GetPage(apiID, pageID string) (json.RawMessage, error)
	CreatePage(apiID string, body json.RawMessage) (json.RawMessage, error)
	UpdatePage(apiID, pageID string, body json.RawMessage) (json.RawMessage, error)
	DeletePage(apiID, pageID string) error
	PublishPage(apiID, pageID string) (json.RawMessage, error)
	UnpublishPage(apiID, pageID string) (json.RawMessage, error)
}

func (s *service) ListPages(apiID string, page, perPage int) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{"page": client.Itoa(page), "perPage": client.Itoa(perPage)})

	data, err := s.client.Get(s.v2(fmt.Sprintf("apis/%s/pages?%s", apiID, q)))
	if err != nil {
		return nil, fmt.Errorf("page list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) GetPage(apiID, pageID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.v2(fmt.Sprintf("apis/%s/pages/%s", apiID, pageID)))
	if err != nil {
		return nil, err
	}

	return raw(data), nil
}

func (s *service) CreatePage(apiID string, body json.RawMessage) (json.RawMessage, error) {
	if err := s.requireWrite(); err != nil {
		return nil, err
	}

	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/pages", apiID)), body)
	if err != nil {
		return nil, fmt.Errorf("page creation failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) UpdatePage(apiID, pageID string, body json.RawMessage) (json.RawMessage, error) {
	if err := s.requireWrite(); err != nil {
		return nil, err
	}

	data, err := s.client.Put(s.v2(fmt.Sprintf("apis/%s/pages/%s", apiID, pageID)), body)
	if err != nil {
		return nil, fmt.Errorf("page update failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) DeletePage(apiID, pageID string) error {
	if err := s.requireWrite(); err != nil {
		return err
	}

	if err := s.client.Delete(s.v2(fmt.Sprintf("apis/%s/pages/%s", apiID, pageID))); err != nil {
		return fmt.Errorf("page deletion failed: %w", err)
	}

	return nil
}

func (s *service) PublishPage(apiID, pageID string) (json.RawMessage, error) {
	if err := s.requireWrite(); err != nil {
		return nil, err
	}

	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/pages/%s/_publish", apiID, pageID)), nil)
	if err != nil {
		return nil, fmt.Errorf("page publish failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) UnpublishPage(apiID, pageID string) (json.RawMessage, error) {
	if err := s.requireWrite(); err != nil {
		return nil, err
	}

	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/pages/%s/_unpublish", apiID, pageID)), nil)
	if err != nil {
		return nil, fmt.Errorf("page unpublish failed: %w", err)
	}

	return raw(data), nil
}
