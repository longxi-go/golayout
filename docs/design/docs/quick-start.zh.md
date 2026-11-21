# 快速开始

使用 go.mod 指定的 Go 版本，当前最低为 Go 1.27.0。在本仓库运行：

```sh
make setup
make format
make check-format
make check
```

## 检查或重写显式文件

v0.1.0 发布后可用 `go install github.com/longxi-go/golayout/cmd/golayout@v0.1.0` 安装，
或在工具模块中运行 `go get -tool github.com/longxi-go/golayout/cmd/golayout@v0.1.0`，
通过 `go tool golayout` 调用。在本仓库中可直接运行：

```sh
go run ./cmd/golayout -lang go1.27 path/to/file.go
go run ./cmd/golayout -w -lang go1.27 path/to/file.go
```

第一条命令只读检查，存在格式差异时以非零状态退出；`-w` 才执行写入。
命令至少接收一个显式文件，仓库文件发现仍由调用方负责。

## 接入生成器

导入 `github.com/longxi-go/golayout/pkg/format`：

```go
formatted, err := format.Source(source, format.Options{
    LangVersion: "go1.27",
})
```

此处是代码片段，`source` 为生成器拥有的字节切片。处理错误后再验证或写出 `formatted`，
输入字节不会被修改。实际目标语言应从使用方 go.mod 推导；示例字面量仅展示选项，
不能用生成器自身工具链版本代替目标应用版本。

公开 API 应在类型验证及输出前调用。源码检查使用同一模块版本，并保留使用方自己的文件排除策略。
golayout 不需要其他同级模块，生态联调可参考可选的[本地工作区说明](development.md#local-workspace)。
