package apim

import (
	"encoding/json"
	"fmt"
	"github.com/gravitee-io/gio-cli/internal/client"
)

// ListApplicationsParams holds parameters for listing applications.
type ListApplicationsParams struct {
	Query   string
	Status  string
	Order   string
	Page    int
	PerPage int
}

// ApplicationService defines application-related operations (V1 API).
type ApplicationService interface {
	ListApplications(params ListApplicationsParams) (*PaginatedResponse, error)
	GetApplication(appID string) (json.RawMessage, error)
	CreateApplication(body json.RawMessage) (json.RawMessage, error)
	DeleteApplication(appID string) error
}

func (s *service) ListApplications(params ListApplicationsParams) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{
		"page": client.Itoa(params.Page), "size": client.Itoa(params.PerPage),
		"query": params.Query, "status": params.Status, "order": params.Order,
	})

	data, err := s.client.Get(s.v1(fmt.Sprintf("applications/_paged?%s", q)))
	if err != nil {
		return nil, fmt.Errorf("application list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) GetApplication(appID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.v1(fmt.Sprintf("applications/%s", appID)))
	if err != nil {
		return nil, err
	}

	return raw(data), nil
}

func (s *service) CreateApplication(body json.RawMessage) (json.RawMessage, error) {
	if err := s.requireWrite(); err != nil {
		return nil, err
	}

	data, err := s.client.Post(s.v1("applications"), body)
	if err != nil {
		return nil, fmt.Errorf("application creation failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) DeleteApplication(appID string) error {
	if err := s.requireWrite(); err != nil {
		return err
	}

	if err := s.client.Delete(s.v1(fmt.Sprintf("applications/%s", appID))); err != nil {
		return fmt.Errorf("application deletion failed: %w", err)
	}

	return nil
}
