# Roadmap

[English](roadmap.md) | [简体中文](roadmap.zh.md)

## Stage 1: Source formatting and notices

The `pkg/format` and `pkg/header` libraries and CLI provide syntax-based
formatting and independent source-notice updates. They do not load application
types. Formatting and supported notice-refresh forms are deterministic and
idempotent. API, ownership and write behavior are defined in
[architecture](design/docs/architecture.md).

## Stage 2: Repository and generator integration

Repositories and generators can pin the same golayout version. A separate tools
module can isolate CLI dependencies from production packages. The shared scripts
select repository source; generators call the in-memory APIs and retain their
own validation and output policy. Configuration examples are in the
[integration guide](integration.md).

## Stage 3: Policy evolution

New formatting rules require a concrete readability need. Intentional rule or
API changes are published in a new version with their compatibility effects.
Consumers control adoption by upgrading their pinned version; the CLI and
library apply the same policy.

## Stage 4: Versioned documentation

Module tags provide pinnable versions. English release notes describe the
changes, and bilingual manuals are built from the corresponding tagged source.
Published versions and notes are listed in
[GitHub Releases](https://github.com/longxi-go/golayout/releases).
