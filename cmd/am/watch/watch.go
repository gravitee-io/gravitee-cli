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

package watch

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewWatchCmd(f *factory.Factory) *cobra.Command {
	var intervalSec int

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Live dashboard — monitor logins, errors, and audit events in real-time",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			if intervalSec < 1 {
				return fmt.Errorf("--interval must be >= 1")
			}

			domainName := f.Resolved.Domain

			var lastErr string
			refresh := func() {
				data, err := f.Client.Get(cmdutil.AMDomainPath(f, "audits?page=0&size=50"))
				if err != nil {
					lastErr = fmt.Sprintf("[%s] fetch failed: %v", time.Now().Format("15:04:05"), err)
					fmt.Fprint(f.IOStreams.Out, "\033[2J\033[H")
					fmt.Fprintf(f.IOStreams.Out, "Watch error: %s\n", lastErr)
					return
				}
				var resp struct {
					Data []map[string]interface{} `json:"data"`
				}
				if err := json.Unmarshal(data, &resp); err != nil {
					lastErr = fmt.Sprintf("[%s] parse failed: %v", time.Now().Format("15:04:05"), err)
					fmt.Fprint(f.IOStreams.Out, "\033[2J\033[H")
					fmt.Fprintf(f.IOStreams.Out, "Watch error: %s\n", lastErr)
					return
				}
				lastErr = ""
				dashboard := buildDashboardData(resp.Data, domainName, f.Config.Current)
				fmt.Fprint(f.IOStreams.Out, "\033[2J\033[H")
				fmt.Fprint(f.IOStreams.Out, render(dashboard, intervalSec))
			}

			refresh()

			ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
			defer ticker.Stop()

			sig := make(chan os.Signal, 1)
			signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

			for {
				select {
				case <-ticker.C:
					refresh()
				case <-sig:
					fmt.Fprintln(f.IOStreams.Out, "\nStopped.")
					return nil
				}
			}
		},
	}
	cmd.Flags().IntVar(&intervalSec, "interval", 5, "Refresh interval in seconds")
	return cmd
}
