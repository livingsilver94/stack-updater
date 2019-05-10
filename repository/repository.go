package repository

import (
	"bytes"
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"
)

const (
	// UnstableURL is where the Solus unstable repository is download from.
	UnstableURL = "https://packages.getsol.us/unstable/eopkg-index.xml.xz"
)

var httpClient = &http.Client{
	Timeout: time.Second * 20,
}

// Repository represents the Solus repository containing a list of packages.
type Repository struct {
	Packages []Package `xml:"Package"`
}

// GetUnstable fetches the unstable Solus repository and saves it to path.
// path is a file that will be created if it doesn't exist.
func GetUnstable(path string) (*Repository, error) {
	modTime := time.Time{}
	if fileInfo, err := os.Stat(path); err == nil {
		modTime = fileInfo.ModTime()
	}
	if HasBeenUpdated(modTime) {
		archive, err := downloadArchive()
		if err != nil {
			return nil, err
		}

		file, err := createFile(path)
		if err != nil {
			return nil, err
		}

		if err := extractArchive(bytes.NewReader(archive), file); err != nil {
			return nil, err
		}
	}
	return ReadAt(path)
}

// HasBeenUpdated tells if the unstable Solus repository has been modified since lastTime.
// HasBeenUpdated also returns true if a connection error occured.
func HasBeenUpdated(lastTime time.Time) bool {
	req, _ := http.NewRequest(http.MethodHead, UnstableURL, nil)
	req.Header.Add("If-Unmodified-Since", lastTime.Format(time.RFC1123))
	resp, err := httpClient.Do(req)
	if err != nil {
		return true
	}
	defer resp.Body.Close()
	return resp.StatusCode == 412
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

// Package returns a package from the repository with the specified name.
// If no package is found, an nil value is returned.
func (repo *Repository) Package(pkgName string) *Package {
	pkgIndex := sort.Search(len(repo.Packages), func(i int) bool {
		return repo.Packages[i].Name >= pkgName
	})

	if !(pkgIndex < len(repo.Packages) && repo.Packages[pkgIndex].Name == pkgName) {
		return nil
	}
	return &repo.Packages[pkgIndex]
}

func parseXML(xmlDoc io.Reader) (*Repository, error) {
	xmlDecoder := xml.NewDecoder(xmlDoc)
	var repo Repository
	err := xmlDecoder.Decode(&repo)
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

// downloadArchive downloads the Solus unstable repository archive. Internally, it checks
// archive's checksum and returns error if it doesn't match the expected one.
func downloadArchive() ([]byte, error) {
	arResp, err := httpClient.Get(UnstableURL)
	if err != nil {
		return nil, err
	}
	defer arResp.Body.Close()

	shaResp, err := httpClient.Get(UnstableURL + ".sha1sum")
	if err != nil {
		return nil, err
	}
	defer shaResp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(arResp.Body)
	arBytes := make([]byte, len(buf.Bytes()))
	copy(arBytes, buf.Bytes())

	buf.Reset()
	buf.ReadFrom(shaResp.Body)
	shaSum := buf.String()
	if fmt.Sprintf("%x", sha1.Sum(arBytes)) != shaSum {
		return nil, fmt.Errorf("Repository's sha1 doesn't match the expected checksum")
	}
	return arBytes, nil
}

func extractArchive(archive io.Reader, dest io.Writer) error {
	cmd := exec.Command("xz", "-d")
	cmd.Stdin = archive
	cmd.Stdout = dest
	err := cmd.Run()
	return err
}

func createFile(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), os.ModeDir); err != nil {
		return nil, err
	}
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}
