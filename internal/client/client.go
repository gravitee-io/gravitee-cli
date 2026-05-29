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

package client

import (
	"fmt"
	"strings"
)

// GraviteeClient defines the operations for communicating with the Gravitee APIM API.
type GraviteeClient interface {
	Get(path string) ([]byte, error)
	Post(path string, body any) ([]byte, error)
	Put(path string, body any) ([]byte, error)
	Patch(path string, body any) ([]byte, error)
	Delete(path string) error
}

// APIError represents an error returned by the Gravitee API.
type APIError struct {
	Message string
	Status  int
}

func (e *APIError) Error() string {
	return e.Message
}

// MapHTTPError maps an HTTP status code to a user-friendly APIError.
// Every branch includes the server body (when non-empty) and an actionable
// hint (when one applies) - callers need the server detail to diagnose.
func MapHTTPError(status int, body []byte) *APIError {
	switch status {
	case 400:
		return newAPIError(status, "invalid request", body, "")
	case 401:
		return newAPIError(status, "authentication failed", body, "run 'gctl login' to configure your credentials")
	case 403:
		return newAPIError(status, "access denied", body, "check your token permissions for this operation")
	case 404:
		return newAPIError(status, "resource not found", body, "")
	case 409:
		return newAPIError(status, "conflict", body, "")
	default:
		if status >= 500 {
			return newAPIError(status, "server error", body, "try again or check the APIM server status")
		}

		return newAPIError(status, "unexpected error", body, "")
	}
}

func newAPIError(status int, label string, body []byte, hint string) *APIError {
	msg := fmt.Sprintf("%s (HTTP %d)", label, status)

	if trimmed := strings.TrimSpace(string(body)); trimmed != "" {
		msg += ": " + sanitizeBody(body)
	}

	if hint != "" {
		msg += "\nHint: " + hint
	}

	return &APIError{Status: status, Message: msg}
}

const maxBodyLen = 500

func sanitizeBody(body []byte) string {
	s := string(body)
	if len(s) > maxBodyLen {
		return s[:maxBodyLen] + "... (truncated)"
	}

	return s
}
