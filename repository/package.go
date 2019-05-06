package repository

import (
	"path/filepath"

	git "gopkg.in/libgit2/git2go.v26"
)

const (
	// SourceBaseURL is a base URL from which to download package sources
	SourceBaseURL = "https://dev.getsol.us/source/"
)

// Update represent a commit of a Solus package. Every new release is
// distinguished by Release
type Update struct {
	Version string `xml:"Version"`
	Release string `xml:"release,attr"`
}

// Package represents a software package inside the Solus repository
type Package struct {
	Name    string `xml:"Name"`
	Source  *packageSource
	Updates []Update `xml:"History>Update"`
}

// DownloadSources downloads this package's source files to directory.
// Internally, it works by cloning a git repository so that it's possible to manually
// browse into directory and perform usual git operations.
//
// DownloadSources also populate Package.Source field
func (pkg *Package) DownloadSources(directory string) error {
	sourcePath := filepath.Join(directory, pkg.Name)
	_, err := git.Clone(SourceBaseURL+pkg.Name, sourcePath, &git.CloneOptions{})
	if err == nil {
		sources, err := newPackageSource(sourcePath)
		if err == nil {
			pkg.Source = sources
		}
	}
	return err
}

// CurrentVersion returns package's latest version available in the repository
func (pkg *Package) CurrentVersion() string {
	return pkg.Updates[0].Version
}
