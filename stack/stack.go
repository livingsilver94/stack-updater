package stack

import (
    "io"
    "net/http"
)

func PageBody(url string) (io.Reader, error) {
    reqResponse, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    if err != nil {
        return nil, err
    }
    defer reqResponse.Body.Close()
    return reqResponse.Body, nil
} 

type Parser interface {
    FetchPackages() (map[string]string, error)
}
