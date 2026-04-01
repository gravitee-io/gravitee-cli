package page

import (
	"encoding/json"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func pageJSON() map[string]any {
	return map[string]any{
		"id":         "page-1",
		"name":       "Getting Started",
		"apiId":      "api-1",
		"type":       "MARKDOWN",
		"visibility": "PUBLIC",
		"published":  true,
		"parentId":   "folder-1",
		"createdAt":  "2026-03-15T10:00:00Z",
		"updatedAt":  "2026-03-25T14:30:00Z",
	}
}

func TestListPages(t *testing.T) {
	t.Run("returns pages from the API", func(t *testing.T) {
		fake := paginatedPages(map[string]any{
			"id": "page-1", "name": "Getting Started", "type": "MARKDOWN",
			"visibility": "PUBLIC", "published": true,
			"updatedAt": "2026-03-25T14:30:00Z",
		})
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newListCmd(tc.Factory), "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Getting Started")
		testutil.AssertOutputContains(t, tc.Out, "MARKDOWN")
	})
}

func TestGetPage(t *testing.T) {
	t.Run("returns page details", func(t *testing.T) {
		fake := testutil.APIReturningItem(pageJSON())
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newGetCmd(tc.Factory), "page-1", "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Getting Started")
		testutil.AssertOutputContains(t, tc.Out, "MARKDOWN")
		testutil.AssertOutputContains(t, tc.Out, "true")
	})

	t.Run("reports not found error", func(t *testing.T) {
		fake := testutil.APIFailingWith(404, "resource not found (HTTP 404)")
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newGetCmd(tc.Factory), "page-999", "--api", "api-1")

		testutil.AssertErrorContains(t, err, "not found")
	})
}

func TestCreatePage(t *testing.T) {
	t.Run("creates a page from a JSON file", func(t *testing.T) {
		file := writeTempJSON(t, `{"name":"Getting Started","type":"MARKDOWN"}`)
		resp, _ := json.Marshal(pageJSON())
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/pages")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newCreateCmd(tc.Factory), "--api", "api-1", "-f", file)

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Getting Started")
	})

}

func TestDeletePage(t *testing.T) {
	t.Run("deletes the page", func(t *testing.T) {
		fake := &client.FakeClient{
			DeleteFunc: func(path string) error {
				testutil.AssertPathCalled(t, path, "/apis/api-1/pages/page-1")

				return nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newDeleteCmd(tc.Factory), "page-1", "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Page 'page-1' deleted.")
	})

}

func TestPublishPage(t *testing.T) {
	t.Run("publishes the page", func(t *testing.T) {
		resp, _ := json.Marshal(map[string]any{
			"id": "page-1", "name": "Getting Started", "apiId": "api-1",
			"type": "MARKDOWN", "visibility": "PUBLIC", "published": true,
			"createdAt": "2026-03-15T10:00:00Z", "updatedAt": "2026-03-25T14:30:00Z",
		})
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/pages/page-1/_publish")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newPublishCmd(tc.Factory), "page-1", "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "true")
		testutil.AssertOutputContains(t, tc.Out, "Getting Started")
	})

}

func TestUnpublishPage(t *testing.T) {
	t.Run("unpublishes the page", func(t *testing.T) {
		resp, _ := json.Marshal(map[string]any{
			"id": "page-1", "name": "Getting Started", "apiId": "api-1",
			"type": "MARKDOWN", "visibility": "PUBLIC", "published": false,
			"createdAt": "2026-03-15T10:00:00Z", "updatedAt": "2026-03-30T09:00:00Z",
		})
		fake := &client.FakeClient{
			PostFunc: func(path string, _ any) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/apis/api-1/pages/page-1/_unpublish")

				return resp, nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newUnpublishCmd(tc.Factory), "page-1", "--api", "api-1")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "false")
		testutil.AssertOutputContains(t, tc.Out, "Getting Started")
	})

}
