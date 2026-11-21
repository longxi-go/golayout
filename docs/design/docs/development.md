# Development

## Local workspace {#local-workspace}

golayout can be cloned and tested alone; it has no dependency on gosr, ioc or nemo.
Use `GOWORK=off make check` to verify the standalone module.
For optional ecosystem work, clone siblings and create a local workspace:

```sh
go work init ./golayout ./gosr ./gosr/tools ./ioc ./nemo
go work edit -replace=github.com/longxi-go/gosr@v0.0.0=./gosr
go work edit -replace=github.com/longxi-go/gosr/tools@v0.0.0=./gosr/tools
```

Use `go work use` for an existing workspace. Include `gosr/tools` separately.
Keep go.work and go.work.sum local; production go.mod files must not contain
sibling-directory replacements. Run Go checks inside each repository/module.

Until other modules publish their real dependency versions, their workspace
substitutions remain local development aids. They do not prove standalone
release readiness. golayout's own module checks are independent of them.

## Quality gates

Use Go 1.27.4 for full quality checks, with Go 1.27.0 as the minimum compatibility
version. Install the configured golangci-lint version (currently v2.15.0).

```sh
make setup
make format
make check-format
make license-headers
make check-license-headers
make check
```

The full gate includes module tidiness, vet, coverage tests, shuffled race tests,
lint, file layout, header ownership and release-validation tests. Inspect actual
outputs; do not report a hosted CI result from a local run.


## Documentation

Python 3.12 and uv build the manual independently of Go dependencies:

```sh
make docs-setup
make docs-serve
make docs-build
```

The last command uses locked dependencies and a strict bilingual build. Sources
live in `docs/design/docs`; add both `page.md` and `page.zh.md` and update both
navigation trees. Keep API examples consistent with current source. Generated
HTML and the virtual environment stay ignored.

## Releases {#releases}

The Go module version comes from an immutable Git tag, not from the go.mod `go`
directive or the documentation project's version. No Maven/npm registry step is
needed. v0.x releases are pre-v1 and may intentionally change APIs or formatting.

1. Review code and maintained English `docs/releases/vX.Y.Z.md` notes.
2. Run `GOWORK=off make check docs-build validate-release TAG=vX.Y.Z`.
3. Commit and push the reviewed source, then wait for branch Quality to pass.
4. Create an annotated stable tag on that exact commit and push the tag.

Release validates stable SemVer and the matching notes title/body, calls reusable
Quality against tagged source, publishes the GitHub Release, then calls Docs Pages.
Malformed, prerelease and build-metadata tags fail validation; v2+ requires a
deliberate module-path change. The GitHub source archives and Go module proxy are
the distribution; install the CLI with `go install .../cmd/golayout@vX.Y.Z`.
Never move or reuse published tags. A failed publication can rerun its failed jobs;
an existing Release must not be overwritten silently.

## Pages publishing

Branch pushes and PRs build documentation in Quality without deployment. Automatic
Pages deployment happens only after Release publishes a stable tag. A manual
Docs Pages run requires an existing published, non-draft, non-prerelease tag and
checks out that tag explicitly. This direct reusable-workflow call also avoids
depending on a Release event generated with `GITHUB_TOKEN` to start another run.

In repository Settings → Pages, select **GitHub Actions** as source before release.
Allow release tags in the `github-pages` environment rules; allow the default
branch too if manual redeployment is needed. Deployment has Pages/OIDC permissions;
ordinary quality and documentation checks have read-only permissions.

The configured URL is `https://longxi-go.github.io/golayout/`, with Chinese under
`/golayout/zh/`. Published builds derive their base URL from `configure-pages`.
For local/custom preview, set `DOCS_SITE_URL` with a trailing `/`.
Local preview is `http://127.0.0.1:8000/golayout/`.

Creating local workflow files does not enable Pages, publish a tag or prove a
hosted Actions result. Initial publication still requires owner approval.
