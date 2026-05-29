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

	alertcmd "gravitee.io/gctl/cmd/am/alert"
	analyticscmd "gravitee.io/gctl/cmd/am/analytics"
	appcmd "gravitee.io/gctl/cmd/am/app"
	auditcmd "gravitee.io/gctl/cmd/am/audit"
	authcmd "gravitee.io/gctl/cmd/am/auth"
	authdevicenotifiercmd "gravitee.io/gctl/cmd/am/auth-device-notifier"
	authorizationenginecmd "gravitee.io/gctl/cmd/am/authorization-engine"
	botdetectioncmd "gravitee.io/gctl/cmd/am/bot-detection"
	certificatecmd "gravitee.io/gctl/cmd/am/certificate"
	dataplanecmd "gravitee.io/gctl/cmd/am/data-plane"
	deviceidentifiercmd "gravitee.io/gctl/cmd/am/device-identifier"
	dictionarycmd "gravitee.io/gctl/cmd/am/dictionary"
	diffcmd "gravitee.io/gctl/cmd/am/diff"
	domaincmd "gravitee.io/gctl/cmd/am/domain"
	emailcmd "gravitee.io/gctl/cmd/am/email"
	entrypointcmd "gravitee.io/gctl/cmd/am/entrypoint"
	extensiongrantcmd "gravitee.io/gctl/cmd/am/extension-grant"
	factorcmd "gravitee.io/gctl/cmd/am/factor"
	flowcmd "gravitee.io/gctl/cmd/am/flow"
	formcmd "gravitee.io/gctl/cmd/am/form"
	groupcmd "gravitee.io/gctl/cmd/am/group"
	idpcmd "gravitee.io/gctl/cmd/am/idp"
	lintcmd "gravitee.io/gctl/cmd/am/lint"
	membercmd "gravitee.io/gctl/cmd/am/member"
	oidctestcmd "gravitee.io/gctl/cmd/am/oidctest"
	orgcmd "gravitee.io/gctl/cmd/am/org"
	passwordpolicycmd "gravitee.io/gctl/cmd/am/password-policy"
	plugincmd "gravitee.io/gctl/cmd/am/plugin"
	protectedresourcecmd "gravitee.io/gctl/cmd/am/protected-resource"
	reportercmd "gravitee.io/gctl/cmd/am/reporter"
	resourcecmd "gravitee.io/gctl/cmd/am/resource"
	rolecmd "gravitee.io/gctl/cmd/am/role"
	scopecmd "gravitee.io/gctl/cmd/am/scope"
	shellcmd "gravitee.io/gctl/cmd/am/shell"
	supportdumpcmd "gravitee.io/gctl/cmd/am/supportdump"
	themecmd "gravitee.io/gctl/cmd/am/theme"
	tokencmd "gravitee.io/gctl/cmd/am/token"
	tracecmd "gravitee.io/gctl/cmd/am/trace"
	usercmd "gravitee.io/gctl/cmd/am/user"
	watchcmd "gravitee.io/gctl/cmd/am/watch"
	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

func newAMBaseCmd(f *factory.Factory) *cobra.Command {
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

	defaultHelp := cmd.HelpFunc()
	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		_ = cmdutil.SetupConfig(f)
		_ = cmdutil.ResolveProductContext(f, "am")
		if header := cmdutil.ContextHeader(f, "am"); header != "" {
			fmt.Fprint(c.OutOrStdout(), header+"\n")
		}

		defaultHelp(c, args)
	})

	cmd.AddCommand(analyticscmd.NewAnalyticsCmd(f))
	cmd.AddCommand(auditcmd.NewAuditCmd(f))

	return cmd
}

// NewAMCmdRO creates the am command with read-only subcommands only.
func NewAMCmdRO(f *factory.Factory) *cobra.Command {
	cmd := newAMBaseCmd(f)

	cmd.AddCommand(alertcmd.NewAlertCmdRO(f))
	cmd.AddCommand(appcmd.NewAppCmdRO(f))
	cmd.AddCommand(authdevicenotifiercmd.NewAuthDeviceNotifierCmdRO(f))
	cmd.AddCommand(authorizationenginecmd.NewAuthorizationEngineCmdRO(f))
	cmd.AddCommand(botdetectioncmd.NewBotDetectionCmdRO(f))
	cmd.AddCommand(certificatecmd.NewCertificateCmdRO(f))
	cmd.AddCommand(dataplanecmd.NewDataPlaneCmd(f))
	cmd.AddCommand(deviceidentifiercmd.NewDeviceIdentifierCmdRO(f))
	cmd.AddCommand(dictionarycmd.NewDictionaryCmdRO(f))
	cmd.AddCommand(domaincmd.NewDomainCmdRO(f))
	cmd.AddCommand(emailcmd.NewEmailCmdRO(f))
	cmd.AddCommand(entrypointcmd.NewEntrypointCmdRO(f))
	cmd.AddCommand(extensiongrantcmd.NewExtensionGrantCmdRO(f))
	cmd.AddCommand(factorcmd.NewFactorCmdRO(f))
	cmd.AddCommand(flowcmd.NewFlowCmdRO(f))
	cmd.AddCommand(formcmd.NewFormCmdRO(f))
	cmd.AddCommand(groupcmd.NewGroupCmdRO(f))
	cmd.AddCommand(idpcmd.NewIDPCmdRO(f))
	cmd.AddCommand(membercmd.NewMemberCmdRO(f))
	// org is excluded: its sub-subcommands mix read and write operations at multiple
	// nesting levels, making a safe RO variant non-trivial without deeper refactoring.
	cmd.AddCommand(passwordpolicycmd.NewPasswordPolicyCmdRO(f))
	cmd.AddCommand(plugincmd.NewPluginCmdRO(f))
	cmd.AddCommand(protectedresourcecmd.NewProtectedResourceCmdRO(f))
	cmd.AddCommand(reportercmd.NewReporterCmdRO(f))
	cmd.AddCommand(resourcecmd.NewResourceCmdRO(f))
	cmd.AddCommand(rolecmd.NewRoleCmdRO(f))
	cmd.AddCommand(scopecmd.NewScopeCmdRO(f))
	cmd.AddCommand(themecmd.NewThemeCmdRO(f))
	cmd.AddCommand(tokencmd.NewTokenCmdRO(f))
	cmd.AddCommand(usercmd.NewUserCmdRO(f))
	cmd.AddCommand(newHealthCmd(f))
	cmd.AddCommand(newStatusCmd(f))
	cmd.AddCommand(newWhoamiCmd(f))

	return cmd
}

// NewAMCmd creates the am parent command with all AM subcommands.
func NewAMCmd(f *factory.Factory) *cobra.Command {
	cmd := newAMBaseCmd(f)

	cmd.AddCommand(alertcmd.NewAlertCmd(f))
	cmd.AddCommand(appcmd.NewAppCmd(f))
	cmd.AddCommand(authcmd.NewAuthCmd(f))
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
