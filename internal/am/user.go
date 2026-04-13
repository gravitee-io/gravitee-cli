package am

import (
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/client"
)

// ListUsersParams holds parameters for listing users.
type ListUsersParams struct {
	Query   string
	Filter  string
	Page    int
	PerPage int
}

// ListUserAuditsParams holds parameters for listing user audits.
type ListUserAuditsParams struct {
	Type    string
	Status  string
	From    string
	To      string
	Page    int
	PerPage int
}

// UserService defines user-related operations.
type UserService interface {
	ListUsers(domainID string, params ListUsersParams) (*PaginatedResponse, error)
	GetUser(domainID, userID string) (json.RawMessage, error)
	CreateUser(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateUser(domainID, userID string, body json.RawMessage) (json.RawMessage, error)
	DeleteUser(domainID, userID string) error
	UpdateUserStatus(domainID, userID string, body json.RawMessage) (json.RawMessage, error)
	ResetPassword(domainID, userID string, body json.RawMessage) error

	// User consents
	ListUserConsents(domainID, userID string) ([]json.RawMessage, error)
	RevokeUserConsent(domainID, userID, consentID string) error
	RevokeAllUserConsents(domainID, userID string) error

	// User roles
	ListUserRoles(domainID, userID string) (json.RawMessage, error)
	AssignUserRoles(domainID, userID string, body json.RawMessage) error
	RevokeUserRole(domainID, userID, roleID string) error

	// User devices
	ListUserDevices(domainID, userID string) ([]json.RawMessage, error)
	DeleteUserDevice(domainID, userID, deviceID string) error

	// User credentials
	ListUserCredentials(domainID, userID string) ([]json.RawMessage, error)
	GetUserCredential(domainID, userID, credentialID string) (json.RawMessage, error)
	RevokeUserCredential(domainID, userID, credentialID string) error

	// User enrolled factors
	ListUserFactors(domainID, userID string) ([]json.RawMessage, error)
	DeleteUserFactor(domainID, userID, factorID string) error

	// User audits
	ListUserAudits(domainID, userID string, params ListUserAuditsParams) (*PaginatedResponse, error)
	GetUserAudit(domainID, userID, auditID string) (json.RawMessage, error)

	// User identities
	ListUserIdentities(domainID, userID string) ([]json.RawMessage, error)
	UnlinkUserIdentity(domainID, userID, identityID string) error

	// User cert-credentials
	ListUserCertCredentials(domainID, userID string) ([]json.RawMessage, error)
	GetUserCertCredential(domainID, userID, credID string) (json.RawMessage, error)
	EnrollUserCertCredential(domainID, userID string, body json.RawMessage) (json.RawMessage, error)
	RevokeUserCertCredential(domainID, userID, credID string) error

	// Bulk user operations
	BulkUserOperation(domainID string, body json.RawMessage) (json.RawMessage, error)

	// Additional user operations
	SendRegistrationConfirmation(domainID, userID string) error
	UpdateUsername(domainID, userID string, body json.RawMessage) (json.RawMessage, error)
}

func (s *service) ListUsers(domainID string, params ListUsersParams) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{
		"page":   client.Itoa(params.Page),
		"size":   client.Itoa(params.PerPage),
		"q":      params.Query,
		"filter": params.Filter,
	})

	data, err := s.client.Get(s.domainPath(domainID, "users?"+q))
	if err != nil {
		return nil, fmt.Errorf("user list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) GetUser(domainID, userID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("users/%s", userID)))
	if err != nil {
		return nil, fmt.Errorf("user get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateUser(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "users"), body)
	if err != nil {
		return nil, fmt.Errorf("user create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateUser(domainID, userID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("users/%s", userID)), body)
	if err != nil {
		return nil, fmt.Errorf("user update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteUser(domainID, userID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("users/%s", userID)))
	if err != nil {
		return fmt.Errorf("user delete failed: %w", err)
	}

	return nil
}

func (s *service) UpdateUserStatus(domainID, userID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("users/%s/status", userID)), body)
	if err != nil {
		return nil, fmt.Errorf("user status update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) ResetPassword(domainID, userID string, body json.RawMessage) error {
	_, err := s.client.Post(s.domainPath(domainID, fmt.Sprintf("users/%s/resetPassword", userID)), body)
	if err != nil {
		return fmt.Errorf("password reset failed: %w", err)
	}

	return nil
}

// User consents

func (s *service) ListUserConsents(domainID, userID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("users/%s/consents", userID)))
	if err != nil {
		return nil, fmt.Errorf("user consent list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return items, nil
}

func (s *service) RevokeUserConsent(domainID, userID, consentID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("users/%s/consents/%s", userID, consentID)))
	if err != nil {
		return fmt.Errorf("user consent revoke failed: %w", err)
	}

	return nil
}

func (s *service) RevokeAllUserConsents(domainID, userID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("users/%s/consents", userID)))
	if err != nil {
		return fmt.Errorf("user consents revoke failed: %w", err)
	}

	return nil
}

// User roles

func (s *service) ListUserRoles(domainID, userID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("users/%s/roles", userID)))
	if err != nil {
		return nil, fmt.Errorf("user role list failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) AssignUserRoles(domainID, userID string, body json.RawMessage) error {
	_, err := s.client.Post(s.domainPath(domainID, fmt.Sprintf("users/%s/roles", userID)), body)
	if err != nil {
		return fmt.Errorf("user role assign failed: %w", err)
	}

	return nil
}

func (s *service) RevokeUserRole(domainID, userID, roleID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("users/%s/roles/%s", userID, roleID)))
	if err != nil {
		return fmt.Errorf("user role revoke failed: %w", err)
	}

	return nil
}

// User devices

func (s *service) ListUserDevices(domainID, userID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("users/%s/devices", userID)))
	if err != nil {
		return nil, fmt.Errorf("user device list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return items, nil
}

func (s *service) DeleteUserDevice(domainID, userID, deviceID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("users/%s/devices/%s", userID, deviceID)))
	if err != nil {
		return fmt.Errorf("user device delete failed: %w", err)
	}

	return nil
}

// User credentials

func (s *service) ListUserCredentials(domainID, userID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("users/%s/credentials", userID)))
	if err != nil {
		return nil, fmt.Errorf("user credential list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return items, nil
}

func (s *service) GetUserCredential(domainID, userID, credentialID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("users/%s/credentials/%s", userID, credentialID)))
	if err != nil {
		return nil, fmt.Errorf("user credential get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) RevokeUserCredential(domainID, userID, credentialID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("users/%s/credentials/%s", userID, credentialID)))
	if err != nil {
		return fmt.Errorf("user credential revoke failed: %w", err)
	}

	return nil
}

// User enrolled factors

func (s *service) ListUserFactors(domainID, userID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("users/%s/factors", userID)))
	if err != nil {
		return nil, fmt.Errorf("user factor list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return items, nil
}

func (s *service) DeleteUserFactor(domainID, userID, factorID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("users/%s/factors/%s", userID, factorID)))
	if err != nil {
		return fmt.Errorf("user factor delete failed: %w", err)
	}

	return nil
}

// User audits

func (s *service) ListUserAudits(domainID, userID string, params ListUserAuditsParams) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{
		"page":   client.Itoa(params.Page),
		"size":   client.Itoa(params.PerPage),
		"type":   params.Type,
		"status": params.Status,
		"from":   params.From,
		"to":     params.To,
	})

	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("users/%s/audits?%s", userID, q)))
	if err != nil {
		return nil, fmt.Errorf("user audit list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) GetUserAudit(domainID, userID, auditID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("users/%s/audits/%s", userID, auditID)))
	if err != nil {
		return nil, fmt.Errorf("user audit get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

// Additional user operations

func (s *service) SendRegistrationConfirmation(domainID, userID string) error {
	_, err := s.client.Post(s.domainPath(domainID, fmt.Sprintf("users/%s/sendRegistrationConfirmation", userID)), nil)
	if err != nil {
		return fmt.Errorf("send registration confirmation failed: %w", err)
	}

	return nil
}

func (s *service) UpdateUsername(domainID, userID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Patch(s.domainPath(domainID, fmt.Sprintf("users/%s/username", userID)), body)
	if err != nil {
		return nil, fmt.Errorf("username update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

// User identities

func (s *service) ListUserIdentities(domainID, userID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("users/%s/identities", userID)))
	if err != nil {
		return nil, fmt.Errorf("user identity list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return items, nil
}

func (s *service) UnlinkUserIdentity(domainID, userID, identityID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("users/%s/identities/%s", userID, identityID)))
	if err != nil {
		return fmt.Errorf("user identity unlink failed: %w", err)
	}

	return nil
}

// User cert-credentials

func (s *service) ListUserCertCredentials(domainID, userID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("users/%s/cert-credentials", userID)))
	if err != nil {
		return nil, fmt.Errorf("user cert-credential list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return items, nil
}

func (s *service) GetUserCertCredential(domainID, userID, credID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("users/%s/cert-credentials/%s", userID, credID)))
	if err != nil {
		return nil, fmt.Errorf("user cert-credential get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) EnrollUserCertCredential(domainID, userID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, fmt.Sprintf("users/%s/cert-credentials", userID)), body)
	if err != nil {
		return nil, fmt.Errorf("user cert-credential enroll failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) RevokeUserCertCredential(domainID, userID, credID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("users/%s/cert-credentials/%s", userID, credID)))
	if err != nil {
		return fmt.Errorf("user cert-credential revoke failed: %w", err)
	}

	return nil
}

// Bulk user operations

func (s *service) BulkUserOperation(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "users/bulk"), body)
	if err != nil {
		return nil, fmt.Errorf("user bulk operation failed: %w", err)
	}

	return json.RawMessage(data), nil
}
