# Python bindings for enry

Python bindings through cFFI (ABI, out-of-line) for calling enry Go functions exposed by CGo wrapper.

## Build

```
# from python/
$ pushd .. && make shared && popd
$ pip install -r requirements.txt
$ python build_enry.py
```

Builds the Go **shared** library for the CGo wrapper (`libenry.so` / `libenry.dylib`), then generates the CFFI out-of-line Python module (`enry/_c_enry.py`) that provides the Python bindings.

### Why a shared library?

Historically, the Python package shipped a **static** library and used CFFI **API mode** (a compiled extension that links against `libenry`).
That approach relied on Go-generated headers/types (e.g. `GoString`, `GoSlice`, and struct return wrappers) and a locally-built archive at build time, which made builds and cross-platform packaging more fragile.

We now build a Go-built **shared** library (`-buildmode=c-shared`) that is bundled inside the wheel and loaded via CFFI **out-of-line (ABI)** mode.
This makes installation simpler and allows `pip install enry` without requiring a Go toolchain. Source builds require Go and a C compiler; the sdist bundles its Go sources under `_go/`.

**Implementation note:** the shared library is located and loaded at import time in `enry/definitions.py` (see `_load_library()`), which prefers the packaged `enry/libenry.*` shipped in wheels and falls back to local dev build locations.

**Future-proof note:** if we later change the binding strategy (or revisit linking/packaging), we should reevaluate whether a shared library is still the best tradeoff and update this documentation accordingly.

## Installation

### From PyPI (Recommended)

For Python 3.12+, install pre-built wheels:
```bash
pip install enry
```

No Go compiler required! Pre-built wheels are available for:
- **Linux**: x86_64 (manylinux)
- **macOS 12+**: x86_64 (Intel) and arm64 (Apple Silicon)

### From Source

If you need to build from source or use an unsupported platform, you'll need Go installed:
```bash
git clone https://github.com/Esoiya/go-enry.git
cd go-enry
cd python
pip install -e .
```

**Requirements for building:**
- Go 1.26.x (used by release CI)
- GCC or compatible C compiler
- Python 3.12 or later

## Developer: publishing to PyPI

Releases are intended to be published from CI on tag pushes (`python-v*`) using **[PyPI Trusted Publishing (OIDC)](https://docs.pypi.org/trusted-publishers/using-a-publisher/)**.

**Note:** CI publishing via OIDC is **gated on PyPI Trusted Publisher configuration** for the `enry` project (must be set up by a PyPI project owner/maintainer for this repo/workflow).

Until that is configured, you can publish manually using a [PyPI API token](https://pypi.org/help/#apitoken).

### Versioning

The package version is derived by setuptools-scm from `python-vX.Y.Z` tags;
there is no version number to edit in `pyproject.toml`. For example, tagging
`python-v0.3.0` builds version `0.3.0`. Go's `v*` tags are ignored. Untagged
commits receive development versions and do not trigger publication.

Choose the release number and push its tag after the changes are ready. This
automates applying the version, not choosing semantic-version bumps or creating
releases on every merge. CI fetches the full Git history and validates the built
artifact versions against the release tag. Source distributions retain their
version when rebuilt without Git. Build from a Git clone or a published sdist;
unversioned source copies are not supported.

### Manual publish (recommended): upload CI-built artifacts

This mirrors what the CI does (cibuildwheel builds platform wheels + an sdist). You simply upload the produced artifacts yourself.

1) Tag a release (this triggers the workflow):

```bash
git tag python-vX.Y.Z
git push origin python-vX.Y.Z
```

2) Download the workflow artifacts from GitHub Actions:
- the built wheels (wheels-*)
- the source distribution (sdist)

3) Upload with [Twine](https://packaging.python.org/tutorials/packaging-projects/) using a PyPI token:

```bash
python -m pip install --upgrade twine

# from the directory where you downloaded artifacts:
TWINE_USERNAME=__token__ TWINE_PASSWORD='pypi-***' python -m twine upload **/*.whl **/*.tar.gz
```

Notes:
- The token must be created on PyPI by an account with upload permission for the enry project.
- This approach is preferred because wheels must be built per-platform/per-arch (Linux manylinux + macOS x86_64/arm64).
- PyPI token notes: set username to __token__ and password to the token value (including the pypi- prefix).

### Manual publish (local, single-platform only)

For a local validation build, use the default build command: it creates an sdist, then builds the wheel from that archive. Native build failures abort packaging.

```bash
# from repo root
cd python
python -m pip install --upgrade build
python -m build

python -m pip install --upgrade twine
TWINE_USERNAME=__token__ TWINE_PASSWORD='pypi-***' python -m twine upload dist/*
```

This requires Go locally and only produces a wheel for the current OS/arch.

## Usage
```python
import enry

# Detect language by filename and content
language = enry.get_language("example.py", b"print('Hello, world!')")
print(f"Detected language: {language}")
```

### FFI / API design notes

The Python `get_language_by_*` functions return `Guess(language, safe)`.
`safe` is true only when Go reports an unambiguous result. The shared library
provides additive `...WithSafety(..., int* out_safe)` exports; the original
string-only exports remain available for existing C consumers. Returned strings
and string arrays are freed by the bindings, including when decoding fails.

## Supported Python Versions

- Python 3.12+
- CPython only (PyPy not yet supported)

Python 3.12 remains the minimum in `requires-python`. Release CI builds and tests
all stable standard (GIL-enabled) CPython versions supported by its pinned
cibuildwheel release and satisfying that minimum (currently 3.12–3.14).
New versions are selected automatically when supported by an updated
cibuildwheel action; Dependabot proposes these updates, which still need to be
merged. There is no per-version wheel list to maintain. A new Python release
does not itself trigger a PyPI upload: push a new `python-vX.Y.Z` tag to publish
a new enry release with the expanded wheel set.

Python 3.11 and older, prerelease interpreters, free-threaded CPython and PyPy
are outside the tested support matrix. Wheels do not upgrade the user's Python
installation, and the minimum version is never raised automatically.
Older Python installations must use an older compatible enry release.

## Platform Support

- ✅ Linux (x86_64)
- ✅ macOS 12+ (Intel x86_64 and Apple Silicon arm64)
- ❌ Linux ARM/aarch64 (not yet available)

## Known Issues

- Memory leak fixed in version 0.2.0 (see [#36](https://github.com/go-enry/go-enry/issues/36))
- Java bindings still target an older shared-library ABI and need a separate migration.



## Run

Example for single exposed API function is provided.

```
$ python enry.py
```

## TODOs
 - [x] helpers for sending/receiving Go slices to C
 - [x] read `libenry.h` and generate `ffibuilder.cdef(...)` content
 - [x] cover the rest of enry API
 - [x] add `setup.py`
 - [x] build/release automation on CI (publish on pypi)
 - [x] try ABI mode, to avoid dependency on C compiler on install (+perf test?)
