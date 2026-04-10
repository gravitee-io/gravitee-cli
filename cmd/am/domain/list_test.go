package domain

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestListDomains(t *testing.T) {
	t.Run("returns domains from the environment", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(_ string) ([]byte, error) {
				return json.Marshal(map[string]any{
					"data": []map[string]any{
						{"id": "dom-1", "name": "My Domain", "enabled": true},
						{"id": "dom-2", "name": "Other", "enabled": false},
					},
					"totalCount":  2,
					"currentPage": 0,
				})
			},
		}
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "My Domain")
		testutil.AssertOutputContains(t, tc.Out, "Other")
	})

	t.Run("returns full JSON with -o json", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(_ string) ([]byte, error) {
				return json.Marshal(map[string]any{
					"data":        []map[string]any{{"id": "dom-1", "name": "Test"}},
					"totalCount":  1,
					"currentPage": 0,
				})
			},
		}
		tc := testutil.NewFactory(fake)
		tc.Factory.OutputFormat = "json"

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, `"totalCount"`)
	})

	t.Run("rejects invalid token with hint", func(t *testing.T) {
		fake := testutil.APIFailingWith(401, "authentication failed")
		tc := testutil.NewFactory(fake)

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertErrorContains(t, err, "authentication failed")
	})

	t.Run("rejects page zero", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		err := testutil.Execute(newListCmd(tc.Factory), "--page", "0")

		testutil.AssertErrorContains(t, err, "--page must be >= 1")
	})

	t.Run("rejects per-page zero", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)

		err := testutil.Execute(newListCmd(tc.Factory), "--per-page", "0")

		testutil.AssertErrorContains(t, err, "--per-page must be >= 1")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
