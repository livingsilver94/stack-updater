package stack

import (
	"fmt"
	"strings"

	"github.com/livingsilver94/stack-updater/html"
)

const (
	// MATEBaseURL is base URL from which it is possible
	// to build a tarball's URL from MATE.
	MATEBaseURL = "https://pub.mate-desktop.org/releases"
	// MATEFileExtension is the default tarball extension for MATE sources.
	MATEFileExtension = ".tar.xz"
)

// MATEHandler fetches packages from the MATE stack.
type MATEHandler struct {
	BaseURL string
	Version string
}

// NewMATEHandler returns a struct to handle the MATE stack, with a default
// base URL.
func NewMATEHandler(version string) MATEHandler {
	return MATEHandler{Version: version, BaseURL: MATEBaseURL}
}

// FetchPackages returns a list of Package belonging to the MATE stack.
func (mate MATEHandler) FetchPackages() ([]Package, error) {
	sourceURL := fmt.Sprintf("%s/%s", MATEBaseURL, mate.Version)
	sourcePage, err := download(sourceURL)
	if err != nil {
		return nil, err
	}
	entries, err := html.ParseListPage(sourcePage)
	if err != nil {
		return nil, err
	}

	var packages []Package
	for _, entry := range entries {
		if strings.HasSuffix(entry, MATEFileExtension) {
			pkgURL := fmt.Sprintf("%s/%s", sourceURL, entry)
			if pkg, err := PackageFromFilename(entry, pkgURL); err == nil {
				lastPkg := lastElement(packages)
				// MATE lists in the same table multiple possible versions
				// of an application. We want to filter away the older versions and
				// keep only the most recent one.
				if lastPkg != nil && lastPkg.Name == pkg.Name && lastPkg.Version < pkg.Version {
					*lastPkg = pkg
				} else {
					packages = append(packages, pkg)
				}
			}
		}
	}
	return packages, nil
}

func lastElement(slice []Package) *Package {
	length := len(slice)
	if length == 0 {
		return nil
	}
	return &slice[length-1]
}
