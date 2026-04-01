package apim

import (
	"encoding/json"
	"fmt"
)

// PluginService defines plugin-related operations.
type PluginService interface {
	ListPlugins(pluginType string) (json.RawMessage, error)
}

func (s *service) ListPlugins(pluginType string) (json.RawMessage, error) {
	path := s.orgV2(fmt.Sprintf("plugins/%s", pluginType))

	data, err := s.client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("plugin list failed: %w", err)
	}

	return raw(data), nil
}
