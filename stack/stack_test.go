package stack_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/livingsilver94/stack-updater/stack"
)

func TestPackageFromFilename(t *testing.T) {
	type Exp struct {
		in  string
		out stack.Package
	}
	url := "http://test.com"
	var tests = []Exp{
		{"this-is-a-123-TEST-1.1.zip", stack.Package{"this-is-a-123-TEST", "1.1", url}},
		{"this-is-a-123-TEST-1.1.tar.gz", stack.Package{"this-is-a-123-TEST", "1.1", url}},
		{"thisisnotvalid.zip", stack.Package{}},
		{"CaPiTaLs-9.99.verylongext.tar", stack.Package{"CaPiTaLs", "9.99", url}},
	}
	for _, testVal := range tests {
		if ret, _ := stack.PackageFromFilename(testVal.in, url); ret != testVal.out {
			t.Errorf("in: %s\texpected: %v\tgot: %v", testVal.in, testVal.out, ret)
		}
	}
}

func TestKDEFetchPackages(t *testing.T) {
	fakePage :=
		`<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 3.2 Final//EN">
		<html>
		<head>
		<title>Index of /stable/applications/19.04.0/src</title>
		</head>
		<body>
		<h1>Index of /stable/applications/19.04.0/src</h1>
		<ul><li><a href="/stable/applications/19.04.0/"> Parent Directory</a></li>
		<li><a href="akonadi-19.04.0.tar.xz"> akonadi-19.04.0.tar.xz</a></li></ul>`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fakePage)
	}))
	defer server.Close()

	kde := stack.KDEHandler{Bundle: "test", Version: "0.0", BaseURL: server.URL}
	pkgs, err := kde.FetchPackages()
	if err != nil {
		t.Errorf("Could not fetch packages: %v", err)
	}
	expectedPkg := stack.Package{Name: "akonadi", Version: "19.04.0", URL: server.URL + "/test/0.0/src/akonadi-19.04.0.tar.xz"}
	if pkgs[0] != expectedPkg {
		t.Errorf("Expected package: %v\tgot: %v", expectedPkg, pkgs[0])
	}
}
