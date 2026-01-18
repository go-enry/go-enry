"""
Python library calling enry Go implementation through cFFI (ABI, out-of-line) and Cgo.
"""

from enry.definitions import (
    get_language,
    get_language_by_extension,
    get_language_by_filename,
    is_vendor,
    is_generated,
)

# Alias for backward compatibility with previous naming
language = get_language
language_by_extension = get_language_by_extension
language_by_filename = get_language_by_filename

## Tests
from collections import namedtuple

def main():
    TestFile = namedtuple("TestFile", "name, content, lang")
    files = [
        TestFile("Parse.hs", b"", "Haskell"),
        TestFile("some.cpp", b"", "C++"), 
        TestFile("orand.go", b"", "Go"),
        TestFile("type.h", b"", "C"), 
        TestFile(".bashrc", b"", "Shell"),
        TestFile(".gitignore", b"", "Ignore List")
    ]

    # Note: .language is the attribute of the Guess namedtuple
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
