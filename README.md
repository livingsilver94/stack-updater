# stack-updater
[![Go Report Card](https://goreportcard.com/badge/github.com/livingsilver94/stack-updater)](https://goreportcard.com/report/github.com/livingsilver94/stack-updater)&nbsp;

`stack-updater` is a simple utility to ease Solus development. It's been written to automate big software stack updates (meaning pieces of software that come split in many packages), so that maintainers won't need to fetch tarballs manually and bump the release number in every package definition file anymore.

## How it works
`stack-updater` works by parsing the download pages of a chosen software stack. Generally, these HTML pages are composed of a list (`<ul>...items...</ul>`) with tarball URLs so it's fairly easy to extract information from such contents.\
After the list extraction, we check the Solus official repository locally, by reading the proper .xml file, and if a package is both in the repository and in the fetched list, `stack-updater` will download the package definition files and update the `package.yml` file with the new tarball URL.

## License
TBD.