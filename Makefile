.DEFAULT_GOAL := help

GO ?= go
UV ?= uv
PYTHON ?= python3
GOLANGCI_LINT ?= golangci-lint
GOLAYOUT_REPOSITORY ?= .
GOLAYOUT_TOOL_MODULE ?= .
GO_LANG_VERSION ?= go$(shell awk '$$1 == "go" { print $$2; exit }' go.mod)

.PHONY: help setup format check-format test-format test race vet lint verify-mod check check-all build

.PHONY: docs-setup docs-serve docs-build
.PHONY: license-headers check-license-headers test-license-headers
.PHONY: test-release validate-release

help:
	@printf '%s\n' 'Usage: make <target>' '' \
	  '  help              Show this target list (default)' \
	  '  license-headers   Refresh owned Go source license headers' \
	  '  check-license-headers  Check source license headers without writes' \
	  '  test-license-headers   Test shared header tooling and file selection' \
	  '  test-release      Test stable release input validation' \
	  '  validate-release  Validate notes for TAG=vX.Y.Z' \
	  '  docs-setup        Install locked documentation dependencies (Python 3.12)' \
	  '  docs-serve        Preview bilingual documentation locally' \
	  '  docs-build        Strictly build both documentation languages' \
	  '  setup             Verify pinned formatter dependencies' \
	  '  format            Rewrite owned Go source using shared policy' \
	  '  check-format      Check full Go layout policy without writes' \
	  '  test-format       Test shared layout, ownership and failures' \
	  '  test              Run all tests with statement coverage' \
	  '  race              Run shuffled race tests' \
	  '  vet               Run go vet' \
	  '  lint              Run golangci-lint (must be installed)' \
	  '  verify-mod        Check module dependency tidiness' \
	  '  build             Build golayout into .tmp/bin' \
	  '  check             Run all required local quality gates' \
	  '  check-all         Alias for check'

setup:
	@GO="$(GO)" bash scripts/format.sh setup "$(CURDIR)" "$(GOLAYOUT_TOOL_MODULE)" "$(GO_LANG_VERSION)"
license-headers:
	@GO="$(GO)" bash scripts/license_headers.sh write "$(CURDIR)" "$(GOLAYOUT_TOOL_MODULE)"
check-license-headers:
	@GO="$(GO)" bash scripts/license_headers.sh check "$(CURDIR)" "$(GOLAYOUT_TOOL_MODULE)"
test-license-headers:
	@GO="$(GO)" bash scripts/license_headers_test.sh "$(abspath $(GOLAYOUT_TOOL_MODULE))"
format:
	@GO="$(GO)" bash scripts/format.sh write "$(CURDIR)" "$(GOLAYOUT_TOOL_MODULE)" "$(GO_LANG_VERSION)"
check-format:
	@GO="$(GO)" bash scripts/format.sh check "$(CURDIR)" "$(GOLAYOUT_TOOL_MODULE)" "$(GO_LANG_VERSION)"
test-format:
	@GO="$(GO)" bash scripts/format_test.sh "$(CURDIR)" "$(GOLAYOUT_TOOL_MODULE)" "$(GO_LANG_VERSION)"

test:
	$(GO) test -cover ./...
race:
	$(GO) test -race -shuffle=on ./...
vet:
	$(GO) vet ./...
lint:
	$(GOLANGCI_LINT) run --timeout=10m
verify-mod:
	$(GO) mod tidy -diff
build:
	$(GO) build -o .tmp/bin/golayout ./cmd/golayout
test-release:
	$(PYTHON) -m unittest discover -s scripts/tests -v
validate-release:
	$(PYTHON) scripts/release.py --tag "$(TAG)"
check: check-license-headers test-license-headers check-format test-format test-release verify-mod vet test race lint
check-all: check

docs-setup:
	cd docs/design && $(UV) sync --locked
docs-serve: docs-setup
	cd docs/design && $(UV) run --locked mkdocs serve
docs-build: docs-setup
	cd docs/design && $(UV) run --locked mkdocs build --strict

.DEFAULT:
	@printf 'Unknown target: %s. Available targets follow.\n' "$@" >&2
	@$(MAKE) --no-print-directory help
	@exit 2
