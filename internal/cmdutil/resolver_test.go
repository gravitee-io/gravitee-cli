package cmdutil

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/apim"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newAPIMCmdWithFlag(f *factory.Factory) *cobra.Command {
	var apiID string

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().StringVar(&apiID, "api", "", "")
	cmd.PersistentPreRunE = func(c *cobra.Command, _ []string) error {
		return ResolveAPIMFlags(f, c)
	}
	cmd.RunE = func(_ *cobra.Command, _ []string) error { return nil }

	return cmd
}

func TestResolveAPIMFlags(t *testing.T) {
	t.Run("rewrites --api with resolved id", func(t *testing.T) {
		mock := &apim.MockService{
			ResolveAPIFunc: func(v string) (string, error) {
				if v != "/my/api" {
					return "", fmt.Errorf("unexpected input %q", v)
				}

				return "resolved-id", nil
			},
		}

		f := &factory.Factory{}
		f.SetAPIMService(mock)

		cmd := newAPIMCmdWithFlag(f)
		cmd.SetArgs([]string{"--api", "/my/api"})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := cmd.Flags().Lookup("api").Value.String(); got != "resolved-id" {
			t.Fatalf("expected flag rewritten to resolved-id, got %q", got)
		}
	})

	t.Run("no-op when --api is not set", func(t *testing.T) {
		called := false
		mock := &apim.MockService{
			ResolveAPIFunc: func(string) (string, error) {
				called = true

				return "", nil
			},
		}

		f := &factory.Factory{}
		f.SetAPIMService(mock)

		cmd := newAPIMCmdWithFlag(f)
		cmd.SetArgs([]string{})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if called {
			t.Fatal("expected no resolve call when flag is empty")
		}
	})

	t.Run("no-op when command has no --api flag", func(t *testing.T) {
		called := false
		mock := &apim.MockService{
			ResolveAPIFunc: func(string) (string, error) {
				called = true

				return "", nil
			},
		}

		f := &factory.Factory{}
		f.SetAPIMService(mock)

		cmd := &cobra.Command{Use: "noflag"}
		cmd.RunE = func(_ *cobra.Command, _ []string) error { return nil }
		cmd.PersistentPreRunE = func(c *cobra.Command, _ []string) error {
			return ResolveAPIMFlags(f, c)
		}

		if err := cmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if called {
			t.Fatal("expected no resolve call when flag is absent")
		}
	})

	t.Run("propagates resolver error", func(t *testing.T) {
		mock := &apim.MockService{
			ResolveAPIFunc: func(string) (string, error) {
				return "", fmt.Errorf("boom")
			},
		}

		f := &factory.Factory{}
		f.SetAPIMService(mock)

		cmd := newAPIMCmdWithFlag(f)
		cmd.SetArgs([]string{"--api", "/x"})
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true

		err := cmd.Execute()
		if err == nil || err.Error() != "boom" {
			t.Fatalf("expected boom error, got: %v", err)
		}
	})
}
