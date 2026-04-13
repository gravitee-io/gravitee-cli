package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTPClient implements GraviteeClient using net/http.
type HTTPClient struct {
	httpClient *http.Client
	debugOut   io.Writer
	baseURL    string
	token      string
	debug      bool
}

// HTTPClientConfig holds the configuration for creating an HTTPClient.
type HTTPClientConfig struct {
	DebugOut io.Writer
	BaseURL  string
	Token    string
	Debug    bool
}

// NewHTTPClient creates a new HTTPClient.
func NewHTTPClient(cfg HTTPClientConfig) *HTTPClient {
	return &HTTPClient{
		baseURL:    strings.TrimRight(cfg.BaseURL, "/"),
		token:      cfg.Token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		debug:      cfg.Debug,
		debugOut:   cfg.DebugOut,
	}
}

// V2Path builds a v2 management API path.
func V2Path(envID string, path string) string {
	return fmt.Sprintf("/management/v2/environments/%s/%s", envID, strings.TrimLeft(path, "/"))
}

// V1Path builds a v1 management API path.
func V1Path(orgID string, envID string, path string) string {
	return fmt.Sprintf("/management/organizations/%s/environments/%s/%s", orgID, envID, strings.TrimLeft(path, "/"))
}

func (c *HTTPClient) Get(path string) ([]byte, error) {
	return c.doRequest(http.MethodGet, path, nil)
}

func (c *HTTPClient) Post(path string, body any) ([]byte, error) {
	return c.doRequest(http.MethodPost, path, body)
}

func (c *HTTPClient) Put(path string, body any) ([]byte, error) {
	return c.doRequest(http.MethodPut, path, body)
}

func (c *HTTPClient) Patch(path string, body any) ([]byte, error) {
	return c.doRequest(http.MethodPatch, path, body)
}

func (c *HTTPClient) Delete(path string) error {
	_, err := c.doRequest(http.MethodDelete, path, nil)

	return err
}

func (c *HTTPClient) doRequest(method, path string, body any) ([]byte, error) {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(context.Background(), method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	if c.debug && c.debugOut != nil {
		c.logRequest(method, path)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if c.debug && c.debugOut != nil {
		c.logResponse(resp.StatusCode)
	}

	if resp.StatusCode >= 400 {
		return nil, MapHTTPError(resp.StatusCode, respBody)
	}

	return respBody, nil
}

func (c *HTTPClient) logRequest(method, path string) {
	maskedToken := maskToken(c.token)
	fmt.Fprintf(c.debugOut, "> %s %s\n> Authorization: Bearer %s\n>\n", method, path, maskedToken)
}

func (c *HTTPClient) logResponse(status int) {
	fmt.Fprintf(c.debugOut, "< HTTP %d\n<\n", status)
}

func maskToken(token string) string {
	if len(token) <= 3 {
		return "***"
	}

	return strings.Repeat("*", len(token)-3) + token[len(token)-3:]
}
