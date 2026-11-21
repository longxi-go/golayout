# 快速开始

要求 Go 1.27.0 或更高版本。

## 安装 CLI

```sh
go install github.com/longxi-go/golayout/cmd/golayout@v0.1.4
```

将 `GOBIN` 加入 `PATH`；未设置 `GOBIN` 时，使用 `$(go env GOPATH)/bin`。
也可以在模块内将 CLI 声明为模块工具：

```sh
go get -tool github.com/longxi-go/golayout/cmd/golayout@v0.1.4
go tool golayout -lang go1.27 path/to/file.go
```

## 检查或重写文件

使用 CLI 时，需指定一个或多个文件路径：

```sh
golayout -lang go1.27 path/to/file.go
golayout -w -lang go1.27 path/to/file.go
```

第一条命令只读检查。存在差异或输入无效时返回非零退出状态；`-w` 才写入修改。
不接受目录或符号链接。所有输入均在开始写入前完成校验；后续文件写入失败时，
先前文件可能已经替换。详见[文件写入行为](architecture.md#cli-and-file-writes)。

## 使用库

在已有 Go 模块中添加格式化包：

```sh
go get github.com/longxi-go/golayout/pkg/format@v0.1.4
```

将以下程序保存为 `main.go`，执行 `go run .`：

```go
package main

import (
	"fmt"
	"log"

	"github.com/longxi-go/golayout/pkg/format"
)

func main() {
	source := []byte("package sample\nfunc Ready() bool { return true }\n")
	formatted, err := format.Source(source, format.Options{
		LangVersion: "go1.27",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(formatted))
}
```

输出中的 `Ready` 函数采用多行布局。`Source` 不修改输入字节。
生成器应从目标模块的 `go.mod` 推导 `LangVersion`，处理错误后再验证并保存结果。

## 维护源码文件头

创建纯文本文件 `codequality/HEADER`，写入项目的版权主体和许可证声明，例如：

```text
Copyright © 2026 the example authors.

SPDX-License-Identifier: Apache-2.0
```

第一行必须是版权行，不带 Go 注释标记。通过以下命令检查或应用文件头：

```sh
golayout -header codequality/HEADER path/to/file.go
golayout -w -header codequality/HEADER path/to/file.go
```

默认只读检查。文件头模式保留代码体，拒绝覆盖其他版权主体的声明。
它不能与 `-lang` 合用，格式化需要单独执行。生成器可在格式化前调用 `pkg/header.Source`。

支持的文件头形式、模板校验、失败行为及语法格式化的限制见 [API 约定](architecture.md#in-memory-apis)。
