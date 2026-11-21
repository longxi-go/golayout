# Roadmap

[English](roadmap.md) | [简体中文](roadmap.zh.md)

## Stage 1: Shared source tooling

Available: `pkg/format`, `pkg/header`, the CLI, focused tests and executable
examples. Quality checks cover Go 1.27 compatibility, formatting, source notices,
module dependencies, vet, unit tests, race tests and lint.

Outcome: repositories and generators share one implementation without loading
application types. Formatting is deterministic and idempotent; source notices
are maintained independently.

## Stage 2: Consumer release adoption

Use one pinned version across consumer source gates and generation steps.
Keep CLI tooling outside dependency-free contracts. Verify each consumer's
published dependencies with `GOWORK=off` before its release.

Outcome: each consumer builds independently and applies the same shared policy.
Adoption and release readiness are verified per repository; a local workspace
build is insufficient.

## Stage 3: Policy maintenance

Require a concrete readability requirement and regression coverage for each new
rule. Version intentional policy changes, test representative consumer fixtures
and document compatibility effects.

Outcome: rule changes remain reviewable and consumers control when they upgrade.

## Stage 4: Release and documentation delivery

Available: stable-tag validation, maintained English notes, reusable tagged-source
quality checks, GitHub Releases and deployment of bilingual documentation from
the published tag.

Outcome: consumers can pin published module versions, and Pages describes released
behavior. Every release verifies its own commit, notes and deployed site.
