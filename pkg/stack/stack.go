package stack

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func PageBody(url string) ([]byte, error) {
	reqResponse, err := http.Get(url)
	defer reqResponse.Body.Close()

	httpCode := reqResponse.StatusCode
	if err != nil || httpCode < 200 || httpCode >= 300 {
		return nil, fmt.Errorf("Cannot fetch page at address %s", url)
	}
	return ioutil.ReadAll(reqResponse.Body)
}

func PackageFromFilename(filename, url string) Package {
	extFinder := regexp.MustCompile("(\\.[a-zA-Z]+)+")
	// Remove the extension (usually .tar.xz) from filename
	cleanName := filename[:extFinder.FindStringIndex(filename)[0]]
	lastDash := strings.LastIndex(cleanName, "-")
	return Package{cleanName[:lastDash], cleanName[lastDash+1:], url}
}

type Parser interface {
	FetchPackages() ([]Package, error)
}

type Package struct {
	Name    string
	Version string
	URL     string
}
