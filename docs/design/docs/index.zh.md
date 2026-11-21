# 统一 Go 代码布局策略

**golayout** 将基于语法的结构化布局与固定版本 gofumpt 组合。
库直接格式化源码字节，不加载应用类型；CLI 检查显式文件，仅在请求写入时修改文件。

[开始使用](quick-start.md){ .md-button .md-button--primary }
[阅读架构](architecture.md){ .md-button }

| 边界 | 职责 |
| --- | --- |
| `pkg/format` | 公开的内存 Source API 与目标语言选项 |
| `internal/layout` | 基于 AST 的多行代码块和声明间距 |
| `internal/command` | 显式文件校验、检查与替换 |
| `cmd/golayout` | 轻量命令入口 |
| 使用方 Makefile/脚本 | 文件所有权、发现和排除策略 |

非空 struct/interface、函数体和结构体形式字面量展开为多行，空函数、空代码块、`struct{}` / `interface{}`
及简单标量集合保留规范形式。格式化工具不负责决定包边界或按职责拆分文件。

gosr、ioc 和 nemo 共同使用此实现：生成器调用库，源码门禁调用 CLI，避免在各工程复制规则导致漂移。
