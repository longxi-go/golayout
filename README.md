# golayout

[English](README.md) | [简体中文](README.zh_CN.md) | [Online manual](https://longxi-go.github.io/golayout/)

golayout provides shared Go source formatting and license-header tooling for
repositories and code generators. `pkg/format` combines AST layout rules with
gofumpt v0.12.0. `pkg/header` maintains source notices as a separate operation.
Neither API loads application packages.

Requires **Go 1.27.0 or later**. This pre-v1 module may change APIs or formatting
rules in future releases; pin a version for reproducible checks.

## Install and use

Install the CLI and check explicit files:

```sh
go install github.com/longxi-go/golayout/cmd/golayout@v0.1.1
golayout -lang go1.27 file.go
golayout -w -lang go1.27 file.go
```

The default mode is read-only. Differences or invalid input produce a nonzero
exit status; `-w` writes changes. The installation directory must be on `PATH`.

To pin the CLI as a Go module tool instead:

```sh
go get -tool github.com/longxi-go/golayout/cmd/golayout@v0.1.1
go tool golayout -lang go1.27 file.go
```

For generators, add `github.com/longxi-go/golayout/pkg/format@v0.1.1` with
`go get`, then call the in-memory API:

```go
formatted, err := format.Source(source, format.Options{
    LangVersion: "go1.27",
})
```

Here, `source` is a byte slice. Handle `err` before validating or writing the
result. Derive `LangVersion` from the target module's `go.mod`; the value above
is an example. See the [quick start](docs/design/docs/quick-start.md) for a
complete program.

## Formatting policy

- Expand nonempty function bodies, struct and interface declarations, and
  control-flow blocks into multiple lines.
- Separate top-level declarations and grouped type declarations.
- Expand named and inferred composite literals without resolving their types.
  Explicit map, slice and array literals expand when their values include
  composite or function literals; simple scalar collections follow gofumpt.
- Permit canonical empty forms such as `struct{}`, `interface{}`, `select {}`
  and empty functions.

Formatting does not choose package boundaries, split files or enforce every
readability convention. See the [architecture](docs/design/docs/architecture.md)
for the exact scope.

## Source license headers

`codequality/HEADER` supplies the plain-text notice for owned Go files:

```sh
golayout -header codequality/HEADER file.go
golayout -w -header codequality/HEADER file.go
```

Header mode preserves generated markers, build constraints, package documentation
and the code body. It rejects replacement of a different copyright owner and
does not format code. Run header and formatting operations separately.

Generators can call `pkg/header.Source` with an explicitly selected template
before formatting and validating output. File discovery, template loading and
file ownership remain the caller's responsibility.

## Development

golayout can be cloned and tested independently; no workspace or sibling module
is required. With the [development prerequisites](docs/design/docs/development.md#quality-gates)
installed, run:

```sh
make setup
make license-headers format
GOWORK=off make check
make docs-build
```

The repository gates select Git-tracked and eligible untracked source, excluding
`.tmp`, `.cache`, vendor files and symlinks. Formatting also includes complete
`testdata/**/*.go.tmpl` fixtures; header refreshes exclude testdata.
Consumer setup is covered by the [integration guide](docs/integration.md).

## Releases

Git tags define module versions; the `go` directive in `go.mod` defines the
minimum Go version. Branch pushes and pull requests run Quality. A `v*` tag
starts Release, which validates stable `vX.Y.Z` input and maintained notes,
checks the tagged source, publishes a GitHub Release and deploys its documentation.
Published tags must not be moved or reused.

[GitHub Releases](https://github.com/longxi-go/golayout/releases) contains version
history. Each version has English notes under [`docs/releases`](docs/releases);
the [release guide](docs/design/docs/development.md#releases) describes the process.
This section documents that process rather than accumulating version entries.

## Documentation

- [English manual](https://longxi-go.github.io/golayout/)
- [简体中文手册](https://longxi-go.github.io/golayout/zh/)
- [Site maintenance](docs/design/README.md)
- [Roadmap](docs/roadmap.md)

Pages serves documentation from a published stable release. Branch changes are
built by CI and become public when a release deploys them. Use `make docs-serve`
to preview local changes.
