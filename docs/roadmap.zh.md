# 路线图

[English](roadmap.md) | [简体中文](roadmap.zh.md)

## 阶段一：源码格式化与文件头

`pkg/format`、`pkg/header` 库及 CLI 提供基于语法的格式化和独立的源码文件头更新，
无需加载应用类型。格式化及支持的文件头刷新形式具有确定性和幂等性。
API、所有权及写入行为见[架构说明](design/docs/architecture.zh.md)。

## 阶段二：仓库与生成器接入

仓库和生成器可固定同一个 golayout 版本，通过独立 tools 模块将 CLI 依赖与生产包隔离。
共享脚本负责仓库源码选择；生成器调用内存 API，并保留自身的验证和输出策略。
配置示例见[接入指南](integration.zh.md)。

## 阶段三：策略演进

新增格式规则需要具体的可读性需求。有意调整规则或 API 时，通过新版本发布并说明兼容性影响。
使用方通过升级固定版本控制接入时机，CLI 与库应用相同的策略。

## 阶段四：版本化文档

模块 tag 提供可固定的版本，英文版本说明记录变更，双语手册从对应 tag 的源码构建。
已发布版本及其说明见[GitHub Releases](https://github.com/longxi-go/golayout/releases)。
