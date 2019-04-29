package main

import (
	"os"

	"github.com/livingsilver94/stack-updater/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
