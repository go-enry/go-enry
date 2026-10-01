"""Exercise the real package configuration without compiling native code."""
import os
from pathlib import Path
import shutil
import subprocess
import tomllib
import tempfile
import unittest

from build.env import DefaultIsolatedEnv
from packaging.version import Version


class VersioningTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        # Use the declared backend requirements once for all metadata fixtures.
        # No native build or runtime/test dependency installation is needed here.
        project = Path(__file__).resolve().parents[1] / "pyproject.toml"
        requirements = tomllib.loads(project.read_text())["build-system"]["requires"]
        cls.backend = cls.enterClassContext(DefaultIsolatedEnv())
        cls.backend.install(requirements)

    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.package = self.root / "python"
        self.package.mkdir()
        source = Path(__file__).resolve().parents[1]
        for name in ("pyproject.toml", "setup.py", "build_enry.py", "README.md", "LICENSE"):
            shutil.copy2(source / name, self.package / name)
        (self.package / "enry").mkdir()
        (self.package / "enry" / "__init__.py").touch()
        # Isolate fixtures from CI checkout state and any local version overrides.
        self.env = {k: v for k, v in os.environ.items()
                    if not k.startswith(("GIT_", "SETUPTOOLS_SCM_"))}
        self.git("init", "-q")
        self.git("add", ".")
        self.git("-c", "user.name=Test", "-c", "user.email=test@example.invalid",
                 "commit", "-qm", "fixture")

    def git(self, *args):
        return subprocess.run(["git", *args], cwd=self.root, env=self.env,
                              check=True, capture_output=True, text=True)

    def version(self):
        result = subprocess.run([self.backend.python_executable, "setup.py", "--version"],
                                cwd=self.package, env=self.env, check=True,
                                capture_output=True, text=True)
        return Version(result.stdout.strip().splitlines()[-1])

    def test_python_tag_ignores_go_tags(self):
        self.git("tag", "python-v1.2.3")
        self.git("tag", "v99.0.0")
        self.assertEqual(self.version(), Version("1.2.3"))

    def test_prerelease_package_tag(self):
        self.git("tag", "python-v1.2.3rc1")
        self.assertEqual(self.version(), Version("1.2.3rc1"))

    def test_commits_after_tag_are_development_versions(self):
        self.git("tag", "python-v1.2.3")
        (self.root / "change").touch()
        self.git("add", ".")
        self.git("-c", "user.name=Test", "-c", "user.email=test@example.invalid",
                 "commit", "-qm", "next change")
        self.assertTrue(self.version().is_devrelease)
        self.assertGreater(self.version(), Version("1.2.3"))

    def test_dirty_release_checkout_is_not_a_release(self):
        self.git("tag", "python-v1.2.3")
        with (self.package / "README.md").open("a") as output:
            output.write("\nchanged\n")
        self.assertTrue(self.version().is_devrelease)

    def test_sdist_metadata_without_git(self):
        self.git("tag", "python-v1.2.3")
        # sdist writes this metadata into its root; rebuilds must retain it.
        (self.package / "PKG-INFO").write_text(
            "Metadata-Version: 2.1\nName: enry\nVersion: 1.2.3\n")
        shutil.rmtree(self.root / ".git")
        self.assertEqual(self.version(), Version("1.2.3"))


if __name__ == "__main__":
    unittest.main()
