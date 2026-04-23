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
)

// MockService implements Service with injectable functions for testing.
type MockService struct {
	ListAPIsFunc           func(ListAPIsParams) (*PaginatedResponse, error)
	ResolveAPIFunc         func(string) (string, error)
	GetAPIFunc             func(string) (json.RawMessage, error)
	CreateAPIFunc          func(json.RawMessage) (json.RawMessage, error)
	UpdateAPIFunc          func(string, json.RawMessage) (json.RawMessage, error)
	DeleteAPIFunc          func(string, bool) error
	StartAPIFunc           func(string) error
	StopAPIFunc            func(string) error
	DeployAPIFunc          func(string, string) error
	ImportAPIFunc          func(json.RawMessage) (json.RawMessage, error)
	ExportAPIFunc          func(string, []string) (json.RawMessage, error)
	RollbackAPIFunc        func(string, string) error
	GetAnalyticsFunc       func(string, AnalyticsParams) (json.RawMessage, error)
	GetHealthFunc          func(string, string) (json.RawMessage, error)
	ListLogsFunc           func(string, ListAPILogsParams) (*PaginatedResponse, error)
	GetLogFunc             func(string, string) (json.RawMessage, error)
	ListPlansFunc          func(string, ListPlansParams) (*PaginatedResponse, error)
	GetPlanFunc            func(string, string) (json.RawMessage, error)
	CreatePlanFunc         func(string, json.RawMessage) (json.RawMessage, error)
	UpdatePlanFunc         func(string, string, json.RawMessage) (json.RawMessage, error)
	DeletePlanFunc         func(string, string) error
	PublishPlanFunc        func(string, string) (json.RawMessage, error)
	DeprecatePlanFunc      func(string, string) (json.RawMessage, error)
	ClosePlanFunc          func(string, string) (json.RawMessage, error)
	ListSubscriptionsFunc  func(string, ListSubscriptionsParams) (*PaginatedResponse, error)
	GetSubscriptionFunc    func(string, string) (json.RawMessage, error)
	CreateSubscriptionFunc func(string, CreateSubscriptionBody) (json.RawMessage, error)
	AcceptSubFunc          func(string, string, AcceptSubscriptionBody) (json.RawMessage, error)
	RejectSubFunc          func(string, string, string) (json.RawMessage, error)
	PauseSubFunc           func(string, string) (json.RawMessage, error)
	ResumeSubFunc          func(string, string) (json.RawMessage, error)
	CloseSubFunc           func(string, string) (json.RawMessage, error)
	TransferSubFunc        func(string, string, string) (json.RawMessage, error)
	ListAPIKeysFunc        func(string, string, int, int) (*PaginatedResponse, error)
	RenewAPIKeyFunc        func(string, string) (json.RawMessage, error)
	RevokeAPIKeyFunc       func(string, string, string) error
	ReactivateAPIKeyFunc   func(string, string, string) (json.RawMessage, error)
	ListMembersFunc        func(string, int, int) (*PaginatedResponse, error)
	AddMemberFunc          func(string, string, string) (json.RawMessage, error)
	RemoveMemberFunc       func(string, string) error
	ListPagesFunc          func(string, string) ([]json.RawMessage, error)
	GetPageFunc            func(string, string) (json.RawMessage, error)
	CreatePageFunc         func(string, json.RawMessage) (json.RawMessage, error)
	UpdatePageFunc         func(string, string, json.RawMessage) (json.RawMessage, error)
	DeletePageFunc         func(string, string) error
	PublishPageFunc        func(string, string) (json.RawMessage, error)
	UnpublishPageFunc      func(string, string) (json.RawMessage, error)
	ListMetadataFunc       func(string, int, int) (*PaginatedResponse, error)
	CreateMetadataFunc     func(string, json.RawMessage) (json.RawMessage, error)
	UpdateMetadataFunc     func(string, string, json.RawMessage) (json.RawMessage, error)
	DeleteMetadataFunc     func(string, string) error
	ListAppsFunc           func(ListApplicationsParams) (*PaginatedResponse, error)
	GetAppFunc             func(string) (json.RawMessage, error)
	CreateAppFunc          func(json.RawMessage) (json.RawMessage, error)
	UpdateAppFunc          func(string, json.RawMessage) (json.RawMessage, error)
	DeleteAppFunc          func(string) error
	ListEnvsFunc           func() (json.RawMessage, error)
	GetEnvFunc             func(string) (json.RawMessage, error)
	ListPluginsFunc        func(string) (json.RawMessage, error)
}

func unexpected(name string) error { return fmt.Errorf("unexpected call: %s", name) }

func (m *MockService) ListAPIs(p ListAPIsParams) (*PaginatedResponse, error) {
	if m.ListAPIsFunc != nil {
		return m.ListAPIsFunc(p)
	}
	return nil, unexpected("ListAPIs")
}
func (m *MockService) ResolveAPI(p string) (string, error) {
	if m.ResolveAPIFunc != nil {
		return m.ResolveAPIFunc(p)
	}
	// default passthrough: tests that don't care about resolution pass UUIDs
	return p, nil
}
func (m *MockService) GetAPI(id string) (json.RawMessage, error) {
	if m.GetAPIFunc != nil {
		return m.GetAPIFunc(id)
	}
	return nil, unexpected("GetAPI")
}
func (m *MockService) CreateAPI(b json.RawMessage) (json.RawMessage, error) {
	if m.CreateAPIFunc != nil {
		return m.CreateAPIFunc(b)
	}
	return nil, unexpected("CreateAPI")
}
func (m *MockService) UpdateAPI(id string, b json.RawMessage) (json.RawMessage, error) {
	if m.UpdateAPIFunc != nil {
		return m.UpdateAPIFunc(id, b)
	}
	return nil, unexpected("UpdateAPI")
}
func (m *MockService) DeleteAPI(id string, cp bool) error {
	if m.DeleteAPIFunc != nil {
		return m.DeleteAPIFunc(id, cp)
	}
	return unexpected("DeleteAPI")
}
func (m *MockService) StartAPI(id string) error {
	if m.StartAPIFunc != nil {
		return m.StartAPIFunc(id)
	}
	return unexpected("StartAPI")
}
func (m *MockService) StopAPI(id string) error {
	if m.StopAPIFunc != nil {
		return m.StopAPIFunc(id)
	}
	return unexpected("StopAPI")
}
func (m *MockService) DeployAPI(id, label string) error {
	if m.DeployAPIFunc != nil {
		return m.DeployAPIFunc(id, label)
	}
	return unexpected("DeployAPI")
}
func (m *MockService) ImportAPI(b json.RawMessage) (json.RawMessage, error) {
	if m.ImportAPIFunc != nil {
		return m.ImportAPIFunc(b)
	}
	return nil, unexpected("ImportAPI")
}
func (m *MockService) ExportAPI(id string, ex []string) (json.RawMessage, error) {
	if m.ExportAPIFunc != nil {
		return m.ExportAPIFunc(id, ex)
	}
	return nil, unexpected("ExportAPI")
}
func (m *MockService) RollbackAPI(id, eid string) error {
	if m.RollbackAPIFunc != nil {
		return m.RollbackAPIFunc(id, eid)
	}
	return unexpected("RollbackAPI")
}
func (m *MockService) GetAPIAnalytics(id string, p AnalyticsParams) (json.RawMessage, error) {
	if m.GetAnalyticsFunc != nil {
		return m.GetAnalyticsFunc(id, p)
	}
	return nil, unexpected("GetAPIAnalytics")
}
func (m *MockService) GetAPIHealth(id, field string) (json.RawMessage, error) {
	if m.GetHealthFunc != nil {
		return m.GetHealthFunc(id, field)
	}
	return nil, unexpected("GetAPIHealth")
}
func (m *MockService) ListAPILogs(id string, p ListAPILogsParams) (*PaginatedResponse, error) {
	if m.ListLogsFunc != nil {
		return m.ListLogsFunc(id, p)
	}
	return nil, unexpected("ListAPILogs")
}
func (m *MockService) GetAPILog(id, rid string) (json.RawMessage, error) {
	if m.GetLogFunc != nil {
		return m.GetLogFunc(id, rid)
	}
	return nil, unexpected("GetAPILog")
}

func (m *MockService) ListPlans(a string, p ListPlansParams) (*PaginatedResponse, error) {
	if m.ListPlansFunc != nil {
		return m.ListPlansFunc(a, p)
	}
	return nil, unexpected("ListPlans")
}
func (m *MockService) GetPlan(a, p string) (json.RawMessage, error) {
	if m.GetPlanFunc != nil {
		return m.GetPlanFunc(a, p)
	}
	return nil, unexpected("GetPlan")
}
func (m *MockService) CreatePlan(a string, b json.RawMessage) (json.RawMessage, error) {
	if m.CreatePlanFunc != nil {
		return m.CreatePlanFunc(a, b)
	}
	return nil, unexpected("CreatePlan")
}
func (m *MockService) UpdatePlan(a, p string, b json.RawMessage) (json.RawMessage, error) {
	if m.UpdatePlanFunc != nil {
		return m.UpdatePlanFunc(a, p, b)
	}
	return nil, unexpected("UpdatePlan")
}
func (m *MockService) DeletePlan(a, p string) error {
	if m.DeletePlanFunc != nil {
		return m.DeletePlanFunc(a, p)
	}
	return unexpected("DeletePlan")
}
func (m *MockService) PublishPlan(a, p string) (json.RawMessage, error) {
	if m.PublishPlanFunc != nil {
		return m.PublishPlanFunc(a, p)
	}
	return nil, unexpected("PublishPlan")
}
func (m *MockService) DeprecatePlan(a, p string) (json.RawMessage, error) {
	if m.DeprecatePlanFunc != nil {
		return m.DeprecatePlanFunc(a, p)
	}
	return nil, unexpected("DeprecatePlan")
}
func (m *MockService) ClosePlan(a, p string) (json.RawMessage, error) {
	if m.ClosePlanFunc != nil {
		return m.ClosePlanFunc(a, p)
	}
	return nil, unexpected("ClosePlan")
}

func (m *MockService) ListSubscriptions(a string, p ListSubscriptionsParams) (*PaginatedResponse, error) {
	if m.ListSubscriptionsFunc != nil {
		return m.ListSubscriptionsFunc(a, p)
	}
	return nil, unexpected("ListSubscriptions")
}
func (m *MockService) GetSubscription(a, s string) (json.RawMessage, error) {
	if m.GetSubscriptionFunc != nil {
		return m.GetSubscriptionFunc(a, s)
	}
	return nil, unexpected("GetSubscription")
}
func (m *MockService) CreateSubscription(a string, b CreateSubscriptionBody) (json.RawMessage, error) {
	if m.CreateSubscriptionFunc != nil {
		return m.CreateSubscriptionFunc(a, b)
	}
	return nil, unexpected("CreateSubscription")
}
func (m *MockService) AcceptSubscription(a, s string, b AcceptSubscriptionBody) (json.RawMessage, error) {
	if m.AcceptSubFunc != nil {
		return m.AcceptSubFunc(a, s, b)
	}
	return nil, unexpected("AcceptSubscription")
}
func (m *MockService) RejectSubscription(a, s, r string) (json.RawMessage, error) {
	if m.RejectSubFunc != nil {
		return m.RejectSubFunc(a, s, r)
	}
	return nil, unexpected("RejectSubscription")
}
func (m *MockService) PauseSubscription(a, s string) (json.RawMessage, error) {
	if m.PauseSubFunc != nil {
		return m.PauseSubFunc(a, s)
	}
	return nil, unexpected("PauseSubscription")
}
func (m *MockService) ResumeSubscription(a, s string) (json.RawMessage, error) {
	if m.ResumeSubFunc != nil {
		return m.ResumeSubFunc(a, s)
	}
	return nil, unexpected("ResumeSubscription")
}
func (m *MockService) CloseSubscription(a, s string) (json.RawMessage, error) {
	if m.CloseSubFunc != nil {
		return m.CloseSubFunc(a, s)
	}
	return nil, unexpected("CloseSubscription")
}
func (m *MockService) TransferSubscription(a, s, p string) (json.RawMessage, error) {
	if m.TransferSubFunc != nil {
		return m.TransferSubFunc(a, s, p)
	}
	return nil, unexpected("TransferSubscription")
}

func (m *MockService) ListAPIKeys(a, s string, pg, pp int) (*PaginatedResponse, error) {
	if m.ListAPIKeysFunc != nil {
		return m.ListAPIKeysFunc(a, s, pg, pp)
	}
	return nil, unexpected("ListAPIKeys")
}
func (m *MockService) RenewAPIKey(a, s string) (json.RawMessage, error) {
	if m.RenewAPIKeyFunc != nil {
		return m.RenewAPIKeyFunc(a, s)
	}
	return nil, unexpected("RenewAPIKey")
}
func (m *MockService) RevokeAPIKey(a, s, k string) error {
	if m.RevokeAPIKeyFunc != nil {
		return m.RevokeAPIKeyFunc(a, s, k)
	}
	return unexpected("RevokeAPIKey")
}
func (m *MockService) ReactivateAPIKey(a, s, k string) (json.RawMessage, error) {
	if m.ReactivateAPIKeyFunc != nil {
		return m.ReactivateAPIKeyFunc(a, s, k)
	}
	return nil, unexpected("ReactivateAPIKey")
}

func (m *MockService) ListMembers(a string, pg, pp int) (*PaginatedResponse, error) {
	if m.ListMembersFunc != nil {
		return m.ListMembersFunc(a, pg, pp)
	}
	return nil, unexpected("ListMembers")
}
func (m *MockService) AddMember(a, u, r string) (json.RawMessage, error) {
	if m.AddMemberFunc != nil {
		return m.AddMemberFunc(a, u, r)
	}
	return nil, unexpected("AddMember")
}
func (m *MockService) RemoveMember(a, mid string) error {
	if m.RemoveMemberFunc != nil {
		return m.RemoveMemberFunc(a, mid)
	}
	return unexpected("RemoveMember")
}

func (m *MockService) ListPages(a, parent string) ([]json.RawMessage, error) {
	if m.ListPagesFunc != nil {
		return m.ListPagesFunc(a, parent)
	}
	return nil, unexpected("ListPages")
}
func (m *MockService) GetPage(a, p string) (json.RawMessage, error) {
	if m.GetPageFunc != nil {
		return m.GetPageFunc(a, p)
	}
	return nil, unexpected("GetPage")
}
func (m *MockService) CreatePage(a string, b json.RawMessage) (json.RawMessage, error) {
	if m.CreatePageFunc != nil {
		return m.CreatePageFunc(a, b)
	}
	return nil, unexpected("CreatePage")
}
func (m *MockService) UpdatePage(a, p string, b json.RawMessage) (json.RawMessage, error) {
	if m.UpdatePageFunc != nil {
		return m.UpdatePageFunc(a, p, b)
	}
	return nil, unexpected("UpdatePage")
}
func (m *MockService) DeletePage(a, p string) error {
	if m.DeletePageFunc != nil {
		return m.DeletePageFunc(a, p)
	}
	return unexpected("DeletePage")
}
func (m *MockService) PublishPage(a, p string) (json.RawMessage, error) {
	if m.PublishPageFunc != nil {
		return m.PublishPageFunc(a, p)
	}
	return nil, unexpected("PublishPage")
}
func (m *MockService) UnpublishPage(a, p string) (json.RawMessage, error) {
	if m.UnpublishPageFunc != nil {
		return m.UnpublishPageFunc(a, p)
	}
	return nil, unexpected("UnpublishPage")
}

func (m *MockService) ListMetadata(a string, pg, pp int) (*PaginatedResponse, error) {
	if m.ListMetadataFunc != nil {
		return m.ListMetadataFunc(a, pg, pp)
	}
	return nil, unexpected("ListMetadata")
}
func (m *MockService) CreateMetadata(a string, b json.RawMessage) (json.RawMessage, error) {
	if m.CreateMetadataFunc != nil {
		return m.CreateMetadataFunc(a, b)
	}
	return nil, unexpected("CreateMetadata")
}
func (m *MockService) UpdateMetadata(a, k string, b json.RawMessage) (json.RawMessage, error) {
	if m.UpdateMetadataFunc != nil {
		return m.UpdateMetadataFunc(a, k, b)
	}
	return nil, unexpected("UpdateMetadata")
}
func (m *MockService) DeleteMetadata(a, k string) error {
	if m.DeleteMetadataFunc != nil {
		return m.DeleteMetadataFunc(a, k)
	}
	return unexpected("DeleteMetadata")
}

func (m *MockService) ListApplications(p ListApplicationsParams) (*PaginatedResponse, error) {
	if m.ListAppsFunc != nil {
		return m.ListAppsFunc(p)
	}
	return nil, unexpected("ListApplications")
}
func (m *MockService) GetApplication(id string) (json.RawMessage, error) {
	if m.GetAppFunc != nil {
		return m.GetAppFunc(id)
	}
	return nil, unexpected("GetApplication")
}
func (m *MockService) CreateApplication(b json.RawMessage) (json.RawMessage, error) {
	if m.CreateAppFunc != nil {
		return m.CreateAppFunc(b)
	}
	return nil, unexpected("CreateApplication")
}
func (m *MockService) UpdateApplication(id string, b json.RawMessage) (json.RawMessage, error) {
	if m.UpdateAppFunc != nil {
		return m.UpdateAppFunc(id, b)
	}
	return nil, unexpected("UpdateApplication")
}
func (m *MockService) DeleteApplication(id string) error {
	if m.DeleteAppFunc != nil {
		return m.DeleteAppFunc(id)
	}
	return unexpected("DeleteApplication")
}

func (m *MockService) ListEnvironments() (json.RawMessage, error) {
	if m.ListEnvsFunc != nil {
		return m.ListEnvsFunc()
	}
	return nil, unexpected("ListEnvironments")
}
func (m *MockService) GetEnvironment(id string) (json.RawMessage, error) {
	if m.GetEnvFunc != nil {
		return m.GetEnvFunc(id)
	}
	return nil, unexpected("GetEnvironment")
}

func (m *MockService) ListPlugins(t string) (json.RawMessage, error) {
	if m.ListPluginsFunc != nil {
		return m.ListPluginsFunc(t)
	}
	return nil, unexpected("ListPlugins")
}
