package kde

import (
	"fmt"
	"github.com/livingsilver94/stack_updater/stack"
	"golang.org/x/net/html"
	"io"
	"strings"
)

const BaseURL = "https://cdn.download.kde.org/stable"

type KDEStack struct {
	Bundle  string
	Version string
}

func (kde KDEStack) FetchPackages() ([]stack.Package, error) {
	page, _ := kde.packagesPage()
	fmt.Println(kde.ParsePage(page))
	return nil, nil
}

func (KDEStack) ParsePage(page []byte) ([]string, error) {
	var pkgList = make([]string, 0, 20)
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

func (kde KDEStack) packagesPage() ([]byte, error) {
	fullURL := fmt.Sprintf("%s/%s/%s/src", BaseURL, kde.Bundle, kde.Version)
	if page, err := stack.PageBody(fullURL); err == nil {
		return page, nil
	}
	// Remove "src"
	fullURL = fullURL[:len(fullURL)-3]
	if page, err := stack.PageBody(fullURL); err == nil {
		return page, nil
	}
	return nil, nil
}

func enclosedString(s, leftToken, rightToken string) string {
	leftIndex := strings.Index(s, leftToken)
	rightIndex := strings.LastIndex(s, rightToken)
	if leftIndex < 0 || rightIndex < 0 {
		return ""
	}
	return s[leftIndex+len(leftToken) : rightIndex]
}
