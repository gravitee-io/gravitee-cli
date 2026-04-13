package am

import (
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/client"
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
