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
	"strings"

	"github.com/gravitee-io/gio-cli/internal/client"
)

// ListSubscriptionsParams holds parameters for listing subscriptions.
type ListSubscriptionsParams struct {
	Statuses []string
	PlanID   string
	AppID    string
	Page     int
	PerPage  int
}

// CreateSubscriptionBody holds the body for creating a subscription.
type CreateSubscriptionBody struct {
	PlanID       string `json:"planId"`
	AppID        string `json:"applicationId"`
	CustomAPIKey string `json:"customApiKey,omitempty"`
}

// AcceptSubscriptionBody holds the body for accepting a subscription.
type AcceptSubscriptionBody struct {
	Reason       string `json:"reason,omitempty"`
	StartingAt   string `json:"startingAt,omitempty"`
	EndingAt     string `json:"endingAt,omitempty"`
	CustomAPIKey string `json:"customApiKey,omitempty"`
}

// SubscriptionService defines subscription-related operations.
type SubscriptionService interface {
	ListSubscriptions(apiID string, params ListSubscriptionsParams) (*PaginatedResponse, error)
	GetSubscription(apiID, subID string) (json.RawMessage, error)
	CreateSubscription(apiID string, body CreateSubscriptionBody) (json.RawMessage, error)
	AcceptSubscription(apiID, subID string, body AcceptSubscriptionBody) (json.RawMessage, error)
	RejectSubscription(apiID, subID string, reason string) (json.RawMessage, error)
	PauseSubscription(apiID, subID string) (json.RawMessage, error)
	ResumeSubscription(apiID, subID string) (json.RawMessage, error)
	CloseSubscription(apiID, subID string) (json.RawMessage, error)
	TransferSubscription(apiID, subID, planID string) (json.RawMessage, error)
}

func (s *service) ListSubscriptions(apiID string, p ListSubscriptionsParams) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{
		"page": client.Itoa(p.Page), "perPage": client.Itoa(p.PerPage),
		"statuses": strings.Join(p.Statuses, ","),
		"planIds":  p.PlanID, "applicationIds": p.AppID,
	})

	data, err := s.client.Get(s.v2(fmt.Sprintf("apis/%s/subscriptions?%s", apiID, q)))
	if err != nil {
		return nil, fmt.Errorf("subscription list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) GetSubscription(apiID, subID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.v2(fmt.Sprintf("apis/%s/subscriptions/%s", apiID, subID)))
	if err != nil {
		return nil, err
	}

	return raw(data), nil
}

func (s *service) CreateSubscription(apiID string, body CreateSubscriptionBody) (json.RawMessage, error) {
	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/subscriptions", apiID)), body)
	if err != nil {
		return nil, fmt.Errorf("subscription creation failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) AcceptSubscription(apiID, subID string, body AcceptSubscriptionBody) (json.RawMessage, error) {
	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/subscriptions/%s/_accept", apiID, subID)), body)
	if err != nil {
		return nil, fmt.Errorf("subscription accept failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) RejectSubscription(apiID, subID, reason string) (json.RawMessage, error) {
	// The API requires a non-null body even when the reason is empty.
	body := map[string]string{"reason": reason}

	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/subscriptions/%s/_reject", apiID, subID)), body)
	if err != nil {
		return nil, fmt.Errorf("subscription reject failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) PauseSubscription(apiID, subID string) (json.RawMessage, error) {
	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/subscriptions/%s/_pause", apiID, subID)), nil)
	if err != nil {
		return nil, fmt.Errorf("subscription pause failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) ResumeSubscription(apiID, subID string) (json.RawMessage, error) {
	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/subscriptions/%s/_resume", apiID, subID)), nil)
	if err != nil {
		return nil, fmt.Errorf("subscription resume failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) CloseSubscription(apiID, subID string) (json.RawMessage, error) {
	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/subscriptions/%s/_close", apiID, subID)), nil)
	if err != nil {
		return nil, fmt.Errorf("subscription close failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) TransferSubscription(apiID, subID, planID string) (json.RawMessage, error) {
	body := map[string]string{"planId": planID}

	data, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/subscriptions/%s/_transfer", apiID, subID)), body)
	if err != nil {
		return nil, fmt.Errorf("subscription transfer failed: %w", err)
	}

	return raw(data), nil
}
