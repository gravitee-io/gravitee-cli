package environment

import (
	"encoding/json"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func environmentList(items ...map[string]string) *client.FakeClient {
	return &client.FakeClient{
		GetFunc: func(_ string) ([]byte, error) {
			data, _ := json.Marshal(items)

			return data, nil
		},
	}
}
