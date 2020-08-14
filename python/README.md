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
 - [ ] automate reading `libenry.h` and generating `ffibuilder.cdef(...)` content from it
 - [x] cover the rest of enry API
 - [x] add `setup.py`
 - [ ] build/release automation on CI (publish on pypi)
 - [ ] try ABI mode, to avoid dependency on C compiler on install (+perf test?)

 ## (experimental) generate bindings using gopy
> [gopy](https://github.com/go-python/gopy) generates (and compiles) a CPython extension module from a go package.

```
# install deps
python3 -m pip install pybindgen
go get golang.org/x/tools/cmd/goimports
go get github.com/go-python/gopy

# generate & build
gopy gen -output=out github.com/go-enry/go-enry/v2
cd out
make

# use
python -c 'import enry; enry.IsVendor("vendors/")'
```