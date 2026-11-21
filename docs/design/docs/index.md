# One Go layout policy

**golayout** combines syntax-aware structural layout with pinned gofumpt. The
library formats source bytes without loading application types; the CLI checks
explicit files and rewrites them only when requested.

[Get started](quick-start.md){ .md-button .md-button--primary }
[Read the architecture](architecture.md){ .md-button }

| Boundary | Responsibility |
| --- | --- |
| `pkg/format` | Public in-memory Source API and target language options |
| `internal/layout` | AST-based multiline bodies and declaration spacing |
| `internal/command` | Explicit file validation, checking and replacement |
| `cmd/golayout` | Thin command entry point |
| Consumer Makefile/scripts | File ownership, discovery and exclusions |

Nonempty structs/interfaces, function bodies and struct-like literals expand to
readable multiline layouts. Empty functions, blocks, `struct{}` / `interface{}` and simple scalar
collections retain canonical forms. The formatter does not choose architectural
package boundaries or split responsibilities into files.

gosr, ioc and nemo consume this implementation. Generators use the library; source
gates use the CLI. Copying layout rules into consumers would let policies drift.
