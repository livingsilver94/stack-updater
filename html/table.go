package html

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// ParseListPage returns all the entries in a HTML table or list.
func ParseListPage(page io.ReadCloser) ([]string, error) {
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

// parseList is actually the core logic of ParseListPage.
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
