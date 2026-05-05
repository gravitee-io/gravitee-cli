package domain

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newImportCmd(f *factory.Factory) *cobra.Command {
	var targetDomainID string
	cmd := &cobra.Command{
		Use:   "import <file>",
		Short: "Import domain configuration from a JSON export file",
		Example: `  gio am domain import domain-export.json
  gio am domain import domain-export.json --target existing-domain-id`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireAMContext(f); err != nil {
				return err
			}
			return runImport(f, args[0], targetDomainID)
		},
	}
	cmd.Flags().StringVar(&targetDomainID, "target", "", "Target domain ID (creates new domain if not set)")
	return cmd
}

func runImport(f *factory.Factory, file, targetDomainID string) error {
	raw, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var exportData map[string]json.RawMessage
	if parseErr := json.Unmarshal(raw, &exportData); parseErr != nil {
		return fmt.Errorf("failed to parse export file: %w", parseErr)
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if targetDomainID == "" {
		var domainObj map[string]interface{}
		if err := json.Unmarshal(exportData["domain"], &domainObj); err != nil {
			return fmt.Errorf("failed to parse domain in export: %w", err)
		}
		body, err := json.Marshal(map[string]interface{}{
			"name":        cmdutil.StringField(domainObj, "name"),
			"description": cmdutil.StringField(domainObj, "description"),
		})
		if err != nil {
			return err
		}
		created, err := f.Client.Post(cmdutil.AMEnvPath(f, "domains"), body)
		if err != nil {
			return fmt.Errorf("failed to create domain: %w", err)
		}
		var newDomain map[string]interface{}
		if err := json.Unmarshal(created, &newDomain); err != nil {
			return fmt.Errorf("failed to parse CreateDomain response: %w", err)
		}
		targetDomainID = cmdutil.StringField(newDomain, "id")
		if targetDomainID == "" {
			return fmt.Errorf("CreateDomain response did not include an ID")
		}
		p.PrintMessage("Created domain '%s'.", targetDomainID)
	}

	imported, skipped := 0, 0
	add := func(i, s int) { imported += i; skipped += s }

	add(importItems(f, exportData, "scopes", targetDomainID, "scopes"))
	add(importItems(f, exportData, "roles", targetDomainID, "roles"))
	add(importItems(f, exportData, "groups", targetDomainID, "groups"))
	add(importItems(f, exportData, "applications", targetDomainID, "applications"))

	p.PrintMessage("Import complete: %d imported, %d skipped.", imported, skipped)
	return nil
}

// importItems creates resources from a JSON array in exportData. Returns (imported, skipped) counts.
func importItems(f *factory.Factory, exportData map[string]json.RawMessage, key, domainID, resource string) (int, int) {
	raw, ok := exportData[key]
	if !ok || len(raw) == 0 {
		return 0, 0
	}
	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil {
		return 0, 1
	}
	imported := 0
	for _, item := range items {
		path := cmdutil.AMDomainPathFor(f, domainID, resource)
		if _, err := f.Client.Post(path, item); err == nil {
			imported++
		}
	}
	return imported, 0
}
