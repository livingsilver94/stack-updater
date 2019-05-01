package cmd

import (
	"fmt"
	"strings"

	"github.com/livingsilver94/stack-updater/stack"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <reponame>[:bundle] <version>",
	Short: "Update a stack",
	Long: `Download and update package definition files beloning to the selected stack.
Note: the ":bundle" part of a repository name is needed when a stack is split in multiple bundles, e.g. KDE.`,
	Args: checkUpdateArgs,
	Run:  updateStack,
}

func init() {
	RootCmd.AddCommand(updateCmd)
}

func checkUpdateArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Not enough arguments provided")
	}
	stackArgs := strings.Split(strings.ToLower(args[0]), ":")
	if _, err := stack.SupportedStackString(stackArgs[0]); err != nil {
		return fmt.Errorf("%s is not a supported stack. Choose any from %s", stackArgs[0], stack.SupportedStackStrings())
	}
	if len(stackArgs) > 1 {
		// We also have a bundle to sanitize
		if stackArgs[1] == "" {
			return fmt.Errorf("You should not use \":\" if you don't mean to specify a bundle")
		}
	}
	return nil
}

func updateStack(cmd *cobra.Command, args []string) {

}
