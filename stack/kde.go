package stack

import (
	"fmt"
	"io"
	"strings"

	"github.com/livingsilver94/stack-updater/html"
)

const (
	// KDEBaseURL is base URL from which it is possible
	// to build a tarball's URL from KDE
	KDEBaseURL = "https://download.kde.org/stable"
	// KDEFileExtension is the default tarball extension for KDE sources
	KDEFileExtension = ".tar.xz"
)

// KDEHandler returns a list of packages by parsing an HTML
// page from the KDE project. Since KDE is split among various bundles,
// this struct keeps only track of one of them. If you need to handle multiple
// bundles, you'll need to instantiate multiple KDE objects.
type KDEHandler struct {
	BaseURL string
	Bundle  string
	Version string
}

// NewKDEHandler returns a struct to handle the KDE stack, with a default
// base URL.
func NewKDEHandler(bundle, version string) KDEHandler {
	return KDEHandler{Bundle: bundle, Version: version, BaseURL: KDEBaseURL}
}

// FetchPackages returns a list of Package objects belonging to the bundle
// and version specified in the KDE handler.
func (kde KDEHandler) FetchPackages() ([]Package, error) {
	pageURL, pageData, err := kde.findCorrectPage()
	if err != nil {
		return nil, err
	}
	files, err := html.ParseListPage(pageData)
	if err != nil {
		return nil, err
	}

	var packages []Package
	for _, file := range files {
		if strings.HasSuffix(file, KDEFileExtension) {
			pkgURL := fmt.Sprintf("%s/%s", pageURL, file)
			if pkg, err := PackageFromFilename(file, pkgURL); err == nil {
				packages = append(packages, pkg)
			}
		}
	}
	return packages, nil
}

func (kde KDEHandler) findCorrectPage() (string, io.ReadCloser, error) {
	urlPatterns := []string{"%s/%s/%s/src", "%s/%s/%s"}
	for _, url := range urlPatterns {
		fullURL := fmt.Sprintf(url, kde.BaseURL, kde.Bundle, kde.Version)
		if page, err := download(fullURL); err == nil {
			return fullURL, page, nil
		}
	}
	return "", nil, fmt.Errorf("Cannot find %s, version %s", kde.Bundle, kde.Version)
}
