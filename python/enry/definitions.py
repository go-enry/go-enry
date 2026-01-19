"""
Python library calling enry Go implementation through cFFI (ABI, out-of-line) and Cgo.
"""
import platform
from pathlib import Path
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
    Locates and opens the Go shared library using Path.
    Handles cross-platform naming and package vs development layouts.
    """
    # Path to the directory containing this file (python/enry/)
    pkg_dir = Path(__file__).resolve().parent
    
    system = platform.system().lower()
    
    if system == "darwin":
        lib_name = "libenry.dylib"
    elif system == "windows":
        lib_name = "libenry.dll"
    else:
        lib_name = "libenry.so"

    # 1. Search in the current package directory (standard for wheels)
    lib_path = pkg_dir / lib_name

    # 2. Fallback: Local dev environment (.shared/os-arch/)
    if not lib_path.exists():
        machine = platform.machine().lower()
        # Map for Go-style architecture names
        go_arch = "amd64" if machine in ["x86_64", "amd64"] else "arm64"
        
        # Look up two levels (repo root) then into .shared
        fallback_path = pkg_dir.parents[1] / ".shared" / f"{system}-{go_arch}" / lib_name
        if fallback_path.exists():
            lib_path = fallback_path

    # 3. Final Fallback: Current Working Directory
    if not lib_path.exists():
        lib_path = Path.cwd() / lib_name

    try:
        # ffi.dlopen requires a string path
        return ffi.dlopen(str(lib_path))
    except OSError as e:
        raise ImportError(
            f"Could not load the enry shared library at {lib_path}.\n"
            f"System: {system}, Architecture: {platform.machine()}\n"
            "Ensure 'make shared' was run and the library is in the enry/ directory."
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