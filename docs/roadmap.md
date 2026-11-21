# Roadmap

[English](roadmap.md) | [简体中文](roadmap.zh.md)

## Stage 1: Shared source tooling

Available: `pkg/format`, `pkg/header`, the CLI, focused tests and executable
examples. Quality checks cover Go 1.27 compatibility, formatting, source notices,
module dependencies, vet, unit tests, race tests and lint.

Repositories and generators share this implementation without loading application
types. Formatting is deterministic and idempotent; source notices are maintained
independently.

## Stage 2: Consumer release adoption

Use one pinned version across consumer source gates and generation steps.
Keep CLI tooling outside dependency-free contracts. Verify each consumer's
published dependencies with `GOWORK=off` before its release.

A consumer completes this stage when it builds independently and applies the
shared policy. Verify adoption and release readiness per repository; a local
workspace build is insufficient.

## Stage 3: Policy maintenance

Each new rule needs a concrete readability requirement and regression coverage.
Version intentional policy changes, test representative consumer fixtures and
document compatibility effects so changes can be reviewed and consumers can
choose when to upgrade.

## Stage 4: Release and documentation delivery

Available: stable-tag validation, maintained English notes, reusable tagged-source
quality checks, GitHub Releases and deployment of bilingual documentation from
the published tag.

Consumers can pin published module versions, and Pages describes released
behavior. Every release verifies its own commit, notes and deployed site.
