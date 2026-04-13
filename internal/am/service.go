package am

import (
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
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
