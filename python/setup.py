"""Build native code only for wheels/editable installs, never for metadata."""
import os
import platform
import shutil
import subprocess
from pathlib import Path

from setuptools import Distribution, setup
from setuptools.command.build_ext import build_ext
from setuptools.command.build_py import build_py
from setuptools.command.sdist import sdist

BASE = Path(__file__).resolve().parent


def go_sources():
    # A published sdist carries its own copy of the Go module.
    return BASE / "_go" if (BASE / "_go" / "go.mod").is_file() else BASE.parent


def build_go_library(destination):
    system = platform.system()
    if system not in {"Darwin", "Linux"}:
        raise RuntimeError("enry native builds support Linux and macOS only")
    architecture = {"x86_64": "amd64", "amd64": "amd64", "arm64": "arm64",
                    "aarch64": "arm64"}.get(platform.machine().lower())
    if architecture is None:
        raise RuntimeError(f"Unsupported architecture: {platform.machine()}")
    destination = Path(destination).resolve()
    destination.mkdir(parents=True, exist_ok=True)
    library = destination / ("libenry.dylib" if system == "Darwin" else "libenry.so")
    env = dict(os.environ, CGO_ENABLED="1", GOARCH=architecture,
               GOOS="darwin" if system == "Darwin" else "linux")
    # Always rebuild: a pre-existing file may be stale or from another platform.
    subprocess.run(["go", "build", "-trimpath", "-buildmode=c-shared", "-o",
                    str(library), "./shared"], cwd=go_sources(), env=env, check=True)
    library.with_suffix(".h").unlink(missing_ok=True)


class BinaryDistribution(Distribution):
    def has_ext_modules(self):
        return True


class BuildPy(build_py):
    def run(self):
        super().run()
        build_go_library(Path(self.build_lib) / "enry")


class BuildExt(build_ext):
    def run(self):
        super().run()
        if self.inplace:
            build_go_library(BASE / "enry")


class SourceDistribution(sdist):
    def make_release_tree(self, base_dir, files):
        super().make_release_tree(base_dir, files)
        root = go_sources()
        target = Path(base_dir) / "_go"
        sources = list(root.glob("*.go"))
        for directory in ("data", "regex", "internal/tokenizer", "shared"):
            sources.extend((root / directory).rglob("*"))
        sources.extend(root / name for name in ("go.mod", "go.sum", "LICENSE"))
        for source in sources:
            if not source.is_file() or source.name.endswith("_test.go"):
                continue
            if source.suffix not in {".go", ".c", ".h"} and source.name not in {"go.mod", "go.sum", "LICENSE"}:
                continue
            output = target / source.relative_to(root)
            output.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy2(source, output)


setup(
    distclass=BinaryDistribution,
    cmdclass={"build_py": BuildPy, "build_ext": BuildExt, "sdist": SourceDistribution},
    cffi_modules=["build_enry.py:ffibuilder"],
)
