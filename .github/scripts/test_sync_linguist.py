import importlib.util
import os
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch

spec = importlib.util.spec_from_file_location("sync_linguist", Path(__file__).with_name("sync_linguist.py"))
sync = importlib.util.module_from_spec(spec)
spec.loader.exec_module(sync)


class SyncTests(unittest.TestCase):
    def test_rejects_shell_syntax_and_nonstable_tags(self):
        for tag in ("v9.5.0; echo bad", "$(id)", "`id`", "v9.5.0\nchanged=true", "--help", "v9.5.0-rc1", "master"):
            with self.subTest(tag=tag), self.assertRaises(ValueError):
                sync.validate_tag(tag)
        self.assertEqual(sync.validate_tag("v9.5.0"), "v9.5.0")

    def test_resolves_annotated_and_lightweight_tags(self):
        tag_sha, commit_sha = "a" * 40, "b" * 40
        refs = f"{tag_sha}\trefs/tags/v9.5.0\n{commit_sha}\trefs/tags/v9.5.0^{{}}\n"
        self.assertEqual(sync.resolve_commit(refs, "v9.5.0"), commit_sha)
        self.assertEqual(sync.resolve_commit(refs.splitlines()[0], "v9.5.0"), tag_sha)
        with self.assertRaises(ValueError):
            sync.resolve_commit(refs, "v1.0.0")

    def test_unchanged_release_is_a_noop(self):
        sha = "a" * 40
        with tempfile.TemporaryDirectory() as tmp:
            fixture, output = Path(tmp) / "generator.go", Path(tmp) / "output"
            fixture.write_text(f'const commit = "{sha}"\n')
            with patch.object(sync, "GENERATOR", fixture), patch.dict(os.environ, {"LINGUIST_TAG": "v9.5.0", "GITHUB_OUTPUT": str(output)}), patch.object(sync.subprocess, "check_output", return_value=f"{sha}\trefs/tags/v9.5.0\n") as run:
                sync.check()
                self.assertIn("changed=false\n", output.read_text())
                self.assertEqual(run.call_count, 1)
                self.assertIsInstance(run.call_args.args[0], list)

    def test_invalid_input_never_reaches_git(self):
        with patch.dict(os.environ, {"LINGUIST_TAG": "$(touch /tmp/injected)"}), patch.object(sync.subprocess, "check_output") as run:
            with self.assertRaises(ValueError):
                sync.check()
            run.assert_not_called()

    def test_updates_only_expected_metadata(self):
        with tempfile.TemporaryDirectory() as tmp:
            old_cwd = Path.cwd()
            try:
                os.chdir(tmp)
                sync.GENERATOR.parent.mkdir(parents=True)
                sync.GENERATOR.write_text('const commit = "' + 'a' * 40 + '"\n')
                Path("README.md").write_text("Linguist version **v9.5.0**.\n")
                with patch.dict(os.environ, {"LINGUIST_TAG": "v9.6.0", "LINGUIST_COMMIT": "b" * 40}):
                    sync.update()
                self.assertIn("b" * 40, sync.GENERATOR.read_text())
                self.assertEqual(Path("README.md").read_text(), "Linguist version **v9.6.0**.\n")
            finally:
                os.chdir(old_cwd)
