// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package am

import (
	"fmt"

	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/config"
)

// Service defines all AM management operations.
type Service interface {
	DomainService
	ApplicationService
	UserService
	RoleService
	ScopeService
	IdentityProviderService
	CertificateService
	FactorService
	GroupService
	FlowService
	FormService
	EmailService
	ThemeService
	PasswordPolicyService
	AuditService
	MemberService
	ExtensionGrantService
	ResourceService
	ReporterService
	BotDetectionService
	DeviceIdentifierService
	AuthDeviceNotifierService
	AuthorizationEngineService
	ProtectedResourceService
	AnalyticsService
	EntrypointService
	DictionaryService
	AlertService
	OrganizationService
	OrgReporterService
	OrgFormService
	OrgIdentityProviderService
	OrgEntrypointService
	OrgTagService
	OrgUserTokenService
}

// service is the concrete implementation backed by an HTTP client.
type service struct {
	client   client.GraviteeClient
	resolved *config.ResolvedContext
}

// NewService creates a new AM service.
func NewService(c client.GraviteeClient, r *config.ResolvedContext) Service {
	return &service{client: c, resolved: r}
}

func (s *service) basePath(path string) string {
	return fmt.Sprintf("/management/organizations/%s/environments/%s/%s", s.resolved.Org, s.resolved.Env, path)
}

func (s *service) orgPath(path string) string {
	return fmt.Sprintf("/management/organizations/%s/%s", s.resolved.Org, path)
}

func (s *service) domainPath(domainID, path string) string {
	return s.basePath(fmt.Sprintf("domains/%s/%s", domainID, path))
}
