package domain

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestGetDomain(t *testing.T) {
	t.Run("returns domain details", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "domains/dom-1")

				return json.Marshal(map[string]any{
					"id": "dom-1", "name": "My Domain", "enabled": true, "description": "A domain",
				})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newGetCmd(tc.Factory), "dom-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Domain")
		testutil.AssertOutputContains(t, tc.Out, "dom-1")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(_ string) ([]byte, error) {
				return json.Marshal(map[string]any{"id": "dom-1", "name": "Test"})
			},
		}
		tc := testutil.NewFactory(fake)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newGetCmd(tc.Factory), "dom-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("requires domain ID", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		err := testutil.Execute(newGetCmd(tc.Factory))

		testutil.AssertErrorContains(t, err, "domain ID or --hrid is required")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newGetCmd(tc.Factory), "dom-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
