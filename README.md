# pd

Browse Python community discussions on discuss.python.org

`pd` is a single pure-Go binary. It speaks to pythondiscuss over plain
HTTPS, shapes the responses into clean records, and pipes into the rest of your
tools. No API key, nothing to run alongside it.

## Install

```bash
go install github.com/tamnd/pythondiscuss-cli/cmd/pd@latest
```

Or grab a prebuilt binary from the [releases](https://github.com/tamnd/pythondiscuss-cli/releases), or run
the container image:

```bash
docker run --rm ghcr.io/tamnd/pd:latest --help
```

## Usage

```bash
pd --help
pd version
```

This is a fresh scaffold. The command tree starts with `version`; build out the
real commands in `cli/` on top of the `pythondiscuss` library package.

## Development

```
cmd/pd/   thin main, wires cli.Root into fang
cli/                 the cobra command tree
pythondiscuss/                the library: HTTP client and data models
docs/                tago documentation site
```

```bash
make build      # ./bin/pd
make test       # go test ./...
make vet        # go vet ./...
```

## Releasing

Push a version tag and GitHub Actions runs GoReleaser, which builds the
archives, Linux packages, the multi-arch GHCR image, checksums, SBOMs, and a
cosign signature:

```bash
git tag v0.1.0
git push --tags
```

The Homebrew and Scoop steps self-disable until their tokens exist, so the first
release works with no extra secrets.

## License

Apache-2.0. See [LICENSE](LICENSE).
