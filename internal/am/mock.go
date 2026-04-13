package am

import (
	"encoding/json"
	"fmt"
)

// MockService implements Service with injectable functions for testing.
type MockService struct {
	// Domain
	ListDomainsFunc                     func(ListDomainsParams) (*PaginatedResponse, error)
	GetDomainFunc                       func(string) (json.RawMessage, error)
	GetDomainByHRIDFunc                 func(string) (json.RawMessage, error)
	CreateDomainFunc                    func(json.RawMessage) (json.RawMessage, error)
	UpdateDomainFunc                    func(string, json.RawMessage) (json.RawMessage, error)
	PatchDomainFunc                     func(string, json.RawMessage) (json.RawMessage, error)
	DeleteDomainFunc                    func(string) error
	UpdateDomainCertificateSettingsFunc func(string, json.RawMessage) (json.RawMessage, error)
	ListDataPlanesFunc                  func() (json.RawMessage, error)

	// Application
	ListApplicationsFunc  func(string, ListApplicationsParams) (*PaginatedResponse, error)
	GetApplicationFunc    func(string, string) (json.RawMessage, error)
	CreateApplicationFunc func(string, json.RawMessage) (json.RawMessage, error)
	UpdateApplicationFunc func(string, string, json.RawMessage) (json.RawMessage, error)
	PatchApplicationFunc  func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteApplicationFunc func(string, string) error

	// Application sub-resources
	ListAppSecretsFunc  func(string, string) (json.RawMessage, error)
	CreateAppSecretFunc func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteAppSecretFunc func(string, string, string) error

	ListAppMembersFunc  func(string, string) (json.RawMessage, error)
	AddAppMemberFunc    func(string, string, json.RawMessage) (json.RawMessage, error)
	RemoveAppMemberFunc func(string, string, string) error

	ListAppFlowsFunc   func(string, string) ([]json.RawMessage, error)
	GetAppFlowFunc     func(string, string, string) (json.RawMessage, error)
	UpdateAppFlowsFunc func(string, string, json.RawMessage) (json.RawMessage, error)

	GetAppEmailFunc    func(string, string, string) (json.RawMessage, error)
	CreateAppEmailFunc func(string, string, json.RawMessage) (json.RawMessage, error)
	UpdateAppEmailFunc func(string, string, string, json.RawMessage) (json.RawMessage, error)
	DeleteAppEmailFunc func(string, string, string) error

	GetAppFormFunc    func(string, string, string) (json.RawMessage, error)
	CreateAppFormFunc func(string, string, json.RawMessage) (json.RawMessage, error)
	UpdateAppFormFunc func(string, string, string, json.RawMessage) (json.RawMessage, error)
	DeleteAppFormFunc func(string, string, string) error

	ListAppResourcesFunc func(string, string) (json.RawMessage, error)
	GetAppResourceFunc   func(string, string, string) (json.RawMessage, error)

	GetAppAnalyticsFunc         func(string, string, AnalyticsParams) (json.RawMessage, error)
	ChangeAppTypeFunc           func(string, string, json.RawMessage) (json.RawMessage, error)
	ListAppResourcePoliciesFunc func(string, string, string) (json.RawMessage, error)
	GetAppResourcePolicyFunc    func(string, string, string, string) (json.RawMessage, error)
	GetAppMemberPermissionsFunc func(string, string) (json.RawMessage, error)

	// User
	ListUsersFunc        func(string, ListUsersParams) (*PaginatedResponse, error)
	GetUserFunc          func(string, string) (json.RawMessage, error)
	CreateUserFunc       func(string, json.RawMessage) (json.RawMessage, error)
	UpdateUserFunc       func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteUserFunc       func(string, string) error
	UpdateUserStatusFunc func(string, string, json.RawMessage) (json.RawMessage, error)
	ResetPasswordFunc    func(string, string, json.RawMessage) error

	// User sub-resources
	ListUserConsentsFunc             func(string, string) ([]json.RawMessage, error)
	RevokeUserConsentFunc            func(string, string, string) error
	RevokeAllUserConsentsFunc        func(string, string) error
	ListUserRolesFunc                func(string, string) (json.RawMessage, error)
	AssignUserRolesFunc              func(string, string, json.RawMessage) error
	RevokeUserRoleFunc               func(string, string, string) error
	ListUserDevicesFunc              func(string, string) ([]json.RawMessage, error)
	DeleteUserDeviceFunc             func(string, string, string) error
	ListUserCredentialsFunc          func(string, string) ([]json.RawMessage, error)
	GetUserCredentialFunc            func(string, string, string) (json.RawMessage, error)
	RevokeUserCredentialFunc         func(string, string, string) error
	ListUserFactorsFunc              func(string, string) ([]json.RawMessage, error)
	DeleteUserFactorFunc             func(string, string, string) error
	ListUserAuditsFunc               func(string, string, ListUserAuditsParams) (*PaginatedResponse, error)
	GetUserAuditFunc                 func(string, string, string) (json.RawMessage, error)
	SendRegistrationConfirmationFunc func(string, string) error
	UpdateUsernameFunc               func(string, string, json.RawMessage) (json.RawMessage, error)
	ListUserIdentitiesFunc           func(string, string) ([]json.RawMessage, error)
	UnlinkUserIdentityFunc           func(string, string, string) error
	ListUserCertCredentialsFunc      func(string, string) ([]json.RawMessage, error)
	GetUserCertCredentialFunc        func(string, string, string) (json.RawMessage, error)
	EnrollUserCertCredentialFunc     func(string, string, json.RawMessage) (json.RawMessage, error)
	RevokeUserCertCredentialFunc     func(string, string, string) error
	BulkUserOperationFunc            func(string, json.RawMessage) (json.RawMessage, error)

	// Role
	ListRolesFunc  func(string, ListRolesParams) (*PaginatedResponse, error)
	GetRoleFunc    func(string, string) (json.RawMessage, error)
	CreateRoleFunc func(string, json.RawMessage) (json.RawMessage, error)
	UpdateRoleFunc func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteRoleFunc func(string, string) error

	// Scope
	ListScopesFunc  func(string, ListScopesParams) (*PaginatedResponse, error)
	GetScopeFunc    func(string, string) (json.RawMessage, error)
	CreateScopeFunc func(string, json.RawMessage) (json.RawMessage, error)
	UpdateScopeFunc func(string, string, json.RawMessage) (json.RawMessage, error)
	PatchScopeFunc  func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteScopeFunc func(string, string) error

	// IdentityProvider
	ListIdentityProvidersFunc   func(string, bool) ([]json.RawMessage, error)
	GetIdentityProviderFunc     func(string, string) (json.RawMessage, error)
	CreateIdentityProviderFunc  func(string, json.RawMessage) (json.RawMessage, error)
	UpdateIdentityProviderFunc  func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteIdentityProviderFunc  func(string, string) error
	UpdateIDPPasswordPolicyFunc func(string, string, json.RawMessage) (json.RawMessage, error)

	// Certificate
	ListCertificatesFunc   func(string) ([]json.RawMessage, error)
	GetCertificateFunc     func(string, string) (json.RawMessage, error)
	CreateCertificateFunc  func(string, json.RawMessage) (json.RawMessage, error)
	UpdateCertificateFunc  func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteCertificateFunc  func(string, string) error
	GetCertificateKeyFunc  func(string, string) (json.RawMessage, error)
	GetCertificateKeysFunc func(string, string) (json.RawMessage, error)
	RotateCertificatesFunc func(string) (json.RawMessage, error)

	// Factor
	ListFactorsFunc  func(string) ([]json.RawMessage, error)
	GetFactorFunc    func(string, string) (json.RawMessage, error)
	CreateFactorFunc func(string, json.RawMessage) (json.RawMessage, error)
	UpdateFactorFunc func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteFactorFunc func(string, string) error

	// Group
	ListGroupsFunc        func(string, ListGroupsParams) (*PaginatedResponse, error)
	GetGroupFunc          func(string, string) (json.RawMessage, error)
	CreateGroupFunc       func(string, json.RawMessage) (json.RawMessage, error)
	UpdateGroupFunc       func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteGroupFunc       func(string, string) error
	ListGroupMembersFunc  func(string, string) (json.RawMessage, error)
	AddGroupMemberFunc    func(string, string, string) error
	RemoveGroupMemberFunc func(string, string, string) error
	ListGroupRolesFunc    func(string, string) (json.RawMessage, error)
	AssignGroupRolesFunc  func(string, string, json.RawMessage) (json.RawMessage, error)
	RevokeGroupRoleFunc   func(string, string, string) error

	// Flow
	ListFlowsFunc   func(string) ([]json.RawMessage, error)
	GetFlowFunc     func(string, string) (json.RawMessage, error)
	UpdateFlowsFunc func(string, json.RawMessage) (json.RawMessage, error)

	// Form
	ListFormsFunc  func(string) ([]json.RawMessage, error)
	GetFormFunc    func(string, string) (json.RawMessage, error)
	CreateFormFunc func(string, json.RawMessage) (json.RawMessage, error)
	UpdateFormFunc func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteFormFunc func(string, string) error

	// Email
	ListEmailsFunc  func(string) ([]json.RawMessage, error)
	GetEmailFunc    func(string, string) (json.RawMessage, error)
	CreateEmailFunc func(string, json.RawMessage) (json.RawMessage, error)
	UpdateEmailFunc func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteEmailFunc func(string, string) error

	// Theme
	ListThemesFunc  func(string) ([]json.RawMessage, error)
	GetThemeFunc    func(string, string) (json.RawMessage, error)
	CreateThemeFunc func(string, json.RawMessage) (json.RawMessage, error)
	UpdateThemeFunc func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteThemeFunc func(string, string) error

	// PasswordPolicy
	ListPasswordPoliciesFunc     func(string) ([]json.RawMessage, error)
	GetPasswordPolicyFunc        func(string, string) (json.RawMessage, error)
	CreatePasswordPolicyFunc     func(string, json.RawMessage) (json.RawMessage, error)
	UpdatePasswordPolicyFunc     func(string, string, json.RawMessage) (json.RawMessage, error)
	DeletePasswordPolicyFunc     func(string, string) error
	GetActivePasswordPolicyFunc  func(string) (json.RawMessage, error)
	SetDefaultPasswordPolicyFunc func(string, string) (json.RawMessage, error)
	EvaluatePasswordPolicyFunc   func(string, string, json.RawMessage) (json.RawMessage, error)

	// Audit
	ListAuditsFunc func(string, ListAuditsParams) (*PaginatedResponse, error)
	GetAuditFunc   func(string, string) (json.RawMessage, error)

	// Member
	ListMembersFunc          func(string) (json.RawMessage, error)
	AddMemberFunc            func(string, json.RawMessage) (json.RawMessage, error)
	RemoveMemberFunc         func(string, string) error
	GetMemberPermissionsFunc func(string) (json.RawMessage, error)

	// ExtensionGrant
	ListExtensionGrantsFunc  func(string) ([]json.RawMessage, error)
	GetExtensionGrantFunc    func(string, string) (json.RawMessage, error)
	CreateExtensionGrantFunc func(string, json.RawMessage) (json.RawMessage, error)
	UpdateExtensionGrantFunc func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteExtensionGrantFunc func(string, string) error

	// Resource
	ListResourcesFunc  func(string) ([]json.RawMessage, error)
	GetResourceFunc    func(string, string) (json.RawMessage, error)
	CreateResourceFunc func(string, json.RawMessage) (json.RawMessage, error)
	UpdateResourceFunc func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteResourceFunc func(string, string) error

	// Reporter
	ListReportersFunc  func(string) ([]json.RawMessage, error)
	GetReporterFunc    func(string, string) (json.RawMessage, error)
	CreateReporterFunc func(string, json.RawMessage) (json.RawMessage, error)
	UpdateReporterFunc func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteReporterFunc func(string, string) error

	// BotDetection
	ListBotDetectionsFunc  func(string) ([]json.RawMessage, error)
	GetBotDetectionFunc    func(string, string) (json.RawMessage, error)
	CreateBotDetectionFunc func(string, json.RawMessage) (json.RawMessage, error)
	UpdateBotDetectionFunc func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteBotDetectionFunc func(string, string) error

	// DeviceIdentifier
	ListDeviceIdentifiersFunc  func(string) ([]json.RawMessage, error)
	GetDeviceIdentifierFunc    func(string, string) (json.RawMessage, error)
	CreateDeviceIdentifierFunc func(string, json.RawMessage) (json.RawMessage, error)
	UpdateDeviceIdentifierFunc func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteDeviceIdentifierFunc func(string, string) error

	// AuthDeviceNotifier
	ListAuthDeviceNotifiersFunc  func(string) ([]json.RawMessage, error)
	GetAuthDeviceNotifierFunc    func(string, string) (json.RawMessage, error)
	CreateAuthDeviceNotifierFunc func(string, json.RawMessage) (json.RawMessage, error)
	UpdateAuthDeviceNotifierFunc func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteAuthDeviceNotifierFunc func(string, string) error

	// AuthorizationEngine
	ListAuthorizationEnginesFunc  func(string) ([]json.RawMessage, error)
	GetAuthorizationEngineFunc    func(string, string) (json.RawMessage, error)
	UpdateAuthorizationEngineFunc func(string, string, json.RawMessage) (json.RawMessage, error)

	// ProtectedResource
	ListProtectedResourcesFunc        func(string) ([]json.RawMessage, error)
	GetProtectedResourceFunc          func(string, string) (json.RawMessage, error)
	CreateProtectedResourceFunc       func(string, json.RawMessage) (json.RawMessage, error)
	UpdateProtectedResourceFunc       func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteProtectedResourceFunc       func(string, string) error
	ListProtectedResourceMembersFunc  func(string, string) (json.RawMessage, error)
	RemoveProtectedResourceMemberFunc func(string, string, string) error
	ListProtectedResourceSecretsFunc  func(string, string) (json.RawMessage, error)

	// Analytics
	GetAnalyticsFunc func(string, AnalyticsParams) (json.RawMessage, error)

	// Entrypoint
	GetEntrypointsFunc   func(string) (json.RawMessage, error)
	CreateEntrypointFunc func(json.RawMessage) (json.RawMessage, error)
	UpdateEntrypointFunc func(string, json.RawMessage) (json.RawMessage, error)
	DeleteEntrypointFunc func(string) error

	// Dictionary
	ListDictionariesFunc        func(string) ([]json.RawMessage, error)
	GetDictionaryFunc           func(string, string) (json.RawMessage, error)
	CreateDictionaryFunc        func(string, json.RawMessage) (json.RawMessage, error)
	UpdateDictionaryFunc        func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteDictionaryFunc        func(string, string) error
	ListDictionaryEntriesFunc   func(string, string) (json.RawMessage, error)
	UpdateDictionaryEntriesFunc func(string, string, json.RawMessage) (json.RawMessage, error)

	// Alert
	ListAlertNotifiersFunc  func(string) ([]json.RawMessage, error)
	GetAlertNotifierFunc    func(string, string) (json.RawMessage, error)
	CreateAlertNotifierFunc func(string, json.RawMessage) (json.RawMessage, error)
	UpdateAlertNotifierFunc func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteAlertNotifierFunc func(string, string) error
	GetAlertTriggersFunc    func(string) (json.RawMessage, error)
	UpdateAlertTriggersFunc func(string, json.RawMessage) (json.RawMessage, error)

	// Organization users
	ListOrgUsersFunc         func(ListOrgUsersParams) (*PaginatedResponse, error)
	GetOrgUserFunc           func(string) (json.RawMessage, error)
	CreateOrgUserFunc        func(json.RawMessage) (json.RawMessage, error)
	UpdateOrgUserFunc        func(string, json.RawMessage) (json.RawMessage, error)
	DeleteOrgUserFunc        func(string) error
	ResetOrgUserPasswordFunc func(string, json.RawMessage) error
	UpdateOrgUserStatusFunc  func(string, json.RawMessage) (json.RawMessage, error)
	UpdateOrgUsernameFunc    func(string, json.RawMessage) (json.RawMessage, error)
	BulkOrgUserOperationFunc func(json.RawMessage) (json.RawMessage, error)

	// Organization groups
	ListOrgGroupsFunc  func() (json.RawMessage, error)
	GetOrgGroupFunc    func(string) (json.RawMessage, error)
	CreateOrgGroupFunc func(json.RawMessage) (json.RawMessage, error)
	UpdateOrgGroupFunc func(string, json.RawMessage) (json.RawMessage, error)
	DeleteOrgGroupFunc func(string) error

	// Organization roles
	ListOrgRolesFunc  func() ([]json.RawMessage, error)
	GetOrgRoleFunc    func(string) (json.RawMessage, error)
	CreateOrgRoleFunc func(json.RawMessage) (json.RawMessage, error)
	UpdateOrgRoleFunc func(string, json.RawMessage) (json.RawMessage, error)
	DeleteOrgRoleFunc func(string) error

	// Organization settings
	GetOrgSettingsFunc   func() (json.RawMessage, error)
	PatchOrgSettingsFunc func(json.RawMessage) (json.RawMessage, error)

	// Organization members
	ListOrgMembersFunc  func() (json.RawMessage, error)
	AddOrgMemberFunc    func(json.RawMessage) (json.RawMessage, error)
	RemoveOrgMemberFunc func(string) error

	// Organization audits
	ListOrgAuditsFunc func(ListOrgAuditsParams) (*PaginatedResponse, error)
	GetOrgAuditFunc   func(string) (json.RawMessage, error)

	// Organization reporters
	ListOrgReportersFunc  func() ([]json.RawMessage, error)
	GetOrgReporterFunc    func(string) (json.RawMessage, error)
	CreateOrgReporterFunc func(json.RawMessage) (json.RawMessage, error)
	UpdateOrgReporterFunc func(string, json.RawMessage) (json.RawMessage, error)
	DeleteOrgReporterFunc func(string) error

	// Organization forms
	GetOrgFormFunc    func(string) (json.RawMessage, error)
	CreateOrgFormFunc func(json.RawMessage) (json.RawMessage, error)
	UpdateOrgFormFunc func(string, json.RawMessage) (json.RawMessage, error)
	DeleteOrgFormFunc func(string) error

	// Organization identity providers
	ListOrgIdentityProvidersFunc  func() ([]json.RawMessage, error)
	GetOrgIdentityProviderFunc    func(string) (json.RawMessage, error)
	CreateOrgIdentityProviderFunc func(json.RawMessage) (json.RawMessage, error)
	UpdateOrgIdentityProviderFunc func(string, json.RawMessage) (json.RawMessage, error)
	DeleteOrgIdentityProviderFunc func(string) error

	// Organization entrypoints
	ListOrgEntrypointsFunc  func() ([]json.RawMessage, error)
	GetOrgEntrypointFunc    func(string) (json.RawMessage, error)
	CreateOrgEntrypointFunc func(json.RawMessage) (json.RawMessage, error)
	UpdateOrgEntrypointFunc func(string, json.RawMessage) (json.RawMessage, error)
	DeleteOrgEntrypointFunc func(string) error

	// Organization tags
	ListOrgTagsFunc  func() ([]json.RawMessage, error)
	GetOrgTagFunc    func(string) (json.RawMessage, error)
	CreateOrgTagFunc func(json.RawMessage) (json.RawMessage, error)
	UpdateOrgTagFunc func(string, json.RawMessage) (json.RawMessage, error)
	DeleteOrgTagFunc func(string) error

	// Organization user tokens
	ListOrgUserTokensFunc  func(string) (json.RawMessage, error)
	CreateOrgUserTokenFunc func(string, json.RawMessage) (json.RawMessage, error)
	RevokeOrgUserTokenFunc func(string, string) error
}

func unexpected(name string) error { return fmt.Errorf("unexpected call: %s", name) }

// Domain

func (m *MockService) ListDomains(p ListDomainsParams) (*PaginatedResponse, error) {
	if m.ListDomainsFunc != nil {
		return m.ListDomainsFunc(p)
	}

	return nil, unexpected("ListDomains")
}

func (m *MockService) GetDomain(id string) (json.RawMessage, error) {
	if m.GetDomainFunc != nil {
		return m.GetDomainFunc(id)
	}

	return nil, unexpected("GetDomain")
}

func (m *MockService) CreateDomain(body json.RawMessage) (json.RawMessage, error) {
	if m.CreateDomainFunc != nil {
		return m.CreateDomainFunc(body)
	}

	return nil, unexpected("CreateDomain")
}

func (m *MockService) UpdateDomain(id string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateDomainFunc != nil {
		return m.UpdateDomainFunc(id, body)
	}

	return nil, unexpected("UpdateDomain")
}

func (m *MockService) PatchDomain(id string, body json.RawMessage) (json.RawMessage, error) {
	if m.PatchDomainFunc != nil {
		return m.PatchDomainFunc(id, body)
	}

	return nil, unexpected("PatchDomain")
}

func (m *MockService) DeleteDomain(id string) error {
	if m.DeleteDomainFunc != nil {
		return m.DeleteDomainFunc(id)
	}

	return unexpected("DeleteDomain")
}

func (m *MockService) GetDomainByHRID(hrid string) (json.RawMessage, error) {
	if m.GetDomainByHRIDFunc != nil {
		return m.GetDomainByHRIDFunc(hrid)
	}

	return nil, unexpected("GetDomainByHRID")
}

func (m *MockService) UpdateDomainCertificateSettings(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateDomainCertificateSettingsFunc != nil {
		return m.UpdateDomainCertificateSettingsFunc(domainID, body)
	}

	return nil, unexpected("UpdateDomainCertificateSettings")
}

func (m *MockService) ListDataPlanes() (json.RawMessage, error) {
	if m.ListDataPlanesFunc != nil {
		return m.ListDataPlanesFunc()
	}

	return nil, unexpected("ListDataPlanes")
}

// Application

func (m *MockService) ListApplications(domainID string, p ListApplicationsParams) (*PaginatedResponse, error) {
	if m.ListApplicationsFunc != nil {
		return m.ListApplicationsFunc(domainID, p)
	}

	return nil, unexpected("ListApplications")
}

func (m *MockService) GetApplication(domainID, appID string) (json.RawMessage, error) {
	if m.GetApplicationFunc != nil {
		return m.GetApplicationFunc(domainID, appID)
	}

	return nil, unexpected("GetApplication")
}

func (m *MockService) CreateApplication(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateApplicationFunc != nil {
		return m.CreateApplicationFunc(domainID, body)
	}

	return nil, unexpected("CreateApplication")
}

func (m *MockService) UpdateApplication(domainID, appID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateApplicationFunc != nil {
		return m.UpdateApplicationFunc(domainID, appID, body)
	}

	return nil, unexpected("UpdateApplication")
}

func (m *MockService) PatchApplication(domainID, appID string, body json.RawMessage) (json.RawMessage, error) {
	if m.PatchApplicationFunc != nil {
		return m.PatchApplicationFunc(domainID, appID, body)
	}

	return nil, unexpected("PatchApplication")
}

func (m *MockService) DeleteApplication(domainID, appID string) error {
	if m.DeleteApplicationFunc != nil {
		return m.DeleteApplicationFunc(domainID, appID)
	}

	return unexpected("DeleteApplication")
}

// Application sub-resources

func (m *MockService) ListAppSecrets(domainID, appID string) (json.RawMessage, error) {
	if m.ListAppSecretsFunc != nil {
		return m.ListAppSecretsFunc(domainID, appID)
	}

	return nil, unexpected("ListAppSecrets")
}

func (m *MockService) CreateAppSecret(domainID, appID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateAppSecretFunc != nil {
		return m.CreateAppSecretFunc(domainID, appID, body)
	}

	return nil, unexpected("CreateAppSecret")
}

func (m *MockService) DeleteAppSecret(domainID, appID, secretID string) error {
	if m.DeleteAppSecretFunc != nil {
		return m.DeleteAppSecretFunc(domainID, appID, secretID)
	}

	return unexpected("DeleteAppSecret")
}

func (m *MockService) ListAppMembers(domainID, appID string) (json.RawMessage, error) {
	if m.ListAppMembersFunc != nil {
		return m.ListAppMembersFunc(domainID, appID)
	}

	return nil, unexpected("ListAppMembers")
}

func (m *MockService) AddAppMember(domainID, appID string, body json.RawMessage) (json.RawMessage, error) {
	if m.AddAppMemberFunc != nil {
		return m.AddAppMemberFunc(domainID, appID, body)
	}

	return nil, unexpected("AddAppMember")
}

func (m *MockService) RemoveAppMember(domainID, appID, memberID string) error {
	if m.RemoveAppMemberFunc != nil {
		return m.RemoveAppMemberFunc(domainID, appID, memberID)
	}

	return unexpected("RemoveAppMember")
}

func (m *MockService) ListAppFlows(domainID, appID string) ([]json.RawMessage, error) {
	if m.ListAppFlowsFunc != nil {
		return m.ListAppFlowsFunc(domainID, appID)
	}

	return nil, unexpected("ListAppFlows")
}

func (m *MockService) GetAppFlow(domainID, appID, flowID string) (json.RawMessage, error) {
	if m.GetAppFlowFunc != nil {
		return m.GetAppFlowFunc(domainID, appID, flowID)
	}

	return nil, unexpected("GetAppFlow")
}

func (m *MockService) UpdateAppFlows(domainID, appID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateAppFlowsFunc != nil {
		return m.UpdateAppFlowsFunc(domainID, appID, body)
	}

	return nil, unexpected("UpdateAppFlows")
}

func (m *MockService) GetAppEmail(domainID, appID, template string) (json.RawMessage, error) {
	if m.GetAppEmailFunc != nil {
		return m.GetAppEmailFunc(domainID, appID, template)
	}

	return nil, unexpected("GetAppEmail")
}

func (m *MockService) CreateAppEmail(domainID, appID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateAppEmailFunc != nil {
		return m.CreateAppEmailFunc(domainID, appID, body)
	}

	return nil, unexpected("CreateAppEmail")
}

func (m *MockService) UpdateAppEmail(domainID, appID, emailID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateAppEmailFunc != nil {
		return m.UpdateAppEmailFunc(domainID, appID, emailID, body)
	}

	return nil, unexpected("UpdateAppEmail")
}

func (m *MockService) DeleteAppEmail(domainID, appID, emailID string) error {
	if m.DeleteAppEmailFunc != nil {
		return m.DeleteAppEmailFunc(domainID, appID, emailID)
	}

	return unexpected("DeleteAppEmail")
}

func (m *MockService) GetAppForm(domainID, appID, template string) (json.RawMessage, error) {
	if m.GetAppFormFunc != nil {
		return m.GetAppFormFunc(domainID, appID, template)
	}

	return nil, unexpected("GetAppForm")
}

func (m *MockService) CreateAppForm(domainID, appID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateAppFormFunc != nil {
		return m.CreateAppFormFunc(domainID, appID, body)
	}

	return nil, unexpected("CreateAppForm")
}

func (m *MockService) UpdateAppForm(domainID, appID, formID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateAppFormFunc != nil {
		return m.UpdateAppFormFunc(domainID, appID, formID, body)
	}

	return nil, unexpected("UpdateAppForm")
}

func (m *MockService) DeleteAppForm(domainID, appID, formID string) error {
	if m.DeleteAppFormFunc != nil {
		return m.DeleteAppFormFunc(domainID, appID, formID)
	}

	return unexpected("DeleteAppForm")
}

func (m *MockService) ListAppResources(domainID, appID string) (json.RawMessage, error) {
	if m.ListAppResourcesFunc != nil {
		return m.ListAppResourcesFunc(domainID, appID)
	}

	return nil, unexpected("ListAppResources")
}

func (m *MockService) GetAppResource(domainID, appID, resourceID string) (json.RawMessage, error) {
	if m.GetAppResourceFunc != nil {
		return m.GetAppResourceFunc(domainID, appID, resourceID)
	}

	return nil, unexpected("GetAppResource")
}

func (m *MockService) GetAppAnalytics(domainID, appID string, params AnalyticsParams) (json.RawMessage, error) {
	if m.GetAppAnalyticsFunc != nil {
		return m.GetAppAnalyticsFunc(domainID, appID, params)
	}

	return nil, unexpected("GetAppAnalytics")
}

func (m *MockService) ChangeAppType(domainID, appID string, body json.RawMessage) (json.RawMessage, error) {
	if m.ChangeAppTypeFunc != nil {
		return m.ChangeAppTypeFunc(domainID, appID, body)
	}

	return nil, unexpected("ChangeAppType")
}

func (m *MockService) ListAppResourcePolicies(domainID, appID, resourceID string) (json.RawMessage, error) {
	if m.ListAppResourcePoliciesFunc != nil {
		return m.ListAppResourcePoliciesFunc(domainID, appID, resourceID)
	}

	return nil, unexpected("ListAppResourcePolicies")
}

func (m *MockService) GetAppResourcePolicy(domainID, appID, resourceID, policyID string) (json.RawMessage, error) {
	if m.GetAppResourcePolicyFunc != nil {
		return m.GetAppResourcePolicyFunc(domainID, appID, resourceID, policyID)
	}

	return nil, unexpected("GetAppResourcePolicy")
}

func (m *MockService) GetAppMemberPermissions(domainID, appID string) (json.RawMessage, error) {
	if m.GetAppMemberPermissionsFunc != nil {
		return m.GetAppMemberPermissionsFunc(domainID, appID)
	}

	return nil, unexpected("GetAppMemberPermissions")
}

// User

func (m *MockService) ListUsers(domainID string, p ListUsersParams) (*PaginatedResponse, error) {
	if m.ListUsersFunc != nil {
		return m.ListUsersFunc(domainID, p)
	}

	return nil, unexpected("ListUsers")
}

func (m *MockService) GetUser(domainID, userID string) (json.RawMessage, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(domainID, userID)
	}

	return nil, unexpected("GetUser")
}

func (m *MockService) CreateUser(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(domainID, body)
	}

	return nil, unexpected("CreateUser")
}

func (m *MockService) UpdateUser(domainID, userID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(domainID, userID, body)
	}

	return nil, unexpected("UpdateUser")
}

func (m *MockService) DeleteUser(domainID, userID string) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(domainID, userID)
	}

	return unexpected("DeleteUser")
}

func (m *MockService) UpdateUserStatus(domainID, userID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateUserStatusFunc != nil {
		return m.UpdateUserStatusFunc(domainID, userID, body)
	}

	return nil, unexpected("UpdateUserStatus")
}

func (m *MockService) ResetPassword(domainID, userID string, body json.RawMessage) error {
	if m.ResetPasswordFunc != nil {
		return m.ResetPasswordFunc(domainID, userID, body)
	}

	return unexpected("ResetPassword")
}

// User sub-resources

func (m *MockService) ListUserConsents(domainID, userID string) ([]json.RawMessage, error) {
	if m.ListUserConsentsFunc != nil {
		return m.ListUserConsentsFunc(domainID, userID)
	}

	return nil, unexpected("ListUserConsents")
}

func (m *MockService) RevokeUserConsent(domainID, userID, consentID string) error {
	if m.RevokeUserConsentFunc != nil {
		return m.RevokeUserConsentFunc(domainID, userID, consentID)
	}

	return unexpected("RevokeUserConsent")
}

func (m *MockService) RevokeAllUserConsents(domainID, userID string) error {
	if m.RevokeAllUserConsentsFunc != nil {
		return m.RevokeAllUserConsentsFunc(domainID, userID)
	}

	return unexpected("RevokeAllUserConsents")
}

func (m *MockService) ListUserRoles(domainID, userID string) (json.RawMessage, error) {
	if m.ListUserRolesFunc != nil {
		return m.ListUserRolesFunc(domainID, userID)
	}

	return nil, unexpected("ListUserRoles")
}

func (m *MockService) AssignUserRoles(domainID, userID string, body json.RawMessage) error {
	if m.AssignUserRolesFunc != nil {
		return m.AssignUserRolesFunc(domainID, userID, body)
	}

	return unexpected("AssignUserRoles")
}

func (m *MockService) RevokeUserRole(domainID, userID, roleID string) error {
	if m.RevokeUserRoleFunc != nil {
		return m.RevokeUserRoleFunc(domainID, userID, roleID)
	}

	return unexpected("RevokeUserRole")
}

func (m *MockService) ListUserDevices(domainID, userID string) ([]json.RawMessage, error) {
	if m.ListUserDevicesFunc != nil {
		return m.ListUserDevicesFunc(domainID, userID)
	}

	return nil, unexpected("ListUserDevices")
}

func (m *MockService) DeleteUserDevice(domainID, userID, deviceID string) error {
	if m.DeleteUserDeviceFunc != nil {
		return m.DeleteUserDeviceFunc(domainID, userID, deviceID)
	}

	return unexpected("DeleteUserDevice")
}

func (m *MockService) ListUserCredentials(domainID, userID string) ([]json.RawMessage, error) {
	if m.ListUserCredentialsFunc != nil {
		return m.ListUserCredentialsFunc(domainID, userID)
	}

	return nil, unexpected("ListUserCredentials")
}

func (m *MockService) GetUserCredential(domainID, userID, credentialID string) (json.RawMessage, error) {
	if m.GetUserCredentialFunc != nil {
		return m.GetUserCredentialFunc(domainID, userID, credentialID)
	}

	return nil, unexpected("GetUserCredential")
}

func (m *MockService) RevokeUserCredential(domainID, userID, credentialID string) error {
	if m.RevokeUserCredentialFunc != nil {
		return m.RevokeUserCredentialFunc(domainID, userID, credentialID)
	}

	return unexpected("RevokeUserCredential")
}

func (m *MockService) ListUserFactors(domainID, userID string) ([]json.RawMessage, error) {
	if m.ListUserFactorsFunc != nil {
		return m.ListUserFactorsFunc(domainID, userID)
	}

	return nil, unexpected("ListUserFactors")
}

func (m *MockService) DeleteUserFactor(domainID, userID, factorID string) error {
	if m.DeleteUserFactorFunc != nil {
		return m.DeleteUserFactorFunc(domainID, userID, factorID)
	}

	return unexpected("DeleteUserFactor")
}

func (m *MockService) ListUserAudits(domainID, userID string, p ListUserAuditsParams) (*PaginatedResponse, error) {
	if m.ListUserAuditsFunc != nil {
		return m.ListUserAuditsFunc(domainID, userID, p)
	}

	return nil, unexpected("ListUserAudits")
}

func (m *MockService) GetUserAudit(domainID, userID, auditID string) (json.RawMessage, error) {
	if m.GetUserAuditFunc != nil {
		return m.GetUserAuditFunc(domainID, userID, auditID)
	}

	return nil, unexpected("GetUserAudit")
}

func (m *MockService) SendRegistrationConfirmation(domainID, userID string) error {
	if m.SendRegistrationConfirmationFunc != nil {
		return m.SendRegistrationConfirmationFunc(domainID, userID)
	}

	return unexpected("SendRegistrationConfirmation")
}

func (m *MockService) UpdateUsername(domainID, userID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateUsernameFunc != nil {
		return m.UpdateUsernameFunc(domainID, userID, body)
	}

	return nil, unexpected("UpdateUsername")
}

func (m *MockService) ListUserIdentities(domainID, userID string) ([]json.RawMessage, error) {
	if m.ListUserIdentitiesFunc != nil {
		return m.ListUserIdentitiesFunc(domainID, userID)
	}

	return nil, unexpected("ListUserIdentities")
}

func (m *MockService) UnlinkUserIdentity(domainID, userID, identityID string) error {
	if m.UnlinkUserIdentityFunc != nil {
		return m.UnlinkUserIdentityFunc(domainID, userID, identityID)
	}

	return unexpected("UnlinkUserIdentity")
}

func (m *MockService) ListUserCertCredentials(domainID, userID string) ([]json.RawMessage, error) {
	if m.ListUserCertCredentialsFunc != nil {
		return m.ListUserCertCredentialsFunc(domainID, userID)
	}

	return nil, unexpected("ListUserCertCredentials")
}

func (m *MockService) GetUserCertCredential(domainID, userID, credID string) (json.RawMessage, error) {
	if m.GetUserCertCredentialFunc != nil {
		return m.GetUserCertCredentialFunc(domainID, userID, credID)
	}

	return nil, unexpected("GetUserCertCredential")
}

func (m *MockService) EnrollUserCertCredential(domainID, userID string, body json.RawMessage) (json.RawMessage, error) {
	if m.EnrollUserCertCredentialFunc != nil {
		return m.EnrollUserCertCredentialFunc(domainID, userID, body)
	}

	return nil, unexpected("EnrollUserCertCredential")
}

func (m *MockService) RevokeUserCertCredential(domainID, userID, credID string) error {
	if m.RevokeUserCertCredentialFunc != nil {
		return m.RevokeUserCertCredentialFunc(domainID, userID, credID)
	}

	return unexpected("RevokeUserCertCredential")
}

func (m *MockService) BulkUserOperation(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.BulkUserOperationFunc != nil {
		return m.BulkUserOperationFunc(domainID, body)
	}

	return nil, unexpected("BulkUserOperation")
}

// Role

func (m *MockService) ListRoles(domainID string, p ListRolesParams) (*PaginatedResponse, error) {
	if m.ListRolesFunc != nil {
		return m.ListRolesFunc(domainID, p)
	}

	return nil, unexpected("ListRoles")
}

func (m *MockService) GetRole(domainID, roleID string) (json.RawMessage, error) {
	if m.GetRoleFunc != nil {
		return m.GetRoleFunc(domainID, roleID)
	}

	return nil, unexpected("GetRole")
}

func (m *MockService) CreateRole(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateRoleFunc != nil {
		return m.CreateRoleFunc(domainID, body)
	}

	return nil, unexpected("CreateRole")
}

func (m *MockService) UpdateRole(domainID, roleID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateRoleFunc != nil {
		return m.UpdateRoleFunc(domainID, roleID, body)
	}

	return nil, unexpected("UpdateRole")
}

func (m *MockService) DeleteRole(domainID, roleID string) error {
	if m.DeleteRoleFunc != nil {
		return m.DeleteRoleFunc(domainID, roleID)
	}

	return unexpected("DeleteRole")
}

// Scope

func (m *MockService) ListScopes(domainID string, p ListScopesParams) (*PaginatedResponse, error) {
	if m.ListScopesFunc != nil {
		return m.ListScopesFunc(domainID, p)
	}

	return nil, unexpected("ListScopes")
}

func (m *MockService) GetScope(domainID, scopeID string) (json.RawMessage, error) {
	if m.GetScopeFunc != nil {
		return m.GetScopeFunc(domainID, scopeID)
	}

	return nil, unexpected("GetScope")
}

func (m *MockService) CreateScope(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateScopeFunc != nil {
		return m.CreateScopeFunc(domainID, body)
	}

	return nil, unexpected("CreateScope")
}

func (m *MockService) UpdateScope(domainID, scopeID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateScopeFunc != nil {
		return m.UpdateScopeFunc(domainID, scopeID, body)
	}

	return nil, unexpected("UpdateScope")
}

func (m *MockService) PatchScope(domainID, scopeID string, body json.RawMessage) (json.RawMessage, error) {
	if m.PatchScopeFunc != nil {
		return m.PatchScopeFunc(domainID, scopeID, body)
	}

	return nil, unexpected("PatchScope")
}

func (m *MockService) DeleteScope(domainID, scopeID string) error {
	if m.DeleteScopeFunc != nil {
		return m.DeleteScopeFunc(domainID, scopeID)
	}

	return unexpected("DeleteScope")
}

// IdentityProvider

func (m *MockService) ListIdentityProviders(domainID string, userProvider bool) ([]json.RawMessage, error) {
	if m.ListIdentityProvidersFunc != nil {
		return m.ListIdentityProvidersFunc(domainID, userProvider)
	}

	return nil, unexpected("ListIdentityProviders")
}

func (m *MockService) GetIdentityProvider(domainID, idpID string) (json.RawMessage, error) {
	if m.GetIdentityProviderFunc != nil {
		return m.GetIdentityProviderFunc(domainID, idpID)
	}

	return nil, unexpected("GetIdentityProvider")
}

func (m *MockService) CreateIdentityProvider(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateIdentityProviderFunc != nil {
		return m.CreateIdentityProviderFunc(domainID, body)
	}

	return nil, unexpected("CreateIdentityProvider")
}

func (m *MockService) UpdateIdentityProvider(domainID, idpID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateIdentityProviderFunc != nil {
		return m.UpdateIdentityProviderFunc(domainID, idpID, body)
	}

	return nil, unexpected("UpdateIdentityProvider")
}

func (m *MockService) DeleteIdentityProvider(domainID, idpID string) error {
	if m.DeleteIdentityProviderFunc != nil {
		return m.DeleteIdentityProviderFunc(domainID, idpID)
	}

	return unexpected("DeleteIdentityProvider")
}

func (m *MockService) UpdateIDPPasswordPolicy(domainID, idpID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateIDPPasswordPolicyFunc != nil {
		return m.UpdateIDPPasswordPolicyFunc(domainID, idpID, body)
	}

	return nil, unexpected("UpdateIDPPasswordPolicy")
}

// Certificate

func (m *MockService) ListCertificates(domainID string) ([]json.RawMessage, error) {
	if m.ListCertificatesFunc != nil {
		return m.ListCertificatesFunc(domainID)
	}

	return nil, unexpected("ListCertificates")
}

func (m *MockService) GetCertificate(domainID, certID string) (json.RawMessage, error) {
	if m.GetCertificateFunc != nil {
		return m.GetCertificateFunc(domainID, certID)
	}

	return nil, unexpected("GetCertificate")
}

func (m *MockService) CreateCertificate(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateCertificateFunc != nil {
		return m.CreateCertificateFunc(domainID, body)
	}

	return nil, unexpected("CreateCertificate")
}

func (m *MockService) UpdateCertificate(domainID, certID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateCertificateFunc != nil {
		return m.UpdateCertificateFunc(domainID, certID, body)
	}

	return nil, unexpected("UpdateCertificate")
}

func (m *MockService) DeleteCertificate(domainID, certID string) error {
	if m.DeleteCertificateFunc != nil {
		return m.DeleteCertificateFunc(domainID, certID)
	}

	return unexpected("DeleteCertificate")
}

func (m *MockService) GetCertificateKey(domainID, certID string) (json.RawMessage, error) {
	if m.GetCertificateKeyFunc != nil {
		return m.GetCertificateKeyFunc(domainID, certID)
	}

	return nil, unexpected("GetCertificateKey")
}

func (m *MockService) GetCertificateKeys(domainID, certID string) (json.RawMessage, error) {
	if m.GetCertificateKeysFunc != nil {
		return m.GetCertificateKeysFunc(domainID, certID)
	}

	return nil, unexpected("GetCertificateKeys")
}

func (m *MockService) RotateCertificates(domainID string) (json.RawMessage, error) {
	if m.RotateCertificatesFunc != nil {
		return m.RotateCertificatesFunc(domainID)
	}

	return nil, unexpected("RotateCertificates")
}

// Factor

func (m *MockService) ListFactors(domainID string) ([]json.RawMessage, error) {
	if m.ListFactorsFunc != nil {
		return m.ListFactorsFunc(domainID)
	}

	return nil, unexpected("ListFactors")
}

func (m *MockService) GetFactor(domainID, factorID string) (json.RawMessage, error) {
	if m.GetFactorFunc != nil {
		return m.GetFactorFunc(domainID, factorID)
	}

	return nil, unexpected("GetFactor")
}

func (m *MockService) CreateFactor(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateFactorFunc != nil {
		return m.CreateFactorFunc(domainID, body)
	}

	return nil, unexpected("CreateFactor")
}

func (m *MockService) UpdateFactor(domainID, factorID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateFactorFunc != nil {
		return m.UpdateFactorFunc(domainID, factorID, body)
	}

	return nil, unexpected("UpdateFactor")
}

func (m *MockService) DeleteFactor(domainID, factorID string) error {
	if m.DeleteFactorFunc != nil {
		return m.DeleteFactorFunc(domainID, factorID)
	}

	return unexpected("DeleteFactor")
}

// Group

func (m *MockService) ListGroups(domainID string, p ListGroupsParams) (*PaginatedResponse, error) {
	if m.ListGroupsFunc != nil {
		return m.ListGroupsFunc(domainID, p)
	}

	return nil, unexpected("ListGroups")
}

func (m *MockService) GetGroup(domainID, groupID string) (json.RawMessage, error) {
	if m.GetGroupFunc != nil {
		return m.GetGroupFunc(domainID, groupID)
	}

	return nil, unexpected("GetGroup")
}

func (m *MockService) CreateGroup(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateGroupFunc != nil {
		return m.CreateGroupFunc(domainID, body)
	}

	return nil, unexpected("CreateGroup")
}

func (m *MockService) UpdateGroup(domainID, groupID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateGroupFunc != nil {
		return m.UpdateGroupFunc(domainID, groupID, body)
	}

	return nil, unexpected("UpdateGroup")
}

func (m *MockService) DeleteGroup(domainID, groupID string) error {
	if m.DeleteGroupFunc != nil {
		return m.DeleteGroupFunc(domainID, groupID)
	}

	return unexpected("DeleteGroup")
}

func (m *MockService) ListGroupMembers(domainID, groupID string) (json.RawMessage, error) {
	if m.ListGroupMembersFunc != nil {
		return m.ListGroupMembersFunc(domainID, groupID)
	}

	return nil, unexpected("ListGroupMembers")
}

func (m *MockService) AddGroupMember(domainID, groupID, memberID string) error {
	if m.AddGroupMemberFunc != nil {
		return m.AddGroupMemberFunc(domainID, groupID, memberID)
	}

	return unexpected("AddGroupMember")
}

func (m *MockService) RemoveGroupMember(domainID, groupID, memberID string) error {
	if m.RemoveGroupMemberFunc != nil {
		return m.RemoveGroupMemberFunc(domainID, groupID, memberID)
	}

	return unexpected("RemoveGroupMember")
}

func (m *MockService) ListGroupRoles(domainID, groupID string) (json.RawMessage, error) {
	if m.ListGroupRolesFunc != nil {
		return m.ListGroupRolesFunc(domainID, groupID)
	}

	return nil, unexpected("ListGroupRoles")
}

func (m *MockService) AssignGroupRoles(domainID, groupID string, body json.RawMessage) (json.RawMessage, error) {
	if m.AssignGroupRolesFunc != nil {
		return m.AssignGroupRolesFunc(domainID, groupID, body)
	}

	return nil, unexpected("AssignGroupRoles")
}

func (m *MockService) RevokeGroupRole(domainID, groupID, roleID string) error {
	if m.RevokeGroupRoleFunc != nil {
		return m.RevokeGroupRoleFunc(domainID, groupID, roleID)
	}

	return unexpected("RevokeGroupRole")
}

// Flow

func (m *MockService) ListFlows(domainID string) ([]json.RawMessage, error) {
	if m.ListFlowsFunc != nil {
		return m.ListFlowsFunc(domainID)
	}

	return nil, unexpected("ListFlows")
}

func (m *MockService) GetFlow(domainID, flowID string) (json.RawMessage, error) {
	if m.GetFlowFunc != nil {
		return m.GetFlowFunc(domainID, flowID)
	}

	return nil, unexpected("GetFlow")
}

func (m *MockService) UpdateFlows(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateFlowsFunc != nil {
		return m.UpdateFlowsFunc(domainID, body)
	}

	return nil, unexpected("UpdateFlows")
}

// Form

func (m *MockService) ListForms(domainID string) ([]json.RawMessage, error) {
	if m.ListFormsFunc != nil {
		return m.ListFormsFunc(domainID)
	}

	return nil, unexpected("ListForms")
}

func (m *MockService) GetForm(domainID, formID string) (json.RawMessage, error) {
	if m.GetFormFunc != nil {
		return m.GetFormFunc(domainID, formID)
	}

	return nil, unexpected("GetForm")
}

func (m *MockService) CreateForm(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateFormFunc != nil {
		return m.CreateFormFunc(domainID, body)
	}

	return nil, unexpected("CreateForm")
}

func (m *MockService) UpdateForm(domainID, formID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateFormFunc != nil {
		return m.UpdateFormFunc(domainID, formID, body)
	}

	return nil, unexpected("UpdateForm")
}

func (m *MockService) DeleteForm(domainID, formID string) error {
	if m.DeleteFormFunc != nil {
		return m.DeleteFormFunc(domainID, formID)
	}

	return unexpected("DeleteForm")
}

// Email

func (m *MockService) ListEmails(domainID string) ([]json.RawMessage, error) {
	if m.ListEmailsFunc != nil {
		return m.ListEmailsFunc(domainID)
	}

	return nil, unexpected("ListEmails")
}

func (m *MockService) GetEmail(domainID, emailID string) (json.RawMessage, error) {
	if m.GetEmailFunc != nil {
		return m.GetEmailFunc(domainID, emailID)
	}

	return nil, unexpected("GetEmail")
}

func (m *MockService) CreateEmail(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateEmailFunc != nil {
		return m.CreateEmailFunc(domainID, body)
	}

	return nil, unexpected("CreateEmail")
}

func (m *MockService) UpdateEmail(domainID, emailID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateEmailFunc != nil {
		return m.UpdateEmailFunc(domainID, emailID, body)
	}

	return nil, unexpected("UpdateEmail")
}

func (m *MockService) DeleteEmail(domainID, emailID string) error {
	if m.DeleteEmailFunc != nil {
		return m.DeleteEmailFunc(domainID, emailID)
	}

	return unexpected("DeleteEmail")
}

// Theme

func (m *MockService) ListThemes(domainID string) ([]json.RawMessage, error) {
	if m.ListThemesFunc != nil {
		return m.ListThemesFunc(domainID)
	}

	return nil, unexpected("ListThemes")
}

func (m *MockService) GetTheme(domainID, themeID string) (json.RawMessage, error) {
	if m.GetThemeFunc != nil {
		return m.GetThemeFunc(domainID, themeID)
	}

	return nil, unexpected("GetTheme")
}

func (m *MockService) CreateTheme(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateThemeFunc != nil {
		return m.CreateThemeFunc(domainID, body)
	}

	return nil, unexpected("CreateTheme")
}

func (m *MockService) UpdateTheme(domainID, themeID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateThemeFunc != nil {
		return m.UpdateThemeFunc(domainID, themeID, body)
	}

	return nil, unexpected("UpdateTheme")
}

func (m *MockService) DeleteTheme(domainID, themeID string) error {
	if m.DeleteThemeFunc != nil {
		return m.DeleteThemeFunc(domainID, themeID)
	}

	return unexpected("DeleteTheme")
}

// PasswordPolicy

func (m *MockService) ListPasswordPolicies(domainID string) ([]json.RawMessage, error) {
	if m.ListPasswordPoliciesFunc != nil {
		return m.ListPasswordPoliciesFunc(domainID)
	}

	return nil, unexpected("ListPasswordPolicies")
}

func (m *MockService) GetPasswordPolicy(domainID, policyID string) (json.RawMessage, error) {
	if m.GetPasswordPolicyFunc != nil {
		return m.GetPasswordPolicyFunc(domainID, policyID)
	}

	return nil, unexpected("GetPasswordPolicy")
}

func (m *MockService) CreatePasswordPolicy(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreatePasswordPolicyFunc != nil {
		return m.CreatePasswordPolicyFunc(domainID, body)
	}

	return nil, unexpected("CreatePasswordPolicy")
}

func (m *MockService) UpdatePasswordPolicy(domainID, policyID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdatePasswordPolicyFunc != nil {
		return m.UpdatePasswordPolicyFunc(domainID, policyID, body)
	}

	return nil, unexpected("UpdatePasswordPolicy")
}

func (m *MockService) DeletePasswordPolicy(domainID, policyID string) error {
	if m.DeletePasswordPolicyFunc != nil {
		return m.DeletePasswordPolicyFunc(domainID, policyID)
	}

	return unexpected("DeletePasswordPolicy")
}

func (m *MockService) GetActivePasswordPolicy(domainID string) (json.RawMessage, error) {
	if m.GetActivePasswordPolicyFunc != nil {
		return m.GetActivePasswordPolicyFunc(domainID)
	}

	return nil, unexpected("GetActivePasswordPolicy")
}

func (m *MockService) SetDefaultPasswordPolicy(domainID, policyID string) (json.RawMessage, error) {
	if m.SetDefaultPasswordPolicyFunc != nil {
		return m.SetDefaultPasswordPolicyFunc(domainID, policyID)
	}

	return nil, unexpected("SetDefaultPasswordPolicy")
}

func (m *MockService) EvaluatePasswordPolicy(domainID, policyID string, body json.RawMessage) (json.RawMessage, error) {
	if m.EvaluatePasswordPolicyFunc != nil {
		return m.EvaluatePasswordPolicyFunc(domainID, policyID, body)
	}

	return nil, unexpected("EvaluatePasswordPolicy")
}

// Audit

func (m *MockService) ListAudits(domainID string, p ListAuditsParams) (*PaginatedResponse, error) {
	if m.ListAuditsFunc != nil {
		return m.ListAuditsFunc(domainID, p)
	}

	return nil, unexpected("ListAudits")
}

func (m *MockService) GetAudit(domainID, auditID string) (json.RawMessage, error) {
	if m.GetAuditFunc != nil {
		return m.GetAuditFunc(domainID, auditID)
	}

	return nil, unexpected("GetAudit")
}

// Member

func (m *MockService) ListMembers(domainID string) (json.RawMessage, error) {
	if m.ListMembersFunc != nil {
		return m.ListMembersFunc(domainID)
	}

	return nil, unexpected("ListMembers")
}

func (m *MockService) AddMember(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.AddMemberFunc != nil {
		return m.AddMemberFunc(domainID, body)
	}

	return nil, unexpected("AddMember")
}

func (m *MockService) RemoveMember(domainID, memberID string) error {
	if m.RemoveMemberFunc != nil {
		return m.RemoveMemberFunc(domainID, memberID)
	}

	return unexpected("RemoveMember")
}

func (m *MockService) GetMemberPermissions(domainID string) (json.RawMessage, error) {
	if m.GetMemberPermissionsFunc != nil {
		return m.GetMemberPermissionsFunc(domainID)
	}

	return nil, unexpected("GetMemberPermissions")
}

// ExtensionGrant

func (m *MockService) ListExtensionGrants(domainID string) ([]json.RawMessage, error) {
	if m.ListExtensionGrantsFunc != nil {
		return m.ListExtensionGrantsFunc(domainID)
	}

	return nil, unexpected("ListExtensionGrants")
}

func (m *MockService) GetExtensionGrant(domainID, grantID string) (json.RawMessage, error) {
	if m.GetExtensionGrantFunc != nil {
		return m.GetExtensionGrantFunc(domainID, grantID)
	}

	return nil, unexpected("GetExtensionGrant")
}

func (m *MockService) CreateExtensionGrant(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateExtensionGrantFunc != nil {
		return m.CreateExtensionGrantFunc(domainID, body)
	}

	return nil, unexpected("CreateExtensionGrant")
}

func (m *MockService) UpdateExtensionGrant(domainID, grantID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateExtensionGrantFunc != nil {
		return m.UpdateExtensionGrantFunc(domainID, grantID, body)
	}

	return nil, unexpected("UpdateExtensionGrant")
}

func (m *MockService) DeleteExtensionGrant(domainID, grantID string) error {
	if m.DeleteExtensionGrantFunc != nil {
		return m.DeleteExtensionGrantFunc(domainID, grantID)
	}

	return unexpected("DeleteExtensionGrant")
}

// Resource

func (m *MockService) ListResources(domainID string) ([]json.RawMessage, error) {
	if m.ListResourcesFunc != nil {
		return m.ListResourcesFunc(domainID)
	}

	return nil, unexpected("ListResources")
}

func (m *MockService) GetResource(domainID, resourceID string) (json.RawMessage, error) {
	if m.GetResourceFunc != nil {
		return m.GetResourceFunc(domainID, resourceID)
	}

	return nil, unexpected("GetResource")
}

func (m *MockService) CreateResource(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateResourceFunc != nil {
		return m.CreateResourceFunc(domainID, body)
	}

	return nil, unexpected("CreateResource")
}

func (m *MockService) UpdateResource(domainID, resourceID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateResourceFunc != nil {
		return m.UpdateResourceFunc(domainID, resourceID, body)
	}

	return nil, unexpected("UpdateResource")
}

func (m *MockService) DeleteResource(domainID, resourceID string) error {
	if m.DeleteResourceFunc != nil {
		return m.DeleteResourceFunc(domainID, resourceID)
	}

	return unexpected("DeleteResource")
}

// Reporter

func (m *MockService) ListReporters(domainID string) ([]json.RawMessage, error) {
	if m.ListReportersFunc != nil {
		return m.ListReportersFunc(domainID)
	}

	return nil, unexpected("ListReporters")
}

func (m *MockService) GetReporter(domainID, reporterID string) (json.RawMessage, error) {
	if m.GetReporterFunc != nil {
		return m.GetReporterFunc(domainID, reporterID)
	}

	return nil, unexpected("GetReporter")
}

func (m *MockService) CreateReporter(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateReporterFunc != nil {
		return m.CreateReporterFunc(domainID, body)
	}

	return nil, unexpected("CreateReporter")
}

func (m *MockService) UpdateReporter(domainID, reporterID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateReporterFunc != nil {
		return m.UpdateReporterFunc(domainID, reporterID, body)
	}

	return nil, unexpected("UpdateReporter")
}

func (m *MockService) DeleteReporter(domainID, reporterID string) error {
	if m.DeleteReporterFunc != nil {
		return m.DeleteReporterFunc(domainID, reporterID)
	}

	return unexpected("DeleteReporter")
}

// BotDetection

func (m *MockService) ListBotDetections(domainID string) ([]json.RawMessage, error) {
	if m.ListBotDetectionsFunc != nil {
		return m.ListBotDetectionsFunc(domainID)
	}

	return nil, unexpected("ListBotDetections")
}

func (m *MockService) GetBotDetection(domainID, botDetectionID string) (json.RawMessage, error) {
	if m.GetBotDetectionFunc != nil {
		return m.GetBotDetectionFunc(domainID, botDetectionID)
	}

	return nil, unexpected("GetBotDetection")
}

func (m *MockService) CreateBotDetection(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateBotDetectionFunc != nil {
		return m.CreateBotDetectionFunc(domainID, body)
	}

	return nil, unexpected("CreateBotDetection")
}

func (m *MockService) UpdateBotDetection(domainID, botDetectionID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateBotDetectionFunc != nil {
		return m.UpdateBotDetectionFunc(domainID, botDetectionID, body)
	}

	return nil, unexpected("UpdateBotDetection")
}

func (m *MockService) DeleteBotDetection(domainID, botDetectionID string) error {
	if m.DeleteBotDetectionFunc != nil {
		return m.DeleteBotDetectionFunc(domainID, botDetectionID)
	}

	return unexpected("DeleteBotDetection")
}

// DeviceIdentifier

func (m *MockService) ListDeviceIdentifiers(domainID string) ([]json.RawMessage, error) {
	if m.ListDeviceIdentifiersFunc != nil {
		return m.ListDeviceIdentifiersFunc(domainID)
	}

	return nil, unexpected("ListDeviceIdentifiers")
}

func (m *MockService) GetDeviceIdentifier(domainID, deviceIdentifierID string) (json.RawMessage, error) {
	if m.GetDeviceIdentifierFunc != nil {
		return m.GetDeviceIdentifierFunc(domainID, deviceIdentifierID)
	}

	return nil, unexpected("GetDeviceIdentifier")
}

func (m *MockService) CreateDeviceIdentifier(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateDeviceIdentifierFunc != nil {
		return m.CreateDeviceIdentifierFunc(domainID, body)
	}

	return nil, unexpected("CreateDeviceIdentifier")
}

func (m *MockService) UpdateDeviceIdentifier(domainID, deviceIdentifierID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateDeviceIdentifierFunc != nil {
		return m.UpdateDeviceIdentifierFunc(domainID, deviceIdentifierID, body)
	}

	return nil, unexpected("UpdateDeviceIdentifier")
}

func (m *MockService) DeleteDeviceIdentifier(domainID, deviceIdentifierID string) error {
	if m.DeleteDeviceIdentifierFunc != nil {
		return m.DeleteDeviceIdentifierFunc(domainID, deviceIdentifierID)
	}

	return unexpected("DeleteDeviceIdentifier")
}

// AuthDeviceNotifier

func (m *MockService) ListAuthDeviceNotifiers(domainID string) ([]json.RawMessage, error) {
	if m.ListAuthDeviceNotifiersFunc != nil {
		return m.ListAuthDeviceNotifiersFunc(domainID)
	}

	return nil, unexpected("ListAuthDeviceNotifiers")
}

func (m *MockService) GetAuthDeviceNotifier(domainID, authDeviceNotifierID string) (json.RawMessage, error) {
	if m.GetAuthDeviceNotifierFunc != nil {
		return m.GetAuthDeviceNotifierFunc(domainID, authDeviceNotifierID)
	}

	return nil, unexpected("GetAuthDeviceNotifier")
}

func (m *MockService) CreateAuthDeviceNotifier(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateAuthDeviceNotifierFunc != nil {
		return m.CreateAuthDeviceNotifierFunc(domainID, body)
	}

	return nil, unexpected("CreateAuthDeviceNotifier")
}

func (m *MockService) UpdateAuthDeviceNotifier(domainID, authDeviceNotifierID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateAuthDeviceNotifierFunc != nil {
		return m.UpdateAuthDeviceNotifierFunc(domainID, authDeviceNotifierID, body)
	}

	return nil, unexpected("UpdateAuthDeviceNotifier")
}

func (m *MockService) DeleteAuthDeviceNotifier(domainID, authDeviceNotifierID string) error {
	if m.DeleteAuthDeviceNotifierFunc != nil {
		return m.DeleteAuthDeviceNotifierFunc(domainID, authDeviceNotifierID)
	}

	return unexpected("DeleteAuthDeviceNotifier")
}

// AuthorizationEngine

func (m *MockService) ListAuthorizationEngines(domainID string) ([]json.RawMessage, error) {
	if m.ListAuthorizationEnginesFunc != nil {
		return m.ListAuthorizationEnginesFunc(domainID)
	}

	return nil, unexpected("ListAuthorizationEngines")
}

func (m *MockService) GetAuthorizationEngine(domainID, engineID string) (json.RawMessage, error) {
	if m.GetAuthorizationEngineFunc != nil {
		return m.GetAuthorizationEngineFunc(domainID, engineID)
	}

	return nil, unexpected("GetAuthorizationEngine")
}

func (m *MockService) UpdateAuthorizationEngine(domainID, engineID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateAuthorizationEngineFunc != nil {
		return m.UpdateAuthorizationEngineFunc(domainID, engineID, body)
	}

	return nil, unexpected("UpdateAuthorizationEngine")
}

// ProtectedResource

func (m *MockService) ListProtectedResources(domainID string) ([]json.RawMessage, error) {
	if m.ListProtectedResourcesFunc != nil {
		return m.ListProtectedResourcesFunc(domainID)
	}

	return nil, unexpected("ListProtectedResources")
}

func (m *MockService) GetProtectedResource(domainID, protectedResourceID string) (json.RawMessage, error) {
	if m.GetProtectedResourceFunc != nil {
		return m.GetProtectedResourceFunc(domainID, protectedResourceID)
	}

	return nil, unexpected("GetProtectedResource")
}

func (m *MockService) CreateProtectedResource(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateProtectedResourceFunc != nil {
		return m.CreateProtectedResourceFunc(domainID, body)
	}

	return nil, unexpected("CreateProtectedResource")
}

func (m *MockService) UpdateProtectedResource(domainID, protectedResourceID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateProtectedResourceFunc != nil {
		return m.UpdateProtectedResourceFunc(domainID, protectedResourceID, body)
	}

	return nil, unexpected("UpdateProtectedResource")
}

func (m *MockService) DeleteProtectedResource(domainID, protectedResourceID string) error {
	if m.DeleteProtectedResourceFunc != nil {
		return m.DeleteProtectedResourceFunc(domainID, protectedResourceID)
	}

	return unexpected("DeleteProtectedResource")
}

func (m *MockService) ListProtectedResourceMembers(domainID, prID string) (json.RawMessage, error) {
	if m.ListProtectedResourceMembersFunc != nil {
		return m.ListProtectedResourceMembersFunc(domainID, prID)
	}

	return nil, unexpected("ListProtectedResourceMembers")
}

func (m *MockService) RemoveProtectedResourceMember(domainID, prID, memberID string) error {
	if m.RemoveProtectedResourceMemberFunc != nil {
		return m.RemoveProtectedResourceMemberFunc(domainID, prID, memberID)
	}

	return unexpected("RemoveProtectedResourceMember")
}

func (m *MockService) ListProtectedResourceSecrets(domainID, prID string) (json.RawMessage, error) {
	if m.ListProtectedResourceSecretsFunc != nil {
		return m.ListProtectedResourceSecretsFunc(domainID, prID)
	}

	return nil, unexpected("ListProtectedResourceSecrets")
}

// Analytics

func (m *MockService) GetAnalytics(domainID string, params AnalyticsParams) (json.RawMessage, error) {
	if m.GetAnalyticsFunc != nil {
		return m.GetAnalyticsFunc(domainID, params)
	}

	return nil, unexpected("GetAnalytics")
}

// Entrypoint

func (m *MockService) GetEntrypoints(domainID string) (json.RawMessage, error) {
	if m.GetEntrypointsFunc != nil {
		return m.GetEntrypointsFunc(domainID)
	}

	return nil, unexpected("GetEntrypoints")
}

func (m *MockService) CreateEntrypoint(body json.RawMessage) (json.RawMessage, error) {
	if m.CreateEntrypointFunc != nil {
		return m.CreateEntrypointFunc(body)
	}

	return nil, unexpected("CreateEntrypoint")
}

func (m *MockService) UpdateEntrypoint(entrypointID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateEntrypointFunc != nil {
		return m.UpdateEntrypointFunc(entrypointID, body)
	}

	return nil, unexpected("UpdateEntrypoint")
}

func (m *MockService) DeleteEntrypoint(entrypointID string) error {
	if m.DeleteEntrypointFunc != nil {
		return m.DeleteEntrypointFunc(entrypointID)
	}

	return unexpected("DeleteEntrypoint")
}

// Dictionary

func (m *MockService) ListDictionaries(domainID string) ([]json.RawMessage, error) {
	if m.ListDictionariesFunc != nil {
		return m.ListDictionariesFunc(domainID)
	}

	return nil, unexpected("ListDictionaries")
}

func (m *MockService) GetDictionary(domainID, dictID string) (json.RawMessage, error) {
	if m.GetDictionaryFunc != nil {
		return m.GetDictionaryFunc(domainID, dictID)
	}

	return nil, unexpected("GetDictionary")
}

func (m *MockService) CreateDictionary(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateDictionaryFunc != nil {
		return m.CreateDictionaryFunc(domainID, body)
	}

	return nil, unexpected("CreateDictionary")
}

func (m *MockService) UpdateDictionary(domainID, dictID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateDictionaryFunc != nil {
		return m.UpdateDictionaryFunc(domainID, dictID, body)
	}

	return nil, unexpected("UpdateDictionary")
}

func (m *MockService) DeleteDictionary(domainID, dictID string) error {
	if m.DeleteDictionaryFunc != nil {
		return m.DeleteDictionaryFunc(domainID, dictID)
	}

	return unexpected("DeleteDictionary")
}

func (m *MockService) ListDictionaryEntries(domainID, dictID string) (json.RawMessage, error) {
	if m.ListDictionaryEntriesFunc != nil {
		return m.ListDictionaryEntriesFunc(domainID, dictID)
	}

	return nil, unexpected("ListDictionaryEntries")
}

func (m *MockService) UpdateDictionaryEntries(domainID, dictID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateDictionaryEntriesFunc != nil {
		return m.UpdateDictionaryEntriesFunc(domainID, dictID, body)
	}

	return nil, unexpected("UpdateDictionaryEntries")
}

// Alert

func (m *MockService) ListAlertNotifiers(domainID string) ([]json.RawMessage, error) {
	if m.ListAlertNotifiersFunc != nil {
		return m.ListAlertNotifiersFunc(domainID)
	}

	return nil, unexpected("ListAlertNotifiers")
}

func (m *MockService) GetAlertNotifier(domainID, notifierID string) (json.RawMessage, error) {
	if m.GetAlertNotifierFunc != nil {
		return m.GetAlertNotifierFunc(domainID, notifierID)
	}

	return nil, unexpected("GetAlertNotifier")
}

func (m *MockService) CreateAlertNotifier(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateAlertNotifierFunc != nil {
		return m.CreateAlertNotifierFunc(domainID, body)
	}

	return nil, unexpected("CreateAlertNotifier")
}

func (m *MockService) UpdateAlertNotifier(domainID, notifierID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateAlertNotifierFunc != nil {
		return m.UpdateAlertNotifierFunc(domainID, notifierID, body)
	}

	return nil, unexpected("UpdateAlertNotifier")
}

func (m *MockService) DeleteAlertNotifier(domainID, notifierID string) error {
	if m.DeleteAlertNotifierFunc != nil {
		return m.DeleteAlertNotifierFunc(domainID, notifierID)
	}

	return unexpected("DeleteAlertNotifier")
}

func (m *MockService) GetAlertTriggers(domainID string) (json.RawMessage, error) {
	if m.GetAlertTriggersFunc != nil {
		return m.GetAlertTriggersFunc(domainID)
	}

	return nil, unexpected("GetAlertTriggers")
}

func (m *MockService) UpdateAlertTriggers(domainID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateAlertTriggersFunc != nil {
		return m.UpdateAlertTriggersFunc(domainID, body)
	}

	return nil, unexpected("UpdateAlertTriggers")
}

// Organization users

func (m *MockService) ListOrgUsers(p ListOrgUsersParams) (*PaginatedResponse, error) {
	if m.ListOrgUsersFunc != nil {
		return m.ListOrgUsersFunc(p)
	}

	return nil, unexpected("ListOrgUsers")
}

func (m *MockService) GetOrgUser(userID string) (json.RawMessage, error) {
	if m.GetOrgUserFunc != nil {
		return m.GetOrgUserFunc(userID)
	}

	return nil, unexpected("GetOrgUser")
}

func (m *MockService) CreateOrgUser(body json.RawMessage) (json.RawMessage, error) {
	if m.CreateOrgUserFunc != nil {
		return m.CreateOrgUserFunc(body)
	}

	return nil, unexpected("CreateOrgUser")
}

func (m *MockService) UpdateOrgUser(userID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateOrgUserFunc != nil {
		return m.UpdateOrgUserFunc(userID, body)
	}

	return nil, unexpected("UpdateOrgUser")
}

func (m *MockService) DeleteOrgUser(userID string) error {
	if m.DeleteOrgUserFunc != nil {
		return m.DeleteOrgUserFunc(userID)
	}

	return unexpected("DeleteOrgUser")
}

func (m *MockService) ResetOrgUserPassword(userID string, body json.RawMessage) error {
	if m.ResetOrgUserPasswordFunc != nil {
		return m.ResetOrgUserPasswordFunc(userID, body)
	}

	return unexpected("ResetOrgUserPassword")
}

func (m *MockService) UpdateOrgUserStatus(userID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateOrgUserStatusFunc != nil {
		return m.UpdateOrgUserStatusFunc(userID, body)
	}

	return nil, unexpected("UpdateOrgUserStatus")
}

func (m *MockService) UpdateOrgUsername(userID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateOrgUsernameFunc != nil {
		return m.UpdateOrgUsernameFunc(userID, body)
	}

	return nil, unexpected("UpdateOrgUsername")
}

func (m *MockService) BulkOrgUserOperation(body json.RawMessage) (json.RawMessage, error) {
	if m.BulkOrgUserOperationFunc != nil {
		return m.BulkOrgUserOperationFunc(body)
	}

	return nil, unexpected("BulkOrgUserOperation")
}

// Organization groups

func (m *MockService) ListOrgGroups() (json.RawMessage, error) {
	if m.ListOrgGroupsFunc != nil {
		return m.ListOrgGroupsFunc()
	}

	return nil, unexpected("ListOrgGroups")
}

func (m *MockService) GetOrgGroup(groupID string) (json.RawMessage, error) {
	if m.GetOrgGroupFunc != nil {
		return m.GetOrgGroupFunc(groupID)
	}

	return nil, unexpected("GetOrgGroup")
}

func (m *MockService) CreateOrgGroup(body json.RawMessage) (json.RawMessage, error) {
	if m.CreateOrgGroupFunc != nil {
		return m.CreateOrgGroupFunc(body)
	}

	return nil, unexpected("CreateOrgGroup")
}

func (m *MockService) UpdateOrgGroup(groupID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateOrgGroupFunc != nil {
		return m.UpdateOrgGroupFunc(groupID, body)
	}

	return nil, unexpected("UpdateOrgGroup")
}

func (m *MockService) DeleteOrgGroup(groupID string) error {
	if m.DeleteOrgGroupFunc != nil {
		return m.DeleteOrgGroupFunc(groupID)
	}

	return unexpected("DeleteOrgGroup")
}

// Organization roles

func (m *MockService) ListOrgRoles() ([]json.RawMessage, error) {
	if m.ListOrgRolesFunc != nil {
		return m.ListOrgRolesFunc()
	}

	return nil, unexpected("ListOrgRoles")
}

func (m *MockService) GetOrgRole(roleID string) (json.RawMessage, error) {
	if m.GetOrgRoleFunc != nil {
		return m.GetOrgRoleFunc(roleID)
	}

	return nil, unexpected("GetOrgRole")
}

func (m *MockService) CreateOrgRole(body json.RawMessage) (json.RawMessage, error) {
	if m.CreateOrgRoleFunc != nil {
		return m.CreateOrgRoleFunc(body)
	}

	return nil, unexpected("CreateOrgRole")
}

func (m *MockService) UpdateOrgRole(roleID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateOrgRoleFunc != nil {
		return m.UpdateOrgRoleFunc(roleID, body)
	}

	return nil, unexpected("UpdateOrgRole")
}

func (m *MockService) DeleteOrgRole(roleID string) error {
	if m.DeleteOrgRoleFunc != nil {
		return m.DeleteOrgRoleFunc(roleID)
	}

	return unexpected("DeleteOrgRole")
}

// Organization settings

func (m *MockService) GetOrgSettings() (json.RawMessage, error) {
	if m.GetOrgSettingsFunc != nil {
		return m.GetOrgSettingsFunc()
	}

	return nil, unexpected("GetOrgSettings")
}

func (m *MockService) PatchOrgSettings(body json.RawMessage) (json.RawMessage, error) {
	if m.PatchOrgSettingsFunc != nil {
		return m.PatchOrgSettingsFunc(body)
	}

	return nil, unexpected("PatchOrgSettings")
}

// Organization members

func (m *MockService) ListOrgMembers() (json.RawMessage, error) {
	if m.ListOrgMembersFunc != nil {
		return m.ListOrgMembersFunc()
	}

	return nil, unexpected("ListOrgMembers")
}

func (m *MockService) AddOrgMember(body json.RawMessage) (json.RawMessage, error) {
	if m.AddOrgMemberFunc != nil {
		return m.AddOrgMemberFunc(body)
	}

	return nil, unexpected("AddOrgMember")
}

func (m *MockService) RemoveOrgMember(memberID string) error {
	if m.RemoveOrgMemberFunc != nil {
		return m.RemoveOrgMemberFunc(memberID)
	}

	return unexpected("RemoveOrgMember")
}

// Organization audits

func (m *MockService) ListOrgAudits(p ListOrgAuditsParams) (*PaginatedResponse, error) {
	if m.ListOrgAuditsFunc != nil {
		return m.ListOrgAuditsFunc(p)
	}

	return nil, unexpected("ListOrgAudits")
}

func (m *MockService) GetOrgAudit(auditID string) (json.RawMessage, error) {
	if m.GetOrgAuditFunc != nil {
		return m.GetOrgAuditFunc(auditID)
	}

	return nil, unexpected("GetOrgAudit")
}

// Organization reporters

func (m *MockService) ListOrgReporters() ([]json.RawMessage, error) {
	if m.ListOrgReportersFunc != nil {
		return m.ListOrgReportersFunc()
	}

	return nil, unexpected("ListOrgReporters")
}

func (m *MockService) GetOrgReporter(reporterID string) (json.RawMessage, error) {
	if m.GetOrgReporterFunc != nil {
		return m.GetOrgReporterFunc(reporterID)
	}

	return nil, unexpected("GetOrgReporter")
}

func (m *MockService) CreateOrgReporter(body json.RawMessage) (json.RawMessage, error) {
	if m.CreateOrgReporterFunc != nil {
		return m.CreateOrgReporterFunc(body)
	}

	return nil, unexpected("CreateOrgReporter")
}

func (m *MockService) UpdateOrgReporter(reporterID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateOrgReporterFunc != nil {
		return m.UpdateOrgReporterFunc(reporterID, body)
	}

	return nil, unexpected("UpdateOrgReporter")
}

func (m *MockService) DeleteOrgReporter(reporterID string) error {
	if m.DeleteOrgReporterFunc != nil {
		return m.DeleteOrgReporterFunc(reporterID)
	}

	return unexpected("DeleteOrgReporter")
}

// Organization forms

func (m *MockService) GetOrgForm(template string) (json.RawMessage, error) {
	if m.GetOrgFormFunc != nil {
		return m.GetOrgFormFunc(template)
	}

	return nil, unexpected("GetOrgForm")
}

func (m *MockService) CreateOrgForm(body json.RawMessage) (json.RawMessage, error) {
	if m.CreateOrgFormFunc != nil {
		return m.CreateOrgFormFunc(body)
	}

	return nil, unexpected("CreateOrgForm")
}

func (m *MockService) UpdateOrgForm(formID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateOrgFormFunc != nil {
		return m.UpdateOrgFormFunc(formID, body)
	}

	return nil, unexpected("UpdateOrgForm")
}

func (m *MockService) DeleteOrgForm(formID string) error {
	if m.DeleteOrgFormFunc != nil {
		return m.DeleteOrgFormFunc(formID)
	}

	return unexpected("DeleteOrgForm")
}

// Organization identity providers

func (m *MockService) ListOrgIdentityProviders() ([]json.RawMessage, error) {
	if m.ListOrgIdentityProvidersFunc != nil {
		return m.ListOrgIdentityProvidersFunc()
	}

	return nil, unexpected("ListOrgIdentityProviders")
}

func (m *MockService) GetOrgIdentityProvider(idpID string) (json.RawMessage, error) {
	if m.GetOrgIdentityProviderFunc != nil {
		return m.GetOrgIdentityProviderFunc(idpID)
	}

	return nil, unexpected("GetOrgIdentityProvider")
}

func (m *MockService) CreateOrgIdentityProvider(body json.RawMessage) (json.RawMessage, error) {
	if m.CreateOrgIdentityProviderFunc != nil {
		return m.CreateOrgIdentityProviderFunc(body)
	}

	return nil, unexpected("CreateOrgIdentityProvider")
}

func (m *MockService) UpdateOrgIdentityProvider(idpID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateOrgIdentityProviderFunc != nil {
		return m.UpdateOrgIdentityProviderFunc(idpID, body)
	}

	return nil, unexpected("UpdateOrgIdentityProvider")
}

func (m *MockService) DeleteOrgIdentityProvider(idpID string) error {
	if m.DeleteOrgIdentityProviderFunc != nil {
		return m.DeleteOrgIdentityProviderFunc(idpID)
	}

	return unexpected("DeleteOrgIdentityProvider")
}

// Organization entrypoints

func (m *MockService) ListOrgEntrypoints() ([]json.RawMessage, error) {
	if m.ListOrgEntrypointsFunc != nil {
		return m.ListOrgEntrypointsFunc()
	}

	return nil, unexpected("ListOrgEntrypoints")
}

func (m *MockService) GetOrgEntrypoint(entrypointID string) (json.RawMessage, error) {
	if m.GetOrgEntrypointFunc != nil {
		return m.GetOrgEntrypointFunc(entrypointID)
	}

	return nil, unexpected("GetOrgEntrypoint")
}

func (m *MockService) CreateOrgEntrypoint(body json.RawMessage) (json.RawMessage, error) {
	if m.CreateOrgEntrypointFunc != nil {
		return m.CreateOrgEntrypointFunc(body)
	}

	return nil, unexpected("CreateOrgEntrypoint")
}

func (m *MockService) UpdateOrgEntrypoint(entrypointID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateOrgEntrypointFunc != nil {
		return m.UpdateOrgEntrypointFunc(entrypointID, body)
	}

	return nil, unexpected("UpdateOrgEntrypoint")
}

func (m *MockService) DeleteOrgEntrypoint(entrypointID string) error {
	if m.DeleteOrgEntrypointFunc != nil {
		return m.DeleteOrgEntrypointFunc(entrypointID)
	}

	return unexpected("DeleteOrgEntrypoint")
}

// Organization tags

func (m *MockService) ListOrgTags() ([]json.RawMessage, error) {
	if m.ListOrgTagsFunc != nil {
		return m.ListOrgTagsFunc()
	}

	return nil, unexpected("ListOrgTags")
}

func (m *MockService) GetOrgTag(tagID string) (json.RawMessage, error) {
	if m.GetOrgTagFunc != nil {
		return m.GetOrgTagFunc(tagID)
	}

	return nil, unexpected("GetOrgTag")
}

func (m *MockService) CreateOrgTag(body json.RawMessage) (json.RawMessage, error) {
	if m.CreateOrgTagFunc != nil {
		return m.CreateOrgTagFunc(body)
	}

	return nil, unexpected("CreateOrgTag")
}

func (m *MockService) UpdateOrgTag(tagID string, body json.RawMessage) (json.RawMessage, error) {
	if m.UpdateOrgTagFunc != nil {
		return m.UpdateOrgTagFunc(tagID, body)
	}

	return nil, unexpected("UpdateOrgTag")
}

func (m *MockService) DeleteOrgTag(tagID string) error {
	if m.DeleteOrgTagFunc != nil {
		return m.DeleteOrgTagFunc(tagID)
	}

	return unexpected("DeleteOrgTag")
}

// Organization user tokens

func (m *MockService) ListOrgUserTokens(userID string) (json.RawMessage, error) {
	if m.ListOrgUserTokensFunc != nil {
		return m.ListOrgUserTokensFunc(userID)
	}

	return nil, unexpected("ListOrgUserTokens")
}

func (m *MockService) CreateOrgUserToken(userID string, body json.RawMessage) (json.RawMessage, error) {
	if m.CreateOrgUserTokenFunc != nil {
		return m.CreateOrgUserTokenFunc(userID, body)
	}

	return nil, unexpected("CreateOrgUserToken")
}

func (m *MockService) RevokeOrgUserToken(userID, tokenID string) error {
	if m.RevokeOrgUserTokenFunc != nil {
		return m.RevokeOrgUserTokenFunc(userID, tokenID)
	}

	return unexpected("RevokeOrgUserToken")
}
