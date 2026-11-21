# Integrating golayout

[English](integration.md) | [简体中文](integration.zh.md)

## Pin the tool and library

Use an immutable published version for both source checks and generated output.
In the module that owns the tool dependency:

```sh
go get -tool github.com/longxi-go/golayout/cmd/golayout@v0.1.1
```

Generators import `github.com/longxi-go/golayout/pkg/format`; they must declare
the same module version. Keep developer-specific directory replacements out of
production `go.mod`. A parent `go.work` can select a local checkout during development.

For dependency-free contract modules such as gosr, put CLI tooling in a separate
`tools` module. Include that module explicitly in local workspaces.

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
`setup` verifies dependencies; `check-format` does not rewrite source or module files.

The language option comes from the consumer's `go.mod`, not the formatter's
build toolchain. Keep generated-source formatting on the same version. CI must
check out shared tooling at the immutable revision required by the consumer;
include both modules' dependency files in cache keys when using sibling checkouts.

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
