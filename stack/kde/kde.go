package kde

import (
	"fmt"
	"github.com/livingsilver94/stack_updater/stack"
	"golang.org/x/net/html"
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

func (KDEStack) ParsePage(page []byte) ([]string, error) {
	pkgList := make([]string, 0, 20)
	pageString := enclosedString(string(page), "<ul>", "</ul>")
	tokenizer := html.NewTokenizer(strings.NewReader(pageString))
	// We are parsing lines like <li><a href="filename"> filename</a></li>
    tokenType := tokenizer.Next()
	for ;tokenType != html.ErrorToken; tokenType = tokenizer.Next() {
        if tokenType == html.StartTagToken {
            tokenizer.Next(); tokenizer.Next()
            token := tokenizer.Token()
			pkgList = append(pkgList, strings.Trim(token.Data, " "))
        }
	}
	return pkgList, nil
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
	return s[leftIndex+len(leftToken) : rightIndex]
}
