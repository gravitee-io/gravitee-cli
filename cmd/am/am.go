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

	alertcmd "github.com/gravitee-io/gio-cli/cmd/am/alert"
	analyticscmd "github.com/gravitee-io/gio-cli/cmd/am/analytics"
	appcmd "github.com/gravitee-io/gio-cli/cmd/am/app"
	auditcmd "github.com/gravitee-io/gio-cli/cmd/am/audit"
	authdevicenotifiercmd "github.com/gravitee-io/gio-cli/cmd/am/auth-device-notifier"
	authorizationenginecmd "github.com/gravitee-io/gio-cli/cmd/am/authorization-engine"
	botdetectioncmd "github.com/gravitee-io/gio-cli/cmd/am/bot-detection"
	certificatecmd "github.com/gravitee-io/gio-cli/cmd/am/certificate"
	dataplanecmd "github.com/gravitee-io/gio-cli/cmd/am/data-plane"
	deviceidentifiercmd "github.com/gravitee-io/gio-cli/cmd/am/device-identifier"
	dictionarycmd "github.com/gravitee-io/gio-cli/cmd/am/dictionary"
	domaincmd "github.com/gravitee-io/gio-cli/cmd/am/domain"
	emailcmd "github.com/gravitee-io/gio-cli/cmd/am/email"
	entrypointcmd "github.com/gravitee-io/gio-cli/cmd/am/entrypoint"
	extensiongrantcmd "github.com/gravitee-io/gio-cli/cmd/am/extension-grant"
	factorcmd "github.com/gravitee-io/gio-cli/cmd/am/factor"
	flowcmd "github.com/gravitee-io/gio-cli/cmd/am/flow"
	formcmd "github.com/gravitee-io/gio-cli/cmd/am/form"
	groupcmd "github.com/gravitee-io/gio-cli/cmd/am/group"
	idpcmd "github.com/gravitee-io/gio-cli/cmd/am/idp"
	membercmd "github.com/gravitee-io/gio-cli/cmd/am/member"
	orgcmd "github.com/gravitee-io/gio-cli/cmd/am/org"
	passwordpolicycmd "github.com/gravitee-io/gio-cli/cmd/am/password-policy"
	protectedresourcecmd "github.com/gravitee-io/gio-cli/cmd/am/protected-resource"
	reportercmd "github.com/gravitee-io/gio-cli/cmd/am/reporter"
	resourcecmd "github.com/gravitee-io/gio-cli/cmd/am/resource"
	rolecmd "github.com/gravitee-io/gio-cli/cmd/am/role"
	scopecmd "github.com/gravitee-io/gio-cli/cmd/am/scope"
	themecmd "github.com/gravitee-io/gio-cli/cmd/am/theme"
	usercmd "github.com/gravitee-io/gio-cli/cmd/am/user"
	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewAMCmd creates the am parent command with all AM subcommands.
func NewAMCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "am",
		Short: "Gravitee Access Management",
		Long:  "Manage Gravitee AM resources: domains, applications, users, identity providers, and more.",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.SetupConfig(f); err != nil {
				return err
			}
			return cmdutil.ResolveProductContext(f, "am")
		},
	}

	// Override help to show context info.
	defaultHelp := cmd.HelpFunc()
	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		_ = cmdutil.SetupConfig(f)
		_ = cmdutil.ResolveProductContext(f, "am")
		if header := cmdutil.ContextHeader(f, "am"); header != "" {
			fmt.Fprint(c.OutOrStdout(), header+"\n")
		}

		defaultHelp(c, args)
	})

	cmd.AddCommand(alertcmd.NewAlertCmd(f))
	cmd.AddCommand(analyticscmd.NewAnalyticsCmd(f))
	cmd.AddCommand(appcmd.NewAppCmd(f))
	cmd.AddCommand(auditcmd.NewAuditCmd(f))
	cmd.AddCommand(authdevicenotifiercmd.NewAuthDeviceNotifierCmd(f))
	cmd.AddCommand(authorizationenginecmd.NewAuthorizationEngineCmd(f))
	cmd.AddCommand(botdetectioncmd.NewBotDetectionCmd(f))
	cmd.AddCommand(certificatecmd.NewCertificateCmd(f))
	cmd.AddCommand(dataplanecmd.NewDataPlaneCmd(f))
	cmd.AddCommand(deviceidentifiercmd.NewDeviceIdentifierCmd(f))
	cmd.AddCommand(dictionarycmd.NewDictionaryCmd(f))
	cmd.AddCommand(domaincmd.NewDomainCmd(f))
	cmd.AddCommand(emailcmd.NewEmailCmd(f))
	cmd.AddCommand(entrypointcmd.NewEntrypointCmd(f))
	cmd.AddCommand(extensiongrantcmd.NewExtensionGrantCmd(f))
	cmd.AddCommand(factorcmd.NewFactorCmd(f))
	cmd.AddCommand(flowcmd.NewFlowCmd(f))
	cmd.AddCommand(formcmd.NewFormCmd(f))
	cmd.AddCommand(groupcmd.NewGroupCmd(f))
	cmd.AddCommand(idpcmd.NewIDPCmd(f))
	cmd.AddCommand(membercmd.NewMemberCmd(f))
	cmd.AddCommand(passwordpolicycmd.NewPasswordPolicyCmd(f))
	cmd.AddCommand(protectedresourcecmd.NewProtectedResourceCmd(f))
	cmd.AddCommand(reportercmd.NewReporterCmd(f))
	cmd.AddCommand(resourcecmd.NewResourceCmd(f))
	cmd.AddCommand(rolecmd.NewRoleCmd(f))
	cmd.AddCommand(scopecmd.NewScopeCmd(f))
	cmd.AddCommand(themecmd.NewThemeCmd(f))
	cmd.AddCommand(usercmd.NewUserCmd(f))
	cmd.AddCommand(orgcmd.NewOrgCmd(f))
	cmd.AddCommand(newLogoutCmd(f))

	return cmd
}
