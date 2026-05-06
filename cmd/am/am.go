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
	diffcmd "github.com/gravitee-io/gio-cli/cmd/am/diff"
	domaincmd "github.com/gravitee-io/gio-cli/cmd/am/domain"
	emailcmd "github.com/gravitee-io/gio-cli/cmd/am/email"
	entrypointcmd "github.com/gravitee-io/gio-cli/cmd/am/entrypoint"
	extensiongrantcmd "github.com/gravitee-io/gio-cli/cmd/am/extension-grant"
	factorcmd "github.com/gravitee-io/gio-cli/cmd/am/factor"
	flowcmd "github.com/gravitee-io/gio-cli/cmd/am/flow"
	formcmd "github.com/gravitee-io/gio-cli/cmd/am/form"
	groupcmd "github.com/gravitee-io/gio-cli/cmd/am/group"
	idpcmd "github.com/gravitee-io/gio-cli/cmd/am/idp"
	lintcmd "github.com/gravitee-io/gio-cli/cmd/am/lint"
	membercmd "github.com/gravitee-io/gio-cli/cmd/am/member"
	oidctestcmd "github.com/gravitee-io/gio-cli/cmd/am/oidctest"
	orgcmd "github.com/gravitee-io/gio-cli/cmd/am/org"
	passwordpolicycmd "github.com/gravitee-io/gio-cli/cmd/am/password-policy"
	plugincmd "github.com/gravitee-io/gio-cli/cmd/am/plugin"
	protectedresourcecmd "github.com/gravitee-io/gio-cli/cmd/am/protected-resource"
	reportercmd "github.com/gravitee-io/gio-cli/cmd/am/reporter"
	resourcecmd "github.com/gravitee-io/gio-cli/cmd/am/resource"
	rolecmd "github.com/gravitee-io/gio-cli/cmd/am/role"
	scopecmd "github.com/gravitee-io/gio-cli/cmd/am/scope"
	shellcmd "github.com/gravitee-io/gio-cli/cmd/am/shell"
	supportdumpcmd "github.com/gravitee-io/gio-cli/cmd/am/supportdump"
	themecmd "github.com/gravitee-io/gio-cli/cmd/am/theme"
	tokencmd "github.com/gravitee-io/gio-cli/cmd/am/token"
	tracecmd "github.com/gravitee-io/gio-cli/cmd/am/trace"
	usercmd "github.com/gravitee-io/gio-cli/cmd/am/user"
	watchcmd "github.com/gravitee-io/gio-cli/cmd/am/watch"
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
	cmd.AddCommand(diffcmd.NewDiffCmd(f))
	cmd.AddCommand(domaincmd.NewDomainCmd(f))
	cmd.AddCommand(emailcmd.NewEmailCmd(f))
	cmd.AddCommand(entrypointcmd.NewEntrypointCmd(f))
	cmd.AddCommand(extensiongrantcmd.NewExtensionGrantCmd(f))
	cmd.AddCommand(factorcmd.NewFactorCmd(f))
	cmd.AddCommand(flowcmd.NewFlowCmd(f))
	cmd.AddCommand(formcmd.NewFormCmd(f))
	cmd.AddCommand(groupcmd.NewGroupCmd(f))
	cmd.AddCommand(idpcmd.NewIDPCmd(f))
	cmd.AddCommand(lintcmd.NewLintCmd(f))
	cmd.AddCommand(membercmd.NewMemberCmd(f))
	cmd.AddCommand(oidctestcmd.NewTestCmd(f))
	cmd.AddCommand(passwordpolicycmd.NewPasswordPolicyCmd(f))
	cmd.AddCommand(plugincmd.NewPluginCmd(f))
	cmd.AddCommand(protectedresourcecmd.NewProtectedResourceCmd(f))
	cmd.AddCommand(reportercmd.NewReporterCmd(f))
	cmd.AddCommand(resourcecmd.NewResourceCmd(f))
	cmd.AddCommand(rolecmd.NewRoleCmd(f))
	cmd.AddCommand(scopecmd.NewScopeCmd(f))
	cmd.AddCommand(newSetCmd(f))
	cmd.AddCommand(supportdumpcmd.NewSupportDumpCmd(f))
	cmd.AddCommand(themecmd.NewThemeCmd(f))
	cmd.AddCommand(tokencmd.NewTokenCmd(f))
	cmd.AddCommand(tracecmd.NewTraceCmd(f))
	cmd.AddCommand(usercmd.NewUserCmd(f))
	cmd.AddCommand(watchcmd.NewWatchCmd(f))
	cmd.AddCommand(orgcmd.NewOrgCmd(f))
	cmd.AddCommand(newLogoutCmd(f))
	cmd.AddCommand(newWhoamiCmd(f))
	cmd.AddCommand(newStatusCmd(f))
	cmd.AddCommand(newHealthCmd(f))
	cmd.AddCommand(newDoctorCmd(f))
	cmd.AddCommand(shellcmd.NewShellCmd(f, cmd))

	return cmd
}
