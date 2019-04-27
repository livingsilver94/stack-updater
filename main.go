package main

import (
	"fmt"
	"strings"

	"github.com/livingsilver94/stack-updater/repository"
	"github.com/livingsilver94/stack-updater/stack"
	"gopkg.in/alecthomas/kingpin.v2"
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
			parser = stack.NewKDE(*bundle, *version)
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
	stackPackages, _ := parser.FetchPackages()
	repo := repository.ReadRepository()
	for _, stackPkg := range stackPackages {
		if repoPkg, err := repo.Package(stackPkg.Name); err == nil {
			if stackPkg.Version >= "1.0.0" {
				repoPkg.DownloadSources("./pacchetti")
				repoPkg.Source.UpdateVersion("TEST")
				repoPkg.Source.UpdateRelease("LOL")
				repoPkg.Source.UpdateSource("https://example.com", "ABCdef")
				repoPkg.Source.Write()
			}
		}
	}
}
