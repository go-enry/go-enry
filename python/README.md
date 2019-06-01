# Python bindings for enry

Python bingings thoug cFFI (API, out-of-line) for calling enr Go functions though CGo wrapper.

## Build

```
$ make static
$ python enry_build.py
```

Will build static library for Cgo wrapper `libenry`, then generate and build `enry.c` 
- a CPython extension that

## Run

Example for single exposed API function is provided.

```
$ python enry.py
```

## TODOs
 - [ ] try ABI mode, to aviod dependency on C compiler on install (+perf test?)
 - [ ] ready `libenry.h` and generate `ffibuilder.cdef` content
 - [ ] cover the rest of enry API
 - [ ] add `setup.py`
 - [ ] build/release automation on CI (publish on pypi)
 - 