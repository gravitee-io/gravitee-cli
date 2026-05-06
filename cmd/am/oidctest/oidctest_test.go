package oidctest

import (
	"testing"
)

func TestDeriveGatewayURL(t *testing.T) {
	cases := []struct {
		mgmtURL  string
		expected string
	}{
		{"http://am.example.com:8093", "http://am.example.com:8092"},
		{"https://am.example.com", "https://am.example.com:8092"},
		{"http://localhost:8093", "http://localhost:8092"},
	}
	for _, tc := range cases {
		got := deriveGatewayURL(tc.mgmtURL)
		if got != tc.expected {
			t.Errorf("deriveGatewayURL(%q) = %q, want %q", tc.mgmtURL, got, tc.expected)
		}
	}
}

func TestDecodeJWT(t *testing.T) {
	// A valid JWT with known payload (header.payload.signature)
	// eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJ1c2VyMSIsImV4cCI6OTk5OTk5OTk5OX0.sig
	//nolint:gosec
	token := "eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJ1c2VyMSIsImV4cCI6OTk5OTk5OTk5OX0.test_signature"
	header, payload, err := decodeJWT(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if header["alg"] != "RS256" {
		t.Errorf("expected alg RS256, got %v", header["alg"])
	}
	if payload["sub"] != "user1" {
		t.Errorf("expected sub user1, got %v", payload["sub"])
	}
}

func TestTruncateToken(t *testing.T) {
	cases := []struct {
		token    string
		maxLen   int
		expected string
	}{
		{"abcdefghijklmnopqrstuvwxyz1234567890", 10, "abcdefghij...(truncated)"},
		{"short", 10, "short"},
		{"exactlyten", 10, "exactlyten"},
	}
	for _, tc := range cases {
		got := truncateToken(tc.token, tc.maxLen)
		if got != tc.expected {
			t.Errorf("truncateToken(%q, %d) = %q, want %q", tc.token, tc.maxLen, got, tc.expected)
		}
	}
}
