package domain

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestCreateDomain(t *testing.T) {
	t.Run("creates a domain with name", func(t *testing.T) {
		fake := &client.FakeClient{
			PostFunc: func(path string, body any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "domains")

				return json.Marshal(map[string]any{
					"id": "new-dom", "name": "My Domain", "enabled": false,
				})
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newCreateCmd(tc.Factory), "--name", "My Domain")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Domain")
		testutil.AssertOutputContains(t, tc.Out, "new-dom")
	})

	t.Run("creates a domain with name and description", func(t *testing.T) {
		fake := &client.FakeClient{
			PostFunc: func(_ string, _ any) ([]byte, error) {
				return json.Marshal(map[string]any{
					"id": "new-dom", "name": "My Domain", "description": "Desc",
				})
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newCreateCmd(tc.Factory), "--name", "My Domain", "--description", "Desc")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Desc")
	})

	t.Run("returns JSON with -o json", func(t *testing.T) {
		fake := &client.FakeClient{
			PostFunc: func(_ string, _ any) ([]byte, error) {
				return json.Marshal(map[string]any{"id": "new-dom", "name": "Test"})
			},
		}
		tc := testutil.NewFactory(fake, false)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newCreateCmd(tc.Factory), "--name", "Test")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"id"`)
	})

	t.Run("includes dataPlaneId in request body", func(t *testing.T) {
		var sentBody map[string]any
		fake := &client.FakeClient{
			PostFunc: func(_ string, body any) ([]byte, error) {
				raw, _ := body.(json.RawMessage)
				_ = json.Unmarshal(raw, &sentBody)

				return json.Marshal(map[string]any{"id": "new", "name": "Test"})
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newCreateCmd(tc.Factory), "--name", "Test")

		testutil.AssertNoError(t, err)
		if sentBody["dataPlaneId"] != "default" {
			t.Errorf("expected dataPlaneId 'default', got %v", sentBody["dataPlaneId"])
		}
	})

	t.Run("requires name flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		err := testutil.Execute(newCreateCmd(tc.Factory))

		testutil.AssertErrorContains(t, err, "required")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newCreateCmd(tc.Factory), "--name", "Test")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
