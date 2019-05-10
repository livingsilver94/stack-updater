package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd is the parent command. It stores application's
// information and works as a dispatcher for subcommands.
var RootCmd = &cobra.Command{
	Use:     "stack-updater",
	Short:   "A stack update helper for Solus",
	Version: "1.0.0",
}
