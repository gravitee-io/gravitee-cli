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
func MapHTTPError(status int, body []byte) *APIError {
	switch status {
	case 401:
		return &APIError{
			Status:  status,
			Message: fmt.Sprintf("authentication failed (HTTP %d)\nHint: run 'gio login' to configure your credentials", status),
		}
	case 403:
		return &APIError{
			Status:  status,
			Message: fmt.Sprintf("access denied (HTTP %d)\nHint: check your token permissions for this operation", status),
		}
	case 404:
		return &APIError{
			Status:  status,
			Message: fmt.Sprintf("resource not found (HTTP %d)", status),
		}
	case 400:
		return &APIError{
			Status:  status,
			Message: fmt.Sprintf("invalid request (HTTP %d): %s", status, sanitizeBody(body)),
		}
	case 409:
		return &APIError{
			Status:  status,
			Message: fmt.Sprintf("conflict (HTTP %d): %s", status, sanitizeBody(body)),
		}
	default:
		if status >= 500 {
			return &APIError{
				Status:  status,
				Message: fmt.Sprintf("server error (HTTP %d)\nHint: try again or check the APIM server status", status),
			}
		}

		return &APIError{
			Status:  status,
			Message: fmt.Sprintf("unexpected error (HTTP %d): %s", status, sanitizeBody(body)),
		}
	}
}

const maxBodyLen = 500

func sanitizeBody(body []byte) string {
	s := string(body)

	lower := strings.ToLower(s)
	if strings.Contains(lower, "token") || strings.Contains(lower, "authorization") {
		return "[redacted: response may contain sensitive data]"
	}

	if len(s) > maxBodyLen {
		return s[:maxBodyLen] + "... (truncated)"
	}

	return s
}
