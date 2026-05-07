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

package trace

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewTraceCmd(f *factory.Factory) *cobra.Command {
	var userArg, appArg string

	cmd := &cobra.Command{
		Use:   "trace",
		Short: "Trace the authentication path for a user and application",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			return runTrace(f, userArg, appArg)
		},
	}
	cmd.Flags().StringVar(&userArg, "user", "", "User to trace (username, email, or ID) (required)")
	cmd.Flags().StringVar(&appArg, "app", "", "Application (name, clientId, or ID) (required)")
	_ = cmd.MarkFlagRequired("user")
	_ = cmd.MarkFlagRequired("app")
	return cmd
}

//nolint:funlen
func runTrace(f *factory.Factory, userArg, appArg string) error {
	user, err := resolveUser(f, userArg)
	if err != nil {
		return fmt.Errorf("user not found %q: %w", userArg, err)
	}
	app, err := resolveApp(f, appArg)
	if err != nil {
		return fmt.Errorf("application not found %q: %w", appArg, err)
	}

	type result struct {
		data []map[string]interface{}
		err  error
	}
	idpCh := make(chan result, 1)
	factorCh := make(chan result, 1)
	flowCh := make(chan result, 1)

	doFetch := func(path string, ch chan<- result) {
		go func() {
			data, err := f.Client.Get(cmdutil.AMDomainPath(f, path))
			if err != nil {
				ch <- result{nil, err}
				return
			}
			var items []map[string]interface{}
			if err := json.Unmarshal(data, &items); err != nil {
				ch <- result{nil, fmt.Errorf("decode %s: %w", path, err)}
				return
			}
			ch <- result{items, nil}
		}()
	}
	doFetch("identities", idpCh)
	doFetch("factors", factorCh)
	doFetch("flows", flowCh)

	idpResult := <-idpCh
	factorResult := <-factorCh
	flowResult := <-flowCh

	// Log non-fatal fetch errors as warnings (don't block the trace)
	out := f.IOStreams.Out
	if idpResult.err != nil {
		fmt.Fprintf(f.IOStreams.Err, "warning: could not fetch identity providers: %v\n", idpResult.err)
	}
	if factorResult.err != nil {
		fmt.Fprintf(f.IOStreams.Err, "warning: could not fetch factors: %v\n", factorResult.err)
	}
	if flowResult.err != nil {
		fmt.Fprintf(f.IOStreams.Err, "warning: could not fetch flows: %v\n", flowResult.err)
	}

	idps := idpResult.data
	factors := factorResult.data
	flows := flowResult.data

	steps := []TraceStep{
		checkUserStatus(user),
		checkIdpMatch(user, app, idps),
		checkGrantTypes(app),
		checkMfa(user, factors),
		checkFlows(flows),
		checkConsent(app),
		checkTokenConfig(app),
	}
	verdict := buildVerdict(steps)

	userLabel := cmdutil.StringField(user, "email")
	if userLabel == "" {
		userLabel = cmdutil.StringField(user, "username")
	}
	appLabel := cmdutil.StringField(app, "name")
	fmt.Fprintf(out, "\nAuth flow trace: %s -> %s\n\n", userLabel, appLabel)
	for _, step := range steps {
		tag := "[OK]  "
		if step.Status == "warn" {
			tag = "[WARN]"
		} else if step.Status == "block" {
			tag = "[FAIL]"
		}
		fmt.Fprintf(out, "  %s %-16s %s\n", tag, step.Label, step.Detail)
	}
	fmt.Fprintln(out)
	if verdict.CanAuthenticate {
		fmt.Fprintf(out, "  Verdict: [OK]   %s\n\n", verdict.Reason)
	} else {
		fmt.Fprintf(out, "  Verdict: [FAIL] %s\n\n", verdict.Reason)
	}
	return nil
}

func resolveUser(f *factory.Factory, userArg string) (map[string]interface{}, error) {
	data, err := f.Client.Get(cmdutil.AMDomainPath(f, "users/"+userArg))
	if err == nil {
		var user map[string]interface{}
		if json.Unmarshal(data, &user) == nil && user["id"] != nil {
			return user, nil
		}
	}
	data, err = f.Client.Get(cmdutil.AMDomainPath(f, "users?q="+url.QueryEscape(userArg)+"&page=0&size=1"))
	if err != nil {
		return nil, err
	}
	var resp struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil || len(resp.Data) == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return resp.Data[0], nil
}

func resolveApp(f *factory.Factory, appArg string) (map[string]interface{}, error) {
	data, err := f.Client.Get(cmdutil.AMDomainPath(f, "applications/"+appArg))
	if err == nil {
		var app map[string]interface{}
		if json.Unmarshal(data, &app) == nil && app["id"] != nil {
			return app, nil
		}
	}
	data, err = f.Client.Get(cmdutil.AMDomainPath(f, "applications?q="+url.QueryEscape(appArg)+"&page=0&size=1"))
	if err != nil {
		return nil, err
	}
	var resp struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil || len(resp.Data) == 0 {
		return nil, fmt.Errorf("application not found")
	}
	return resp.Data[0], nil
}
