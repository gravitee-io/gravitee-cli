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

package apim

import (
	"encoding/json"
	"fmt"

	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/config"
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
