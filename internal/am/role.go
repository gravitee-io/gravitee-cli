package am

import (
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/client"
)

// ListRolesParams holds parameters for listing roles.
type ListRolesParams struct {
	Query   string
	Page    int
	PerPage int
}

// RoleService defines role-related operations.
type RoleService interface {
	ListRoles(domainID string, params ListRolesParams) (*PaginatedResponse, error)
	GetRole(domainID, roleID string) (json.RawMessage, error)
	CreateRole(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateRole(domainID, roleID string, body json.RawMessage) (json.RawMessage, error)
	DeleteRole(domainID, roleID string) error
}

func (s *service) ListRoles(domainID string, params ListRolesParams) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{
		"page": client.Itoa(params.Page),
		"size": client.Itoa(params.PerPage),
		"q":    params.Query,
	})

	data, err := s.client.Get(s.domainPath(domainID, "roles?"+q))
	if err != nil {
		return nil, fmt.Errorf("role list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) GetRole(domainID, roleID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("roles/%s", roleID)))
	if err != nil {
		return nil, fmt.Errorf("role get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateRole(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "roles"), body)
	if err != nil {
		return nil, fmt.Errorf("role create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateRole(domainID, roleID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("roles/%s", roleID)), body)
	if err != nil {
		return nil, fmt.Errorf("role update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteRole(domainID, roleID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("roles/%s", roleID)))
	if err != nil {
		return fmt.Errorf("role delete failed: %w", err)
	}

	return nil
}
