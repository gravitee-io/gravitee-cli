package apim

import (
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
)

// Service defines all APIM management operations.
type Service interface {
	APIService
	PlanService
	SubscriptionService
	APIKeyService
	MemberService
	PageService
	MetadataService
	ApplicationService
	EnvironmentService
	PluginService
}

// service is the concrete implementation backed by an HTTP client.
type service struct {
	client   client.GraviteeClient
	resolved *config.ResolvedContext
}

// NewService creates a new APIM service.
func NewService(c client.GraviteeClient, r *config.ResolvedContext) Service {
	return &service{client: c, resolved: r}
}

func (s *service) v2(path string) string {
	return client.V2Path(s.resolved.Env, path)
}

func (s *service) v1(path string) string {
	return client.V1Path(s.resolved.Org, s.resolved.Env, path)
}

func (s *service) orgV2(path string) string {
	return fmt.Sprintf("/management/v2/organizations/%s/%s", s.resolved.Org, path)
}

func i64toa(n int64) string {
	return fmt.Sprintf("%d", n)
}

func raw(data []byte) json.RawMessage {
	return data
}
