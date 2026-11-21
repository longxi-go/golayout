# Integrating golayout

[English](integration.md) | [简体中文](integration.zh.md)

## Pin the tool and library

Use an immutable published version for both source checks and generated output.
In the module that owns the tool dependency:

```sh
go get -tool github.com/longxi-go/golayout/cmd/golayout@v0.1.4
```

Generators import the library from the same module version:

```sh
go get github.com/longxi-go/golayout/pkg/format@v0.1.4
```

If your production module must remain free of tool dependencies, create a
separate module for the CLI. From the repository root:

```sh
mkdir -p tools
cd tools
go mod init example.com/project/tools
go get -tool github.com/longxi-go/golayout/cmd/golayout@v0.1.4
```

Replace `example.com/project/tools` with your module path. Commit its `go.mod`
and `go.sum`; generators declare the library dependency in their own module.
Keep developer-specific directory replacements out of production `go.mod`.

## Delegate repository checks

The following Makefile fragment uses shared scripts from the selected module.
Run it at the consuming repository's root. Set `GOLAYOUT_TOOL_MODULE=tools` when
the tool is declared in a nested tools module:

```makefile
GO ?= go
GOLAYOUT_TOOL_MODULE ?= .
GOLAYOUT_DIR := $(shell $(GO) -C "$(GOLAYOUT_TOOL_MODULE)" list -m \
    -f '{{.Dir}}' github.com/longxi-go/golayout)
GO_LANG_VERSION := go$(shell awk '$$1 == "go" { print $$2; exit }' go.mod)

.PHONY: setup format check-format

setup:
	@GO="$(GO)" bash "$(GOLAYOUT_DIR)/scripts/format.sh" setup \
	    "$(CURDIR)" "$(GOLAYOUT_TOOL_MODULE)" "$(GO_LANG_VERSION)"

format:
	@GO="$(GO)" bash "$(GOLAYOUT_DIR)/scripts/format.sh" write \
	    "$(CURDIR)" "$(GOLAYOUT_TOOL_MODULE)" "$(GO_LANG_VERSION)"

check-format:
	@GO="$(GO)" bash "$(GOLAYOUT_DIR)/scripts/format.sh" check \
	    "$(CURDIR)" "$(GOLAYOUT_TOOL_MODULE)" "$(GO_LANG_VERSION)"
```

The script selects Git-tracked and eligible untracked Go files, excludes `.tmp`,
`.cache`, vendor files and symlinks, and includes complete `testdata/**/*.go.tmpl`
fixtures. Other testdata files and production generator templates stay excluded.
`setup` verifies dependencies without installing a binary; `check-format` does
not rewrite source or module files. These commands require Git, Bash, Make and
the Go version required by the selected golayout release.

The language option comes from the consumer's `go.mod`, not the formatter's
build toolchain. Use the same language version when formatting generated source.

## CI and local development

CI resolves the golayout version declared by the tool module through Go's
checksum-verified module mechanism. No sibling checkout is required. From the
consuming repository root, download dependencies before resolving script paths:

```sh
GOWORK=off go mod download
GOWORK=off make setup check-format
```

For a nested tools module, use:

```sh
GOWORK=off go -C tools mod download
GOWORK=off make setup check-format GOLAYOUT_TOOL_MODULE=tools
```

Include the tool module's `go.mod` and `go.sum` in CI cache keys. Validate generated
output through the generator's own checks as well.

A local `go.work` may select a different golayout checkout. Include nested tool
modules explicitly when using a workspace. `GOWORK=off` makes the commands above
use the dependencies declared in the tool module instead.

## Format generated output

Call `pkg/format.Source` before the generator's type validation and file write.
Pass source bytes and the target module's language version. Handle errors before
using the result. Do not invoke the CLI as a generation subprocess or copy
private formatting code into a consumer.

For source notices, load the target repository's template explicitly and call
`pkg/header.Source` before formatting. The caller decides whether a template is
required and how an absent template is handled; golayout does not select ownership.
Repository-wide header scripts read the consumer's `codequality/HEADER` and exclude
all testdata. Their check and write modes are separate from formatting.

## Maintain readable fixtures

Store substantial Go test programs as per-scenario `testdata/<scenario>/*.go.tmpl`
files. Read static fixtures directly; use a strict template renderer only when
substitutions are required. The suffix prevents Go package discovery from compiling
the fixture as host test code.

These files opt into formatting and must parse as Go before template execution.
Place template actions inside strings or comments, or use valid identifier
placeholders. Keep deliberately invalid syntax and formatter-input samples in
`.go.txt` files or short inline literals so formatting cannot erase the test case.
