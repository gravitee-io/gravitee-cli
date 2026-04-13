package domain

import (
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestDeleteDomain(t *testing.T) {
	t.Run("deletes a domain", func(t *testing.T) {
		fake := &client.FakeClient{
			DeleteFunc: func(path string) error {
				testutil.AssertPathCalled(t, path, "domains/dom-1")

				return nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newDeleteCmd(tc.Factory), "dom-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Domain 'dom-1' deleted.")
	})

	t.Run("returns error on API failure", func(t *testing.T) {
		fake := &client.FakeClient{
			DeleteFunc: func(_ string) error {
				return &client.APIError{Status: 404, Message: "resource not found (HTTP 404)"}
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newDeleteCmd(tc.Factory), "dom-1")

		testutil.AssertErrorContains(t, err, "not found")
	})

	t.Run("requires domain ID argument", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		err := testutil.Execute(newDeleteCmd(tc.Factory))

		testutil.AssertErrorContains(t, err, "accepts 1 arg")
	})

	t.Run("requires a configured context", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)
		tc.Factory.Resolved = nil

		err := testutil.Execute(newDeleteCmd(tc.Factory), "dom-1")

		testutil.AssertErrorContains(t, err, "no context configured")
	})
}
