package repository

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

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

// Repository represents the Solus repository containing a list of packages
type Repository struct {
	Packages []Package `xml:"Package"`
}

// ReadRepository initializes a new Repository by reading
// the Solus repository from a default filepath
func ReadRepository() *Repository {
	return ReadRepositoryAt("/var/lib/eopkg/index/Solus/eopkg-index.xml")
}

// ReadRepositoryAt initializes a new Repository by reading
// the Solus repository from the given filepath
func ReadRepositoryAt(path string) *Repository {
	if xmlFile, err := os.Open(path); err == nil {
		defer xmlFile.Close()

		if fileBytes, err := ioutil.ReadAll(xmlFile); err == nil {
			var repo Repository
			xml.Unmarshal(fileBytes, &repo)
			return &repo
		}
	}
	return nil
}

// Package returns a package from the repository with the specified name.
// If package is not found, an error is returned along with an empty Package object
func (repo *Repository) Package(pkgName string) (Package, error) {
	pkgIndex := sort.Search(len(repo.Packages), func(i int) bool {
		return repo.Packages[i].Name >= pkgName
	})

	if pkgIndex < len(repo.Packages) && repo.Packages[pkgIndex].Name == pkgName {
		return repo.Packages[pkgIndex], nil
	}
	return Package{}, fmt.Errorf("No package named %s in this repository", pkgName)
}
