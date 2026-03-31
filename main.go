package main

import (
	"os"

	"github.com/gravitee-io/gio-cli/cmd"
)

// version is set at build time via -ldflags.
var version = "dev"

func main() {
	os.Exit(cmd.Execute(version))
}
