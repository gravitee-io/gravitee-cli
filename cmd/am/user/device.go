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

package user

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newDeviceCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var userID string

	cmd := &cobra.Command{
		Use:   "device",
		Short: "Manage user devices",
	}

	cmd.PersistentFlags().StringVar(&userID, "user-id", "", "User ID (required)")
	_ = cmd.MarkPersistentFlagRequired("user-id")

	cmd.AddCommand(newDeviceListCmd(f, domainID, &userID))
	cmd.AddCommand(newDeviceDeleteCmd(f, domainID, &userID))

	return cmd
}

func newDeviceListCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List user devices",
		Example: `  gio am user device list --domain my-domain --user-id user-1
  gio am user device list --domain my-domain --user-id user-1 -o json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runDeviceList(f, *domainID, *userID)
		},
	}
}

func runDeviceList(f *factory.Factory, domainID, userID string) error {
	items, err := f.AM().ListUserDevices(domainID, userID)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(items)
	}

	return p.PrintList(items, deviceColumns())
}

func newDeviceDeleteCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <deviceID>",
		Short:   "Delete a user device",
		Example: `  gio am user device delete device-1 --domain my-domain --user-id user-1`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().DeleteUserDevice(*domainID, *userID, args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Device '%s' deleted.", args[0])

			return nil
		},
	}
}

func deviceColumns() []printer.Column {
	return []printer.Column{
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Device ID", Value: func(i any) string { return cmdutil.StringField(i, "deviceId") }},
		{Name: "Type", Value: func(i any) string { return cmdutil.StringField(i, "type") }},
		{Name: "Client", Value: func(i any) string { return cmdutil.StringField(i, "client") }},
	}
}
