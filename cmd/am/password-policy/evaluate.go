package passwordpolicy

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newEvaluateCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var password string

	cmd := &cobra.Command{
		Use:     "evaluate <policyID> --password <password>",
		Short:   "Evaluate a password against a policy",
		Example: `  gio am password-policy evaluate policy-123 --domain my-domain --password "test123"`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body := map[string]any{"password": password}
			raw, _ := json.Marshal(body)

			data, err := f.AM().EvaluatePasswordPolicy(*domainID, args[0], json.RawMessage(raw))
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

	cmd.Flags().StringVar(&password, "password", "", "Password to evaluate (required)")
	_ = cmd.MarkFlagRequired("password")

	return cmd
}
