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
	return KDEHandler{Bundle: bundle, Version: version, BaseURL: "https://download.kde.org/stable"}
}

// FetchPackages returns a list of Package objects belonging to the bundle
// and version specified in the KDE handler.
func (kde KDEHandler) FetchPackages() ([]Package, error) {
	fileExtension := ".tar.xz"

	pageURL, pageData, err := kde.findCorrectPage()
	if err != nil {
		return nil, err
	}
	files, err := parsePage(pageData)
	if err != nil {
		return nil, err
	}

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

func parsePage(page io.ReadCloser) ([]string, error) {
	defer page.Close()

	tokenizer := html.NewTokenizer(page)
	for tokenType := tokenizer.Next(); tokenType != html.ErrorToken; tokenType = tokenizer.Next() {
		token := tokenizer.Token()
		if tokenType == html.StartTagToken && (token.DataAtom == atom.Ul || token.DataAtom == atom.Table) {
			// We found the list
			return parseList(tokenizer, token.DataAtom), nil
		}
	}
	// We couldn't find the list
	return nil, fmt.Errorf("Couldn't find a list or table in this page")
}

func parseList(tokenizer *html.Tokenizer, HTMLTag atom.Atom) []string {
	var pkgList []string
LOOP:
	for {
		switch tokenizer.Next() {
		case html.StartTagToken:
			{
				listElement := tokenizer.Token()
				if listElement.DataAtom == atom.A {
					attribs := listElement.Attr
					for i := range attribs {
						if attribs[i].Key == "href" {
							pkgList = append(pkgList, strings.TrimSpace(attribs[i].Val))
							break
						}
					}
				}
			}
		case html.EndTagToken:
			{
				t := tokenizer.Token()
				if t.DataAtom == HTMLTag {
					break LOOP
				}
			}
		}
	}
	return pkgList
}
