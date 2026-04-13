package am

import (
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/client"
)

// ListOrgUsersParams holds parameters for listing organization users.
type ListOrgUsersParams struct {
	Page    int
	PerPage int
}

// ListOrgAuditsParams holds parameters for listing organization audits.
type ListOrgAuditsParams struct {
	Type    string
	Status  string
	From    string
	To      string
	Page    int
	PerPage int
}

// OrganizationService defines organization-level operations.
type OrganizationService interface {
	// Org users
	ListOrgUsers(params ListOrgUsersParams) (*PaginatedResponse, error)
	GetOrgUser(userID string) (json.RawMessage, error)
	CreateOrgUser(body json.RawMessage) (json.RawMessage, error)
	UpdateOrgUser(userID string, body json.RawMessage) (json.RawMessage, error)
	DeleteOrgUser(userID string) error
	ResetOrgUserPassword(userID string, body json.RawMessage) error
	UpdateOrgUserStatus(userID string, body json.RawMessage) (json.RawMessage, error)
	UpdateOrgUsername(userID string, body json.RawMessage) (json.RawMessage, error)
	BulkOrgUserOperation(body json.RawMessage) (json.RawMessage, error)

	// Org groups
	ListOrgGroups() (json.RawMessage, error)
	GetOrgGroup(groupID string) (json.RawMessage, error)
	CreateOrgGroup(body json.RawMessage) (json.RawMessage, error)
	UpdateOrgGroup(groupID string, body json.RawMessage) (json.RawMessage, error)
	DeleteOrgGroup(groupID string) error

	// Org roles
	ListOrgRoles() ([]json.RawMessage, error)
	GetOrgRole(roleID string) (json.RawMessage, error)
	CreateOrgRole(body json.RawMessage) (json.RawMessage, error)
	UpdateOrgRole(roleID string, body json.RawMessage) (json.RawMessage, error)
	DeleteOrgRole(roleID string) error

	// Org settings
	GetOrgSettings() (json.RawMessage, error)
	PatchOrgSettings(body json.RawMessage) (json.RawMessage, error)

	// Org members
	ListOrgMembers() (json.RawMessage, error)
	AddOrgMember(body json.RawMessage) (json.RawMessage, error)
	RemoveOrgMember(memberID string) error

	// Org audits
	ListOrgAudits(params ListOrgAuditsParams) (*PaginatedResponse, error)
	GetOrgAudit(auditID string) (json.RawMessage, error)
}

// Org users

func (s *service) ListOrgUsers(params ListOrgUsersParams) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{
		"page": client.Itoa(params.Page),
		"size": client.Itoa(params.PerPage),
	})

	data, err := s.client.Get(s.orgPath("users?" + q))
	if err != nil {
		return nil, fmt.Errorf("org user list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) GetOrgUser(userID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.orgPath(fmt.Sprintf("users/%s", userID)))
	if err != nil {
		return nil, fmt.Errorf("org user get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateOrgUser(body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.orgPath("users"), body)
	if err != nil {
		return nil, fmt.Errorf("org user create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateOrgUser(userID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.orgPath(fmt.Sprintf("users/%s", userID)), body)
	if err != nil {
		return nil, fmt.Errorf("org user update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteOrgUser(userID string) error {
	err := s.client.Delete(s.orgPath(fmt.Sprintf("users/%s", userID)))
	if err != nil {
		return fmt.Errorf("org user delete failed: %w", err)
	}

	return nil
}

func (s *service) ResetOrgUserPassword(userID string, body json.RawMessage) error {
	_, err := s.client.Post(s.orgPath(fmt.Sprintf("users/%s/resetPassword", userID)), body)
	if err != nil {
		return fmt.Errorf("org user password reset failed: %w", err)
	}

	return nil
}

func (s *service) UpdateOrgUserStatus(userID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.orgPath(fmt.Sprintf("users/%s/status", userID)), body)
	if err != nil {
		return nil, fmt.Errorf("org user status update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateOrgUsername(userID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Patch(s.orgPath(fmt.Sprintf("users/%s/username", userID)), body)
	if err != nil {
		return nil, fmt.Errorf("org username update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) BulkOrgUserOperation(body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.orgPath("users/bulk"), body)
	if err != nil {
		return nil, fmt.Errorf("org user bulk operation failed: %w", err)
	}

	return json.RawMessage(data), nil
}

// Org groups

func (s *service) ListOrgGroups() (json.RawMessage, error) {
	data, err := s.client.Get(s.orgPath("groups"))
	if err != nil {
		return nil, fmt.Errorf("org group list failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) GetOrgGroup(groupID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.orgPath(fmt.Sprintf("groups/%s", groupID)))
	if err != nil {
		return nil, fmt.Errorf("org group get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateOrgGroup(body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.orgPath("groups"), body)
	if err != nil {
		return nil, fmt.Errorf("org group create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateOrgGroup(groupID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.orgPath(fmt.Sprintf("groups/%s", groupID)), body)
	if err != nil {
		return nil, fmt.Errorf("org group update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteOrgGroup(groupID string) error {
	err := s.client.Delete(s.orgPath(fmt.Sprintf("groups/%s", groupID)))
	if err != nil {
		return fmt.Errorf("org group delete failed: %w", err)
	}

	return nil
}

// Org roles

func (s *service) ListOrgRoles() ([]json.RawMessage, error) {
	data, err := s.client.Get(s.orgPath("roles"))
	if err != nil {
		return nil, fmt.Errorf("org role list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return items, nil
}

func (s *service) GetOrgRole(roleID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.orgPath(fmt.Sprintf("roles/%s", roleID)))
	if err != nil {
		return nil, fmt.Errorf("org role get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateOrgRole(body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.orgPath("roles"), body)
	if err != nil {
		return nil, fmt.Errorf("org role create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateOrgRole(roleID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.orgPath(fmt.Sprintf("roles/%s", roleID)), body)
	if err != nil {
		return nil, fmt.Errorf("org role update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteOrgRole(roleID string) error {
	err := s.client.Delete(s.orgPath(fmt.Sprintf("roles/%s", roleID)))
	if err != nil {
		return fmt.Errorf("org role delete failed: %w", err)
	}

	return nil
}

// Org settings

func (s *service) GetOrgSettings() (json.RawMessage, error) {
	data, err := s.client.Get(s.orgPath("settings"))
	if err != nil {
		return nil, fmt.Errorf("org settings get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) PatchOrgSettings(body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Patch(s.orgPath("settings"), body)
	if err != nil {
		return nil, fmt.Errorf("org settings patch failed: %w", err)
	}

	return json.RawMessage(data), nil
}

// Org members

func (s *service) ListOrgMembers() (json.RawMessage, error) {
	data, err := s.client.Get(s.orgPath("members"))
	if err != nil {
		return nil, fmt.Errorf("org member list failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) AddOrgMember(body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.orgPath("members"), body)
	if err != nil {
		return nil, fmt.Errorf("org member add failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) RemoveOrgMember(memberID string) error {
	err := s.client.Delete(s.orgPath(fmt.Sprintf("members/%s", memberID)))
	if err != nil {
		return fmt.Errorf("org member remove failed: %w", err)
	}

	return nil
}

// Org audits

func (s *service) ListOrgAudits(params ListOrgAuditsParams) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{
		"page":   client.Itoa(params.Page),
		"size":   client.Itoa(params.PerPage),
		"type":   params.Type,
		"status": params.Status,
		"from":   params.From,
		"to":     params.To,
	})

	data, err := s.client.Get(s.orgPath("audits?" + q))
	if err != nil {
		return nil, fmt.Errorf("org audit list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) GetOrgAudit(auditID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.orgPath(fmt.Sprintf("audits/%s", auditID)))
	if err != nil {
		return nil, fmt.Errorf("org audit get failed: %w", err)
	}

	return json.RawMessage(data), nil
}
