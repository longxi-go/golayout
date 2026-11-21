# golayout

`github.com/longxi-go/golayout` provides one formatting policy for Go tools,
generated source and application repositories. It adds AST-based structural
layout before pinned gofumpt, since gofumpt has no blanket multiline rule.

Requires Go 1.27.0 or newer. The module version is defined by Git tags such as
`v0.1.0`; the `go` directive is the minimum Go version, not the library version.
Nonempty function bodies, structs and control blocks expand to readable layouts.
Empty forms such as `struct{}`, `interface{}`, `select {}` and `func() {}` may stay
on one line. This pre-v1 project may change APIs and formatting policy in later releases.

## Install

After v0.1.0 is published, install the CLI or declare it as a module tool:

```sh
go install github.com/longxi-go/golayout/cmd/golayout@v0.1.0
go get -tool github.com/longxi-go/golayout/cmd/golayout@v0.1.0
```

Use `go get github.com/longxi-go/golayout@v0.1.0` for the library. Repository
contributors can clone golayout alone and run `GOWORK=off make check`; no sibling
module or local `replace` is required.

The library formats source bytes without loading application types:

```go
formatted, err := format.Source(source, format.Options{
    LangVersion: "go1.27",
})
```

The CLI checks explicit files by default and writes only with `-w`:

```sh
go tool golayout -lang go1.27 file.go
go tool golayout -w -lang go1.27 file.go
```

Repository Makefiles own file discovery and exclusions. See
[the architecture](docs/design/docs/architecture.md), [integration guide](docs/integration.md)
and [roadmap](docs/roadmap.md) before integrating another project. Do not copy
the implementation into consumers.

## Go source license headers

Keep the plain-text notice in `codequality/HEADER`. Refresh source-file comments
with `make license-headers`; use `make check-license-headers` for a read-only check.
Both use the shared CLI, which also accepts explicit files:

```sh
go tool golayout -header codequality/HEADER file.go
go tool golayout -w -header codequality/HEADER file.go
```

Header mode preserves the code body and does not format it. It updates existing
notices only for the template's copyright owner, without guessing current years.
Generated markers stay first; build constraints and package documentation remain.
Bulk selection uses Git-owned files, including eligible untracked files, and
excludes ignored files, snapshots, caches, vendor, testdata and symlinks.

Generators use `pkg/header.Source` with their target module's optional
`codequality/HEADER`, before formatting and validation. An absent template leaves
generated notices unchanged; a malformed template fails generation. The generator
distribution's copyright never supplies an application's notice.

## Local ecosystem development

golayout itself needs no workspace. For optional work across ecosystem repositories,
clone them into sibling directories and select the local modules in a parent go.work:

```sh
go work init ./golayout ./gosr ./gosr/tools ./ioc ./nemo
```

Until the other modules publish real dependency versions, map their development
placeholders explicitly in that local workspace:

```sh
go work edit -replace=github.com/longxi-go/gosr@v0.0.0=./gosr
go work edit -replace=github.com/longxi-go/gosr/tools@v0.0.0=./gosr/tools
```

Use `go work use` to extend an existing workspace. Include gosr/tools separately;
it has its own go.mod. Keep go.work and go.work.sum local. Run checks inside each
module, rather than `go test ./...` from the parent directory.

Workspace builds do not establish standalone release readiness for other modules.
Keep release checks independent of the local workspace using `GOWORK=off`.

## Releases

Branch pushes and pull requests run quality checks and a strict bilingual docs
build. A pushed stable `vX.Y.Z` tag starts Release: validate its maintained notes,
run quality checks against the exact tagged source, publish the GitHub Release,
then deploy that tag's documentation to Pages. Invalid/prerelease tags do not publish.
Pages can also redeploy an existing published stable release through a manual run.

Go modules are distributed by tags and Go module proxies; a GitHub Release is the
human-readable release record. No separate package registry or version field in
go.mod is needed. Published tags must not be moved or reused. See the
[v0.1.0 notes](docs/releases/v0.1.0.md) and [release guide](docs/design/docs/development.md#releases).

## Documentation site

The bilingual manual lives in [docs/design](docs/design/README.md).
Run `make docs-serve` to preview or `make docs-build` for the strict build.
The configured Pages destination is <https://longxi-go.github.io/golayout/>;
it becomes available after a stable release is published and Pages is enabled.
