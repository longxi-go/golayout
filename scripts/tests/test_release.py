"""Stable Go module releases require a matching tag and reviewed release notes."""

from pathlib import Path
import sys
import tempfile
import unittest

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
import release


class ReleaseTests(unittest.TestCase):
    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory()
        self.addCleanup(self.temporary.cleanup)
        self.root = Path(self.temporary.name)
        self.notes = self.root / "docs/releases/v0.1.0.md"
        self.notes.parent.mkdir(parents=True)
        self.notes.write_text("# golayout v0.1.0\n\nInitial stable release.\n", encoding="utf-8")

    def test_accepts_stable_semver_with_matching_notes(self):
        self.assertEqual("0.1.0", release.validate_tag(self.root, "v0.1.0"))

    def test_rejects_nonstable_malformed_or_path_like_tags(self):
        for tag in ["0.1.0", "v00.1.0", "v0.01.0", "v0.1.00", "v0.1", "v0.1.0-rc.1", "v0.1.0+build", "v0.1.0\n", "v2.0.0", "v0.1.0/../../main"]:
            with self.subTest(tag=tag), self.assertRaises(ValueError):
                release.validate_tag(self.root, tag)

    def test_rejects_missing_blank_or_mismatched_release_notes(self):
        self.notes.unlink()
        with self.assertRaises(ValueError):
            release.validate_tag(self.root, "v0.1.0")
        self.notes.write_text(" \n", encoding="utf-8")
        with self.assertRaises(ValueError):
            release.validate_tag(self.root, "v0.1.0")
        self.notes.write_text("# golayout v0.2.0\n\nWrong version.\n", encoding="utf-8")
        with self.assertRaises(ValueError):
            release.validate_tag(self.root, "v0.1.0")
        self.notes.write_text("# golayout v0.1.0\n", encoding="utf-8")
        with self.assertRaises(ValueError):
            release.validate_tag(self.root, "v0.1.0")


if __name__ == "__main__":
    unittest.main()
