package domain

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

// amPaginatedResponse is the AM paginated API response format (0-based pagination).
type amPaginatedResponse struct {
	Data        []json.RawMessage `json:"data"`
	CurrentPage int               `json:"currentPage"`
	TotalCount  int               `json:"totalCount"`
}

func newExportCmd(f *factory.Factory) *cobra.Command {
	var file string
	cmd := &cobra.Command{
		Use:     "export <domainId>",
		Short:   "Export domain configuration to JSON",
		Example: `  gio am domain export abc-123 -f domain-export.json`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireAMContext(f); err != nil {
				return err
			}
			return runExport(f, args[0], file)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "Output file path (default: stdout)")
	return cmd
}

func runExport(f *factory.Factory, domainID, file string) error {
	export, err := exportToMemory(f, domainID)
	if err != nil {
		return err
	}
	out, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return err
	}
	if file != "" {
		return os.WriteFile(file, out, 0600)
	}
	fmt.Fprintln(f.IOStreams.Out, string(out))
	return nil
}

func exportToMemory(f *factory.Factory, domainID string) (map[string]json.RawMessage, error) {
	domainPath := cmdutil.AMEnvPath(f, fmt.Sprintf("domains/%s", domainID))
	domainData, err := f.Client.Get(domainPath)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch domain: %w", err)
	}

	jobs := buildExportJobs(f, domainID)
	results, err := runExportJobs(jobs)
	if err != nil {
		return nil, err
	}

	return map[string]json.RawMessage{
		"domain":            domainData,
		"applications":      results["applications"],
		"identityProviders": results["identityProviders"],
		"roles":             results["roles"],
		"scopes":            results["scopes"],
		"factors":           results["factors"],
		"groups":            results["groups"],
		"flows":             results["flows"],
		"certificates":      results["certificates"],
	}, nil
}

func buildExportJobs(f *factory.Factory, domainID string) []struct {
	key string
	fn  func() (json.RawMessage, error)
} {
	return []struct {
		key string
		fn  func() (json.RawMessage, error)
	}{
		{"applications", func() (json.RawMessage, error) {
			return fetchAllPaginated(f, domainID, "applications")
		}},
		{"identityProviders", func() (json.RawMessage, error) {
			data, err := f.Client.Get(cmdutil.AMDomainPathFor(f, domainID, "identityProviders"))
			if err != nil {
				return nil, err
			}
			return json.RawMessage(data), nil
		}},
		{"roles", func() (json.RawMessage, error) {
			return fetchAllPaginated(f, domainID, "roles")
		}},
		{"scopes", func() (json.RawMessage, error) {
			return fetchAllPaginated(f, domainID, "scopes")
		}},
		{"factors", func() (json.RawMessage, error) {
			data, err := f.Client.Get(cmdutil.AMDomainPathFor(f, domainID, "factors"))
			if err != nil {
				return nil, err
			}
			return json.RawMessage(data), nil
		}},
		{"groups", func() (json.RawMessage, error) {
			return fetchAllPaginated(f, domainID, "groups")
		}},
		{"flows", func() (json.RawMessage, error) {
			data, err := f.Client.Get(cmdutil.AMDomainPathFor(f, domainID, "flows"))
			if err != nil {
				return nil, err
			}
			return json.RawMessage(data), nil
		}},
		{"certificates", func() (json.RawMessage, error) {
			data, err := f.Client.Get(cmdutil.AMDomainPathFor(f, domainID, "certificates"))
			if err != nil {
				return nil, err
			}
			return json.RawMessage(data), nil
		}},
	}
}

func runExportJobs(jobs []struct {
	key string
	fn  func() (json.RawMessage, error)
}) (map[string]json.RawMessage, error) {
	results := make(map[string]json.RawMessage)
	var mu sync.Mutex
	var wg sync.WaitGroup
	var firstErr error

	for _, job := range jobs {
		wg.Add(1)
		go func(j struct {
			key string
			fn  func() (json.RawMessage, error)
		}) {
			defer wg.Done()
			data, err := j.fn()
			mu.Lock()
			defer mu.Unlock()
			if err != nil && firstErr == nil {
				firstErr = fmt.Errorf("fetch %s: %w", j.key, err)
				return
			}
			results[j.key] = data
		}(job)
	}
	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}
	return results, nil
}

func fetchAllPaginated(f *factory.Factory, domainID, resource string) (json.RawMessage, error) {
	var all []json.RawMessage
	size := 100
	for page := 0; page <= 1000; page++ {
		q := url.Values{}
		q.Set("page", strconv.Itoa(page))
		q.Set("size", strconv.Itoa(size))
		path := cmdutil.AMDomainPathFor(f, domainID, resource+"?"+q.Encode())
		data, err := f.Client.Get(path)
		if err != nil {
			return nil, err
		}
		var resp amPaginatedResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse %s response: %w", resource, err)
		}
		all = append(all, resp.Data...)
		if len(all) >= resp.TotalCount || len(resp.Data) < size {
			break
		}
	}
	b, err := json.Marshal(all)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}
