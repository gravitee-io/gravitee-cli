package supportdump

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

// NewSupportDumpCmd creates the support-dump command for collecting diagnostic information.
func NewSupportDumpCmd(f *factory.Factory) *cobra.Command {
	var outputFile string
	var allDomains bool
	var noAudit bool
	var auditSize int
	var includeUsers bool
	var noRedact bool

	cmd := &cobra.Command{
		Use:   "support-dump",
		Short: "Generate a comprehensive support diagnostic dump",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runSupportDump(f, outputFile, allDomains, noAudit, auditSize, includeUsers, noRedact)
		},
	}
	cmd.Flags().StringVarP(&outputFile, "file", "f", "", "Output file path (default: stdout)")
	cmd.Flags().BoolVar(&allDomains, "all-domains", false, "Dump all domains in the environment")
	cmd.Flags().BoolVar(&noAudit, "no-audit", false, "Skip audit logs")
	cmd.Flags().IntVar(&auditSize, "audit-size", 100, "Number of recent audit events to include")
	cmd.Flags().BoolVar(&includeUsers, "include-users", false, "Include user list (contains PII)")
	cmd.Flags().BoolVar(&noRedact, "no-redact", false, "Disable secret redaction")
	return cmd
}

func runSupportDump(f *factory.Factory, outputFile string, allDomains, noAudit bool, auditSize int, includeUsers, noRedact bool) error {
	if err := cmdutil.RequireAMContext(f); err != nil {
		return err
	}
	shouldRedact := !noRedact

	domainIDs, err := resolveDomainIDs(f, allDomains)
	if err != nil {
		return err
	}

	output := buildDumpOutput(f, domainIDs, shouldRedact, includeUsers, noAudit)

	if allDomains {
		domains := make([]interface{}, 0, len(domainIDs))
		for _, domainID := range domainIDs {
			sections, errs := collectDomain(f, domainID, !noAudit, auditSize, includeUsers)
			entry := map[string]interface{}{"domainId": domainID}
			for k, v := range sections {
				entry[k] = v
			}
			if len(errs) > 0 {
				entry["_errors"] = errs
			}
			domains = append(domains, entry)
		}
		output["domains"] = domains
	} else {
		sections, errs := collectDomain(f, domainIDs[0], !noAudit, auditSize, includeUsers)
		for k, v := range sections {
			output[k] = v
		}
		if len(errs) > 0 {
			output["_errors"] = errs
		}
	}

	if shouldRedact {
		redacted, ok := redactSecrets(output).(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to redact secrets")
		}
		output = redacted
	}

	return writeDumpOutput(f, outputFile, output)
}

func resolveDomainIDs(f *factory.Factory, allDomains bool) ([]string, error) {
	if allDomains {
		data, err := f.Client.Get(cmdutil.AMEnvPath(f, "domains?page=0&size=1000"))
		if err != nil {
			return nil, err
		}
		var resp struct {
			Data []map[string]interface{} `json:"data"`
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, err
		}
		var ids []string
		for _, d := range resp.Data {
			if id, ok := d["id"].(string); ok {
				ids = append(ids, id)
			}
		}
		return ids, nil
	}

	if err := cmdutil.RequireAMDomain(f); err != nil {
		return nil, err
	}
	return []string{f.Resolved.Domain}, nil
}

func buildDumpOutput(f *factory.Factory, domainIDs []string, shouldRedact bool, includeUsers bool, noAudit bool) map[string]interface{} {
	return map[string]interface{}{
		"_metadata": map[string]interface{}{
			"serverUrl":       f.Resolved.URL,
			"organizationId":  f.Resolved.Org,
			"environmentId":   f.Resolved.Env,
			"secretsRedacted": shouldRedact,
			"includesUsers":   includeUsers,
			"includesAudit":   !noAudit,
			"domainCount":     len(domainIDs),
		},
	}
}

func writeDumpOutput(f *factory.Factory, outputFile string, output map[string]interface{}) error {
	jsonBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	if outputFile != "" {
		if err := os.WriteFile(outputFile, jsonBytes, 0600); err != nil {
			return err
		}
		fmt.Fprintf(f.IOStreams.Out, "Support dump written to %s\n", outputFile)
	} else {
		fmt.Fprintln(f.IOStreams.Out, string(jsonBytes))
	}
	return nil
}

func collectDomain(f *factory.Factory, domainID string, includeAudit bool, auditSize int, includeUsers bool) (map[string]interface{}, []string) {
	sections := make(map[string]interface{})
	var errs []string

	get := func(label, path string, dest *interface{}) {
		data, err := f.Client.Get(path)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", label, err))
			return
		}
		var v interface{}
		if err := json.Unmarshal(data, &v); err != nil {
			errs = append(errs, fmt.Sprintf("%s: parse error", label))
			return
		}
		*dest = v
	}

	domainPath := func(p string) string {
		return cmdutil.AMDomainPathFor(f, domainID, p)
	}

	var domain, apps, idps, certs, flows, factors, roles, scopes, groups, members, users, audits interface{}
	get("domain", domainPath(""), &domain)
	get("applications", domainPath("applications?page=0&size=1000"), &apps)
	get("identities", domainPath("identities"), &idps)
	get("certificates", domainPath("certificates"), &certs)
	get("flows", domainPath("flows"), &flows)
	get("factors", domainPath("factors"), &factors)
	get("roles", domainPath("roles?page=0&size=1000"), &roles)
	get("scopes", domainPath("scopes?page=0&size=1000"), &scopes)
	get("groups", domainPath("groups?page=0&size=1000"), &groups)
	get("members", domainPath("members"), &members)

	if includeUsers {
		get("users", domainPath("users?page=0&size=1000"), &users)
	}
	if includeAudit {
		get("audits", domainPath(fmt.Sprintf("audits?page=0&size=%d", auditSize)), &audits)
	}

	setIfNotNil := func(key string, val interface{}) {
		if val != nil {
			sections[key] = val
		}
	}
	setIfNotNil("domain", domain)
	setIfNotNil("applications", extractData(apps))
	setIfNotNil("identityProviders", idps)
	setIfNotNil("certificates", certs)
	setIfNotNil("flows", flows)
	setIfNotNil("factors", factors)
	setIfNotNil("roles", extractData(roles))
	setIfNotNil("scopes", extractData(scopes))
	setIfNotNil("groups", extractData(groups))
	setIfNotNil("members", members)
	if includeUsers {
		setIfNotNil("users", extractData(users))
	}
	if includeAudit {
		setIfNotNil("recentAudits", extractData(audits))
	}

	return sections, errs
}

func extractData(v interface{}) interface{} {
	if m, ok := v.(map[string]interface{}); ok {
		if data, ok := m["data"]; ok {
			return data
		}
	}
	return v
}
