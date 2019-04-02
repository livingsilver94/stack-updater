package stack

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"strings"
)

const (
	BaseURL       = "https://cdn.download.kde.org/stable"
	FileExtension = ".tar.xz"
)

type KDE struct {
	Bundle  string
	Version string
}

func (kde KDE) FetchPackages() ([]Package, error) {
	if pageURL, pageData, err := kde.packagesPage(); err == nil {
		if files, err := kde.ParsePage(pageData); err == nil {
			var packages []Package
			for _, file := range files {
				if strings.HasSuffix(file, FileExtension) {
					pkgURL := fmt.Sprintf("%s/%s", pageURL, file)
					if pkg, err := PackageFromFilename(file, pkgURL); err == nil{
						packages = append(packages, pkg)
					}
				}
			}
			return packages, nil
		}
	}
	return nil, fmt.Errorf("Cannot fetch packages")
}

func (KDE) ParsePage(page []byte) ([]string, error) {
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

func (kde KDE) packagesPage() (string, []byte, error) {
	urlPatterns := []string{"%s/%s/%s/src", "%s/%s/%s"}
	for _, url := range urlPatterns {
		fullURL := fmt.Sprintf(url, BaseURL, kde.Bundle, kde.Version)
		if page, err := PageBody(fullURL); err == nil {
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
