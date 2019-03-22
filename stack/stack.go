package stack

import (
    "net/http"
    "io/ioutil"
)

func PageBody(url string) ([]byte, error) {
    reqResponse, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer reqResponse.Body.Close()

    return ioutil.ReadAll(reqResponse.Body)
} 

type Parser interface {
    FetchPackages() (map[string]string, error)
}
