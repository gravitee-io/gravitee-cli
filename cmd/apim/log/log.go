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

package log

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/apim"
	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

// NewLogCmd creates the log parent command.
func NewLogCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log",
		Short: "Manage API request logs",
		Args:  cobra.NoArgs,
	}

	cmdutil.AddOutputFlags(cmd, f)
	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newGetCmd(f))

	return cmd
}

func newListCmd(f *factory.Factory) *cobra.Command {
	var (
		apiID          string
		applicationIDs []string
		planIDs        []string
		methods        []string
		from           int64
		to             int64
		page           int
		perPage        int
		all            bool
		follow         bool
	)

	cmd := &cobra.Command{
		Use:   "list --api <apiId>",
		Short: "List API connection logs",
		Example: `  gio apim log list --api /my/api --from 1711497600000 --to 1711584000000
  gio apim log list --api /my/api --methods GET --per-page 20
  gio apim log list --api /my/api --follow`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runList(f, apiID, applicationIDs, planIDs, methods, from, to, page, perPage, all, follow)
		},
	}

	cmdutil.AddAPIFlag(cmd, &apiID)
	cmd.Flags().Int64Var(&from, "from", 0, "Start timestamp (epoch millis)")
	cmd.Flags().Int64Var(&to, "to", 0, "End timestamp (epoch millis)")
	cmd.Flags().StringArrayVar(&applicationIDs, "application-ids", nil, "Filter by application IDs")
	cmd.Flags().StringArrayVar(&planIDs, "plan-ids", nil, "Filter by plan IDs")
	cmd.Flags().StringArrayVar(&methods, "methods", nil, "Filter by HTTP methods")
	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&perPage, "per-page", 10, "Results per page")
	cmd.Flags().BoolVar(&all, "all", false, "Fetch all pages")
	cmd.Flags().BoolVar(&follow, "follow", false, "Poll for new logs continuously (Ctrl+C to stop)")

	return cmd
}

func newGetCmd(f *factory.Factory) *cobra.Command {
	var apiID string

	cmd := &cobra.Command{
		Use:     "get <requestId> --api <apiId>",
		Short:   "Get details of a specific request log",
		Example: `  gio apim log get req-aaaa-bbbb-cccc-dddd-eeeeeeee --api /my/api`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			resolvedID, err := f.APIM().ResolveAPI(apiID)
			if err != nil {
				return err
			}

			data, err := f.APIM().GetAPILog(resolvedID, args[0])
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			if printer.IsStructured(f.OutputFormat) {
				return p.PrintDetail(data)
			}

			return printLogDetail(p, data)
		},
	}

	cmdutil.AddAPIFlag(cmd, &apiID)

	return cmd
}

func runList(f *factory.Factory, apiID string, applicationIDs, planIDs, methods []string, from, to int64, page, perPage int, all, follow bool) error {
	if err := cmdutil.RequireContext(f); err != nil {
		return err
	}

	resolvedID, err := f.APIM().ResolveAPI(apiID)
	if err != nil {
		return err
	}

	if follow {
		return runFollow(f, resolvedID, apim.ListAPILogsParams{
			ApplicationIDs: applicationIDs,
			PlanIDs:        planIDs,
			Methods:        methods,
			From:           from,
			Page:           1,
			PerPage:        100,
		})
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	params := func(pg int) apim.ListAPILogsParams {
		return apim.ListAPILogsParams{
			ApplicationIDs: applicationIDs,
			PlanIDs:        planIDs,
			Methods:        methods,
			From:           from,
			To:             to,
			Page:           pg,
			PerPage:        perPage,
		}
	}

	if all {
		var allData []json.RawMessage
		allData, err = apim.FetchAllPages(func(pg int) (*apim.PaginatedResponse, error) {
			return f.APIM().ListAPILogs(resolvedID, params(pg))
		})
		if err != nil {
			return err
		}

		if printer.IsStructured(f.OutputFormat) {
			return p.PrintDetail(allData)
		}

		if err = p.PrintList(allData, logColumns()); err != nil {
			return err
		}

		if len(allData) > 0 {
			p.PrintHint("Showing %d results.", len(allData))
		}

		return nil
	}

	resp, err := f.APIM().ListAPILogs(resolvedID, params(page))
	if err != nil {
		return err
	}

	if printer.IsStructured(f.OutputFormat) {
		return p.PrintDetail(resp)
	}

	if err := p.PrintList(resp.Data, logColumns()); err != nil {
		return err
	}

	pg := resp.Pagination
	cmdutil.PrintPaginationHint(p, pg.Page, pg.PerPage, pg.PageCount, pg.TotalCount, pg.PageItemsCount, all)

	return nil
}

func runFollow(f *factory.Factory, apiID string, baseParams apim.ListAPILogsParams) error {
	seen := make(map[string]struct{})

	w := tabwriter.NewWriter(f.IOStreams.Out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tREQUEST ID\tMETHOD\tSTATUS\tPATH")
	_ = w.Flush()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sig:
			return nil
		case <-ticker.C:
			resp, err := f.APIM().ListAPILogs(apiID, baseParams)
			if err != nil {
				return err
			}

			var newEntries []map[string]any

			for _, raw := range resp.Data {
				var entry map[string]any
				if err := json.Unmarshal(raw, &entry); err != nil {
					continue
				}

				id, _ := entry["requestId"].(string)
				if _, already := seen[id]; already {
					continue
				}

				seen[id] = struct{}{}
				newEntries = append(newEntries, entry)
			}

			for i := len(newEntries) - 1; i >= 0; i-- {
				e := newEntries[i]
				fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%s\n",
					cmdutil.TimestampField(e, "timestamp"),
					cmdutil.StringField(e, "requestId"),
					cmdutil.StringField(e, "method"),
					cmdutil.StringField(e, "status"),
					cmdutil.StringField(e, "uri"),
				)
			}

			_ = w.Flush()
		}
	}
}

func printLogDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, field := range []struct{ label, key string }{
		{"Request ID", "requestId"},
		{"Timestamp", "timestamp"},
		{"Method", "method"},
		{"Path", "uri"},
		{"Status", "status"},
		{"Response Time", "gatewayResponseTime"},
		{"Application", "application"},
		{"Plan", "plan"},
	} {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-17s%v", field.label+":", v)
		}
	}

	return nil
}

func logColumns() []printer.Column {
	return []printer.Column{
		{Name: "Timestamp", Value: func(i any) string { return cmdutil.TimestampField(i, "timestamp") }},
		{Name: "Request ID", Value: func(i any) string { return cmdutil.StringField(i, "requestId") }},
		{Name: "Method", Value: func(i any) string { return cmdutil.StringField(i, "method") }},
		{Name: "Status", Value: func(i any) string { return cmdutil.StringField(i, "status") }},
		{Name: "Path", Value: func(i any) string { return cmdutil.StringField(i, "uri") }},
	}
}
