package repository

import (
	"bytes"
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	git "gopkg.in/libgit2/git2go.v26"
)

const (
	// SourceBaseURL is a base URL from which to download package sources
	SourceBaseURL = "https://dev.getsol.us/source/"
	UnstableURL   = "https://packages.getsol.us/unstable/eopkg-index.xml.xz"
)

// Update represent a commit of a Solus package. Every new release is
// distinguished by Release
type Update struct {
	Version string `xml:"Version"`
	Release string `xml:"release,attr"`
}

// Package represents a software package inside the Solus repository
type Package struct {
	Name    string `xml:"Name"`
	Source  *packageSource
	Updates []Update `xml:"History>Update"`
}

// DownloadSources downloads this package's source files to directory.
// Internally, it works by cloning a git repository so that it's possible to manually
// browse into directory and perform usual git operations.
//
// DownloadSources also populate Package.Source field
func (pkg *Package) DownloadSources(directory string) error {
	sourcePath := filepath.Join(directory, pkg.Name)
	_, err := git.Clone(SourceBaseURL+pkg.Name, sourcePath, &git.CloneOptions{})
	if err == nil {
		sources, err := newPackageSource(sourcePath)
		if err == nil {
			pkg.Source = sources
		}
	}
	return err
}

// CurrentVersion returns package's latest version available in the repository
func (pkg *Package) CurrentVersion() string {
	return pkg.Updates[0].Version
}

// Repository represents the Solus repository containing a list of packages
type Repository struct {
	Packages []Package `xml:"Package"`
}

func GetUnstable(path string) (*Repository, error) {
	fileInfo, err := os.Stat(path)
	modTime := time.Time{}
	if err == nil {
		modTime = fileInfo.ModTime()
	}
	repoBody, repoSha, err := downloadArchive(modTime)
	if err != nil {
		return nil, err
	}
	if repoBody != nil {
		defer repoBody.Close()
		buf := bytes.Buffer{}
		buf.ReadFrom(repoBody)
		repoBytes := buf.Bytes()
		if fmt.Sprintf("%x", sha1.Sum(repoBytes)) != repoSha {
			return nil, fmt.Errorf("Repository's sha1 doesn't match the expected checksum")
		}
		destFile, _ := os.Create(path)
		extractArchive(bytes.NewReader(repoBytes), destFile)
	}
	return ReadAt(path)
}

// ReadAt initializes a new Repository by reading
// the Solus repository from path.
func ReadAt(path string) (*Repository, error) {
	xmlFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()
	return parseXML(xmlFile)
}

func parseXML(xmlDoc io.Reader) (*Repository, error) {
	fileBytes, err := ioutil.ReadAll(xmlDoc)
	if err != nil {
		return nil, err
	}
	var repo Repository
	xml.Unmarshal(fileBytes, &repo)
	return &repo, nil
}

func downloadArchive(modTime time.Time) (archive io.ReadCloser, sha1 string, err error) {
	client := &http.Client{}
	arReq, _ := http.NewRequest("GET", UnstableURL, nil)
	arReq.Header.Add("If-Modified-Since", modTime.Format("Wed, 21 Oct 2015 07:28:00 GMT"))
	arResp, err := client.Do(arReq)
	if err != nil {
		return
	}
	if arResp.StatusCode != 304 {
		shaResp, err := client.Get(UnstableURL + ".sha1sum")
		if err != nil {
			return archive, sha1, err
		}
		defer shaResp.Body.Close()

		archive = arResp.Body
		buf := bytes.Buffer{}
		buf.ReadFrom(shaResp.Body)
		sha1 = buf.String()
		return archive, sha1, nil
	}
	arResp.Body.Close()
	return nil, "", nil
}

func extractArchive(archive io.Reader, dest io.Writer) error {
	cmd := exec.Command("xz", "-d")
	cmd.Stdin = archive
	cmd.Stdout = dest
	err := cmd.Run()
	return err
}

// Package returns a package from the repository with the specified name.
// If no package is found, an nil value is returned
func (repo *Repository) Package(pkgName string) *Package {
	pkgIndex := sort.Search(len(repo.Packages), func(i int) bool {
		return repo.Packages[i].Name >= pkgName
	})

	if !(pkgIndex < len(repo.Packages) && repo.Packages[pkgIndex].Name == pkgName) {
		return nil
	}
	return &repo.Packages[pkgIndex]
}
