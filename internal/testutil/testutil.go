package testutil

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// TestContext holds a Factory and captured output buffers for testing.
type TestContext struct {
	Factory *factory.Factory
	Out     *bytes.Buffer
}

// NewFactory creates a Factory configured for testing with the given client and read-only setting.
func NewFactory(c client.GraviteeClient, readOnly bool) *TestContext {
	out := &bytes.Buffer{}

	return &TestContext{
		Factory: &factory.Factory{
			Config: &config.Config{
				Current: "test",
				Contexts: map[string]*config.Context{
					"test": {
						Org: "DEFAULT", Env: "DEFAULT", ReadOnly: readOnly,
						APIM: &config.ProductConfig{URL: "https://test.com", Token: "tok"},
						AM:   &config.ProductConfig{URL: "https://test.com", Token: "tok"},
					},
				},
			},
			Resolved: &config.ResolvedContext{
				Name: "test", URL: "https://test.com", Token: "tok",
				Org: "DEFAULT", Env: "DEFAULT", ReadOnly: readOnly,
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

// --- Fake API builders ---

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

// APIFailingWith creates a FakeClient whose Get returns an API error.
func APIFailingWith(status int, message string) *client.FakeClient {
	return &client.FakeClient{
		GetFunc: func(_ string) ([]byte, error) {
			return nil, &client.APIError{Status: status, Message: message}
		},
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

// --- Command execution ---

// Execute runs a cobra command with the given args and silences usage/error output.
func Execute(cmd *cobra.Command, args ...string) error {
	cmd.SetArgs(args)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd.Execute()
}

// --- Assertions ---

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
