package stack

import (
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"strings"
)

type KDE struct {
	BaseURL string
	Bundle  string
	Version string
}

// NewKDE returns a struct to handle the KDE stack, with a default
// working URL.
//
// The KDE stack is split among bundles (applications, frameworks...)
// and every bundle has its own version, so these parameters are required.
func NewKDE(bundle, version string) KDE {
	return KDE{Bundle: bundle, Version: version, BaseURL: "https://cdn.download.kde.org/stable"}
}

func (kde KDE) FetchPackages() ([]Package, error) {
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

func (KDE) parsePage(page io.ReadCloser) ([]string, error) {
	var pkgList []string

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
							return pkgList, nil
						}
					}
				}
			}
		}
	}
	// We couldn't find the list
	return pkgList, fmt.Errorf("Couldn't find a list in this page")
}

func (kde KDE) findCorrectPage() (string, io.ReadCloser, error) {
	urlPatterns := []string{"%s/%s/%s/src", "%s/%s/%s"}
	for _, url := range urlPatterns {
		fullURL := fmt.Sprintf(url, kde.BaseURL, kde.Bundle, kde.Version)
		if page, err := pageBody(fullURL); err == nil {
			return fullURL, page, nil
		}
	}
	return "", nil, fmt.Errorf("Cannot find %s, version %s", kde.Bundle, kde.Version)
}
