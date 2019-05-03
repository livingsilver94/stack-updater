package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "stack-updater",
	Short: "A stack update helper for Solus",
}
