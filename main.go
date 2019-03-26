package main

import (
	"fmt"
	"github.com/livingsilver94/stack_updater/repository"
	"github.com/livingsilver94/stack_updater/stack"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
)

var (
	app       = kingpin.New("Stack Updater", "Simplify software stack updates for Solus")
	stackname = kingpin.Arg("stack", "The stack you want to update").Required().String()
	version   = kingpin.Arg("version", "Version of the stack/bundle you want to update at").Required().String()
	bundle    = kingpin.Arg("bundle", "KDE bundle to update").Default("kde").String()
	dryRun    = kingpin.Flag("dry-run", "List what will be updated without touching any file").Short('d').Bool()

	parser stack.Parser
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
	switch strings.ToLower(*stackname) {
	case "kde":
		{
			parser = stack.KDE{*bundle, *version}
		}
	case "mate":
		{
			fmt.Println("To implement")
		}
	default:
		{
			fmt.Println("Nope")
			return
		}
	}
	packages, _ := parser.FetchPackages()
	repo := repository.ReadRepository()
	fmt.Println(repo.Package(packages[3].Name))
}
