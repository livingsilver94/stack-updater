package repository

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"sort"
)

const (
	Filepath = "/var/lib/eopkg/index/Solus/eopkg-index.xml"
)

type Update struct {
	Version string `xml:"Version"`
	Release string `xml:"release,attr"`
}

type Package struct {
	Name    string   `xml:"Name"`
	Updates []Update `xml:"History>Update"`
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
