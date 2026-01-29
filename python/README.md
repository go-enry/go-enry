# Python bindings for enry

Python bindings through cFFI (ABI, out-of-line) for calling enry Go functions exposed by CGo wrapper.

## Build

```
$ pushd .. && make shared && popd
$ pip install -r requirements.txt
$ python build_enry.py
```

Will build a static library for Cgo wrapper `libenry`, then generate and build `enry.c` - a CPython extension that provides actual bindings.

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
