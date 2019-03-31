package repository

import (
	"encoding/xml"
	"gopkg.in/libgit2/git2go.v26"
	"io/ioutil"
	"os"
	"sort"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	Filepath          = "/var/lib/eopkg/index/Solus/eopkg-index.xml"
	SourceBaseURL     = "https://dev.getsol.us/source/"
	PkgDefinitionFile = "package.yml"
)

// This is a ugly yet necessary hack since we don't have a package.yml linter *yet*
// A package.yml file in fact is not a fully YAML-compliant file and we cannot parse it as a YAML
var yamlEditor = regexp.MustCompile(`^(?P<key>[[:alnum:]]+)(?P<separator> *: *)(?P<value>[[:print:]]+)$`)

type packageSource struct {
	path string
	definition string
}

func newPackageSource(path string) (*packageSource, error) {
	defFile, err := ioutil.ReadFile(filepath.Join(path, PkgDefinitionFile))
	if err == nil {
		return &packageSource{path, string(defFile)}, err
	}
	return nil, err
}

func (source *packageSource) Entry(key string) string {
	for _, line := range source.readDefLines() {
		if matches := yamlEditor.FindStringSubmatch(line); matches != nil && matches[1] == key {
			return matches[3]
		}
	}
	return ""
}

func (source *packageSource) UpdateEntry(key, value string) error {
	var builder strings.Builder
	for _, line := range source.readDefLines() {
		if matches := yamlEditor.FindStringSubmatch(line); matches != nil && matches[1] == key {
			line = yamlEditor.ReplaceAllString(line, "${key}${separator}" + value)
		}
		builder.WriteString(line+"\n")
	}
	source.definition = builder.String()
	return source.writeDefFile()
}

func (source *packageSource) UpdateEntries(defs map[string]string) error {
	var builder strings.Builder
	for _, line := range source.readDefLines() {
		for key, value := range defs {
			if matches := yamlEditor.FindStringSubmatch(line); matches != nil && matches[1] == key {
				line = yamlEditor.ReplaceAllString(line, "${key}${separator}" + value)
				break
			}
		}
		builder.WriteString(line+"\n")
	}
	source.definition = builder.String()
	return source.writeDefFile()
}

func (source *packageSource) writeDefFile() error {
	return ioutil.WriteFile(filepath.Join(source.path, PkgDefinitionFile), []byte(source.definition), 0644)
}

func (source *packageSource) readDefLines() []string {
	return strings.Split(source.definition, "\n")
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
