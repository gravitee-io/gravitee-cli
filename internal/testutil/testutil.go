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

package testutil

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/config"
	"gravitee.io/gctl/internal/factory"
)

// TestContext holds a Factory and captured output buffers for testing.
type TestContext struct {
	Factory *factory.Factory
	Out     *bytes.Buffer
}

// NewAMTestFactory creates a Factory configured for AM testing.
func NewAMTestFactory(c client.GraviteeClient, domainID string) *TestContext {
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}

	return &TestContext{
		Factory: &factory.Factory{
			Config: &config.Config{
				Current: "am-test",
				Contexts: map[string]*config.Context{
					"am-test": {
						Org:    "DEFAULT",
						Env:    "DEFAULT",
						Type:   "am",
						Domain: domainID,
						AM:     &config.ProductConfig{URL: "https://am-test.company.com", Token: "am-test-token"},
					},
				},
			},
			Resolved: &config.ResolvedContext{
				Name: "am-test", URL: "https://am-test.company.com", Token: "am-test-token",
				Org: "DEFAULT", Env: "DEFAULT",
				Type: "am", Domain: domainID,
			},
			Client:       c,
			IOStreams:    factory.IOStreams{Out: out, Err: errOut, In: &bytes.Buffer{}},
			OutputFormat: "table",
		},
		Out: out,
	}
}

// NewFactory creates a Factory configured for testing with the given client.
func NewFactory(c client.GraviteeClient) *TestContext {
	out := &bytes.Buffer{}

	return &TestContext{
		Factory: &factory.Factory{
			Config: &config.Config{
				Current: "test",
				Contexts: map[string]*config.Context{
					"test": {
						Org:  "DEFAULT",
						Env:  "DEFAULT",
						APIM: &config.ProductConfig{URL: "https://test.com", Token: "tok"},
						AM:   &config.ProductConfig{URL: "https://test.com", Token: "tok"},
					},
				},
			},
			Resolved: &config.ResolvedContext{
				Name: "test", URL: "https://test.com", Token: "tok",
				Org: "DEFAULT", Env: "DEFAULT",
			},
			Client:       c,
			IOStreams:    factory.IOStreams{Out: out, Err: &bytes.Buffer{}, In: &bytes.Buffer{}},
			OutputFormat: "table",
		},
		Out: out,
	}
}

// NoOpClient is a FakeClient with no configured methods. Any call will error.
var NoOpClient = client.FakeClient{}

// APIReturning creates a FakeClient whose Get returns a paginated response with the given items.
func APIReturning(items []map[string]any) *client.FakeClient {
	return &client.FakeClient{
		GetFunc: func(_ string) ([]byte, error) {
			resp := map[string]any{
				"data":       items,
				"pagination": map[string]int{"page": 1, "perPage": 10, "pageCount": 1, "totalCount": len(items), "pageItemsCount": len(items)},
			}

			data, _ := json.Marshal(resp)

			return data, nil
		},
	}
}

// APIReturningItem creates a FakeClient whose Get returns a single item.
func APIReturningItem(item map[string]any) *client.FakeClient {
	return &client.FakeClient{
		GetFunc: func(_ string) ([]byte, error) {
			data, _ := json.Marshal(item)

			return data, nil
		},
	}
}

// APIFailingWith creates a FakeClient whose Get and Post return an API error.
func APIFailingWith(status int, message string) *client.FakeClient {
	err := &client.APIError{Status: status, Message: message}

	return &client.FakeClient{
		GetFunc:  func(_ string) ([]byte, error) { return nil, err },
		PostFunc: func(_ string, _ any) ([]byte, error) { return nil, err },
	}
}

// PostSucceeding creates a FakeClient whose Post returns the given response.
func PostSucceeding(resp map[string]any) *client.FakeClient {
	return &client.FakeClient{
		PostFunc: func(_ string, _ any) ([]byte, error) {
			data, _ := json.Marshal(resp)

			return data, nil
		},
	}
}

// PostFailingWith creates a FakeClient whose Post returns an API error.
func PostFailingWith(status int, message string) *client.FakeClient {
	return &client.FakeClient{
		PostFunc: func(_ string, _ any) ([]byte, error) {
			return nil, &client.APIError{Status: status, Message: message}
		},
	}
}

// DeleteSucceeding creates a FakeClient whose Delete succeeds.
func DeleteSucceeding() *client.FakeClient {
	return &client.FakeClient{
		DeleteFunc: func(_ string) error { return nil },
	}
}

// Execute runs a cobra command with the given args and silences usage/error output.
func Execute(cmd *cobra.Command, args ...string) error {
	cmd.SetArgs(args)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd.Execute()
}

// AssertNoError fails the test if err is not nil.
func AssertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// AssertErrorContains fails the test if err is nil or doesn't contain the expected substring.
func AssertErrorContains(t *testing.T, err error, expected string) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected error containing %q, got nil", expected)
	}

	if !strings.Contains(err.Error(), expected) {
		t.Errorf("expected error containing %q, got: %v", expected, err)
	}
}

// AssertOutputContains fails the test if the buffer doesn't contain the expected substring.
func AssertOutputContains(t *testing.T, out *bytes.Buffer, expected string) {
	t.Helper()

	if !strings.Contains(out.String(), expected) {
		t.Errorf("expected output containing %q, got: %s", expected, out.String())
	}
}

// AssertPathCalled fails the test if the path doesn't contain the expected substring.
func AssertPathCalled(t *testing.T, actual, expected string) {
	t.Helper()

	if !strings.Contains(actual, expected) {
		t.Errorf("expected path containing %q, got: %s", expected, actual)
	}
}
