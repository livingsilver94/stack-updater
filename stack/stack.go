package stack

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func pageBody(url string) ([]byte, error) {
	reqResponse, err := http.Get(url)
	defer reqResponse.Body.Close()

	httpCode := reqResponse.StatusCode
	if err != nil || httpCode < 200 || httpCode >= 300 {
		return nil, fmt.Errorf("Cannot fetch page at address %s", url)
	}
	return ioutil.ReadAll(reqResponse.Body)
}

func PackageFromFilename(filename, url string) (Package, error) {
	extFinder := regexp.MustCompile("(\\.[a-zA-Z]+)+")
	// Remove the extension (usually .tar.xz) from filename
	if indexes := extFinder.FindStringIndex(filename); indexes != nil {
		cleanName := filename[:indexes[0]]
		if lastDash := strings.LastIndex(cleanName, "-"); lastDash >= 0 {
			// Make sure filename has at least a dash (to separate name from version)
			return Package{cleanName[:lastDash], cleanName[lastDash+1:], url}, nil
		}
	}
	return Package{}, fmt.Errorf("Filename is not valid: %s", filename)
}

type Parser interface {
	FetchPackages() ([]Package, error)
}

type Package struct {
	Name    string
	Version string
	URL     string
}
