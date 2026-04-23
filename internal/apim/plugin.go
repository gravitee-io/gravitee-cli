// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
