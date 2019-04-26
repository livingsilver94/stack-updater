package repository

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	PkgDefinitionFile = "package.yml"
)

// This is a ugly yet necessary hack since we don't have a package.yml linter *yet*
// A package.yml file in fact is not a fully YAML-compliant file and we cannot parse it as a YAML
var lineMatcher = regexp.MustCompile(`^(?P<key>[[:alnum:]]+)(?P<separator> *: *)(?P<value>[[:print:]]+)$`)

type packageSource struct {
	path       string
	definition string
}

func newPackageSource(path string) (*packageSource, error) {
	defFile, err := ioutil.ReadFile(filepath.Join(path, PkgDefinitionFile))
	return &packageSource{path, string(defFile)}, err
}

func (source *packageSource) Version() string {
	return source.singleLineEntry("version")
}

func (source *packageSource) Release() string {
	return source.singleLineEntry("release")
}

func (source *packageSource) UpdateVersion(value string) {
	source.updateEntry("version", value)
}

func (source *packageSource) UpdateRelease(value string) {
	source.updateEntry("release", value)
}

func (source *packageSource) UpdateSource(url, sha256 string) {
	var builder strings.Builder
	var newLine string
	lines := source.readDefLines()
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "source") {
			builder.WriteString(lines[i] + "\n")
			newLine = fmt.Sprintf("    - %s : %s", url, sha256)
			i++
		} else {
			newLine = lines[i]
		}
		builder.WriteString(newLine + "\n")
	}
	source.definition = builder.String()
}

func (source *packageSource) Write() error {
	return ioutil.WriteFile(filepath.Join(source.path, PkgDefinitionFile), []byte(source.definition), 0644)
}

func (source *packageSource) singleLineEntry(key string) string {
	for _, line := range source.readDefLines() {
		// matches[1] means <key> group
		// matches[3] means <value> group
		if matches := lineMatcher.FindStringSubmatch(line); matches != nil && matches[1] == key {
			return matches[3]
		}
	}
	return ""
}

func (source *packageSource) updateEntry(key, value string) {
	var builder strings.Builder
	for _, ymlLine := range source.readDefLines() {
		if strings.HasPrefix(ymlLine, key) {
			ymlLine = lineMatcher.ReplaceAllString(ymlLine, "${key}${separator}"+value)
		}
		builder.WriteString(ymlLine + "\n")
	}
	source.definition = builder.String()
}

func (source *packageSource) readDefLines() []string {
	return strings.Split(source.definition, "\n")
}
