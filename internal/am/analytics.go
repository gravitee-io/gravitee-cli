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

	"gravitee.io/gctl/internal/client"
)

// AnalyticsParams holds query parameters for the analytics endpoint.
type AnalyticsParams struct {
	Type     string
	Field    string
	From     string
	To       string
	Interval string
	Size     int
}

// AnalyticsService defines analytics-related operations.
type AnalyticsService interface {
	GetAnalytics(domainID string, params AnalyticsParams) (json.RawMessage, error)
}

func (s *service) GetAnalytics(domainID string, params AnalyticsParams) (json.RawMessage, error) {
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

	path := s.domainPath(domainID, "analytics")
	if q != "" {
		path += "?" + q
	}

	data, err := s.client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("analytics get failed: %w", err)
	}

	return json.RawMessage(data), nil
}
