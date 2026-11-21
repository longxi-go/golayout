# Documentation site

[English](README.md) | [简体中文](README.zh.md)

The [published manual](https://longxi-go.github.io/golayout/) contains English and
Simplified Chinese documentation for a released version. MkDocs Material builds
both languages from `docs/`; Python 3.12 dependencies are locked in `uv.lock`.

## Build and preview

Install Python 3.12 and uv. From the repository root:

```sh
make docs-setup
make docs-serve
make docs-build
```

Preview English at `http://127.0.0.1:8000/golayout/` and Chinese at `/golayout/zh/`.
`docs-build` performs a strict build. To build from this directory directly:

```sh
uv sync --locked
uv run --locked mkdocs build --strict
```

Set `DOCS_SITE_URL` with a trailing slash to override the base URL. Keep `site/`
and `.venv/` untracked.

## Maintain both languages

Each page has an English `page.md` and a Chinese `page.zh.md`. Update both
navigation trees in `mkdocs.yml` when adding pages. Keep commands, code examples,
requirements and limitations equivalent. Review technical accuracy separately
from build and link validation.

## Publish

Quality builds both languages on branch pushes and pull requests. Release
publishes a stable GitHub Release, then deploys documentation from its tag.
Manual Docs Pages runs redeploy an existing published stable release.

See [Development](docs/development.md) for prerequisites, environment policies and
the release checklist. Repository README and supplementary guides are maintained
outside the site; version-specific English notes remain in `docs/releases/` at
the repository root.
