package stack

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func pageBody(url string) (io.ReadCloser, error) {
	reqResponse, err := http.Get(url)

	httpCode := reqResponse.StatusCode
	if err != nil || httpCode < 200 || httpCode >= 300 {
		return nil, fmt.Errorf("Cannot fetch page at address %s", url)
	}
	return reqResponse.Body, nil
}

// PackageFromFilename parses a package's filename and returns a stack.Package
// instance if all was successfull. It expects a filename structured in
// this way: `pkgname-version.some.extension`.
//
// PackageFromFilename also takes a `url` argument since there's no way to get it
// from a filename.
func PackageFromFilename(filename, url string) (Package, error) {
	extFinder := regexp.MustCompile("(\\.[a-zA-Z]+)+")
	// Remove the extension (usually .tar.xz) from filename
	if indexes := extFinder.FindStringIndex(filename); indexes != nil {
		cleanName := filename[:indexes[0]]
		// Make sure filename has at least a dash (to separate name from version)
		if lastDash := strings.LastIndex(cleanName, "-"); lastDash >= 0 {
			name, version := cleanName[:lastDash], cleanName[lastDash+1:]
			return Package{Name: name, Version: version, URL: url}, nil
		}
	}
	return Package{}, fmt.Errorf("Filename is not valid: %s", filename)
}

// Parser is an interface representing the ability to build a list of Package
// from the information that a struct has.
type Parser interface {
	// FetchPackages returns a list of Package structs.
	FetchPackages() ([]Package, error)
}

// Package represents a piece of software from a certain stack.
type Package struct {
	Name    string
	Version string
	URL     string
}
