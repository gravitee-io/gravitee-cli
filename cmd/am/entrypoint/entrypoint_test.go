package entrypoint

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

// --- Get ---

func TestGetEntrypoints(t *testing.T) {
	t.Run("returns entrypoint data", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			GetEntrypointsFunc: func(domainID string) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				return json.Marshal(map[string]any{
					"url": "https://example.com",
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewEntrypointCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "example.com")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		mock := &am.MockService{
			GetEntrypointsFunc: func(_ string) (json.RawMessage, error) {
				return json.Marshal(map[string]any{"url": "https://example.com"})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewEntrypointCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "-o", "json", "get")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"url"`)
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		cmd := NewEntrypointCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "get")

		testutil.AssertErrorContains(t, err, "no context configured")
	})

	t.Run("requires domain flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		cmd := NewEntrypointCmd(tc.Factory)
		err := testutil.Execute(cmd, "get")

		testutil.AssertErrorContains(t, err, "required")
	})
}
