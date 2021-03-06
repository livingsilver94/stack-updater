package repository_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/livingsilver94/stack-updater/repository"
)

func createRepo() *repository.Repository {
	repo, _ := repository.ReadAt("../test_data/repository.xml")
	return repo
}

func TestReadRepository(t *testing.T) {
	expectedPkgs := []repository.Package{
		{Name: "pkg1", Source: nil, Updates: []repository.Update{{Version: "0.0.1", Release: 1}}},
		{Name: "pkg2", Source: nil, Updates: []repository.Update{{Version: "0.0.2", Release: 2}}},
	}

	repo := createRepo()
	if !cmp.Equal(repo.Packages, expectedPkgs) {
		t.Errorf("Expected: %v\tGot: %v", expectedPkgs, repo.Packages)
	}
}

func TestFindPackageInRepo(t *testing.T) {
	repo := createRepo()
	if pkg1 := repo.Package("pkg1"); pkg1.Name != "pkg1" {
		t.Errorf("Expected pkg1 but got %v", pkg1.Name)
	}

	if repo.Package("nonexistent") != nil {
		t.Errorf("Expected a nil value but got something different")
	}
}
