package repository

import (
	"encoding/xml"
	"os"
	"sort"
	"path/filepath"
	"gopkg.in/libgit2/git2go.v26"
	"io/ioutil"
)

const (
	Filepath          = "/var/lib/eopkg/index/Solus/eopkg-index.xml"
	SourceBaseURL     = "https://dev.getsol.us/source/"
)

type update struct {
	version string `xml:"Version"`
	release string `xml:"release,attr"`
}

type Package struct {
	Name    string `xml:"Name"`
	Source  *packageSource
	updates []update `xml:"History>Update"`
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
	return pkg.updates[0].version
}

type Repository struct {
	Packages []Package `xml:"Package"`
}

func ReadRepository() *Repository {
	if xmlFile, err := os.Open(Filepath); err == nil {
		defer xmlFile.Close()
		if fileBytes, err := ioutil.ReadAll(xmlFile); err == nil {
			var repo Repository
			xml.Unmarshal(fileBytes, &repo)
			return &repo
		}
	}
	return nil
}

func (repo *Repository) Package(pkgName string) Package {
	var pkg Package
	pkgIndex := sort.Search(len(repo.Packages), func(i int) bool {
		return repo.Packages[i].Name >= pkgName
	})
	if pkgIndex >= 0 {
		pkg = repo.Packages[pkgIndex]
	}
	return pkg
}
