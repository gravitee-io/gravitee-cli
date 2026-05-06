package lint

import (
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

type lintContext struct {
	apps    []map[string]interface{}
	certs   []map[string]interface{}
	factors []map[string]interface{}
	scopes  []map[string]interface{}
}

// ruleCount returns how many rules runAllRules executes — kept in sync by reading the slice once.
const ruleCount = 14

func NewLintCmd(f *factory.Factory) *cobra.Command {
	var ci bool

	cmd := &cobra.Command{
		Use:   "lint",
		Short: fmt.Sprintf("Run security audit rules against the current domain (%d rules, scored 0-10)", ruleCount),
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			ctx, err := collectLintData(f)
			if err != nil {
				return err
			}
			findings := runAllRules(ctx)
			score := calculateScore(findings)

			out := f.IOStreams.Out
			if len(findings) == 0 {
				fmt.Fprintf(out, "No findings. Score: 10/10\n")
				return nil
			}
			for _, finding := range findings {
				fmt.Fprintf(out, "  [%-8s] %-30s %-30s %s\n",
					finding.Severity, finding.Rule, finding.Resource, finding.Message)
			}
			fmt.Fprintf(out, "\nScore: %d/10 (%d critical, %d warning)\n",
				score, countBySeverity(findings, "critical"), countBySeverity(findings, "warning"))

			if ci && countBySeverity(findings, "critical") > 0 {
				return fmt.Errorf("lint failed: critical findings present")
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&ci, "ci", false, "Exit with code 1 if any critical findings are found")
	return cmd
}

// collectLintData fetches data needed by lint rules. The applications fetch is
// fatal (rules can't run without it). Other fetches degrade gracefully but the
// user gets a stderr warning so they know findings may be incomplete.
func collectLintData(f *factory.Factory) (lintContext, error) {
	var ctx lintContext
	appsData, err := f.Client.Get(cmdutil.AMDomainPath(f, "applications?page=0&size=1000"))
	if err != nil {
		return ctx, err
	}
	var appsResp struct {
		Data []map[string]interface{} `json:"data"`
	}
	if parseErr := json.Unmarshal(appsData, &appsResp); parseErr != nil {
		return ctx, parseErr
	}
	ctx.apps = appsResp.Data

	warn := func(section string, err error) {
		fmt.Fprintf(f.IOStreams.Err, "warning: failed to load %s: %v — related rules may report no findings\n", section, err)
	}

	if data, certErr := f.Client.Get(cmdutil.AMDomainPath(f, "certificates")); certErr != nil {
		warn("certificates", certErr)
	} else if parseErr := json.Unmarshal(data, &ctx.certs); parseErr != nil {
		warn("certificates", parseErr)
	}

	if data, factorErr := f.Client.Get(cmdutil.AMDomainPath(f, "factors")); factorErr != nil {
		warn("factors", factorErr)
	} else if parseErr := json.Unmarshal(data, &ctx.factors); parseErr != nil {
		warn("factors", parseErr)
	}

	if data, scopeErr := f.Client.Get(cmdutil.AMDomainPath(f, "scopes?page=0&size=1000")); scopeErr != nil {
		warn("scopes", scopeErr)
	} else {
		var scopesResp struct {
			Data []map[string]interface{} `json:"data"`
		}
		if parseErr := json.Unmarshal(data, &scopesResp); parseErr != nil {
			warn("scopes", parseErr)
		} else {
			ctx.scopes = scopesResp.Data
		}
	}
	return ctx, nil
}

func runAllRules(ctx lintContext) []LintFinding {
	var all []LintFinding
	all = append(all, ruleImplicitGrant(ctx.apps)...)
	all = append(all, ruleNoPkce(ctx.apps)...)
	all = append(all, ruleLongTokenLifetime(ctx.apps)...)
	all = append(all, ruleLongRefreshLifetime(ctx.apps)...)
	all = append(all, ruleNoIdp(ctx.apps)...)
	all = append(all, ruleLocalhostRedirect(ctx.apps)...)
	all = append(all, ruleHttpRedirect(ctx.apps)...)
	all = append(all, ruleWildcardRedirect(ctx.apps)...)
	all = append(all, ruleAppDisabled(ctx.apps)...)
	all = append(all, ruleNoFactors(ctx.apps, ctx.factors)...)
	all = append(all, ruleCertExpiry(ctx.certs)...)
	all = append(all, ruleUnusedScope(ctx.apps, ctx.scopes)...)
	all = append(all, rulePasswordGrantNoMfa(ctx.apps, ctx.factors)...)
	all = append(all, ruleEmptyDomain(ctx.apps)...)
	return all
}

func countBySeverity(findings []LintFinding, severity string) int {
	count := 0
	for _, f := range findings {
		if f.Severity == severity {
			count++
		}
	}
	return count
}
