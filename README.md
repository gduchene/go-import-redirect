# go-import-redirect

A thing that tells `go get` where to get Go packages and redirects users
to https://godoc.org. It runs on AWS Lambda, but can easily be adapted
to run elsewhere.

You need to set up the following environment variables for it to work:

* `DEST` for the base URL that will be used to build the repository URL,
* `PREFIX` for the prefix that must be removed from your package name,
  e.g. `golang.org/x/` for `golang.org/x/image`, and
* `VCS` for the type of VCS you are using, e.g. `git`.
