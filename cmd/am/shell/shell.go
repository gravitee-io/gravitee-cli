package shell

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func splitArgs(input string) []string {
	var args []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)

	for _, ch := range input {
		switch {
		case inQuote:
			if ch == quoteChar {
				inQuote = false
			} else {
				current.WriteRune(ch)
			}
		case ch == '"' || ch == '\'':
			inQuote = true
			quoteChar = ch
		case unicode.IsSpace(ch):
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(ch)
		}
	}
	if inQuote {
		// unterminated quote: treat remaining as plain token
		if current.Len() > 0 {
			args = append(args, current.String())
		}
	} else if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}

func buildPrompt(workspace, domain string) string {
	if workspace == "" {
		return "[not-configured] am> "
	}
	domainLabel := domain
	if domainLabel == "" {
		domainLabel = "(no domain)"
	}
	return fmt.Sprintf("[%s:%s] am> ", workspace, domainLabel)
}

// NewShellCmd creates the interactive shell command that dispatches to the parent am command tree.
func NewShellCmd(f *factory.Factory, parent *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:     "shell",
		Aliases: []string{"interactive"},
		Short:   "Start an interactive shell session",
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runShell(f, parent)
		},
	}
}

func runShell(f *factory.Factory, parent *cobra.Command) error {
	out := f.IOStreams.Out
	fmt.Fprintln(out, "\nGravitee AM CLI - Interactive Shell")
	fmt.Fprintln(out, "Type commands without the 'am' prefix. Type 'help' for available commands, 'exit' to quit.")
	fmt.Fprintln(out)

	scanner := bufio.NewScanner(f.IOStreams.In)
	return shellLoop(out, scanner, f, parent)
}

func shellLoop(out io.Writer, scanner *bufio.Scanner, f *factory.Factory, parent *cobra.Command) error {
	for {
		workspace := ""
		domain := ""
		if f.Config != nil {
			workspace = f.Config.Current
		}
		if f.Resolved != nil {
			domain = f.Resolved.Domain
		}
		fmt.Fprint(out, buildPrompt(workspace, domain))

		if !scanner.Scan() {
			fmt.Fprintln(out, "\nGoodbye!")
			if err := scanner.Err(); err != nil {
				return err
			}
			return nil
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if processShellCommand(out, line, f, parent) {
			return nil
		}
	}
}

func processShellCommand(out io.Writer, line string, f *factory.Factory, parent *cobra.Command) bool {
	switch line {
	case "exit", "quit":
		fmt.Fprintln(out, "Goodbye!")
		return true
	case "clear":
		fmt.Fprint(out, "\033[2J\033[H")
		return false
	case "help":
		_ = parent.Help()
		return false
	default:
		args := splitArgs(line)
		parent.SetArgs(args)
		if err := parent.Execute(); err != nil {
			fmt.Fprintf(f.IOStreams.Err, "Error: %v\n", err)
		}
		// Reset flag values back to defaults so the next invocation isn't
		// polluted by state from this one. Cobra mutates *Command in place
		// (`-o json` followed by `list` would otherwise stay JSON).
		resetFlags(parent)
		return false
	}
}

// resetFlags walks the command tree and resets all flags' values back to their
// declared defaults. Without this, persistent flags like --output set in one
// shell iteration leak into the next.
func resetFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		flag.Changed = false
		_ = flag.Value.Set(flag.DefValue)
	})
	cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		flag.Changed = false
		_ = flag.Value.Set(flag.DefValue)
	})
	for _, sub := range cmd.Commands() {
		resetFlags(sub)
	}
}
