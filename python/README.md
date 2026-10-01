# Python bindings for enry

Python bindings through cFFI (ABI, out-of-line) for calling enry Go functions exposed by CGo wrapper.

## Build

```bash
# from python/
python -m pip install -e .
```

The editable install builds the Go shared library and CFFI bindings automatically.
Build requirements and runtime dependencies come from `pyproject.toml`.

### Why a shared library?

Historically, the Python package shipped a **static** library and used CFFI **API mode** (a compiled extension that links against `libenry`).
That approach relied on Go-generated headers/types (e.g. `GoString`, `GoSlice`, and struct return wrappers) and a locally-built archive at build time, which made builds and cross-platform packaging more fragile.

We now build a Go-built **shared** library (`-buildmode=c-shared`) that is bundled inside the wheel and loaded via CFFI **out-of-line (ABI)** mode.
This makes installation simpler and allows `pip install enry-python` without requiring a Go toolchain. Source builds require Go and a C compiler; the sdist bundles its Go sources under `_go/`.

**Implementation note:** the shared library is located and loaded at import time in `enry/definitions.py` (see `_load_library()`), which prefers the packaged `enry/libenry.*` shipped in wheels and falls back to local dev build locations.

**Future-proof note:** if we later change the binding strategy (or revisit linking/packaging), we should reevaluate whether a shared library is still the best tradeoff and update this documentation accordingly.

## Installation

### From PyPI (Recommended)

For Python 3.12+, install pre-built wheels:
```bash
pip install enry-python
```

No Go compiler required! Pre-built wheels are available for:
- **Linux**: x86_64 (manylinux)
- **macOS 12+**: x86_64 (Intel) and arm64 (Apple Silicon)

### From Source

If you need to build from source or use an unsupported platform, you'll need Go installed:
```bash
git clone https://github.com/go-enry/go-enry.git
cd go-enry
cd python
pip install -e .
```

**Requirements for building:**
- Go 1.26.x (used by release CI)
- GCC or compatible C compiler
- Python 3.12 or later

## Developer: publishing to PyPI

Releases are published by the tagged wheel workflow using **[PyPI Trusted Publishing (OIDC)](https://docs.pypi.org/trusted-publishers/using-a-publisher/)**.

**Note:** CI publishing via OIDC is **gated on PyPI Trusted Publisher configuration** for the `enry-python` project (must be set up by a PyPI project owner/maintainer for this repo/workflow).

Until that is configured, you can publish manually using a [PyPI API token](https://pypi.org/help/#apitoken).

### Versioning

The package version is derived by setuptools-scm from `python-vX.Y.Z` tags;
there is no version number to edit in `pyproject.toml`. For example, tagging
`python-v0.3.0` builds version `0.3.0`. Go's `v*` tags are ignored. Untagged
commits receive development versions and do not trigger publication.

### Automated releases

The **Prepare Python Release** workflow runs after pushes to `master`. Release
Please opens or updates a release PR with the next version and changelog.
Review and merge that PR when ready to publish: the workflow then creates the
`python-vX.Y.Z` tag and GitHub release, and explicitly starts **Build Python
Wheels** at that tag. All wheel and source tests must pass before PyPI upload.
Changes to bundled Go/native detector code, Go dependencies or language data
request at least a **minor** Python release, including non-conventional upstream
merges and `chore(data)` Linguist syncs. Multiple changes since the last Python
release are collected into one release PR. For example, native changes after
`0.3.0` propose `0.4.0`. Breaking changes retain their normal versioning precedence.

Like Sync Linguist, this automation opens/updates a PR for review; it does not
merge it automatically. Merging the release PR authorizes publication. Java,
CLI, generator-source, native-test-only, docs and general CI changes do not
request a Python release. Python source/packaging and its wheel workflow retain
conventional-commit versioning. The policy lives in `.github/release/` and is
tested with the pinned Release Please library; dependency updates are lockfile-backed.

Use conventional commit titles when merging changes:

- `fix(python): ...` requests a patch release.
- `feat(python): ...` requests a minor release.
- `feat(python)!: ...` or a `BREAKING CHANGE:` footer requests a breaking release;
  before 1.0 this increments the minor version, and from 1.0 it increments major.
- `docs`, `chore` and other non-release types do not request a release themselves.

For squash merges, put the conventional title and any breaking-change footer
in the squash commit. Review the proposed version before merging the release PR.

`.github/.release-please-manifest.json` records the last released version for
Release Please; the bot updates it. It is not the Python package's version source.
The root package configuration includes Go detector and native binding changes.
The migration boundary (`last-release-sha`) starts the first changelog at the
Python refresh, excluding imported upstream history; later releases stop at
their own release commits.
The Go release strategy writes only the changelog, leaving setuptools-scm to
supply Python metadata from the tag. CI verifies artifact versions against that
tag, and sdists retain their version when rebuilt without Git.

The workflow uses the built-in `GITHUB_TOKEN`, with no personal access token.
Enable **Allow GitHub Actions to create and approve pull requests** in the
repository's Actions settings. Bot-created PRs and tags do not trigger normal
PR/push workflows with this token; the explicit wheel workflow dispatch handles
release validation. If branch rules require PR checks, provide those checks
before merging (for example, by using a GitHub App token for Release Please).
Do not bypass required checks. This repository's automation does not configure
an App or bypass branch protection.

If the dispatch job fails, rerun that failed job. To retry the complete tagged
build, run:

```bash
gh workflow run python-wheels.yml --repo go-enry/go-enry --ref python-vX.Y.Z
```

This command can publish to PyPI once tests pass. Reuse the existing tag when
retrying a failed build; never move a published tag. A GitHub release indicates
tag creation, not successful PyPI publication—check **Build Python Wheels**.
Build from a Git clone or published sdist; unversioned source copies are unsupported.

### Manual fallback: upload CI-built artifacts

This mirrors what the CI does (cibuildwheel builds platform wheels + an sdist). You simply upload the produced artifacts yourself.

1) Prefer merging the automated release PR. If automation is unavailable, create
an unused release tag manually (this triggers the workflow):

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
- The token must be created on PyPI by an account with upload permission for the enry-python project.
- This approach is preferred because wheels must be built per-platform/per-arch (Linux manylinux + macOS x86_64/arm64).
- PyPI token notes: set username to __token__ and password to the token value (including the pypi- prefix).

### Manual publish (local, single-platform only)

For a local validation build, use the default build command: it creates an sdist, then builds the wheel from that archive. Native build failures abort packaging.

```bash
# from repo root
cd python
python -m pip install --upgrade pip
python -m pip install --group ci
python -m build

TWINE_USERNAME=__token__ TWINE_PASSWORD='pypi-***' python -m twine upload dist/*
```

This requires Go locally and only produces a wheel for the current OS/arch.

## Development dependencies

`python/pyproject.toml` is the source of truth: `[build-system].requires` supplies
isolated build environments, `[project].dependencies` supplies runtime packages,
and `[dependency-groups]` defines `test` and `ci` tools. Requirements files are
no longer needed. CFFI is declared for both build and runtime because both use it.

```bash
# from python/; dependency groups require pip 25.1 or later
python -m pip install --upgrade pip
python -m pip install --group test -e .
python -m pytest tests -q
python -m pip install --group ci
python -m unittest discover -s packaging_tests -v
```

The workflows use these same groups; cibuildwheel reads `test-groups` directly.
The SHA-pinned cibuildwheel action owns its tool version, updated by Dependabot,
so there is no second cibuildwheel pin in a requirements file. Runner operating
systems and Go versions remain workflow settings.

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
does not itself trigger a PyPI upload: merge the next release PR to publish
a new enry release with the expanded wheel set. If a tooling update needs a
release, use a `fix(python): ...` commit describing the newly supported wheels.

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
- Java bindings use the same C ABI; see `../java/README.md` for build instructions.



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
