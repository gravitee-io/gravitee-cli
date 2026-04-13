package app

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newSecretCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var appID string

	cmd := &cobra.Command{
		Use:   "secret",
		Short: "Manage application secrets",
	}

	cmd.PersistentFlags().StringVar(&appID, "app-id", "", "Application ID (required)")
	_ = cmd.MarkPersistentFlagRequired("app-id")

	cmd.AddCommand(newSecretListCmd(f, domainID, &appID))
	cmd.AddCommand(newSecretCreateCmd(f, domainID, &appID))
	cmd.AddCommand(newSecretDeleteCmd(f, domainID, &appID))

	return cmd
}

func newSecretListCmd(f *factory.Factory, domainID, appID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List application secrets",
		Example: `  gio am app secret list --domain my-domain --app-id my-app
  gio am app secret list --domain my-domain --app-id my-app -o json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().ListAppSecrets(*domainID, *appID)
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
}

func newSecretCreateCmd(f *factory.Factory, domainID, appID *string) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create an application secret",
		Example: `  gio am app secret create --domain my-domain --app-id my-app --name "my secret"`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body, _ := json.Marshal(map[string]any{"name": name})

			data, err := f.AM().CreateAppSecret(*domainID, *appID, json.RawMessage(body))
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

			p.PrintMessage("Secret created successfully.")

			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Secret name (required)")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func newSecretDeleteCmd(f *factory.Factory, domainID, appID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <secretID>",
		Short:   "Delete an application secret",
		Example: `  gio am app secret delete my-secret-id --domain my-domain --app-id my-app`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().DeleteAppSecret(*domainID, *appID, args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Secret '%s' deleted.", args[0])

			return nil
		},
	}
}
