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
	"fmt"

	"gravitee.io/gctl/internal/client"
)

// MetadataService defines metadata-related operations.
// Note: V2 API only supports listing metadata. Create/update/delete are not available in V2.
type MetadataService interface {
	ListMetadata(apiID string, page, perPage int) (*PaginatedResponse, error)
}

func (s *service) ListMetadata(apiID string, page, perPage int) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{"page": client.Itoa(page), "perPage": client.Itoa(perPage)})

	data, err := s.client.Get(s.v2(fmt.Sprintf("apis/%s/metadata?%s", apiID, q)))
	if err != nil {
		return nil, fmt.Errorf("metadata list failed: %w", err)
	}

	return parsePaginatedResponse(data)
}
