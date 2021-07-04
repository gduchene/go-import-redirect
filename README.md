# go-import-redirect

A thing that tells `go get` where to get Go packages and redirects users
to https://godoc.org.

You need to set up the following environment variables for it to work:

* `DEST` for the base URL that will be used to build the repository URL,
* `PREFIX` for the prefix that must be removed from your package name,
  e.g. `golang.org/x/` for `golang.org/x/image`, and
* `VCS` for the type of VCS you are using, e.g. `git`.

There are three modes of operation:

1. **As a systemd service.** In this mode, it is recommended to enable
   the companion systemd socket and configure that. systemd will only
   start the service when needed. If you are not using the companion
   systemd socket but rather the systemd service directly, you will need
   to set the `ADDR` environment variable to tell go-import-redirect how
   to listen for incoming connections. If `ADDR` is set and starts with
   `/`, go-import-redirect will treat it as a UNIX socket path. If
   `ADDR` is not set, it will default to `:8080`. See
   https://golang.org/pkg/net/#Dial for more details.

2. **Inside a Docker container.** In this mode, the port 8080 will be
   exposed and no other environment variable besides the ones above need
   to be set.

3. **As an AWS Lambda.** This is very similar to Docker. You will need
   to set the environment variables above. Note that this version needs
   to be compiled as a Linux binary and with the `aws` build tag set.
