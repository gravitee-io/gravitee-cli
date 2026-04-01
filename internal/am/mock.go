package am

import (
	"encoding/json"
	"fmt"
)

// MockService implements Service with injectable functions for testing.
type MockService struct {
	ListDomainsFunc func(ListDomainsParams) (*PaginatedResponse, error)
	GetDomainFunc   func(string) (json.RawMessage, error)
}

func unexpected(name string) error { return fmt.Errorf("unexpected call: %s", name) }

func (m *MockService) ListDomains(p ListDomainsParams) (*PaginatedResponse, error) {
	if m.ListDomainsFunc != nil {
		return m.ListDomainsFunc(p)
	}

	return nil, unexpected("ListDomains")
}

func (m *MockService) GetDomain(id string) (json.RawMessage, error) {
	if m.GetDomainFunc != nil {
		return m.GetDomainFunc(id)
	}

	return nil, unexpected("GetDomain")
}
