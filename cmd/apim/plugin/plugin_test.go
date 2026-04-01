package plugin

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/testutil"
)

func TestListPlugins(t *testing.T) {
	t.Run("returns plugins filtered by type", func(t *testing.T) {
		fake := pluginList(
			map[string]string{"id": "rate-limit", "name": "Rate Limiting", "version": "4.5.0", "description": "Rate limiting policy"},
			map[string]string{"id": "api-key", "name": "API Key", "version": "4.5.0", "description": "API key validation policy"},
		)
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newListCmd(tc.Factory), "--type", "policies")

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Rate Limiting")
		testutil.AssertOutputContains(t, tc.Out, "API Key")
	})

	t.Run("calls the correct API path for a given type", func(t *testing.T) {
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				testutil.AssertPathCalled(t, path, "/management/v2/organizations/DEFAULT/plugins/policies")

				data, _ := json.Marshal([]map[string]string{})

				return data, nil
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newListCmd(tc.Factory), "--type", "policies")

		testutil.AssertNoError(t, err)
	})

	t.Run("hides TYPE column when --type is set", func(t *testing.T) {
		fake := pluginList(
			map[string]string{"id": "rate-limit", "name": "Rate Limiting", "version": "4.5.0"},
		)
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newListCmd(tc.Factory), "--type", "policies")

		testutil.AssertNoError(t, err)
		lines := strings.Split(strings.TrimSpace(tc.Out.String()), "\n")
		if len(lines) > 0 && strings.Contains(lines[0], "TYPE") {
			t.Errorf("TYPE column should not appear when --type is set, got header: %s", lines[0])
		}
	})

	t.Run("returns all plugin types when no --type is set", func(t *testing.T) {
		endpointsData, _ := json.Marshal([]map[string]string{
			{"id": "kafka", "name": "Kafka", "version": "4.5.0", "description": "Kafka endpoint connector"},
		})
		entrypointsData, _ := json.Marshal([]map[string]string{
			{"id": "http-proxy", "name": "HTTP Proxy", "version": "4.5.0", "description": "HTTP proxy entrypoint"},
		})
		policiesData, _ := json.Marshal([]map[string]string{
			{"id": "rate-limit", "name": "Rate Limiting", "version": "4.5.0", "description": "Rate limiting policy"},
		})
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				switch path {
				case "/management/v2/organizations/DEFAULT/plugins/endpoints":
					return endpointsData, nil
				case "/management/v2/organizations/DEFAULT/plugins/entrypoints":
					return entrypointsData, nil
				case "/management/v2/organizations/DEFAULT/plugins/policies":
					return policiesData, nil
				default:
					t.Errorf("unexpected path: %s", path)

					return nil, nil
				}
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "Kafka")
		testutil.AssertOutputContains(t, tc.Out, "HTTP Proxy")
		testutil.AssertOutputContains(t, tc.Out, "Rate Limiting")
	})

	t.Run("shows TYPE column when no --type is set", func(t *testing.T) {
		endpointsData, _ := json.Marshal([]map[string]string{
			{"id": "kafka", "name": "Kafka", "version": "4.5.0"},
		})
		entrypointsData, _ := json.Marshal([]map[string]string{})
		policiesData, _ := json.Marshal([]map[string]string{})
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				switch path {
				case "/management/v2/organizations/DEFAULT/plugins/endpoints":
					return endpointsData, nil
				case "/management/v2/organizations/DEFAULT/plugins/entrypoints":
					return entrypointsData, nil
				case "/management/v2/organizations/DEFAULT/plugins/policies":
					return policiesData, nil
				default:
					return nil, nil
				}
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertNoError(t, err)
		lines := strings.Split(strings.TrimSpace(tc.Out.String()), "\n")
		if len(lines) > 0 && !strings.Contains(lines[0], "TYPE") {
			t.Errorf("TYPE column should appear when --type is not set, got header: %s", lines[0])
		}
	})

	t.Run("uses singular type labels", func(t *testing.T) {
		endpointsData, _ := json.Marshal([]map[string]string{
			{"id": "kafka", "name": "Kafka", "version": "4.5.0"},
		})
		entrypointsData, _ := json.Marshal([]map[string]string{
			{"id": "http-proxy", "name": "HTTP Proxy", "version": "4.5.0"},
		})
		policiesData, _ := json.Marshal([]map[string]string{
			{"id": "rate-limit", "name": "Rate Limiting", "version": "4.5.0"},
		})
		fake := &client.FakeClient{
			GetFunc: func(path string) ([]byte, error) {
				switch path {
				case "/management/v2/organizations/DEFAULT/plugins/endpoints":
					return endpointsData, nil
				case "/management/v2/organizations/DEFAULT/plugins/entrypoints":
					return entrypointsData, nil
				case "/management/v2/organizations/DEFAULT/plugins/policies":
					return policiesData, nil
				default:
					return nil, nil
				}
			},
		}
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newListCmd(tc.Factory))

		testutil.AssertNoError(t, err)
		testutil.AssertOutputContains(t, tc.Out, "endpoint")
		testutil.AssertOutputContains(t, tc.Out, "entrypoint")
		testutil.AssertOutputContains(t, tc.Out, "policy")
	})

	t.Run("rejects invalid token with hint", func(t *testing.T) {
		fake := testutil.APIFailingWith(401, "authentication failed (HTTP 401)")
		tc := testutil.NewFactory(fake, false)

		err := testutil.Execute(newListCmd(tc.Factory), "--type", "policies")

		testutil.AssertErrorContains(t, err, "authentication failed")
	})

	t.Run("rejects invalid type flag", func(t *testing.T) {
		tc := testutil.NewFactory(&testutil.NoOpClient, false)

		err := testutil.Execute(newListCmd(tc.Factory), "--type", "connectors")

		testutil.AssertErrorContains(t, err, "invalid value 'connectors' for flag --type")
	})
}
