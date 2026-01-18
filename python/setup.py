import subprocess
import shutil
import sys
import os
from pathlib import Path
from setuptools import setup, Distribution

class BinaryDistribution(Distribution):
    def has_ext_modules(self):
        return True

def build_go_library():
    # Only build if the library doesn't exist
    # This allows GitHub Actions to 'pre-build' if desired, 
    # but provides a fallback for cibuildwheel.
    base_dir = Path(__file__).resolve().parent
    enry_dir = base_dir / "enry"
    lib_exists = any(enry_dir.glob("libenry.*"))
    
    if not lib_exists:
        print("Building Go shared library...")
        # Path to the root where Makefile lives
        root_dir = base_dir.parent
        try:
            # Explicitly call make shared
            subprocess.check_call(["make", "shared"], cwd=root_dir)
            # The Makefile likely puts the lib in shared/
            # We need to make sure it's in python/enry/
            for lib in (root_dir / "shared").glob("libenry.*"):
                shutil.copy(lib, enry_dir)
        except Exception as e:
            print(f"Warning: Go build failed: {e}. If this is a wheel build, ensure Go is installed.")

# Trigger the build before setup
build_go_library()

setup(
    distclass=BinaryDistribution,
    cffi_modules=["build_enry.py:ffibuilder"],
)