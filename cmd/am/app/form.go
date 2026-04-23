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

package app

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newAppFormCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var appID string

	cmd := &cobra.Command{
		Use:   "form",
		Short: "Manage application form templates",
	}

	cmd.PersistentFlags().StringVar(&appID, "app-id", "", "Application ID (required)")
	_ = cmd.MarkPersistentFlagRequired("app-id")

	cmd.AddCommand(newAppFormGetCmd(f, domainID, &appID))
	cmd.AddCommand(newAppFormCreateCmd(f, domainID, &appID))
	cmd.AddCommand(newAppFormUpdateCmd(f, domainID, &appID))
	cmd.AddCommand(newAppFormDeleteCmd(f, domainID, &appID))

	return cmd
}

func newAppFormGetCmd(f *factory.Factory, domainID, appID *string) *cobra.Command {
	var template string

	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get an application form template",
		Example: `  gio am app form get --domain my-domain --app-id my-app --template LOGIN`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().GetAppForm(*domainID, *appID, template)
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			return p.PrintDetail(data)
		},
	}

	cmd.Flags().StringVar(&template, "template", "", "Form template name (required)")
	_ = cmd.MarkFlagRequired("template")

	return cmd
}

func newAppFormCreateCmd(f *factory.Factory, domainID, appID *string) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:     "create --file <form.json>",
		Short:   "Create an application form template",
		Example: `  gio am app form create --domain my-domain --app-id my-app --file form.json`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body, err := cmdutil.ReadJSONFile(file)
			if err != nil {
				return err
			}

			data, err := f.AM().CreateAppForm(*domainID, *appID, body)
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			if f.OutputFormat != printer.FormatTable {
				return p.PrintDetail(data)
			}

			p.PrintMessage("Form template created successfully.")

			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func newAppFormUpdateCmd(f *factory.Factory, domainID, appID *string) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:     "update <formID> --file <form.json>",
		Short:   "Update an application form template",
		Example: `  gio am app form update form-1 --domain my-domain --app-id my-app --file form.json`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body, err := cmdutil.ReadJSONFile(file)
			if err != nil {
				return err
			}

			data, err := f.AM().UpdateAppForm(*domainID, *appID, args[0], body)
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			if f.OutputFormat != printer.FormatTable {
				return p.PrintDetail(data)
			}

			p.PrintMessage("Form template updated successfully.")

			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func newAppFormDeleteCmd(f *factory.Factory, domainID, appID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <formID>",
		Short:   "Delete an application form template",
		Example: `  gio am app form delete form-1 --domain my-domain --app-id my-app`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().DeleteAppForm(*domainID, *appID, args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Form template '%s' deleted.", args[0])

			return nil
		},
	}
}
