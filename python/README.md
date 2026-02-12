# Python bindings for enry

Python bindings through cFFI (ABI, out-of-line) for calling enry Go functions exposed by CGo wrapper.

## Build

```
# from python/
$ pushd .. && make shared && popd
$ pip install -r requirements.txt
$ python build_enry.py
```

Builds the Go **shared** library for the CGo wrapper (`libenry.so` / `libenry.dylib`), then generates and builds the CFFI out-of-line module (`enry.c`) that provides the Python bindings.

### Why a shared library?

Historically, the Python package shipped a **static** library and used CFFI **API mode** (a compiled extension that links against `libenry`).
That approach relied on Go-generated headers/types (e.g. `GoString`, `GoSlice`, and struct return wrappers) and a locally-built archive at build time, which made builds and cross-platform packaging more fragile.

We now build a Go-built **shared** library (`-buildmode=c-shared`) that is bundled inside the wheel and loaded via CFFI **out-of-line (ABI)** mode.
This makes installation simpler and allows `pip install enry` without requiring a Go toolchain.

**Implementation note:** the shared library is located and loaded at import time in `enry/definitions.py` (see `_load_library()`), which prefers the packaged `enry/libenry.*` shipped in wheels and falls back to local dev build locations.

**Future-proof note:** if we later change the binding strategy (or revisit linking/packaging), we should reevaluate whether a shared library is still the best tradeoff and update this documentation accordingly.

## Installation

### From PyPI (Recommended)

For Python 3.9+, install pre-built wheels:
```bash
pip install enry
```

No Go compiler required! Pre-built wheels are available for:
- **Linux**: x86_64 (manylinux)
- **macOS**: x86_64 (Intel) and arm64 (Apple Silicon)

### From Source

If you need to build from source or use an unsupported platform, you'll need Go installed:
```bash
git clone https://github.com/go-enry/go-enry.git
cd go-enry
make shared
cd python
pip install -e .
```

**Requirements for building:**
- Go 1.21 or later
- GCC or compatible C compiler
- Python 3.9 or later

## Developer: publishing to PyPI

Releases are intended to be published from CI on tag pushes (`python-v*`) using **[PyPI Trusted Publishing (OIDC)](https://docs.pypi.org/trusted-publishers/using-a-publisher/)**.

**Note:** CI publishing via OIDC is **gated on PyPI Trusted Publisher configuration** for the `enry` project (must be set up by a PyPI project owner/maintainer for this repo/workflow).

Until that is configured, you can publish manually using a [PyPI API token](https://pypi.org/help/#apitoken).

### Versioning

Before tagging a release, bump the Python package version in `python/pyproject.toml` (`[project].version`) and commit it.
The git tag should match the package version (e.g. `version = "0.2.1"` and tag `python-v0.2.1`).

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

If you only need to build and upload artifacts for your current machine, you should build the Go shared library first (mirrors the CI’s make shared step):

```bash
# from repo root: builds and copies libenry.* into python/enry/
make shared

cd python
cp ../LICENSE .
python -m pip install --upgrade build
python -m build --sdist --wheel

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

Some upstream Go `enry` functions return `(value, safe)` where `safe` indicates whether the result is considered unambiguous.
The shared library (`libenry`) exports used by these Python bindings intentionally return **only the primary string value** and do not currently expose the `safe` flag.

**Rationale:** keeping the C ABI to simple primitives (`char*` in/out + explicit free) avoids returning structs/tuples across the language boundary and reduces ABI + memory-ownership pitfalls. It also keeps the shared library broadly usable by non-Python consumers without committing the core ABI to a particular “safe mode” policy.

**Tradeoff:** the Python bindings cannot directly access the Go `safe` signal via the current exports. If we decide we need it, an ABI-friendly extension would be to add parallel exports that surface safety without structs (e.g. `...WithSafety(..., int* out_safe)`), while keeping the existing string-only exports for backwards compatibility.


## Supported Python Versions

- Python 3.9+
- CPython only (PyPy not yet supported)

**Note:** Python 3.6, 3.7 and 3.8 reached end-of-life and are no longer supported. 
Use enry 0.1.1 if you must use these versions (not recommended for security reasons).

## Platform Support

- ✅ Linux (x86_64)
- ✅ macOS (Intel x86_64 and Apple Silicon arm64)
- ❌ Linux ARM/aarch64 (not yet available)

## Known Issues

- Memory leak fixed in version 0.2.0 (see [#36](https://github.com/go-enry/go-enry/issues/36))
- The current shared-library exports return only a string result and do not expose the Go `safe` flag (by design; see “FFI / API design notes” above).



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
