package repository

import (
	"io/ioutil"
	"regexp"
	"strings"
	"path/filepath"
)

const (
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
