"""
Python library calling enry Go implementation trough cFFI (API, out-of-line) and Cgo.
"""
import os
import platform
from typing import List
from enry._c_enry import ffi
from enry.types import Guess
from enry.utils import (
    go_str_to_py, 
    go_str_slice_to_py, 
    prepare_candidates, 
    go_guess_to_py
)

def _load_library():
    """
    Locates and opens the Go shared library based on the operating system.
    This is the core of the ABI mode implementation.
    """
    base_path = os.path.dirname(__file__)
    system = platform.system()
    
    if system == "Darwin":
        lib_name = "libenry.dylib"
    elif system == "Windows":
        lib_name = "libenry.dll"
    else:
        lib_name = "libenry.so"
        
    lib_path = os.path.join(base_path, lib_name)
    
    # Fallback to current directory for local testing
    if not os.path.exists(lib_path):
        lib_path = os.path.join(os.getcwd(), lib_name)

    try:
        return ffi.dlopen(lib_path)
    except OSError as e:
        raise ImportError(
            f"Could not load the enry shared library at {lib_path}. "
            "Ensure the Go library is built and placed in the enry/ directory."
        ) from e

# Load the library globally for this module
lib = _load_library()

# --- Singular API ---

def get_language(filename: str, content: bytes) -> str:
    res = lib.GetLanguage(filename.encode(), content, len(content))
    return go_str_to_py(lib, res)

def get_language_by_content(filename: str, content: bytes) -> Guess:
    res = lib.GetLanguageByContent(filename.encode(), content, len(content))
    return go_guess_to_py(lib, res)

def get_language_by_extension(filename: str) -> Guess:
    res = lib.GetLanguageByExtension(filename.encode())
    return go_guess_to_py(lib, res)

def get_language_by_filename(filename: str) -> Guess:
    res = lib.GetLanguageByFilename(filename.encode())
    return go_guess_to_py(lib, res)

def get_language_by_modeline(content: bytes) -> Guess:
    res = lib.GetLanguageByModeline(content, len(content))
    return go_guess_to_py(lib, res)

def get_language_by_shebang(content: bytes) -> Guess:
    res = lib.GetLanguageByShebang(content, len(content))
    return go_guess_to_py(lib, res)

def get_language_by_emacs_modeline(content: bytes) -> Guess:
    res = lib.GetLanguageByEmacsModeline(content, len(content))
    return go_guess_to_py(lib, res)

def get_language_by_vim_modeline(content: bytes) -> Guess:
    res = lib.GetLanguageByVimModeline(content, len(content))
    return go_guess_to_py(lib, res)

def get_mime_type(path: str, language: str) -> str:
    res = lib.GetMimeType(path.encode(), language.encode())
    return go_str_to_py(lib, res)

def get_color(language: str) -> str:
    res = lib.GetColor(language.encode())
    return go_str_to_py(lib, res)

def get_language_type(language: str) -> str:
    res = lib.GetLanguageType(language.encode())
    return go_str_to_py(lib, res)

# --- Boolean API ---

def is_binary(content: bytes) -> bool:
    return bool(lib.IsBinary(content, len(content)))

def is_vendor(filename: str) -> bool:
    return bool(lib.IsVendor(filename.encode()))

def is_generated(filename: str, content: bytes) -> bool:
    return bool(lib.IsGenerated(filename.encode(), content, len(content)))

def is_configuration(path: str) -> bool:
    return bool(lib.IsConfiguration(path.encode()))

def is_documentation(path: str) -> bool:
    return bool(lib.IsDocumentation(path.encode()))

def is_dot_file(path: str) -> bool:
    return bool(lib.IsDotFile(path.encode()))

def is_image(path: str) -> bool:
    return bool(lib.IsImage(path.encode()))

def is_test(path: str) -> bool:
    return bool(lib.IsTest(path.encode()))

# --- Plural API ---

def get_languages(filename: str, content: bytes) -> List[str]:
    res = lib.GetLanguages(filename.encode(), content, len(content))
    return go_str_slice_to_py(lib, res)

def get_language_extensions(language: str) -> List[str]:
    res = lib.GetLanguageExtensions(language.encode())
    return go_str_slice_to_py(lib, res)

def get_languages_by_filename(filename: str, content: bytes = b"", candidates: List[str] = None) -> List[str]:
    c_cand = prepare_candidates(candidates)
    res = lib.GetLanguagesByFilename(filename.encode(), content, len(content), c_cand)
    return go_str_slice_to_py(lib, res)

def get_languages_by_shebang(filename: str, content: bytes = b"", candidates: List[str] = None) -> List[str]:
    c_cand = prepare_candidates(candidates)
    res = lib.GetLanguagesByShebang(filename.encode(), content, len(content), c_cand)
    return go_str_slice_to_py(lib, res)