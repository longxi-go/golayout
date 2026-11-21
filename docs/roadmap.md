# Roadmap

## Stage 1: Shared structural formatter

Deliver the pure pkg/format library, private AST editor, source-matched CLI,
tests, executable examples and Go 1.27 quality gates. Outcome: one deterministic,
type-loading-free formatting implementation that preserves comments, rejects
invalid syntax and produces idempotent output.

## Stage 2: Ecosystem adoption

Migrate gosr, ioc and nemo Make targets and generation steps to the shared
golayout module. Keep gosr's contract module free of tool dependencies. Outcome:
all ecosystem projects use one policy and CI checks the same immutable revision.

## Stage 3: Maintained shared policy

Version intentional formatting changes, run compatibility against representative
consumer fixtures, and document upgrades. Add no rule without a concrete
readability requirement and regression coverage.

## Stage 4: Standalone module releases

Validate stable tags and English release notes, rerun standalone quality and
bilingual documentation checks on tagged source, publish GitHub Releases, then
deploy documentation from the same published release. Outcome: consumers can
pin a real module tag without local workspace substitutions; Pages describes
released behavior. Initial v0.1.0 publication remains pending owner confirmation.
