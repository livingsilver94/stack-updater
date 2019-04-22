package stack

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"strings"
)

type kdeHandler struct {
	BaseURL string
	Bundle  string
	Version string
}

// NewKDEHandler returns a struct to handle the KDE stack, with a default
// working URL.
//
// The KDE stack is split among bundles (applications, frameworks...)
// and every bundle has its own version, so these parameters are required.
func NewKDEHandler(bundle, version string) kdeHandler {
	return NewKDEHandlerWithURL(bundle, version, "https://cdn.download.kde.org/stable")
}

// NewKDEHandlerWithURL returns a struct to handle the KDE stack, using the provided URL.
//
// The KDE stack is split among bundles (applications, frameworks...)
// and every bundle has its own version, so these parameters are required.
func NewKDEHandlerWithURL(bundle, version, url string) kdeHandler {
	return kdeHandler{BaseURL: url,
		Bundle:  bundle,
		Version: version}
}

func (kde kdeHandler) FetchPackages() ([]Package, error) {
	fileExtension :=".tar.xz"

	if pageURL, pageData, err := kde.packagesPage(); err == nil {
		if files, err := kde.ParsePage(pageData); err == nil {
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

// ParsePage extracts a list of filenames from a given KDE HTML page.
func (kdeHandler) ParsePage(page []byte) ([]string, error) {
	var pkgList []string
	var err error

	pageString := enclosedString(string(page), "<ul>", "</ul>")
	tokenizer := html.NewTokenizer(strings.NewReader(pageString))
loop:
	for {
		// We are parsing lines like <li><a href="filename"> filename</a></li>
		switch tokenizer.Next() {
		case html.StartTagToken:
			{
				tokenizer.Next()
				tokenizer.Next()
				token := tokenizer.Token()
				pkgList = append(pkgList, strings.TrimSpace(token.Data))
			}
		case html.ErrorToken:
			{
				if parseErr := tokenizer.Err(); parseErr != io.EOF {
					err = fmt.Errorf("Cannot parse the page: %v", parseErr)
				}
				break loop
			}
		}
	}
	return pkgList, err
}

func (kde kdeHandler) packagesPage() (string, []byte, error) {
	urlPatterns := []string{"%s/%s/%s/src", "%s/%s/%s"}
	for _, url := range urlPatterns {
		fullURL := fmt.Sprintf(url, kde.BaseURL, kde.Bundle, kde.Version)
		if page, err := pageBody(fullURL); err == nil {
			return fullURL, page, nil
		}
	}
	return "", nil, fmt.Errorf("Cannot find %s, version %s", kde.Bundle, kde.Version)
}

func enclosedString(s, leftToken, rightToken string) string {
	leftIndex := strings.Index(s, leftToken)
	rightIndex := strings.LastIndex(s, rightToken)
	if leftIndex < 0 || rightIndex < 0 {
		return ""
	}
	return s[leftIndex+len(leftToken) : rightIndex]
}
