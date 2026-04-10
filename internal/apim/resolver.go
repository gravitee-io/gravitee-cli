package apim

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ResolveAPI returns the API id for an id or a context path.
// A value starting with "/" is resolved by searching APIs and matching V2
// contextPath or V4 listener paths exactly. Any other value is passed through
// unchanged (the server validates the id).
func (s *service) ResolveAPI(pathOrID string) (string, error) {
	if !strings.HasPrefix(pathOrID, "/") {
		return pathOrID, nil
	}

	items, err := FetchAllPages(func(page int) (*PaginatedResponse, error) {
		return s.ListAPIs(ListAPIsParams{Query: pathOrID, Page: page, PerPage: 100})
	})
	if err != nil {
		return "", fmt.Errorf("resolve context path %q: %w", pathOrID, err)
	}

	var matches []string

	for _, rawItem := range items {
		var item map[string]any
		if err := json.Unmarshal(rawItem, &item); err != nil {
			continue
		}

		if !matchesContextPath(item, pathOrID) {
			continue
		}

		if id, _ := item["id"].(string); id != "" {
			matches = append(matches, id)
		}
	}

	switch len(matches) {
	case 0:
		return "", fmt.Errorf("no API found with context path %q (tip: `gio apim api list` to see available paths)", pathOrID)
	case 1:
		return matches[0], nil
	default:
		return "", fmt.Errorf("multiple APIs (%d) share context path %q - use --api <uuid> to disambiguate", len(matches), pathOrID)
	}
}

// matchesContextPath reports whether any of the API's access paths equals target.
// V2 APIs expose contextPath directly; V4 APIs expose listeners[].paths[].path.
func matchesContextPath(item map[string]any, target string) bool {
	if cp, _ := item["contextPath"].(string); cp == target {
		return true
	}

	listeners, _ := item["listeners"].([]any)
	for _, l := range listeners {
		lm, ok := l.(map[string]any)
		if !ok {
			continue
		}

		paths, _ := lm["paths"].([]any)
		for _, p := range paths {
			pm, ok := p.(map[string]any)
			if !ok {
				continue
			}

			if path, _ := pm["path"].(string); path == target {
				return true
			}
		}
	}

	return false
}
