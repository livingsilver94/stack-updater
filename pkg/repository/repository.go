package repository

import (
	"encoding/xml"
	"gopkg.in/libgit2/git2go.v26"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"sort"
	"path/filepath"
)

const (
	Filepath          = "/var/lib/eopkg/index/Solus/eopkg-index.xml"
	SourceBaseURL     = "https://dev.getsol.us/source/"
	PkgDefinitionFile = "package.yml"
)

type packageSource struct {
	path string
	definitions yaml.MapSlice
}

func newPackageSource(path string) (*packageSource, error) {
	var definitions yaml.MapSlice
	ymlFile, err := ioutil.ReadFile(filepath.Join(path, PkgDefinitionFile))
	if err == nil {
		yaml.Unmarshal(ymlFile, &definitions)
		return &packageSource{path, definitions}, err
	}
	return nil, err
}

func (source *packageSource) ReadEntry(key string) interface{}{
	for i := range source.definitions {
		if source.definitions[i].Key == key {
			return source.definitions[i].Value
		}
	}
	return nil
}

func (source *packageSource) UpdateEntry(key string, value interface{}) error {
	source.updateEntry(key, value)
	return source.writeDefFile()
}

func (source *packageSource) UpdateEntries(defs map[string]interface{}) error {
	for key, value := range defs {
		source.updateEntry(key, value)
	}
	return source.writeDefFile()
}

func (source *packageSource) writeDefFile() error {
	file, err := yaml.Marshal(source.definitions)
	if err == nil {
		err = ioutil.WriteFile(filepath.Join(source.path, PkgDefinitionFile), file, 0644)
	}
	return err
}

func (source *packageSource) updateEntry(key string, value interface{}) {
	for i := range source.definitions {
		if source.definitions[i].Key == key {
			source.definitions[i].Value = value
			break
		}
	}
}


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
