"""
Minimal setup.py that handles building the Go static library and cffi bindings.
All package metadata is in pyproject.toml.
"""
import subprocess
import shutil
import sys
from pathlib import Path
from setuptools import setup
from setuptools.command.build_py import build_py
from setuptools.command.develop import develop
from setuptools.command.build_ext import build_ext


def build_go_static_lib():
    """Build the Go static library using make."""
    if shutil.which("go") is None:
        print(
            "\n" + "="*70 + "\n"
            "ERROR: Go compiler not found!\n\n"
            "To build from source, you need Go installed.\n"
            "Install Go: https://golang.org/dl/\n\n"
            "Or install pre-built wheels instead:\n"
            "  pip install enry\n"
            + "="*70 + "\n",
            file=sys.stderr
        )
        sys.exit(1)
    
    print("Building Go static library (this may take a minute)...")
    try:
        subprocess.check_call(["make", "static"], cwd="..")
        print("✓ Go static library built successfully")
    except subprocess.CalledProcessError as e:
        print(f"\n✗ Failed to build Go static library: {e}", file=sys.stderr)
        sys.exit(1)


def build_cffi_extension():
    """Build the cffi extension by running build_enry.py."""
    print("Building cffi extension...")
    try:
        # Run build_enry.py which creates the cffi bindings
        subprocess.check_call([sys.executable, "build_enry.py"])
        print("✓ cffi extension built successfully")
    except subprocess.CalledProcessError as e:
        print(f"\n✗ Failed to build cffi extension: {e}", file=sys.stderr)
        sys.exit(1)


class BuildPyWithGo(build_py):
    """Build Python package after building Go static library and cffi extension."""
    
    def run(self):
        build_go_static_lib()
        build_cffi_extension()
        super().run()


class DevelopWithGo(develop):
    """Install in development mode after building Go static library and cffi extension."""
    
    def run(self):
        build_go_static_lib()
        build_cffi_extension()
        super().run()


if __name__ == "__main__":
    setup(
        # Reference to cffi builder
        cffi_modules=["build_enry.py:ffibuilder"],
        cmdclass={
            "build_py": BuildPyWithGo,
            "develop": DevelopWithGo,
        }
    )