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
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/client"
)

// ListApplicationsParams holds parameters for listing applications.
type ListApplicationsParams struct {
	Query   string
	Page    int
	PerPage int
}

// ApplicationService defines application-related operations.
type ApplicationService interface {
	ListApplications(domainID string, params ListApplicationsParams) (*PaginatedResponse, error)
	GetApplication(domainID, appID string) (json.RawMessage, error)
	CreateApplication(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdateApplication(domainID, appID string, body json.RawMessage) (json.RawMessage, error)
	PatchApplication(domainID, appID string, body json.RawMessage) (json.RawMessage, error)
	DeleteApplication(domainID, appID string) error

	// Application sub-resources
	ListAppSecrets(domainID, appID string) (json.RawMessage, error)
	CreateAppSecret(domainID, appID string, body json.RawMessage) (json.RawMessage, error)
	DeleteAppSecret(domainID, appID, secretID string) error
	RenewAppSecret(domainID, appID, secretID string) (json.RawMessage, error)

	ListAppMembers(domainID, appID string) (json.RawMessage, error)
	AddAppMember(domainID, appID string, body json.RawMessage) (json.RawMessage, error)
	RemoveAppMember(domainID, appID, memberID string) error

	ListAppFlows(domainID, appID string) ([]json.RawMessage, error)
	GetAppFlow(domainID, appID, flowID string) (json.RawMessage, error)
	UpdateAppFlows(domainID, appID string, body json.RawMessage) (json.RawMessage, error)

	GetAppEmail(domainID, appID, template string) (json.RawMessage, error)
	CreateAppEmail(domainID, appID string, body json.RawMessage) (json.RawMessage, error)
	UpdateAppEmail(domainID, appID, emailID string, body json.RawMessage) (json.RawMessage, error)
	DeleteAppEmail(domainID, appID, emailID string) error

	GetAppForm(domainID, appID, template string) (json.RawMessage, error)
	CreateAppForm(domainID, appID string, body json.RawMessage) (json.RawMessage, error)
	UpdateAppForm(domainID, appID, formID string, body json.RawMessage) (json.RawMessage, error)
	DeleteAppForm(domainID, appID, formID string) error

	ListAppResources(domainID, appID string) (json.RawMessage, error)
	GetAppResource(domainID, appID, resourceID string) (json.RawMessage, error)

	GetAppAnalytics(domainID, appID string, params AnalyticsParams) (json.RawMessage, error)
	ChangeAppType(domainID, appID string, body json.RawMessage) (json.RawMessage, error)

	// Application resource policies
	ListAppResourcePolicies(domainID, appID, resourceID string) (json.RawMessage, error)
	GetAppResourcePolicy(domainID, appID, resourceID, policyID string) (json.RawMessage, error)

	// Application member permissions
	GetAppMemberPermissions(domainID, appID string) (json.RawMessage, error)
}

func (s *service) ListApplications(domainID string, params ListApplicationsParams) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{
		"page": client.Itoa(params.Page),
		"size": client.Itoa(params.PerPage),
		"q":    params.Query,
	})

	data, err := s.client.Get(s.domainPath(domainID, "applications?"+q))
	if err != nil {
		return nil, fmt.Errorf("application list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) GetApplication(domainID, appID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("applications/%s", appID)))
	if err != nil {
		return nil, fmt.Errorf("application get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateApplication(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "applications"), body)
	if err != nil {
		return nil, fmt.Errorf("application create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateApplication(domainID, appID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("applications/%s", appID)), body)
	if err != nil {
		return nil, fmt.Errorf("application update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) PatchApplication(domainID, appID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Patch(s.domainPath(domainID, fmt.Sprintf("applications/%s", appID)), body)
	if err != nil {
		return nil, fmt.Errorf("application patch failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteApplication(domainID, appID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("applications/%s", appID)))
	if err != nil {
		return fmt.Errorf("application delete failed: %w", err)
	}

	return nil
}

// Application Secrets

func (s *service) appPath(domainID, appID, sub string) string {
	return s.domainPath(domainID, fmt.Sprintf("applications/%s/%s", appID, sub))
}

func (s *service) ListAppSecrets(domainID, appID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.appPath(domainID, appID, "secrets"))
	if err != nil {
		return nil, fmt.Errorf("app secret list failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateAppSecret(domainID, appID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.appPath(domainID, appID, "secrets"), body)
	if err != nil {
		return nil, fmt.Errorf("app secret create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteAppSecret(domainID, appID, secretID string) error {
	err := s.client.Delete(s.appPath(domainID, appID, fmt.Sprintf("secrets/%s", secretID)))
	if err != nil {
		return fmt.Errorf("app secret delete failed: %w", err)
	}

	return nil
}

func (s *service) RenewAppSecret(domainID, appID, secretID string) (json.RawMessage, error) {
	data, err := s.client.Post(s.appPath(domainID, appID, fmt.Sprintf("secrets/%s/_renew", secretID)), json.RawMessage(`{}`))
	if err != nil {
		return nil, fmt.Errorf("app secret renew failed: %w", err)
	}

	return json.RawMessage(data), nil
}

// Application Members

func (s *service) ListAppMembers(domainID, appID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.appPath(domainID, appID, "members"))
	if err != nil {
		return nil, fmt.Errorf("app member list failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) AddAppMember(domainID, appID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.appPath(domainID, appID, "members"), body)
	if err != nil {
		return nil, fmt.Errorf("app member add failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) RemoveAppMember(domainID, appID, memberID string) error {
	err := s.client.Delete(s.appPath(domainID, appID, fmt.Sprintf("members/%s", memberID)))
	if err != nil {
		return fmt.Errorf("app member remove failed: %w", err)
	}

	return nil
}

// Application Flows

func (s *service) ListAppFlows(domainID, appID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.appPath(domainID, appID, "flows"))
	if err != nil {
		return nil, fmt.Errorf("app flow list failed: %w", err)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse app flow list: %w", err)
	}

	return items, nil
}

func (s *service) GetAppFlow(domainID, appID, flowID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.appPath(domainID, appID, fmt.Sprintf("flows/%s", flowID)))
	if err != nil {
		return nil, fmt.Errorf("app flow get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateAppFlows(domainID, appID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.appPath(domainID, appID, "flows"), body)
	if err != nil {
		return nil, fmt.Errorf("app flow update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

// Application Emails

func (s *service) GetAppEmail(domainID, appID, template string) (json.RawMessage, error) {
	q := client.BuildQuery(map[string]string{"template": template})

	path := s.appPath(domainID, appID, "emails")
	if q != "" {
		path += "?" + q
	}

	data, err := s.client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("app email get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateAppEmail(domainID, appID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.appPath(domainID, appID, "emails"), body)
	if err != nil {
		return nil, fmt.Errorf("app email create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateAppEmail(domainID, appID, emailID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.appPath(domainID, appID, fmt.Sprintf("emails/%s", emailID)), body)
	if err != nil {
		return nil, fmt.Errorf("app email update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteAppEmail(domainID, appID, emailID string) error {
	err := s.client.Delete(s.appPath(domainID, appID, fmt.Sprintf("emails/%s", emailID)))
	if err != nil {
		return fmt.Errorf("app email delete failed: %w", err)
	}

	return nil
}

// Application Forms

func (s *service) GetAppForm(domainID, appID, template string) (json.RawMessage, error) {
	q := client.BuildQuery(map[string]string{"template": template})

	path := s.appPath(domainID, appID, "forms")
	if q != "" {
		path += "?" + q
	}

	data, err := s.client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("app form get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreateAppForm(domainID, appID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.appPath(domainID, appID, "forms"), body)
	if err != nil {
		return nil, fmt.Errorf("app form create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdateAppForm(domainID, appID, formID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.appPath(domainID, appID, fmt.Sprintf("forms/%s", formID)), body)
	if err != nil {
		return nil, fmt.Errorf("app form update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeleteAppForm(domainID, appID, formID string) error {
	err := s.client.Delete(s.appPath(domainID, appID, fmt.Sprintf("forms/%s", formID)))
	if err != nil {
		return fmt.Errorf("app form delete failed: %w", err)
	}

	return nil
}

// Application Resources

func (s *service) ListAppResources(domainID, appID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.appPath(domainID, appID, "resources"))
	if err != nil {
		return nil, fmt.Errorf("app resource list failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) GetAppResource(domainID, appID, resourceID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.appPath(domainID, appID, fmt.Sprintf("resources/%s", resourceID)))
	if err != nil {
		return nil, fmt.Errorf("app resource get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

// Application Analytics

func (s *service) GetAppAnalytics(domainID, appID string, params AnalyticsParams) (json.RawMessage, error) {
	qp := map[string]string{
		"type":  params.Type,
		"field": params.Field,
		"from":  params.From,
		"to":    params.To,
	}

	if params.Interval != "" {
		qp["interval"] = params.Interval
	}

	if params.Size > 0 {
		qp["size"] = client.Itoa(params.Size)
	}

	q := client.BuildQuery(qp)

	path := s.appPath(domainID, appID, "analytics")
	if q != "" {
		path += "?" + q
	}

	data, err := s.client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("app analytics get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

// Application Type Change

func (s *service) ChangeAppType(domainID, appID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.appPath(domainID, appID, "type"), body)
	if err != nil {
		return nil, fmt.Errorf("app type change failed: %w", err)
	}

	return json.RawMessage(data), nil
}

// Application resource policies

func (s *service) ListAppResourcePolicies(domainID, appID, resourceID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.appPath(domainID, appID, fmt.Sprintf("resources/%s/policies", resourceID)))
	if err != nil {
		return nil, fmt.Errorf("app resource policy list failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) GetAppResourcePolicy(domainID, appID, resourceID, policyID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.appPath(domainID, appID, fmt.Sprintf("resources/%s/policies/%s", resourceID, policyID)))
	if err != nil {
		return nil, fmt.Errorf("app resource policy get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

// Application member permissions

func (s *service) GetAppMemberPermissions(domainID, appID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.appPath(domainID, appID, "members/permissions"))
	if err != nil {
		return nil, fmt.Errorf("app member permissions get failed: %w", err)
	}

	return json.RawMessage(data), nil
}
