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

package audit

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

// --- List ---

func TestListAudits(t *testing.T) {
	t.Run("returns audits", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			ListAuditsFunc: func(domainID string, p am.ListAuditsParams) (*am.PaginatedResponse, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return &am.PaginatedResponse{
					Data: []json.RawMessage{
						json.RawMessage(`{"id":"aud-1","type":"USER_LOGIN","status":"SUCCESS","actor":{"displayName":"admin"},"target":{"displayName":"user1"},"timestamp":"2024-01-01"}`),
						json.RawMessage(`{"id":"aud-2","type":"USER_LOGOUT","status":"SUCCESS","actor":{"displayName":"admin"},"target":{"displayName":"user2"},"timestamp":"2024-01-02"}`),
					},
					TotalCount:  2,
					CurrentPage: 0,
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAuditCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "aud-1")
		testutil.AssertOutputContains(t, tc.Out, "USER_LOGIN")
		testutil.AssertOutputContains(t, tc.Out, "admin")
	})

	t.Run("passes filter flags", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			ListAuditsFunc: func(_ string, p am.ListAuditsParams) (*am.PaginatedResponse, error) {
				if p.Type != "USER_LOGIN" {
					t.Errorf("expected type 'USER_LOGIN', got %q", p.Type)
				}

				if p.Status != "SUCCESS" {
					t.Errorf("expected status 'SUCCESS', got %q", p.Status)
				}

				return &am.PaginatedResponse{
					Data:       []json.RawMessage{},
					TotalCount: 0,
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAuditCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list", "--type", "USER_LOGIN", "--status", "SUCCESS")

		testutil.AssertNoError(t, err)
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			ListAuditsFunc: func(_ string, _ am.ListAuditsParams) (*am.PaginatedResponse, error) {
				return &am.PaginatedResponse{
					Data: []json.RawMessage{
						json.RawMessage(`{"id":"aud-1","type":"USER_LOGIN"}`),
					},
					TotalCount:  1,
					CurrentPage: 0,
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAuditCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "list")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"totalCount"`)
	})

	t.Run("fetches all pages", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		callCount := 0
		mock := &am.MockService{
			ListAuditsFunc: func(_ string, p am.ListAuditsParams) (*am.PaginatedResponse, error) {
				callCount++

				if callCount == 1 {
					return &am.PaginatedResponse{
						Data: []json.RawMessage{
							json.RawMessage(`{"id":"aud-1","type":"LOGIN"}`),
						},
						TotalCount:  2,
						CurrentPage: 0,
					}, nil
				}

				return &am.PaginatedResponse{
					Data: []json.RawMessage{
						json.RawMessage(`{"id":"aud-2","type":"LOGOUT"}`),
					},
					TotalCount:  2,
					CurrentPage: 1,
				}, nil
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAuditCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list", "--all", "--per-page", "1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "aud-1")
		testutil.AssertOutputContains(t, tc.Out, "aud-2")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewAuditCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "list")

		testutil.AssertErrorContains(t, err, "no context configured")
	})

	t.Run("requires domain flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAuditCmd(tc.Factory)
		err := testutil.Execute(cmd, "list")

		testutil.AssertErrorContains(t, err, "required")
	})
}

// --- Get ---

func TestGetAudit(t *testing.T) {
	t.Run("returns audit details", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetAuditFunc: func(domainID, auditID string) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if auditID != "aud-1" {
					t.Errorf("expected auditID 'aud-1', got %q", auditID)
				}

				return json.Marshal(map[string]any{
					"id": "aud-1", "type": "USER_LOGIN", "status": "SUCCESS", "timestamp": "2024-01-01",
					"actor":  map[string]any{"displayName": "admin"},
					"target": map[string]any{"displayName": "user1"},
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAuditCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "aud-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "aud-1")
		testutil.AssertOutputContains(t, tc.Out, "USER_LOGIN")
		testutil.AssertOutputContains(t, tc.Out, "admin")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetAuditFunc: func(_, _ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"id": "aud-1", "type": "USER_LOGIN"})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAuditCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "get", "aud-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires audit ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAuditCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewAuditCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "aud-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})

	t.Run("returns error on API failure", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			GetAuditFunc: func(_, _ string) (json.RawMessage, error) {
				return nil, &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAuditCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get", "aud-1")

		testutil.AssertErrorContains(t, err, "not found")
	})
}

