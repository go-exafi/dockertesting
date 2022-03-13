# `go-dockertesting`

An opinionated extension of
[`ory/dockertest`](https://github.com/ory/dockertest) which will build a
Dockerfile, expose all ports, wait for its health check to pass, and return the
[`dockertest.Resource`](https://pkg.go.dev/github.com/ory/dockertest/v3#Resource)
which was created.  It will also register a cleanup handler which may be run if
the tests aren't killed with C-c.  The resource should also be set to expire
to ensure no dangling containers.

## License

See [LICENSE](./LICENSE) file.
