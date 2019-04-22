package repository

import (
	"encoding/xml"
	"os"
	"sort"
	"path/filepath"
	"gopkg.in/libgit2/git2go.v26"
	"io/ioutil"
	"fmt"
)

const (
	SourceBaseURL     = "https://dev.getsol.us/source/"
)

type Update struct {
	Version string `xml:"Version"`
	Release string `xml:"release,attr"`
}

type Package struct {
	Name    string `xml:"Name"`
	Source  *packageSource
	Updates []Update `xml:"History>Update"`
}

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

func (pkg *Package) CurrentVersion() string {
	return pkg.Updates[0].Version
}

type Repository struct {
	Packages []Package `xml:"Package"`
}

func ReadRepository() *Repository {
	return ReadRepositoryAt("/var/lib/eopkg/index/Solus/eopkg-index.xml")
}

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

func (repo *Repository) Package(pkgName string) (Package, error) {
	pkgIndex := sort.Search(len(repo.Packages), func(i int) bool {
		return repo.Packages[i].Name >= pkgName
	})

	if pkgIndex < len(repo.Packages) && repo.Packages[pkgIndex].Name == pkgName {
		return repo.Packages[pkgIndex], nil
	}
	return Package{}, fmt.Errorf("No package named %s in this repository", pkgName)
}
