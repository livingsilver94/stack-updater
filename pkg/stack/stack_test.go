package stack_test

import (
	"github.com/livingsilver94/stack-updater/pkg/stack"
	"net/http/httptest"
	"testing"
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
		{"fooFoo-9.99.java.tar", stack.Package{"fooFoo", "9.99", url}},
	}
	for _, testVal := range tests {
		if ret, _ := stack.PackageFromFilename(testVal.in, url); ret != testVal.out {
			t.Errorf("in: %s\texpected: %v\tgot: %v", testVal.in, testVal.out, ret)
		}
	}
}
