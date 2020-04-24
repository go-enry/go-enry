# Python bindings for enry

Python bindings through cFFI (API, out-of-line) for calling enry Go functions exposed by CGo wrapper.

## Build

```
$ cd .. && make static
$ python build_enry.py
```

Will build a static library for Cgo wrapper `libenry`, then generate and build `enry.c` - a CPython extension that provides actual bindings.

## Run

Example for single exposed API function is provided.

```
$ python enry.py
```

## TODOs
 - [x] helpers for sending/receiving Go slices to C
 - [ ] read `libenry.h` and generate `ffibuilder.cdef(...)` content
 - [ ] cover the rest of enry API
 - [ ] add `setup.py`
 - [ ] build/release automation on CI (publish on pypi)
 - [ ] try ABI mode, to avoid dependency on C compiler on install (+perf test?)