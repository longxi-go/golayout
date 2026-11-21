# 文档站点

[English](README.md) | [简体中文](README.zh.md)

[在线手册](https://longxi-go.github.io/golayout/zh/)包含已发布版本的英文和简体中文文档。
MkDocs Material 从 `docs/` 构建两个语言版本，使用 Python 3.12，依赖由 `uv.lock` 固定。

## 构建与预览

安装 Python 3.12 和 uv。在仓库根目录执行：

```sh
make docs-setup
make docs-serve
make docs-build
```

英文预览地址为 `http://127.0.0.1:8000/golayout/`，中文位于 `/golayout/zh/`。
`docs-build` 执行严格构建。在当前目录中也可以直接运行：

```sh
uv sync --locked
uv run --locked mkdocs build --strict
```

通过以斜杠结尾的 `DOCS_SITE_URL` 覆盖基础 URL。`site/` 和 `.venv/` 不纳入版本控制。

## 维护双语内容

每页同时提供英文 `page.md` 和中文 `page.zh.md`。新增页面时更新 `mkdocs.yml` 的两份导航。
两种语言的命令、代码示例、使用前提和限制应一致。构建和链接检查不能替代技术内容审查。

## 发布

Quality 在分支推送和 Pull Request 中构建两个语言版本。
Release 发布稳定版本的 GitHub Release 后，部署对应 tag 的文档。
手动执行 Docs Pages 可以重新部署已公开发布的稳定版本。

开发依赖、环境规则和发布清单见[开发指南](docs/development.zh.md)。
仓库 README 和补充指南在站点外维护；各版本的英文说明保存在仓库根目录下的 `docs/releases/`。
