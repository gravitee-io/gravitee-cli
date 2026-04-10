package app

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestUpdateApplication(t *testing.T) {
	t.Run("updates application name", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			PatchApplicationFunc: func(domainID, appID string, body json.RawMessage) (json.RawMessage, error) {
				if domainID != "dom-1" {
					t.Errorf("expected domainID 'dom-1', got %q", domainID)
				}

				if appID != "app-1" {
					t.Errorf("expected appID 'app-1', got %q", appID)
				}

				return json.Marshal(map[string]any{
					"id": "app-1", "name": "Updated", "type": "web", "enabled": true,
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "app-1", "--name", "Updated")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Updated")
	})

	t.Run("updates application enabled status", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		mock := &am.MockService{
			PatchApplicationFunc: func(_ string, _ string, body json.RawMessage) (json.RawMessage, error) {
				var m map[string]any
				_ = json.Unmarshal(body, &m)

				if m["enabled"] != false {
					t.Errorf("expected enabled=false, got %v", m["enabled"])
				}

				return json.Marshal(map[string]any{
					"id": "app-1", "name": "Test", "enabled": false,
				})
			},
		}
		tc.Factory.SetAMService(mock)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "app-1", "--enabled", "false")

		testutil.AssertNoError(t, err)
	})

	t.Run("rejects invalid enabled value", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "app-1", "--enabled", "maybe")

		testutil.AssertErrorContains(t, err, "must be 'true' or 'false'")
	})

	t.Run("requires at least one flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "app-1")

		testutil.AssertErrorContains(t, err, "at least one flag")
	})

	t.Run("requires app ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update")

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		cmd := NewAppCmd(tc.Factory)
		err := testutil.Execute(cmd, "--domain", "dom-1", "update", "app-1", "--name", "Test")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
