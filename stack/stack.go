package stack

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// SupportedStack represents a stack supported by this package.
type SupportedStack int

//go:generate go run github.com/dmarkham/enumer -transform=lower -type=SupportedStack

const (
	// KDE is a supported stack
	KDE SupportedStack = iota
	// MATE is a supported stack
	MATE
)

func download(url string) (io.ReadCloser, error) {
	reqResponse, err := http.Get(url)

	httpCode := reqResponse.StatusCode
	if err != nil || httpCode < 200 || httpCode >= 300 {
		return nil, fmt.Errorf("Cannot fetch page at address %s", url)
	}
	return reqResponse.Body, nil
}

// PackageFromFilename parses a package's filename and returns a stack.Package
// instance if all was successful. It expects a filename structured in
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

// CreateStackHandler constructs and returns a stack handler if supported. It returns nil otherwise.
//
//Note: bundle is ignored if selected stack is known for not to be split in bundles.
func CreateStackHandler(stack SupportedStack, version, bundle string) Handler {
	switch stack {
	case KDE:
		{
			return NewKDEHandler(bundle, version)
		}
	case MATE:
		{
			return NewMATEHandler(version)
		}
	default:
		{
			return nil
		}
	}
}

// Handler is an interface representing the ability to build a list of Package
// from the information that a struct has.
type Handler interface {
	FetchPackages() ([]Package, error)
}

// Package represents a piece of software from a certain stack.
type Package struct {
	Name    string
	Version string
	URL     string
}

// Download downloads the package's tarball
func (pkg Package) Download() (io.ReadCloser, error) {
	return download(pkg.URL)
}
