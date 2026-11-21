"""Validate a stable Go module tag and its maintained English release notes."""

import argparse
from pathlib import Path
import re

STABLE_TAG = re.compile(r"v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)")


def validate_tag(root: Path, tag: str) -> str:
    """Git tags own module versions; go.mod's go directive is a language minimum."""
    match = STABLE_TAG.fullmatch(tag)
    if match is None:
        raise ValueError("Release tags must use stable vX.Y.Z versions without leading zeroes")
    # The current module has no /vN suffix. A v2+ release needs a deliberate
    # module-path migration rather than silently publishing an unusable Go tag.
    if int(match[1]) >= 2:
        raise ValueError("v2+ requires a versioned module path and a release-policy update")
    notes = root / "docs/releases" / f"{tag}.md"
    if not notes.is_file():
        raise ValueError(f"Missing release notes: {notes}")
    text = notes.read_text(encoding="utf-8").strip()
    lines = text.splitlines()
    if not lines or lines[0] != f"# golayout {tag}" or not "\n".join(lines[1:]).strip():
        raise ValueError(f"Release notes must have a matching title and nonempty body: {notes}")
    return tag[1:]


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--root", type=Path, default=Path(__file__).resolve().parents[1])
    parser.add_argument("--tag", required=True)
    args = parser.parse_args()
    try:
        print(validate_tag(args.root, args.tag))
    except (ValueError, OSError) as error:
        parser.error(str(error))


if __name__ == "__main__":
    main()
