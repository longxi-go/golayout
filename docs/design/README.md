# golayout documentation site

MkDocs Material builds English and Simplified Chinese from `docs/`, with local
search and light/dark themes. Dependencies are locked in `uv.lock` and use Python
3.12, following the reference kttp documentation structure.

From the repository root:

```sh
make docs-setup
make docs-serve
make docs-build
```

Preview: <http://127.0.0.1:8000/golayout/>; Chinese: `/golayout/zh/`.
Configured Pages URL: <https://longxi-go.github.io/golayout/>.
The URL describes the deployment destination, not an already deployed site.

Direct commands from this directory:

```sh
uv sync --locked
uv run --locked mkdocs build --strict
```

Override `DOCS_SITE_URL` for a custom domain/base path (include a trailing slash).
Add both `page.md` and `page.zh.md`; update both navigation trees in `mkdocs.yml`.
Never commit `site/` or `.venv/`. The lockfile and Markdown are project sources.

Quality strictly builds documentation for branch pushes and PRs without deployment.
Stable tag publication runs the same checks, creates a GitHub Release, then calls
Docs Pages to deploy that released tag through `github-pages`. Manual redeployment
also requires an existing published stable release. Choose GitHub Actions as the
Pages source before publication and permit release tags in environment rules.

See the [development guide](docs/development.md) for workspace/publication limits.
