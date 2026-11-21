# 接入 golayout

[English](integration.md) | [简体中文](integration.zh.md)

## 固定工具和库版本

源码检查和生成代码应使用同一个已发布且不可变的版本。
在管理工具依赖的模块中执行：

```sh
go get -tool github.com/longxi-go/golayout/cmd/golayout@v0.1.4
```

生成器导入相同模块版本中的库：

```sh
go get github.com/longxi-go/golayout/pkg/format@v0.1.4
```

如果生产模块需要保持无工具依赖，在仓库根目录创建独立模块管理 CLI：

```sh
mkdir -p tools
cd tools
go mod init example.com/project/tools
go get -tool github.com/longxi-go/golayout/cmd/golayout@v0.1.4
```

将 `example.com/project/tools` 替换为实际模块路径，提交其 `go.mod` 和 `go.sum`；
生成器在自身模块中声明库依赖。生产 `go.mod` 不应包含开发者专属的目录替换。

<a id="委托仓库检查"></a>

## 配置仓库检查

以下 Makefile 片段调用已选定模块版本中的共享脚本。将其加入使用方的 Makefile，
在仓库根目录运行。如果工具依赖声明在嵌套 tools 模块中，设置 `GOLAYOUT_TOOL_MODULE=tools`：

```makefile
GO ?= go
GOLAYOUT_TOOL_MODULE ?= .
GOLAYOUT_DIR := $(shell $(GO) -C "$(GOLAYOUT_TOOL_MODULE)" list -m \
    -f '{{.Dir}}' github.com/longxi-go/golayout)
GO_LANG_VERSION := go$(shell awk '$$1 == "go" { print $$2; exit }' go.mod)

.PHONY: setup format check-format

setup:
	@GO="$(GO)" bash "$(GOLAYOUT_DIR)/scripts/format.sh" setup \
	    "$(CURDIR)" "$(GOLAYOUT_TOOL_MODULE)" "$(GO_LANG_VERSION)"

format:
	@GO="$(GO)" bash "$(GOLAYOUT_DIR)/scripts/format.sh" write \
	    "$(CURDIR)" "$(GOLAYOUT_TOOL_MODULE)" "$(GO_LANG_VERSION)"

check-format:
	@GO="$(GO)" bash "$(GOLAYOUT_DIR)/scripts/format.sh" check \
	    "$(CURDIR)" "$(GOLAYOUT_TOOL_MODULE)" "$(GO_LANG_VERSION)"
```

脚本选择 Git 已跟踪文件和符合条件的未跟踪 Go 文件，排除 `.tmp`、`.cache`、vendor 文件和符号链接，
并包含完整的 `testdata/**/*.go.tmpl` 样本。其他 testdata 文件和生产生成器模板不参与格式化。
`setup` 验证依赖，不安装二进制；`check-format` 不重写源码或模块文件。
这些命令需要 Git、Bash、Make，以及所选 golayout 版本要求的 Go 版本。

语言选项取自使用方的 `go.mod`，与构建格式化工具所用的 Go 版本无关。
格式化生成代码时应使用相同的语言版本。

## CI 与本地开发

CI 通过 Go 的校验和验证机制解析工具模块声明的 golayout 版本，不需要同级源码仓库。
在使用方仓库根目录先下载依赖，再解析脚本路径：

```sh
GOWORK=off go mod download
GOWORK=off make setup check-format
```

使用嵌套 tools 模块时，执行：

```sh
GOWORK=off go -C tools mod download
GOWORK=off make setup check-format GOLAYOUT_TOOL_MODULE=tools
```

CI 缓存键应包含工具模块的 `go.mod` 和 `go.sum`。生成代码还需通过生成器自身的检查。

本地 `go.work` 可以选择另一份 golayout 源码，使用工作区时需显式包含嵌套工具模块。
设置 `GOWORK=off` 后，上述命令改用工具模块声明的依赖。

## 格式化生成代码

在生成器执行类型验证和文件写入之前，调用 `pkg/format.Source`，传入源码字节和目标模块的语言版本。
使用结果前须先处理错误；不要通过子进程调用 CLI，也不要复制私有格式化实现。

需要源码文件头时，显式读取目标仓库的模板，在格式化前调用 `pkg/header.Source`。
模板是否必需、缺失时如何处理，由调用方决定；golayout 不决定版权归属。
仓库批量文件头脚本读取使用方的 `codequality/HEADER`，排除全部 testdata。
文件头检查和写入独立于格式化。

## 维护可读测试样本

将较完整的 Go 测试程序保存为按场景组织的 `testdata/<scenario>/*.go.tmpl`。
静态样本直接读取，仅在需要替换变量时使用严格模板渲染。
该后缀使 Go 在查找包时不会将样本作为测试源码编译。

这类样本参与格式化，必须在执行模板之前就能作为 Go 源码解析。
模板动作应置于字符串或注释中，或使用合法标识符占位。
故意无效的语法和格式化输入样本应保存在 `.go.txt` 或简短内联字面量中，避免格式化改变测试条件。
