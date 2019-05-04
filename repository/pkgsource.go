package repository

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	// PkgDefinitionFile is the name of a package definition file, without
	// its base directory
	PkgDefinitionFile = "package.yml"
)

// This is an ugly yet necessary hack since we don't have a package.yml linter *yet*
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

// Version returns version field's value from the package definition
func (source *packageSource) Version() string {
	return source.singleLineEntry("version")
}

// Release returns release field's value from the package definition
func (source *packageSource) Release() int {
	intRelease, _ := strconv.Atoi(source.singleLineEntry("release"))
	return intRelease
}

// UpdateVersion replaces the version field's value in the package definition
func (source *packageSource) UpdateVersion(value string) {
	source.updateEntry("version", value)
}

// UpdateRelease replaces the release field's value in the package definition
func (source *packageSource) UpdateRelease(value int) {
	strRelease := strconv.Itoa(value)
	source.updateEntry("release", strRelease)
}

// UpdateSource replaces the source field's value in the package definition
func (source *packageSource) UpdateSource(url, sha256 string) {
	var builder strings.Builder
	var newLine string
	lines := source.readDefLines()
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "source") {
			builder.WriteString(lines[i] + "\n")
			newLine = fmt.Sprintf("    - %s : %s\n", url, sha256)
			i++
		} else {
			if lines[i] == "" {
				newLine = lines[i]
			} else {
				newLine = lines[i] + "\n"
			}
		}
		builder.WriteString(newLine)
	}
	source.definition = builder.String()
}

// Write writes the package definition to disk
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
		if ymlLine != "" {
			ymlLine += "\n"
		}
		builder.WriteString(ymlLine)
	}
	source.definition = builder.String()
}

func (source *packageSource) readDefLines() []string {
	return strings.Split(source.definition, "\n")
}
