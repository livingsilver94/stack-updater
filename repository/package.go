package repository

import (
	"os"
	"path/filepath"

	git "gopkg.in/libgit2/git2go.v26"
)

const (
	// SourceBaseURL is a base URL from which to download package sources.
	SourceBaseURL = "https://dev.getsol.us/source/"
)

// Update represent a commit of a Solus package.
type Update struct {
	// Upstream version
	Version string `xml:"Version"`
	// Solus package release, incremental.
	Release string `xml:"release,attr"`
}

// Package represents a software package inside the Solus repository.
type Package struct {
	Name    string `xml:"Name"`
	Source  *packageSource
	Updates []Update `xml:"History>Update"`
}

// DownloadSources clones this package's git repository to directory
// and populates Package.Source field.
func (pkg *Package) DownloadSources(directory string) error {
	sourcePath := filepath.Join(directory, pkg.Name)
	if _, err := os.Stat(filepath.Join(sourcePath, PkgDefinitionFile)); os.IsNotExist(err) {
		_, err = git.Clone(SourceBaseURL+pkg.Name, sourcePath, &git.CloneOptions{})
		if err != nil {
			return err
		}
	}
	sources, err := newPackageSource(sourcePath)
	if err != nil {
		return err
	}
	pkg.Source = sources
	return nil
}

// CurrentVersion returns package's latest version available in the repository.
func (pkg *Package) CurrentVersion() string {
	return pkg.Updates[0].Version
}
