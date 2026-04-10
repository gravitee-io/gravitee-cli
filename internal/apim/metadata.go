package apim

import (
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/client"
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
