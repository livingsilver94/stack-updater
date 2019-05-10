package cmd

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/format"
	"github.com/DataDrake/waterlog/level"
	"github.com/livingsilver94/stack-updater/repository"
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
	log.SetLevel(level.Info)
	log.SetFormat(format.Min)

	updateCmd.Flags().StringP("directory", "t", "", "where to store package sources")
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
			return fmt.Errorf(`You should not use ":" if you don't mean to specify a bundle`)
		}
	}
	return nil
}

func updateStack(cmd *cobra.Command, args []string) {
	stackParams := strings.Split(strings.ToLower(args[0]), ":")
	if len(stackParams) == 1 {
		stackParams = append(stackParams, "")
	}
	chosenStack, _ := stack.SupportedStackString(stackParams[0])
	stackHandler := stack.CreateStackHandler(chosenStack, args[1], stackParams[1])
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = "."
	}

	stackPackages, err := stackHandler.FetchPackages()
	if err != nil {
		log.Fatal(err)
	}
	repo, err := repository.GetUnstable(filepath.Join(cacheDir, "stack-updater", "unstable.xml"))
	if err != nil {
		log.Fatal(err)
	}

	for _, stackPkg := range stackPackages {
		repoPkg := repo.Package(stackPkg.Name)
		if repoPkg == nil {
			log.Infof("%s not found in Solus repository\n", stackPkg.Name)
			continue
		}
		if stackPkg.Version > repoPkg.CurrentVersion() {
			log.Printf("Updating %s from %s to %s\n", repoPkg.Name, repoPkg.CurrentVersion(), stackPkg.Version)
			if err := repoPkg.DownloadSources(cmd.Flag("directory").Value.String()); err != nil {
				log.Errorf("%s: %s. Skipping...\n", repoPkg.Name, err)
				continue
			}
			repoPkg.Source.UpdateRelease(repoPkg.Source.Release() + 1)
			repoPkg.Source.UpdateVersion(stackPkg.Version)
			repoPkg.Source.UpdateSource(stackPkg.URL, packageHash(stackPkg))
			if repoPkg.Source.Write() != nil {
				log.Errorf("Couldn't write updated sources for %s\n", repoPkg.Name)
			}
		}
	}
}

func packageHash(pkg stack.Package) string {
	tarball, err := pkg.Download()
	if err != nil {
		return ""
	}
	defer tarball.Close()

	hasher := sha256.New()
	io.Copy(hasher, tarball)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
