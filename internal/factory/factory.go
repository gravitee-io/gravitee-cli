package factory

import (
	"io"
	"os"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
)

// IOStreams abstracts standard I/O for testability.
type IOStreams struct {
	Out io.Writer
	Err io.Writer
	In  io.Reader
}

// DefaultIOStreams returns IOStreams connected to os stdin/stdout/stderr.
func DefaultIOStreams() IOStreams {
	return IOStreams{
		Out: os.Stdout,
		Err: os.Stderr,
		In:  os.Stdin,
	}
}

// Factory is the central dependency injection container passed to all commands.
type Factory struct {
	Config       *config.Config
	Resolved     *config.ResolvedContext
	Client       client.GraviteeClient
	IOStreams    IOStreams
	ConfigPath   string
	OutputFormat string
	Quiet        bool
}
