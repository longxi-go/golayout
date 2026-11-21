# Go 格式化与源码文件头

**golayout** 为 Go 仓库和代码生成器提供统一的格式化库、源码文件头库及 CLI。
要求 Go 1.27.0 或更高版本。

[开始使用](quick-start.md){ .md-button .md-button--primary }
[阅读架构](architecture.md){ .md-button }

| 组件 | 职责 |
| --- | --- |
| `pkg/format` | 使用 AST 布局规则与 gofumpt v0.12.0 格式化源码字节 |
| `pkg/header` | 应用显式文件头模板，不格式化代码 |
| `internal/layout` | 实现结构规则，不加载应用类型 |
| `internal/command` | 校验显式文件、报告差异并写入替换内容 |
| `cmd/golayout` | 提供 CLI 入口 |
| 仓库脚本 | 选择项目拥有的文件并应用排除规则 |

非空函数体、结构体和接口声明、控制流代码块采用多行布局，空形式允许保留紧凑形式。
复合字面量按语法分类处理，详见[格式规则](architecture.md#formatting-rules)。
文件职责划分和人工断行仍属于设计与审查范围。

生成器调用公开库；仓库检查使用 CLI 和共享脚本。两条路径应固定同一个 golayout 版本。
本模块可以独立使用，不依赖 gosr、ioc 或 nemo。
