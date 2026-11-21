# Architecture

## In-memory composition

`pkg/format.Source` validates the optional language version, calls private
structural layout and then pinned gofumpt. Calls are independent and do not
mutate the input. Invalid input returns no formatted output.

```text
source bytes → AST structural layout → pinned gofumpt → formatted bytes
```

Structural rules expand nonempty executable blocks, nonempty declarations and struct-like
literals, and separate declarations. Syntax classification avoids turning every
map/slice into a multiline structure. Formatting is deterministic and idempotent;
it does not resolve imports through application package loading.
Empty forms may remain on one line, including `select {}`, `struct{}`, `interface{}`
and empty functions; no artificial statements or comments are added to expand them.

`pkg/header.Source` is a separate operation for source notices. It uses an explicit
plain-text template, preserves generated markers/build constraints/package docs,
and rejects replacing a different copyright owner. Header mode does not format
the code body; generators apply it before their formatting and validation steps.

## CLI and file ownership

`cmd/golayout` delegates to `internal/command`. The command prepares and validates
explicit files before writing. Its read-only mode reports changed paths and
returns a failure status; write mode applies replacements.

The API does not know Git, vendor directories, generated-file ownership or test
fixture policy. Shared scripts discover eligible source for repository gates;
consumers select their root/tool module and exclude caches and intentional fixtures.
Complete Go `.go.tmpl` fixtures can opt into the same policy.

## Dependency direction

golayout depends on pinned gofumpt, not on gosr, ioc, nemo or application source.
Generators call `pkg/format`; Makefiles delegate to the shared CLI/scripts.
One implementation owns policy changes, so source and generated output cannot
silently follow different rule copies.

Upstream gopls/gofumpt settings alone do not enforce these additional structural
rules. File responsibilities and sensible wrapping remain source-design decisions.
