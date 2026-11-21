# Integrating golayout

Each consumer declares golayout in go.mod (or its nested tools module). Local
multi-repository development selects the sibling module through a parent go.work;
do not put developer-specific local replace directives in production go.mod files.
Use the published `v0.1.0` tag after its release; before publication only local
workspace selection can resolve it. Standalone consumers require the published
tag and matching go.sum checksums. golayout itself builds independently of the
ecosystem workspace. Never move a published tag to change a formatting policy.

```go
require github.com/longxi-go/golayout v0.1.0

tool github.com/longxi-go/golayout/cmd/golayout
```

Go source targets the module language version, such as `go1.27`. A Makefile
should pass an explicit inventory of tracked and unignored Go files, excluding
snapshots, caches, vendor and ordinary testdata. Complete parseable Go sources
under `testdata/**/*.go.tmpl` explicitly opt into the same format/check gates.
Production generator templates outside testdata stay excluded. Delegate
format/check/setup to scripts in the selected golayout module. Consumer Makefiles
resolve its directory with go list -m, supporting workspace and module-cache sources. The check path must not install
consumer-local binaries or rewrite source/module files.

Generators call `github.com/longxi-go/golayout/pkg/format.Source` on output
bytes before validation and persistence. They must not run the CLI as a
subprocess or import a consumer's private formatter implementation.

CI checks out golayout at the same immutable commit/tag as the consumer's
module requirement. For a sibling checkout, place it alongside the consumer
repository and add its `go.mod` to the Go cache dependency paths. If golayout
is private, configure a read token that can access the shared repository.

For gosr, keep formatter tooling in the nested `tools` module. This prevents
the standards/contracts module from depending on developer tooling.

## Readable Go test fixtures

Keep substantial test programs in per-scenario `testdata/.../*.go.tmpl` files,
with normal imports, declarations, comments and multiline blocks. Tests read
static files directly; use strict template rendering only for actual variables.
Go ignores testdata during package discovery, and the .tmpl suffix keeps these
sources from being compiled as part of the host test package.

The shared selector includes this opt-in extension in both root and nested
testdata directories. Fixtures must parse as Go before substitutions; use template
actions inside import strings/comments or explicit valid identifier placeholders
when necessary. Intentionally malformed syntax and formatter-input samples should
keep a different extension (for example .go.txt) or a short inline literal, so a
format pass cannot erase the condition the test is meant to exercise. Neither
production generator templates nor every file in testdata is auto-formatted.
