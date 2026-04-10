package cmdutil

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/factory"
)

// ResolveAPIMFlags rewrites the --api flag (if present and set) with the id
// resolved by the APIM service. Meant to run in PersistentPreRunE so every
// APIM subcommand accepting --api gets path → id resolution for free.
func ResolveAPIMFlags(f *factory.Factory, cmd *cobra.Command) error {
	flag := cmd.Flags().Lookup("api")
	if flag == nil {
		return nil
	}

	val := flag.Value.String()
	if val == "" {
		return nil
	}

	svc := f.APIM()
	if svc == nil {
		return nil
	}

	resolved, err := svc.ResolveAPI(val)
	if err != nil {
		return err
	}

	return flag.Value.Set(resolved)
}
