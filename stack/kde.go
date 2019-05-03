package stack

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
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
	return KDEHandler{Bundle: bundle, Version: version, BaseURL: "https://cdn.download.kde.org/stable"}
}

// FetchPackages returns a list of Package objects belonging to the bundle
// and version specified in the KDE handler.
func (kde KDEHandler) FetchPackages() ([]Package, error) {
	fileExtension := ".tar.xz"

	if pageURL, pageData, err := kde.findCorrectPage(); err == nil {
		if files, err := kde.parsePage(pageData); err == nil {
			var packages []Package
			for _, file := range files {
				if strings.HasSuffix(file, fileExtension) {
					pkgURL := fmt.Sprintf("%s/%s", pageURL, file)
					if pkg, err := PackageFromFilename(file, pkgURL); err == nil {
						packages = append(packages, pkg)
					}
				}
			}
			return packages, nil
		}
	}
	return nil, fmt.Errorf("Cannot fetch packages")
}

func (KDEHandler) parsePage(page io.ReadCloser) ([]string, error) {
	var pkgList []string
	var err error

	doc := html.NewTokenizer(page)
	for tokenType := doc.Next(); tokenType != html.ErrorToken; tokenType = doc.Next() {
		token := doc.Token()
		if tokenType == html.StartTagToken && token.DataAtom == atom.Ul {
			// We found the list
			for {
				switch doc.Next() {
				case html.StartTagToken:
					{
						doc.Next()
						doc.Next()
						pkgList = append(pkgList, strings.TrimSpace(doc.Token().Data))
					}
				case html.EndTagToken:
					{
						if doc.Token().DataAtom == atom.Ul {
							goto RETURN
						}
					}
				}
			}
		}
	}
	// We couldn't find the list
	err = fmt.Errorf("Couldn't find a list in this page")
RETURN:
	page.Close()
	return pkgList, err
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
