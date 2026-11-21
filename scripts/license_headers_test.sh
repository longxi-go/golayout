#!/usr/bin/env bash
# Exercise owned-file discovery and real CLI writes in a disposable repository.
set -euo pipefail
checker="$(cd "$(dirname "$0")" && pwd)/license_headers.sh"
tool_module="${1:?absolute tool module required}"
fixture=$(mktemp -d)
trap 'rm -rf "$fixture"' EXIT
cd "$fixture"
git init -q
mkdir -p codequality .tmp .codex .cache vendor nested/vendor testdata ignored
printf 'ignored/\n' > .gitignore
printf 'Copyright © 2026 the fixture authors.\n' > codequality/HEADER
printf 'package fixture\n' > 'file name.go'
cp 'file name.go' original
for path in .tmp/bad.go .codex/bad.go .cache/bad.go vendor/bad.go nested/vendor/bad.go testdata/bad.go ignored/bad.go; do
    printf 'intentionally invalid source\n' > "$path"
done
git add vendor/bad.go nested/vendor/bad.go testdata/bad.go
printf 'package fixture\n' > deleted.go
git add deleted.go
rm deleted.go
ln -s testdata/bad.go linked.go

if bash "$checker" check "$fixture" "$tool_module" > check.log 2>&1; then
    echo 'FAIL: check accepted a missing source header' >&2
    exit 1
fi
cmp 'file name.go' original
grep -q 'source headers differ' check.log
bash "$checker" write "$fixture" "$tool_module"
bash "$checker" check "$fixture" "$tool_module"
grep -Fqx '// Copyright © 2026 the fixture authors.' 'file name.go'
cp 'file name.go' refreshed
bash "$checker" write "$fixture" "$tool_module"
cmp 'file name.go' refreshed

printf 'package fixture\n' > earlier.go
cp earlier.go before-failure
printf 'package fixture\nfunc Broken(' > broken.go
if bash "$checker" write "$fixture" "$tool_module" > failure.log 2>&1; then
    echo 'FAIL: invalid Go accepted' >&2
    exit 1
fi
cmp earlier.go before-failure
printf 'PASS: source header checks, writes, exclusions, idempotence and preflight\n'
