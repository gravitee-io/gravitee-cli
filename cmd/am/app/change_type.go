package app

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newChangeTypeCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var (
		appID   string
		appType string
	)

	cmd := &cobra.Command{
		Use:   "change-type",
		Short: "Change the type of an application",
		Example: `  gio am app change-type --domain my-domain --app-id my-app --type browser
  gio am app change-type --domain my-domain --app-id my-app --type web`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.ValidateEnum(appType, "type", []string{"web", "native", "browser", "service", "resource_server"}); err != nil {
				return err
			}

			body, _ := json.Marshal(map[string]any{
				"type": strings.ToUpper(appType),
			})

			data, err := f.AM().ChangeAppType(*domainID, appID, json.RawMessage(body))
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

			p.PrintMessage(fmt.Sprintf("Application type changed to '%s'.", strings.ToUpper(appType)))

			return nil
		},
	}

	cmd.Flags().StringVar(&appID, "app-id", "", "Application ID (required)")
	cmd.Flags().StringVar(&appType, "type", "", "New application type: web, native, browser, service, resource_server (required)")
	_ = cmd.MarkFlagRequired("app-id")
	_ = cmd.MarkFlagRequired("type")

	return cmd
}
