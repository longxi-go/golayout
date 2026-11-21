# AGENTS.md

## Project Purpose and Scope

- golayout owns deterministic Go source formatting for the longxi-go ecosystem.
- Public pkg/header separately owns template-driven Go source notices shared by
  repository tooling and generators. Licensing and copyright-owner selection
  remain explicit repository policy, outside source formatting.
- Public pkg/format composes private AST layout and pinned gofumpt. The CLI uses
  that API. Never depend on gosr, ioc, nemo or application type loading.
- Keep file selection outside the library and all command logic out of cmd.
- Follow docs/architecture/design.md and docs/roadmap.md for this extraction.

## Engineering Mandates

### Forbidden

1. MVP-only delivery: presenting temporary, partial or demonstration-grade work as complete production capability.
2. Untested production behavior or placeholder implementations presented as finished work.
3. Significant implementation without a clear design, package boundaries and acceptance criteria.
4. Feature accumulation that creates oversized, tangled or difficult-to-maintain packages.
5. Mixing unrelated responsibilities or leaking implementation details across architectural boundaries.
6. Isolated feature-by-feature delivery without considering the stage roadmap and known dependencies.
7. Designing each next feature only after finishing the previous one.
8. Speculative interfaces, factories, providers, registries or design patterns without an actual architectural purpose.
9. Silent changes to public contracts, configuration, errors or lifecycle behavior without assessing compatibility.
10. False completion claims, hidden verification gaps or weakened tests/linters used to obtain a passing result.

### Expected

- Prioritize correctness, maintainability, architectural coherence, testability and necessary extensibility.
- Prefer the simplest design that fully satisfies known requirements and preserves real extension boundaries.
- Encapsulate state and implementation details. Use composition, interfaces and design patterns where they solve an identified problem.
- Keep dependencies, responsibilities, naming, ownership and side effects explicit.
- Build complete, verifiable capability slices; do not substitute superficial feature counts for usable behavior.
- Simplicity must not excuse incomplete engineering; extensibility must not excuse unnecessary abstraction.

## Framework Design Expectations

- Each package MUST have one cohesive responsibility and a minimal intentional API.
- Before creating or expanding a package, justify its responsibility and dependency direction.
- Avoid dumping grounds such as `common`, `misc`, `helpers`, `utils` or vaguely named managers.
- Define small behavioral interfaces at real shared or consumer boundaries. Do not merely mirror concrete types.
- Prefer composition over inheritance-like embedding hierarchies. Adopt useful Spring IoC ideas through Go-compatible contracts and ownership.
- Shared APIs must remain independent of unnecessary transports, storage engines, runtimes and compiler internals.
- Runtime selection, registration, activation, bootstrapping and cleanup must have explicit owners in the responsible implementation.
- Extension points must specify initialization, lifecycle, failure and cleanup semantics where applicable.
- Default implementations must permit documented alternatives; registration, selection and configuration precedence must be deterministic.
- Registration must be discoverable and documented. Generated `init` may contribute definitions; it MUST NOT construct Beans, perform external I/O or start background work.
- Prefer standard-library solutions when adequate. Justify third-party dependencies by their value, compatibility and maintenance costs.
- Foundational contracts must not depend on higher-level implementations. Avoid import cycles and unnecessary transitive dependencies.
- Do not hide trivial dependencies behind artificial abstraction layers or transplant Java class-loading/inheritance machinery into Go.

## Package Layout and Go Conventions

- Put public APIs and contracts in `pkg/<capability>`; put private implementation in `internal/<responsibility>`.
- Keep root Go code deliberately small, such as package documentation or a thin facade. Do not collect domain logic and helpers there.
- Executables belong in `cmd/<tool>` with thin entry points delegating to cohesive packages.
- Do not create empty directories to fill a template. Respect explicit deferral of existing layout changes; deleted packages are not a layout to restore.
- Every owned Go package MUST contain `doc.go` and `example_test.go`, with meaningful examples that run under `go test`.
- Use precise package names, minimal exports, consistent receiver naming, straightforward control flow and focused functions.
- Follow the Go version in `go.mod`; verify language/tool support rather than assuming older-version restrictions.
- Use generics and code generation when they improve static safety or wiring. Keep their complexity contained; neither is prohibited merely as metaprogramming.
- Format owned source with the pinned gofumpt version through `make format`; gofmt alone is not the formatting gate. Generated Go must meet the same compilation, formatting and quality expectations as handwritten Go.
- Handle errors explicitly; preserve causes with `%w` where appropriate and classify with `errors.Is`/`errors.As`.
- Do not discard meaningful errors, rely on error strings for control flow or use panic for ordinary operational failure.
- Document cleanup failures and provide useful diagnostics without exposing sensitive data.
- Pass `context.Context` through context-aware operations; respect cancellation and deadlines. Do not retain request contexts in long-lived objects.
- Define ownership of goroutines, connections, timers and other resources. Release them on success, failure, cancellation and shutdown as applicable.
- Minimize mutable globals, hidden dependencies and initialization side effects.
- Make concurrency safety and limits explicit. Use synchronization intentionally; avoid races, leaks, unbounded work and concurrency without a useful purpose.

## Go Code Layout and Readability

- Keep each `.go` file focused on one cohesive responsibility. Split distinct
  responsibilities into accurately named files within the same package; package
  visibility does not require putting related capabilities into one large file.
- Compiler `pipeline.go` owns orchestration only. Keep invocation state, step
  execution, each individual `*_step.go`, semantic modeling, rendering and output
  persistence in their respective files. Do not collect step implementations in
  the orchestration file or introduce an unnecessary package for each step.
- Nonempty function/method bodies and control-flow blocks MUST use multiple lines,
  including forwarding methods and generated converter closures.
- Empty functions, control-flow blocks, structs, interfaces and literals may use
  single-line Go forms, including `select {}`, `struct{}` and `interface{}`.
  Do not add artificial comments or force line breaks solely to expand empty forms.
- Separate top-level declarations with a blank line, keeping doc comments attached
  to their declaration. Separate distinct types inside grouped declarations too.
- Every named, non-empty struct declaration MUST use multiple lines, even when
  it contains only one field. Empty structs retain the canonical `struct{}` form
  required by gofmt/gofumpt; do not fight the formatter with artificial comments.
- Non-empty struct literals with multiple fields MUST use multiple lines, with
  one field per line. Prefer that layout for single-field literals when it makes
  construction intent easier to scan.
- Non-trivial map and slice literals MUST use multiple lines when that improves
  readability, especially for nested values or collections of definitions.
- Avoid excessively long lines in function declarations/calls, composite
  literals, boolean expressions and nested expressions. Split at logical
  boundaries rather than relying on a fixed-width formatter.
- Prefer splitting long or multi-step fluent/method chains at call boundaries,
  usually one chained call per line, when it makes the sequence easier to follow.
  Short, clear chains may stay on one line. This is readability advice, not a
  mandatory line-width or call-count rule and not an automatic formatter gate.
  Keep the continuation dot on the preceding line in Go; consider named
  intermediate values when they read better. Apply the same advice to tests
  and generator templates, without mechanically splitting every selector chain.
- Prefer explicit, readable layouts over compact expressions. Do not collapse
  deliberately structured multiline code into less readable single-line forms.
- The public `pkg/format` API is the one supported formatter for consumer generators. The CLI delegates to that same API; never duplicate either pipeline.
- Ecosystem repositories use the CLI from their pinned tool module for source
  files and `pkg/format` for in-memory generated output.
- Empty structs, interfaces and literals retain canonical compact forms. A
  standalone gofumpt invocation, including `-extra`, does not enforce this layout.
- File responsibilities and logical line breaks remain source-design/review
  requirements; a formatter cannot infer appropriate architectural boundaries.
- Apply the same readability requirements to generated Go. Templates own
  structural layout; the formatting stage runs before validation/output.
- After changing Go source, run `make format`, then `make check-format`.

For a chain whose setup sequence is easier to scan across lines, prefer:

```go
t, err := template.New("codegen").
    Option("missingkey=error").
    Funcs(functions).
    ParseFS(source, patterns...)
```

Avoid:

```go
type Config struct{ Port int }
config := Config{Port: 8080}
```

Prefer:

```go
type Config struct {
    Port int
}
config := Config{
    Port: 8080,
}
```

## Multi-Repository Development

- Use a parent `go.work` for local ecosystem development. Include gosr/tools as
  a separate module; parent module selection does not automatically include it.
- Production `go.mod` files declare dependencies and MUST NOT contain local
  sibling-directory replacements. Temporary consumer fixtures may retain explicit
  local replacements to exercise isolated module behavior.
- Keep go.work, go.work.sum, execution state, caches, backups, OpenSpec and
  superpowers planning artifacts outside the publishable change set.
- Unpublished dependency placeholders are for workspace development only. Verify
  standalone module resolution after actual dependency versions are published;
  do not treat a workspace-only pass as proof of standalone release readiness.

## Design and Implementation Workflow

- Before a meaningful change, inspect current code, dependencies, relevant documents, extension seams and compatibility constraints.
- Current user instructions and current source are authoritative; stale documents and cached summaries do not override them.
- Establish the problem, scope/non-goals, architectural placement, public contracts, ownership, failure behavior, testing strategy and acceptance criteria.
- Consider relationships between known capabilities before implementation. Do not speculate beyond established needs.
- Keep design effort proportional; small fixes do not need unnecessary ceremony.
- Use existing suitable abstractions, but do not copy a demonstrably unsuitable pattern or retain a flawed design merely because it already exists.
- Implement against the agreed design with cohesive, reviewable changes. Keep code, tests and documents aligned.
- Reassess a flawed abstraction when evidence warrants it; do not continually patch around architectural defects.
- Distinguish required capability work, necessary supporting work, related optional work and unrelated changes.
- Do not add speculative features or omit essential correctness, testing or integration to reduce scope.
- Evaluate downstream effects of public API, configuration, error and lifecycle changes. Document intentional breaks and migration implications.
- Breaking changes already authorized for this rewrite do not need repeated permission; do not resurrect obsolete APIs or add compatibility shims without a concrete requirement.

## Testing Expectations

- Every production package MUST be covered by focused tests. Contract-only packages require meaningful usage, semantic or compile-time examples rather than a production runtime built just for testing.
- New or modified behavior must have automated tests covering ordinary behavior, edge cases and meaningful failure paths.
- Prefer behavioral, deterministic and isolated tests; use table-driven cases when they improve clarity.
- Do not mirror implementation details, rely on incidental order/globals or add meaningless tests for coverage.
- Test new extension points with default and alternative implementations, registration/selection, lifecycle and relevant invalid-use cases.
- Integration tests must verify interactions owned by this project: contracts in gosr; boot, runtime wiring, generation and cross-module composition in ioc.
- Configuration precedence and external integration tests are required when those capabilities are in the authorized scope, not merely because a generic checklist mentions them.
- Practical bug fixes SHOULD include a regression reproducing the original failure and asserting the corrected behavior.
- Use TDD for new behavior: observe the failing behavioral test before implementing its solution.
- Run race detection for applicable concurrent/shared-state changes; meet configured coverage thresholds without treating percentages as proof of quality.
- A conformance suite must name its version and exact coverage. Do not equate a limited catalog, test fixture or compiling interface with complete conformance.

## Test Go Source Fixtures

- Prefer per-scenario `testdata/<scenario>/*.go.tmpl` files for substantial Go
  programs compiled, generated or inspected by tests. Keep scenario names and
  expected diagnostics in the test table; keep declarations and runtime checks
  readable as normal Go source in the fixture.
- Read static fixtures directly. Use the existing strict template renderer only
  when substitutions are needed; keep view models named and substitutions small.
  Do not add a general fixture framework for simple file reads.
- Complete `.go.tmpl` files under `testdata` opt into the shared golayout
  `make format` / `make check-format` gates. They must parse as Go before template
  execution; keep template actions inside quoted strings or comments when possible.
- Short, clear inline snippets may remain. Intentionally malformed syntax or
  compact inputs that test formatting must retain their exact shape; use clearly
  named `.go.txt` files or brief inline samples rather than opting them into
  formatting. Never normalize away the behavior a negative test must exercise.
- These are readability recommendations, not a requirement to externalize every
  literal. Preserve annotations, build constraints, tags and test assertions when
  moving source into fixtures.

## Documentation Standards

- Prefer diagrams before detailed prose when explaining architecture, dependency
  directions, workflows, states or interactions. Use Mermaid flowcharts, sequence
  diagrams or state diagrams as appropriate; do not force diagrams onto simple facts.
- Keep diagram labels and technical meaning aligned in English and Simplified
  Chinese. Distinguish existing mechanisms, normative obligations and planned work.
  Retain ownership, failure, concurrency and compatibility details in concise prose.
- Validate diagram rendering as well as Markdown builds; a successful static build
  does not prove Mermaid syntax or visual readability.
- Every package MUST have `doc.go` and `example_test.go`; this is not limited to core packages or optional discoverability improvements.
- Package documentation must explain purpose, responsibility, intended use, constraints and extension boundaries.
- Exported contracts, types and important functions MUST have meaningful Go doc comments.
- Source comments, package documentation and API comments MUST be in English.
- Explain why an abstraction exists, its invariants, design choices and prevented misuse. Do not restate source line by line.
- Describe relevant ownership, lifecycle, concurrency, cancellation, failure and compatibility semantics.
- Keep examples compilable, executable and representative of supported behavior. Do not add empty examples merely to satisfy filenames.
- Update documentation with contract changes. Keep comments precise, necessary and consistent with actual behavior.
- Source-level documentation must make architectural seams understandable without requiring every external design document.

## Linting and Quality Gates

- `codequality/HEADER` is the authoritative plain-text Go source-file notice.
  Preserve its copyright owner and year policy; do not infer ownership from Git.
- Use `make license-headers` to refresh owned `.go` files and
  `make check-license-headers` for the read-only gate included in `make check`.
- Shared notice logic lives in `golayout/pkg/header`; generators apply the output
  owner's optional template before formatting/validation, keeping generated markers
  first. Do not duplicate notice-replacement logic in consumers.
- Exclude ignored files, backups, caches, vendor, testdata and symlinks from bulk
  refreshes. A different existing copyright owner requires an explicit decision.

- Pinned gofumpt `v0.12.0` is the final stage; optional `-extra` is disabled.
  Rule/version changes require intentional review and a golayout version change.
- `make setup` MUST verify tool dependencies resolve without installing binaries
  or modifying module files. `make check-format` MUST be read-only.
- Consumer CI MUST check out the immutable golayout revision used by go.mod/tools
  before running setup, formatting checks or generated-code tests.
- Format/check only Git-owned or unignored Go source, including generated files explicitly. Exclude `.tmp`, `.cache`, dependency/vendor trees and ordinary intentional `testdata` samples; opt complete `testdata/**/*.go.tmpl` fixtures into the shared gate. Never recurse blindly over `.`.
- `golangci-lint` is mandatory for code delivery. Run `golangci-lint run --timeout=10m` with zero issues before declaring code complete.
- Satisfy all enabled linters, including `errcheck`, `staticcheck`, `unused` and applicable related checks.
- Do not suppress findings merely to pass validation. Any permitted suppression needs a specific justified reason.
- Migrate deprecated APIs/patterns when an appropriate replacement exists. Retention needs a documented intentional compatibility reason.
- Resolve introduced lint/test failures in the same capability slice. Investigate existing failures and disclose them accurately.
- For code changes, verify formatting of owned Go files, run `go test ./...`, `go vet ./...`, required race/integration/coverage checks and lint.
- Run `make check` when the project defines that target; satisfy any additional applicable project gates.
- Limit formatting/file scans to owned source; ignored snapshots and dependency caches are not current packages.
- Inspect actual output and exit status. Report skipped, unavailable and failed checks with reasons; do not make an unqualified completion claim while required gates fail.
- Do not weaken tests, disable linters or remove gates to obtain success.
- Instruction/documentation-only tasks require diff, requirement and reference validation; they do not require unrelated Go builds/tests.
- Never restore deleted source to make a gate run or claim that an empty source tree proves implementation correctness.

## Delivery Standard

- Design, implementation, tests, examples and documentation must evolve together.
- Applicable OpenSpec artifacts must stay consistent with current code and requirements.
- Satisfy acceptance criteria with coherent boundaries, intentional APIs and documented compatibility implications.
- Do not leave placeholders, unexplained critical TODOs or unresolved material design conflicts presented as complete.
- Verify error handling, resource/lifecycle ownership and applicable concurrency behavior.
- Required unit, integration, extension, regression, race, coverage and lint checks must finish and their results must be inspected.
- Keep source easy to read, extend and replace in parts; do not introduce unrelated functionality or complexity.
- Completion requires current evidence for the delivered revision, not optimism, successful compilation alone or an agent's report.

## Roadmap Governance

- A stage-based roadmap MUST exist before large-scale iteration continues.
- Define each stage's capability bands, dependencies, usable outcome and acceptance criteria.
- OpenSpec changes MUST derive from the roadmap and group work by stage/capability, not isolated feature impulses.
- Complete cohesive stage-capability slices; do not invent one disconnected feature after another.
- If stages are marked complete but their promised usable outcome is absent, the roadmap/decomposition is wrong and MUST be corrected before further implementation.
- For gosr, usability means independently understandable and implementable contracts; for ioc, a working managed/generated IoC path.
- Roadmap updates are required development work. They do not authorize speculative implementation or restoring deleted artifacts during a rules-only task.

## Execution, Verification and Completion

- Work only within current authorization. A rules review does not authorize application code, scaffolding, reference engines, new specs or another repository's implementation.
- Preserve user changes, staging and deletions. Inspect current files before modifying them; never restore rejected designs to simplify a task.
- Continue authorized implementation through verification and completion. Do not substitute a plan, partial handoff or tool-recovery exercise for the deliverable.
- A failed tool/test, unfamiliar package or large task is not by itself a blocker. Inspect evidence, adjust the approach, retry when justified and continue.
- Do not repeat unchanged failures or busy-poll. A real blocker requires missing information, authority or an external dependency with no useful independent work remaining.
- Resolve genuine material ambiguities before dependent work; do not repeat approval questions for settled requirements.
- For permitted parallel work, assign clear path ownership and verification responsibility. An agent report does not establish integrated correctness.
- Keep changes/commits scoped. Do not stage unrelated user deletions, force history changes or publish externally as part of a rules edit.
- Before completion, inspect the intended diff and map multi-area acceptance requirements to current implementation/artifacts and actual verification.
- Reassess verification after relevant changes; do not present stale passes as current evidence.
- Report actual changes, architectural impact, tests, commands/results, documentation and material remaining issues in proportion to the task.
- Every final report must include `Completed`, `Verified`, `Not Verified` and `Known Issues` (or the user's language equivalents). Use None for empty fields.
- Distinguish verified, unverified, partial, unimplemented and genuinely blocked work. State blockers, impact and skipped checks; never invent results or silently lower acceptance criteria.

## Task State and Instruction Protection

- The user deleted earlier `.codex/*` state. Do not restore its content, resume its plans or infer authorization from former records.
- Do not create task-state files for this rules-only edit.
- For future authorized long/multi-phase work, required task-state tracking remains in force. Build fresh records from current scope and evidence.
- Fresh state must identify the task/workspace and responsible role, plus Objective, Current Phase, Completed, In Progress, Remaining, Changed Files, Verification, Decisions and Known Issues.
- Record branch/baseline, relevant verification scope/results/gaps, the first unfinished action and outstanding asynchronous work.
- Update state at meaningful milestones and before compaction/handoff. After resumption, reconcile it with current files/Git evidence before continuing.
- Concurrent tasks require separate records and one writer per shared index; state files do not isolate overlapping code writes.
- Preserve other applicable continuation, verification, completion and configuration safeguards. Do not remove state or weaken safeguards without an explicit user request.
- When editing rules, compare requirements section by section against the originals; retaining headings alone is insufficient.

## Final Engineering Standard

- Architecture before accumulation; responsibilities before abstraction; correctness before convenience.
- Use idiomatic Go and explicit contracts while preserving justified static generation and Spring-inspired IoC semantics.
- Deliver verified capabilities, not feature collections or optimistic reports.
- Maintainability, complete behavior and proportionate process are all required.
