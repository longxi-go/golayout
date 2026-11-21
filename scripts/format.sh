#!/usr/bin/env bash
# Select only repository-owned source; golayout owns all formatting rules.
set -euo pipefail
mode="${1:?setup, check or write required}"
project_root="${2:?consumer project root required}"
tool_module="${3:?consumer tool module directory required}"
language="${4:?Go language version required}"
project_root=$(cd "$project_root" && pwd)
if [[ "$tool_module" != /* ]]; then
    tool_module="$project_root/$tool_module"
fi
tool_module=$(cd "$tool_module" && pwd)
tool=github.com/longxi-go/golayout/cmd/golayout

if [[ "$mode" == setup ]]; then
    "${GO:-go}" -C "$tool_module" list -mod=readonly -deps "$tool" > /dev/null
    printf 'golayout ready\n'
    exit 0
fi

case "$mode" in
    write) args=(-w) ;;
    check) args=() ;;
    *) printf 'Unknown formatting mode: %s\n' "$mode" >&2; exit 2 ;;
esac

cd "$project_root"
cd "$(git rev-parse --show-toplevel)"
manifest=$(mktemp)
trap 'rm -f "$manifest"' EXIT
git ls-files -z --cached --others --exclude-standard -- '*.go' ':(glob)**/testdata/**/*.go.tmpl' > "$manifest"
files=()
while IFS= read -r -d '' path; do
    case "$path" in
        .tmp/*|*/.tmp/*|.cache/*|*/.cache/*|vendor/*|*/vendor/*) continue ;;
        testdata/*|*/testdata/*)
            # The extension explicitly opts complete, parseable fixture source in.
            [[ "$path" == *.go.tmpl ]] || continue
            ;;
    esac
    if [[ -f "$path" && ! -L "$path" ]]; then
        files+=("$PWD/$path")
    fi
done < "$manifest"
if (( ${#files[@]} == 0 )); then
    exit 0
fi

GOFLAGS="${GOFLAGS:-} -mod=readonly" "${GO:-go}" -C "$tool_module" tool golayout "${args[@]}" -lang "$language" "${files[@]}"
