/**
 * main.go
 * gho - Ghost Admin API CLI
 *
 * Ghost Admin API CLI tool with a user experience similar to gog-cli
 */

package main

import (
	"fmt"
	"os"

	"github.com/mtane0412/ghocli/internal/cmd"
)

var (
	// Version information (set via -ldflags at build time)
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Call Execute function to run the command
	err := cmd.Execute(os.Args, cmd.ExecuteOptions{
		Version: buildVersion(),
	})

	// If error exists, output to stderr and exit with appropriate exit code
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(cmd.ExitCode(err))
	}
}

// buildVersion constructs the version string
func buildVersion() string {
	if version == "dev" {
		return "gho dev (commit: " + commit + ", built at: " + date + ")"
	}
	return "gho " + version
}
