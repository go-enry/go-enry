"""
Python library calling enry Go implementation trough cFFI (API, out-of-line) and Cgo.
"""

from _c_enry import ffi, lib

## cgo -> ffi helpers
def py_bytes_to_go(py_bytes: bytes):
    c_bytes = ffi.new("char[]", len(py_bytes))
    go_slice = ffi.new("GoSlice *", [c_bytes, len(py_bytes), len(py_bytes)])
    return (go_slice[0], c_bytes)

def py_str_to_go(py_str: str):
    str_bytes = py_str.encode()
    c_str = ffi.new("char[]", str_bytes)
    go_str = ffi.new("_GoString_ *", [c_str, len(str_bytes)])
    return (go_str[0], c_str)

def go_str_to_py(go_str: str):
    str_len = go_str.n
    if str_len > 0:
        return ffi.unpack(go_str.p, go_str.n).decode()
    return ""

def go_bool_to_py(go_bool: bool):
    return go_bool == 1


## API, TODO(bzz): add docstrings
def language(filename: str, content: bytes) -> str:
    fName, c_str = py_str_to_go(filename)
    fContent, c_bytes = py_bytes_to_go(content)
    guess = lib.GetLanguage(fName, fContent)
    lang = go_str_to_py(guess)
    return lang

def language_by_extension(filename: str) -> str:
    fName, c_str = py_str_to_go(filename)
    guess = lib.GetLanguageByExtension(fName)
    lang = go_str_to_py(guess.r0)
    return lang

def language_by_filename(filename: str) -> str:
    fName, c_str = py_str_to_go(filename)
    guess = lib.GetLanguageByFilename(fName)
    lang = go_str_to_py(guess.r0)
    return lang

def is_vendor(filename: str) -> bool:
    fName, c_str = py_str_to_go(filename)
    guess = lib.IsVendor(fName)
    return go_bool_to_py(guess)

def is_generated(filename: str, content: bytes) -> bool:
    fname, c_str = py_str_to_go(filename)
    fcontent, c_bytes = py_bytes_to_go(content)
    guess = lib.IsGenerated(fname, fcontent)
    return go_bool_to_py(guess)


## Tests
from collections import namedtuple

def main():
    TestFile = namedtuple("TestFile", "name, content, lang")
    files = [
        TestFile("Parse.hs", b"", "Haskell"), TestFile("some.cpp", b"", "C++"), 
        TestFile("orand.go", b"", "Go"), TestFile("type.h", b"", "C"), 
        TestFile(".bashrc", b"", "Shell"), TestFile(".gitignore", b"", "Ignore List")
    ]

    print("\nstrategy: extension")
    for f in files:
        lang = language_by_extension(f.name)
        print("\tfile: {:10s} language: '{}'".format(f.name, lang))

    print("\nstrategy: filename")
    for f in files:
        lang = language_by_filename(f.name)
        print("\tfile: {:10s} language: '{}'".format(f.name, lang))

    print("\ncheck: is vendor?")
    for f in files:
        vendor = is_vendor(f.name)
        print("\tfile: {:10s} vendor: '{}'".format(f.name, vendor))

    print("\nstrategy: all")
    for f in files:
        lang = language(f.name, f.content)
        print("\tfile: {:10s} language: '{}'".format(f.name, lang))
        assert lang == f.lang, "Expected '{}' but got '{}'".format(f.lang, lang)

if __name__ == "__main__":
    main()
