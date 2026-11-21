#!/usr/bin/env bash
# The repository owns file selection; the shared Go library owns notice updates.
set -euo pipefail
mode="${1:?check or write required}"
project_root="${2:?consumer repository root required}"
tool_module="${3:?consumer tool module required}"
project_root=$(cd "$project_root" && git rev-parse --show-toplevel)
if [[ "$tool_module" != /* ]]; then
    tool_module="$project_root/$tool_module"
fi
tool_module=$(cd "$tool_module" && pwd)
case "$mode" in
    write) args=(-w) ;;
    check) args=() ;;
    *) printf 'Unknown source header mode: %s\n' "$mode" >&2; exit 2 ;;
esac
template="$project_root/codequality/HEADER"
if [[ ! -f "$template" ]]; then
    printf 'Missing source header template: %s\n' "$template" >&2
    exit 2
fi
cd "$project_root"
manifest=$(mktemp)
trap 'rm -f "$manifest"' EXIT
git ls-files -z --cached --others --exclude-standard -- '*.go' > "$manifest"
files=()
while IFS= read -r -d '' path; do
    case "$path" in
        .tmp/*|*/.tmp/*|.codex/*|*/.codex/*|.cache/*|*/.cache/*|vendor/*|*/vendor/*|testdata/*|*/testdata/*) continue ;;
    esac
    if [[ -f "$path" && ! -L "$path" ]]; then
        files+=("$PWD/$path")
    fi
done < "$manifest"
if (( ${#files[@]} == 0 )); then
    exit 0
fi
GOFLAGS="${GOFLAGS:-} -mod=readonly" "${GO:-go}" -C "$tool_module" tool golayout "${args[@]}" -header "$template" "${files[@]}"
