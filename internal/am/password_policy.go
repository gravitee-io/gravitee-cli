package am

import (
	"encoding/json"
	"fmt"
)

// PasswordPolicyService defines password policy-related operations.
type PasswordPolicyService interface {
	ListPasswordPolicies(domainID string) ([]json.RawMessage, error)
	GetPasswordPolicy(domainID, policyID string) (json.RawMessage, error)
	CreatePasswordPolicy(domainID string, body json.RawMessage) (json.RawMessage, error)
	UpdatePasswordPolicy(domainID, policyID string, body json.RawMessage) (json.RawMessage, error)
	DeletePasswordPolicy(domainID, policyID string) error
	GetActivePasswordPolicy(domainID string) (json.RawMessage, error)
	SetDefaultPasswordPolicy(domainID, policyID string) (json.RawMessage, error)
	EvaluatePasswordPolicy(domainID, policyID string, body json.RawMessage) (json.RawMessage, error)
}

func (s *service) ListPasswordPolicies(domainID string) ([]json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "password-policies"))
	if err != nil {
		return nil, fmt.Errorf("password policy list failed: %w", err)
	}

	if len(data) == 0 {
		return nil, nil
	}

	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse password policy list: %w", err)
	}

	return items, nil
}

func (s *service) GetPasswordPolicy(domainID, policyID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("password-policies/%s", policyID)))
	if err != nil {
		return nil, fmt.Errorf("password policy get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) CreatePasswordPolicy(domainID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, "password-policies"), body)
	if err != nil {
		return nil, fmt.Errorf("password policy create failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) UpdatePasswordPolicy(domainID, policyID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("password-policies/%s", policyID)), body)
	if err != nil {
		return nil, fmt.Errorf("password policy update failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) DeletePasswordPolicy(domainID, policyID string) error {
	err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("password-policies/%s", policyID)))
	if err != nil {
		return fmt.Errorf("password policy delete failed: %w", err)
	}

	return nil
}

func (s *service) GetActivePasswordPolicy(domainID string) (json.RawMessage, error) {
	data, err := s.client.Get(s.domainPath(domainID, "password-policies/activePolicy"))
	if err != nil {
		return nil, fmt.Errorf("active password policy get failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) SetDefaultPasswordPolicy(domainID, policyID string) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, fmt.Sprintf("password-policies/%s/default", policyID)), nil)
	if err != nil {
		return nil, fmt.Errorf("set default password policy failed: %w", err)
	}

	return json.RawMessage(data), nil
}

func (s *service) EvaluatePasswordPolicy(domainID, policyID string, body json.RawMessage) (json.RawMessage, error) {
	data, err := s.client.Post(s.domainPath(domainID, fmt.Sprintf("password-policies/%s/evaluate", policyID)), body)
	if err != nil {
		return nil, fmt.Errorf("password policy evaluate failed: %w", err)
	}

	return json.RawMessage(data), nil
}
