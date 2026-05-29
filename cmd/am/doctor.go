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

package am

import (
	"fmt"

	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

type doctorCheck struct {
	label  string
	status string // "OK", "WARN", "FAIL"
	detail string
}

func newDoctorCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Run diagnostic checks on the CLI configuration",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			checks := runDoctorChecks(f)
			out := f.IOStreams.Out

			for _, c := range checks {
				fmt.Fprintf(out, "  [%-4s] %-20s %s\n", c.status, c.label, c.detail)
			}

			return nil
		},
	}
}

func runDoctorChecks(f *factory.Factory) []doctorCheck {
	var checks []doctorCheck

	// 1. Config
	if f.Config == nil || len(f.Config.Contexts) == 0 {
		checks = append(checks, doctorCheck{"config", "FAIL", "No contexts configured — run 'gctl login am'"})
		return checks
	}

	checks = append(checks, doctorCheck{"config", "OK", fmt.Sprintf("%d context(s) found", len(f.Config.Contexts))})

	// 2. Current context
	if f.Config.Current == "" {
		checks = append(checks, doctorCheck{"context", "WARN", "No current context set"})
		return checks
	}

	ctx, ok := f.Config.Contexts[f.Config.Current]
	if !ok {
		checks = append(checks, doctorCheck{"context", "FAIL", fmt.Sprintf("Context %q not found", f.Config.Current)})
		return checks
	}

	amURL := ""
	if ctx.AM != nil {
		amURL = ctx.AM.URL
	}

	checks = append(checks, doctorCheck{"context", "OK", fmt.Sprintf("%s @ %s", f.Config.Current, amURL)})

	// 3. Token
	if ctx.AM == nil || ctx.AM.Token == "" {
		checks = append(checks, doctorCheck{"auth", "FAIL", "No AM token stored — run 'gctl login am'"})
		return checks
	}

	checks = append(checks, doctorCheck{"auth", "OK", "Token present"})

	// 4. Domain
	domain := ""
	if f.Resolved != nil {
		domain = f.Resolved.Domain
	}

	if domain == "" {
		checks = append(checks, doctorCheck{"domain", "WARN", "No domain set — run 'gctl am set domain <id>'"})
	} else {
		checks = append(checks, doctorCheck{"domain", "OK", domain})
	}

	// 5. Connectivity
	if err := cmdutil.RequireAMContext(f); err != nil {
		checks = append(checks, doctorCheck{"connect", "WARN", "Skipped — AM context not active"})
	} else if f.Client == nil {
		checks = append(checks, doctorCheck{"connect", "FAIL", "No HTTP client available"})
	} else if _, err := f.Client.Get("/management/user"); err != nil {
		checks = append(checks, doctorCheck{"connect", "FAIL", fmt.Sprintf("Cannot reach AM: %v", err)})
	} else {
		checks = append(checks, doctorCheck{"connect", "OK", "AM management API reachable"})
	}

	return checks
}
