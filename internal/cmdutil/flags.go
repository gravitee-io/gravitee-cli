package cmdutil

import (
	"github.com/spf13/cobra"
)

// AddAPIFlag registers the standard --api flag on cmd, bound to target,
// and marks it required. All APIM subcommands that target a specific API
// should use this so the flag's behavior (description, shortcut, validation)
// lives in one place.
func AddAPIFlag(cmd *cobra.Command, target *string) {
	cmd.Flags().StringVar(target, "api", "",
		"API id or context path (e.g. /my/api) (required)")
	_ = cmd.MarkFlagRequired("api")

	// MarkFlagRequired checks presence, not content: --api "" would pass.
	// Reject it explicitly before RunE so callers never hit the server with
	// an empty id (which silently resolves to wrong endpoints - see B3).
	prev := cmd.PreRunE
	cmd.PreRunE = func(c *cobra.Command, args []string) error {
		if prev != nil {
			if err := prev(c, args); err != nil {
				return err
			}
		}

		return RequireNonEmpty("--api", *target)
	}
}
