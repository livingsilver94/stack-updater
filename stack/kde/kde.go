package kde

import (
	"fmt"
	"github.com/livingsilver94/stack_updater/stack"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"strings"
)

const BaseURL = "https://cdn.download.kde.org/stable"

type KDEStack struct {
	Bundle  string
	Version string
}

func (kde KDEStack) FetchPackages() (map[string]string, error) {
	page, _ := kde.packagesPage()
	fmt.Println(kde.ParsePage(page))
	return nil, nil
}

func (KDEStack) ParsePage(page io.Reader) ([]string, error) {
	pkgList := make([]string, 0, 20)
	var pageString string
	if pageData, err := ioutil.ReadAll(page); err != nil {
		return pkgList, err
	} else {
		pageString = string(pageData)
	}
	pageString = enclosedString(pageString, "<ul>", "</ul>")
	tokenizer := html.NewTokenizer(strings.NewReader(pageString))
	// We are parsing lines like <li><a href="filename"> filename</a></li>
loop:
	for {
		switch tokenizer.Next() {
		case html.StartTagToken:
			{
				tokenizer.Next()
				tokenizer.Next()
				token := tokenizer.Token()
				pkgList = append(pkgList, token.Data)
			}
		case html.ErrorToken:
			{
				break loop
			}
		}
	}
	return pkgList, nil
}

func (kde KDEStack) packagesPage() (io.Reader, error) {
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
	return s[leftIndex+len(leftToken) : rightIndex]
}
