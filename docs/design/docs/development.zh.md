# 开发指南

<a id="quality-gates"></a>

## 质量检查

开发需要 Git、Bash、Make 和 Go。CI 使用 Go 1.27.0 检查兼容性，使用 Go 1.27.4 执行完整质量检查。
安装 `.github/workflows/quality.yml` 指定的 linter 版本：

```sh
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0
```

确保 `golangci-lint` 位于 `PATH` 中，或通过 `GOLANGCI_LINT` 指定可执行文件。
在仓库根目录执行：

```sh
make setup
make license-headers format
GOWORK=off make check
```

`setup` 解析格式化工具依赖，不安装二进制或修改模块文件。
`license-headers` 和 `format` 写入项目拥有的源码。
`check` 对源码执行只读检查，涵盖文件头、布局、模块依赖整洁性、vet、带覆盖率的单元测试、
随机顺序的 race 测试、lint 及工具测试；执行过程可能写入工具缓存。
单独的检查目标见 `make help`。

<a id="local-workspace"></a>

## 本地工作区

golayout 不依赖其他生态仓库，可以单独测试。
需要联调使用方时，在两个工程的父目录创建工作区：

```sh
go work init ./golayout ./consumer
```

将 `./consumer` 替换为实际路径。已有工作区使用 `go work use` 扩展。
独立 tools 模块等嵌套模块需要显式加入。`go.work` 和 `go.work.sum` 仅供本地使用；
生产 `go.mod` 不得包含开发者专属的目录替换。

在各模块内部运行检查。工作区构建不能证明使用方声明的依赖版本可独立解析，
此项检查需要设置 `GOWORK=off`。

## 文档

站点需要 Python 3.12 和 uv。在仓库根目录执行：

```sh
make docs-setup
make docs-serve
make docs-build
```

`docs-build` 使用锁定依赖和严格校验。
英文预览地址为 `http://127.0.0.1:8000/golayout/`，中文位于 `/golayout/zh/`。
站点源码位于 `docs/design/docs`。新增页面需同时提供 `page.md` 和 `page.zh.md`，
并更新 `mkdocs.yml` 的两份导航。生成的 HTML 和虚拟环境不纳入版本控制。

提交文档前，应依据实现审查两个语言版本，核对前提、命令执行位置、API 行为、限制、版本和链接。
构建成功只能证明可渲染，不能证明事实或翻译准确。

<a id="releases"></a>

## 版本发布

Git tag 定义模块版本，`go.mod` 的 `go` 指令定义最低 Go 版本，文档工程版本不控制模块发布。
每个版本在 `docs/releases/vX.Y.Z.md` 保存英文说明，以 `# golayout vX.Y.Z` 开头并包含非空正文。
README 说明使用方法和发布流程；版本历史由这些说明和 GitHub Releases 记录。

1. 添加该版本说明，同步更新两个语言版本的安装示例。
2. 审查完整差异，包括文档准确性和双语一致性。
3. 执行 `GOWORK=off make check docs-build validate-release TAG=vX.Y.Z`，
   将 `vX.Y.Z` 替换为待发布版本。
4. 提交并推送已审查文件，等待该提交的 Quality 检查通过。
5. 在同一提交上创建附注 tag，说明使用 `Release golayout vX.Y.Z`，然后推送 tag。
6. 核对 Release 工作流、已发布说明和 Pages 内容。

推送 `v*` tag 会触发 Release。校验仅接受稳定版本 `vX.Y.Z`，不允许前导零、预发布后缀或构建元数据。
v2 及后续版本需要版本化模块路径和发布策略调整。当前 v0.x 版本可能有意调整 API 或格式规则。

Release 通过可复用的 Quality 检查 tag 源码，发布 GitHub Release，然后部署同一 tag 的文档。
Go 模块代理在 tag 推送后即可解析该版本，不依赖 GitHub Release 是否成功。
已发布 tag 不得移动或复用。工作流失败时应解决原因并重跑失败任务，不得覆盖已有 Release。

## Pages 发布

Pages 提供[英文手册](https://longxi-go.github.io/golayout/)和[中文手册](https://longxi-go.github.io/golayout/zh/)。
分支推送和 Pull Request 只构建文档，不执行部署；Release 发布稳定版本的 GitHub Release 后部署文档。

手动执行 Docs Pages 时，`release_tag` 必须指向已公开发布的稳定版本，工作流会检出该 tag。
自动部署要求 `github-pages` 环境允许 `v*` tag；从 `main` 手动执行时，还需允许该分支。
Pages 的 Source 使用 GitHub Actions。

部署任务具有 Pages 和 OIDC 权限，Quality 仅有仓库只读权限。
Release 直接调用 Docs Pages，不依赖 `GITHUB_TOKEN` 再次产生事件触发部署。

发布构建从 `configure-pages` 获取基础 URL。自定义域名或预览路径可通过 `DOCS_SITE_URL` 覆盖，
值需要以斜杠结尾。
