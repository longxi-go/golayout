# 开发指南

## 本地工作区 {#local-workspace}

golayout 可以单独克隆和验证，不依赖 gosr、ioc 或 nemo。
使用 `GOWORK=off make check` 验证独立模块。可选的生态联调才需要同级目录及本地工作区：

```sh
go work init ./golayout ./gosr ./gosr/tools ./ioc ./nemo
go work edit -replace=github.com/longxi-go/gosr@v0.0.0=./gosr
go work edit -replace=github.com/longxi-go/gosr/tools@v0.0.0=./gosr/tools
```

已有工作区使用 `go work use` 扩展，`gosr/tools` 必须单独纳入。
go.work 和 go.work.sum 保持本地使用，生产 go.mod 不应包含同级目录 replace。
在每个仓库/模块内部执行 Go 检查。

其他模块尚未发布的依赖仍可通过本地工作区替换进行开发，但不能据此声称它们已具备独立发布能力。
golayout 自身的模块检查不受这些占位依赖影响。

## 质量门禁

完整质量检查使用 Go 1.27.4，最低兼容版本为 Go 1.27.0。
本地安装配置要求的 golangci-lint，当前为 v2.15.0。

```sh
make setup
make format
make check-format
make license-headers
make check-license-headers
make check
```

完整门禁包括模块整洁性、vet、覆盖率测试、随机顺序 race、lint、代码布局、文件头及发布校验测试。
检查实际输出；本地执行不等于远程 CI 已通过。


## 文档

Python 3.12 和 uv 可独立于 Go 依赖构建站点：

```sh
make docs-setup
make docs-serve
make docs-build
```

最后一条命令使用锁定依赖严格构建双语内容。源码位于 `docs/design/docs`，新增页面同时提供
`page.md` 和 `page.zh.md`，并更新两份导航。API 示例应与当前源码一致，HTML 产物及虚拟环境保持忽略。

## 版本发布 {#releases}

Go 模块版本由不可变 Git tag 决定，go.mod 中的 `go` 是最低 Go 版本，文档工程版本也不是模块版本。
不需要 Maven/npm 注册步骤。v0.x 为 pre-v1，后续版本可能有意调整 API 或格式策略。

1. 审查代码和英文 `docs/releases/vX.Y.Z.md` 发布说明。
2. 运行 `GOWORK=off make check docs-build validate-release TAG=vX.Y.Z`。
3. 提交并推送已审查源码，等待分支 Quality 成功。
4. 在该提交创建带英文说明的稳定 tag，并推送 tag。

Release 校验稳定 SemVer 和对应说明的标题/正文，复用 Quality 校验 tag 源码，发布 GitHub Release，
然后调用 Docs Pages。非法、预发布、构建元数据 tag 均失败；v2+ 需要显式调整模块路径。
分发方式是 GitHub 源码归档及 Go module proxy；CLI 使用 `go install .../cmd/golayout@vX.Y.Z`。
已发布 tag 不应移动或复用；失败时可重跑失败任务，不能悄悄覆盖已有 Release。

## Pages 发布

分支推送和 PR 只在 Quality 构建文档；自动部署仅在稳定 tag 的 Release 发布成功后执行。
手动 Docs Pages 必须指定已发布、非 draft、非 prerelease 的版本，并显式检出该 tag。
直接调用可复用工作流，避免依赖 `GITHUB_TOKEN` 创建的 Release 事件再次触发流水线。

首次发布前，在 Settings → Pages 将 Source 设为 **GitHub Actions**。
`github-pages` environment 应允许发布 tag；需要手动重部署时也允许默认分支。
部署任务持有 Pages/OIDC 权限，普通检查只有只读权限。

站点配置地址为 `https://longxi-go.github.io/golayout/`，中文位于 `/golayout/zh/`。
发布构建从 `configure-pages` 获取基础 URL；本地或自定义预览可设置以 `/` 结尾的 `DOCS_SITE_URL`。
本地预览地址为 `http://127.0.0.1:8000/golayout/`。

创建本地工作流文件不会自动启用 Pages 或发布 tag，也不能证明线上 Actions 成功。
首次发布仍须获得用户确认。
