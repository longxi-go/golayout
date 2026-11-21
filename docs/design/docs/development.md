# Development

## Quality gates

Development requires Git, Bash, Make and the Go version selected by
`.github/workflows/quality.yml`. The minimum language version is declared in
`go.mod`.
Install the linter version used by `.github/workflows/quality.yml`:

```sh
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0
```

Ensure `golangci-lint` is on `PATH`, or set `GOLANGCI_LINT` to its executable.
From the repository root:

```sh
make setup
make license-headers format
GOWORK=off make check
```

`setup` resolves formatter dependencies without installing a binary or modifying
module files. `license-headers` and `format` write owned source. `check` is the
read-only gate for source notices, layout, module tidiness, vet, unit tests with
coverage, shuffled race tests, lint and tooling tests. It may populate tool caches.
Use `make help` for individual targets.

## Source selection

Repository scripts select Git-tracked and eligible untracked Go files; ignored
untracked files are excluded. Formatting skips `.tmp`, `.cache`, vendor files and
symlinks. Complete, parseable `testdata/**/*.go.tmpl` fixtures opt into formatting;
other testdata and production generator templates are excluded.

Header refreshes read `codequality/HEADER` and additionally exclude `.codex` and
all testdata. Keep third-party source outside the owned inventory. A different
recognized copyright owner causes the header operation to fail.
The CLI itself accepts explicit regular files and does not discover repositories.

## Local workspace

For joint development with a consumer, create a workspace in their parent
directory:

```sh
go work init ./golayout ./consumer
```

Replace `./consumer` with the actual path. Use `go work use` to extend an existing
workspace. Nested modules, such as a separate tools module, must be included
explicitly. Keep `go.work` and `go.work.sum` local; production `go.mod` files must
not contain developer-specific directory replacements.

Run checks inside each module. A workspace build does not verify that a consumer's
declared dependency versions resolve independently; use `GOWORK=off` for that check.

## Documentation

The site requires Python 3.12 and uv. From the repository root:

```sh
make docs-setup
make docs-serve
make docs-build
```

`docs-build` uses locked dependencies and strict validation. Preview the English
site at `http://127.0.0.1:8000/golayout/` and Chinese at `/golayout/zh/`.
Site sources are under `docs/design/docs`. Add both `page.md` and `page.zh.md`,
and update both navigation trees in `mkdocs.yml`. Keep generated HTML and the
virtual environment untracked.
