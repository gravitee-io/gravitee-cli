package apim

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/gravitee-io/gio-cli/internal/client"
)

// ListAPIsParams holds parameters for listing APIs.
type ListAPIsParams struct {
	Query   string
	Status  string
	Page    int
	PerPage int
}

// ListAPILogsParams holds parameters for listing API logs.
type ListAPILogsParams struct {
	ApplicationIDs []string
	PlanIDs        []string
	Methods        []string
	From           int64
	To             int64
	Page           int
	PerPage        int
}

// AnalyticsParams holds parameters for API analytics.
type AnalyticsParams struct {
	Terms        []string
	Field        string
	Type         string
	Ranges       string
	Aggregations string
	Order        string
	Query        string
	From         int64
	To           int64
	Interval     int64
	Size         int
}

// APIService defines API-related operations.
type APIService interface {
	ListAPIs(params ListAPIsParams) (*PaginatedResponse, error)
	ResolveAPI(pathOrID string) (string, error)
	GetAPI(apiID string) (json.RawMessage, error)
	CreateAPI(body json.RawMessage) (json.RawMessage, error)
	UpdateAPI(apiID string, body json.RawMessage) (json.RawMessage, error)
	DeleteAPI(apiID string, closePlans bool) error
	StartAPI(apiID string) error
	StopAPI(apiID string) error
	DeployAPI(apiID string, label string) error
	ImportAPI(body json.RawMessage) (json.RawMessage, error)
	ExportAPI(apiID string, exclude []string) (json.RawMessage, error)
	RollbackAPI(apiID string, eventID string) error
	GetAPIAnalytics(apiID string, params AnalyticsParams) (json.RawMessage, error)
	GetAPIHealth(apiID string, field string) (json.RawMessage, error)
	ListAPILogs(apiID string, params ListAPILogsParams) (*PaginatedResponse, error)
	GetAPILog(apiID string, requestID string) (json.RawMessage, error)
}

func (s *service) ListAPIs(params ListAPIsParams) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{
		"page": client.Itoa(params.Page), "perPage": client.Itoa(params.PerPage),
		"status": params.Status,
	})

	body := map[string]string{}
	if params.Query != "" {
		body["query"] = params.Query
	}

	data, err := s.client.Post(s.v2("apis/_search?"+q), body)
	if err != nil {
		return nil, fmt.Errorf("API list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) GetAPI(apiID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.v2(fmt.Sprintf("apis/%s", apiID)))
	if err != nil {
		return nil, err
	}

	return raw(data), nil
}

func (s *service) CreateAPI(body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.v2("apis"), body)
	if err != nil {
		return nil, fmt.Errorf("API creation failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) UpdateAPI(apiID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.v2(fmt.Sprintf("apis/%s", apiID)), body)
	if err != nil {
		return nil, fmt.Errorf("API update failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) DeleteAPI(apiID string, closePlans bool) error {
	path := s.v2(fmt.Sprintf("apis/%s", apiID))
	if closePlans {
		path += "?closePlans=true"
	}

	if err := s.client.Delete(path); err != nil {
		return fmt.Errorf("API deletion failed: %w", err)
	}

	return nil
}

func (s *service) StartAPI(apiID string) error {
	if _, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/_start", apiID)), nil); err != nil {
		return fmt.Errorf("API start failed: %w", err)
	}

	return nil
}

func (s *service) StopAPI(apiID string) error {
	if _, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/_stop", apiID)), nil); err != nil {
		return fmt.Errorf("API stop failed: %w", err)
	}

	return nil
}

func (s *service) DeployAPI(apiID string, label string) error {
	var body any
	if label != "" {
		body = map[string]string{"deploymentLabel": label}
	}

	if _, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/deployments", apiID)), body); err != nil {
		return fmt.Errorf("API deploy failed: %w", err)
	}

	return nil
}

func (s *service) ImportAPI(body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.v2("apis/_import/definition"), body)
	if err != nil {
		return nil, fmt.Errorf("API import failed: %w", err)
	}

	return raw(data), nil
}

func (s *service) ExportAPI(apiID string, exclude []string) (json.RawMessage, error) {
	path := s.v2(fmt.Sprintf("apis/%s/_export/definition", apiID))

	if len(exclude) > 0 {
		escaped := make([]string, len(exclude))
		for i, e := range exclude {
			escaped[i] = url.QueryEscape(e)
		}

		path += "?excludeAdditionalData=" + strings.Join(escaped, ",")
	}

	data, err := s.client.Get(path)
	if err != nil {
		return nil, err
	}

	return raw(data), nil
}

func (s *service) RollbackAPI(apiID, eventID string) error {
	body := map[string]string{"eventId": eventID}

	if _, err := s.client.Post(s.v2(fmt.Sprintf("apis/%s/_rollback", apiID)), body); err != nil {
		return fmt.Errorf("API rollback failed: %w", err)
	}

	return nil
}

func (s *service) GetAPIAnalytics(apiID string, p AnalyticsParams) (json.RawMessage, error) {
	q := url.Values{}
	// from and to are required by the API - always send them.
	q.Set("from", i64toa(p.From))
	q.Set("to", i64toa(p.To))

	if p.Interval != 0 {
		q.Set("interval", i64toa(p.Interval))
	}

	if p.Field != "" {
		q.Set("field", p.Field)
	}

	if p.Type != "" {
		q.Set("type", p.Type)
	}

	if p.Size != 0 {
		q.Set("size", client.Itoa(p.Size))
	}

	if p.Ranges != "" {
		q.Set("ranges", p.Ranges)
	}

	if p.Aggregations != "" {
		q.Set("aggregations", p.Aggregations)
	}

	if p.Order != "" {
		q.Set("order", p.Order)
	}

	if p.Query != "" {
		q.Set("query", p.Query)
	}

	for _, t := range p.Terms {
		q.Add("terms", t)
	}

	data, err := s.client.Get(s.v2(fmt.Sprintf("apis/%s/analytics?%s", apiID, q.Encode())))
	if err != nil {
		return nil, err
	}

	return raw(data), nil
}

func (s *service) GetAPIHealth(apiID, field string) (json.RawMessage, error) {
	data, err := s.client.Get(s.v2(fmt.Sprintf("apis/%s/health/availability?field=%s", apiID, url.QueryEscape(field))))
	if err != nil {
		return nil, err
	}

	return raw(data), nil
}

func (s *service) ListAPILogs(apiID string, p ListAPILogsParams) (*PaginatedResponse, error) {
	q := url.Values{}
	q.Set("page", client.Itoa(p.Page))
	q.Set("perPage", client.Itoa(p.PerPage))

	if p.From != 0 {
		q.Set("from", i64toa(p.From))
	}

	if p.To != 0 {
		q.Set("to", i64toa(p.To))
	}

	for _, id := range p.ApplicationIDs {
		q.Add("applicationIds", id)
	}

	for _, id := range p.PlanIDs {
		q.Add("planIds", id)
	}

	for _, m := range p.Methods {
		q.Add("methods", m)
	}

	data, err := s.client.Get(s.v2(fmt.Sprintf("apis/%s/logs?%s", apiID, q.Encode())))
	if err != nil {
		return nil, fmt.Errorf("API logs failed: %w", err)
	}

	return parsePaginatedResponse(data)
}

func (s *service) GetAPILog(apiID, requestID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.v2(fmt.Sprintf("apis/%s/logs/%s", apiID, requestID)))
	if err != nil {
		return nil, err
	}

	return raw(data), nil
}
