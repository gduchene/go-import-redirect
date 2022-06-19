# go-import-redirect

A thing that tells `go get` where to get Go packages and redirects users
to https://godoc.org.

## systemd Configuration

`go-import-redirect` takes the following flags:

* `-addr` for the address on which `go-import-redirect` should listen.
  It can either be a normal `IP:PORT` address or an absolute path to a
  UNIX socket that will be created. Defaults to `localhost:8080`. See
  https://golang.org/pkg/net/#Dial for more details.
* `-from` for the regular expression that the import path must match,
  including any capturing group, e.g. `go\\.example\\.com/(.+)`.
* `-to` for the URL that will be used to build the repository URL.
  Capturing groups can be used, e.g. `https://git.example.com/$1`.
* `-vcs` for the type of VCS you are using, e.g. `git`. Defaults to
  `git`.

Additionally, a configuration file can be passed with `-c`. See
`archlinux/go-import-redirect.conf` for an example.

It is recommended to enable the companion systemd socket and customize
it so systemd can start the service when needed and pass the socket to
`go-import-redirect`.

If you do not want to use socket activation, you must override the
`IPAddressDeny` and `RestrictAddressFamilies` unit settings to
appropriate values.

Likewise, you must customize the service definition to pass the right
flag values.

## AWS Lambda Configuration

The AWS Lambda version requires the use of the following environment
variables: `FROM`, `TO`, and `VCS`. Those have the same semantics as the
flags described above.

Note that this version needs to be compiled as a Linux binary and with
the `aws` build tag set.

## Docker Configuration

You can build the Docker image as usual, and pass the `-from`, `-to`,
and `-vcs` flags when you invoke `docker run`. The `-addr` flag is
already set to an appropriate value in the Dockerfile.
