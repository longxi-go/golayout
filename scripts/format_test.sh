#!/usr/bin/env bash
# Exercise repository file selection and the shared CLI.
set -euo pipefail
checker="$(cd "$(dirname "$0")" && pwd)/format.sh"
consumer_root="${1:?consumer project root required}"
tool_module="${2:?consumer tool module directory required}"
language="${3:?Go language version required}"
consumer_root=$(cd "$consumer_root" && pwd)
if [[ "$tool_module" != /* ]]; then
    tool_module="$consumer_root/$tool_module"
fi
tool_module=$(cd "$tool_module" && pwd)
module_before=$(cat "$tool_module/go.mod")
fixture=$(mktemp -d)
trap 'rm -rf "$fixture"' EXIT
cd "$fixture"
git init -q
printf '.tmp/\n.cache/\n' > .gitignore

cat > 'file name.go' <<'EOF'
package fixture

var (
Name = "one"
)
type Config struct{ Port int }
func Port() int { return 8080 }
EOF
cp 'file name.go' original
if bash "$checker" check "$fixture" "$tool_module" "$language" > check.log 2>&1; then
    echo 'FAIL: check accepted unformatted owned source' >&2
    exit 1
fi
cmp 'file name.go' original
grep -q 'go layout differs' check.log
bash "$checker" write "$fixture" "$tool_module" "$language"
grep -Fqx 'var Name = "one"' 'file name.go'
grep -Fqx 'type Config struct {' 'file name.go'
grep -Fqx 'func Port() int {' 'file name.go'
bash "$checker" check "$fixture" "$tool_module" "$language"

mkdir -p .tmp .cache vendor testdata nested/vendor
for path in .tmp/ignored.go .cache/ignored.go vendor/dependency.go testdata/invalid.go nested/vendor/invalid.go; do
    printf 'invalid Go syntax\n' > "$path"
done
git add vendor/dependency.go testdata/invalid.go nested/vendor/invalid.go
printf 'package fixture\n' > deleted.go
git add deleted.go
rm deleted.go
ln -s testdata/invalid.go linked.go
bash "$checker" check "$fixture" "$tool_module" "$language"

# Complete Go fixture templates explicitly opt into layout checks, including
# paths containing spaces. Ordinary malformed fixtures and generator templates
# stay outside this inventory.
mkdir -p testdata/scenario internal/compiler/testdata/scenario templates
printf 'package scenario\ntype Config struct{ Port int }\n' > 'testdata/scenario/source file.go.tmpl'
printf 'package scenario\nfunc Value() int { return 1 }\n' > internal/compiler/testdata/scenario/main.go.tmpl
printf '{{define "intentionally not Go"}}\n' > templates/file.go.tmpl
cp 'testdata/scenario/source file.go.tmpl' fixture-original
git add 'testdata/scenario/source file.go.tmpl'
if bash "$checker" check "$fixture" "$tool_module" "$language" > fixture-check.log 2>&1; then
    echo 'FAIL: check accepted compact opted-in fixture templates' >&2
    exit 1
fi
cmp 'testdata/scenario/source file.go.tmpl' fixture-original
grep -q 'go layout differs' fixture-check.log
bash "$checker" write "$fixture" "$tool_module" "$language"
grep -Fqx 'type Config struct {' 'testdata/scenario/source file.go.tmpl'
grep -Fqx 'func Value() int {' internal/compiler/testdata/scenario/main.go.tmpl
bash "$checker" check "$fixture" "$tool_module" "$language"

bash "$checker" setup "$consumer_root" "$tool_module" "$language"
if GO=/bin/false bash "$checker" setup "$consumer_root" "$tool_module" "$language" > setup.log 2>&1; then
    echo 'FAIL: failed tool preparation was reported successful' >&2
    exit 1
fi
if bash "$checker" setup "$fixture" "$fixture/missing" "$language" > missing.log 2>&1; then
    echo 'FAIL: missing tool module was accepted' >&2
    exit 1
fi
printf 'PASS: shared formatting, readonly rejection, ownership exclusions and setup errors\n'
[[ "$module_before" == "$(cat "$tool_module/go.mod")" ]]
