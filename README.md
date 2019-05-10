# stack-updater
[![Go Report Card](https://goreportcard.com/badge/github.com/livingsilver94/stack-updater)](https://goreportcard.com/report/github.com/livingsilver94/stack-updater)&nbsp;

`stack-updater` is a simple utility to ease Solus development. It's been written to automate big software stack updates (meaning pieces of software that come split in many packages), so that maintainers won't need to fetch tarballs manually and bump the release number in every package definition file anymore.\
`stack-updater` doesn't need to be run on Solus to work.

## How it works
`stack-updater` works by parsing the download page of a chosen software stack. Generally, these HTML pages are composed of a list (`<ul>...items...</ul>`) with tarball URLs so it's fairly easy to extract information from such contents.\
After the list extraction, we check the Solus Unstable repository by downloading and reading the proper .xml file, and if a package is both in the repository and in the fetched list, `stack-updater` if needed will download the package definition files and update the `package.yml` file with new data.

## Examples
The `--help` flag should provide all the necessary information. Anyway, here are some examples on how to run the command (note the `:` to separate stack name from a bundle):
```shell
stack-updater update kde:applications 19.04.0
stack-updater update mate 1.23
stack-updater update kde:frameworks 5.56 -t /destination/directory
```

## Dependencies
 - libgit (dynamically linked)
 - `xz` command

## License
TBD.
