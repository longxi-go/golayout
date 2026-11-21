# 开发指南

<a id="quality-gates"></a>

## 质量检查

开发需要 Git、Bash、Make，以及 `.github/workflows/quality.yml` 选定的 Go 版本。
最低语言版本由 `go.mod` 声明。
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
`license-headers` 和 `format` 修改项目维护的源码。
`check` 对源码执行只读检查，涵盖文件头、布局、模块依赖整洁性、vet、带覆盖率的单元测试、
随机顺序的 race 测试、lint 及工具测试；执行过程可能写入工具缓存。
单独的检查目标见 `make help`。

<a id="source-selection"></a>

## 文件选择

仓库脚本选择 Git 已跟踪文件和符合条件的未跟踪 Go 文件，排除被忽略的未跟踪文件。
格式化跳过 `.tmp`、`.cache`、vendor 文件和符号链接。
完整且可作为 Go 源码解析的 `testdata/**/*.go.tmpl` 样本参与格式化；
其他 testdata 文件和生产生成器模板不参与。

文件头刷新读取 `codequality/HEADER`，并额外排除 `.codex` 和全部 testdata。
第三方源码不应进入项目维护的文件清单。识别到其他版权主体时应处理报错，不得擅自改写归属。
CLI 本身仅接受显式指定的普通文件，不负责发现仓库文件。

<a id="local-workspace"></a>

## 本地工作区

需要联调使用方时，在两个工程的父目录创建工作区：

```sh
go work init ./golayout ./consumer
```

将 `./consumer` 替换为实际路径。已有工作区使用 `go work use` 扩展。
独立 tools 模块等嵌套模块需要显式加入。`go.work` 和 `go.work.sum` 仅供本地使用；
生产 `go.mod` 不得包含开发者专属的目录替换。

在各模块内运行检查。工作区构建不会验证使用方声明的依赖版本能否独立解析，
验证时需设置 `GOWORK=off`。

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
