package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// SupportedStacks is a list of stacks this application can handle
var SupportedStacks = [...]string{"kde"}

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
		return fmt.Errorf("Usage: %s", cmd.Use)
	}
	stack := strings.Split(strings.ToLower(args[0]), ":")
	if !isValidStack(stack[0]) {
		return fmt.Errorf("%s is not a supported stack. Choose any from %s", stack[0], SupportedStacks)
	}
	if len(stack) > 1 {
		// We also have a bundle to sanitize
		if stack[1] == "" {
			return fmt.Errorf("You should not use \":\" if you don't mean to specify a bundle")
		}
	}
	return nil
}

func isValidStack(stack string) bool {
	for _, validBundle := range SupportedStacks {
		if stack == validBundle {
			return true
		}
	}
	return false
}

func updateStack(cmd *cobra.Command, args []string) {

}
