package repository_test

import (
	"github.com/google/go-cmp/cmp"
	repo "github.com/livingsilver94/stack-updater/repository"
	"testing"
)

func TestReadRepository(t *testing.T) {
	expectedPkgs := []repo.Package{
		{Name: "pkg1", Source: nil, Updates: []repo.Update{{Version: "0.0.1", Release: "1"}}},
		{Name: "pkg2", Source: nil, Updates: []repo.Update{{Version: "0.0.2", Release: "2"}}},
	}

	repo := repo.ReadRepositoryAt("../test_data/repository.xml")
	if !cmp.Equal(repo.Packages, expectedPkgs) {
		t.Errorf("Expected: %v\tGot: %v", expectedPkgs, repo.Packages)
	}
}
