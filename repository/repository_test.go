package repository_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/livingsilver94/stack-updater/repository"
	"testing"
)

func createRepo() *repository.Repository {
	return repository.ReadRepositoryAt("../test_data/repository.xml")
}

func TestReadRepository(t *testing.T) {
	expectedPkgs := []repository.Package{
		{Name: "pkg1", Source: nil, Updates: []repository.Update{{Version: "0.0.1", Release: "1"}}},
		{Name: "pkg2", Source: nil, Updates: []repository.Update{{Version: "0.0.2", Release: "2"}}},
	}

	repo := createRepo()
	if !cmp.Equal(repo.Packages, expectedPkgs) {
		t.Errorf("Expected: %v\tGot: %v", expectedPkgs, repo.Packages)
	}
}

func TestFindPackageInRepo(t *testing.T) {
	repo := createRepo()
	if pkg1, err := repo.Package("pkg1"); pkg1.Name != "pkg1" || err != nil {
		t.Errorf("Expected pkg1 with a nil error;\t Got package name \"%v\" with error \"%v\"", pkg1.Name, err)
	}

	if _, err := repo.Package("nonexistent"); err == nil {
		t.Errorf("Expected an error but got no one")
	}
}
