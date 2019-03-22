package stack

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func PageBody(url string) ([]byte, error) {
	reqResponse, err := http.Get(url)
	httpCode := reqResponse.StatusCode
	defer reqResponse.Body.Close()

	if err != nil || httpCode < 200 || httpCode >= 300 {
		return nil, fmt.Errorf("Cannot fetch page at address %s", url)
	}
	return ioutil.ReadAll(reqResponse.Body)
}

type Parser interface {
	FetchPackages() ([]Package, error)
}

type Package struct {
	Name    string
	Version string
	URL     string
}
