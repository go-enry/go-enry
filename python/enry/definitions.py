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
    """
    Return the language of the given file based on the filename and its contents.

    :param filename: name of the file with the extension
    :param content: array of bytes with the contents of the file (the code)
    :return: the guessed language
    """
    res = lib.GetLanguage(filename.encode(), content, len(content))
    return go_str_to_py(lib, res)

def get_language_by_content(filename: str, content: bytes) -> Guess:
    """
    Return detected language by its content.
    If there are more than one possible language, return the first language
    in alphabetical order and safe = False.

    :param filename: path of the file
    :param content: array of bytes with the contents of the file (the code)
    :return: guessed result
    """
    res = lib.GetLanguageByContent(filename.encode(), content, len(content))
    return go_guess_to_py(lib, res)

def get_language_by_extension(filename: str) -> Guess:
    """
    Return detected language by the extension of the filename.
    If there are more than one possible language return the first language
    in alphabetical order and safe = False.

    :param filename: path of the file
    :return: guessed result
    """
    res = lib.GetLanguageByExtension(filename.encode())
    return go_guess_to_py(lib, res)

def get_language_by_filename(filename: str) -> Guess:
    """
    Return detected language by its filename.
    If there are more than one possible language return the first language
    in alphabetical order and safe = False.

    :param filename: path of the file
    :return: guessed result
    """
    res = lib.GetLanguageByFilename(filename.encode())
    return go_guess_to_py(lib, res)

def get_language_by_modeline(content: bytes) -> Guess:
    """
    Return detected language by its modeline.
    If there are more than one possible language return the first language
    in alphabetical order and safe = False.

    :param content: array of bytes with the contents of the file (the code)
    :return: guessed result
    """
    res = lib.GetLanguageByModeline(content, len(content))
    return go_guess_to_py(lib, res)

def get_language_by_shebang(content: bytes) -> Guess:
    """
    Return detected langauge by its shebang.
    If there are more than one possible language return the first language
    in alphabetical order and safe = False.

    :param content: array of bytes with the contents of the file (the code)
    :return: guessed result
    """
    res = lib.GetLanguageByShebang(content, len(content))
    return go_guess_to_py(lib, res)

def get_language_by_emacs_modeline(content: bytes) -> Guess:
    """
    Return detected langauge by its emacs modeline.
    If there are more than one possible language return the first language
    in alphabetical order and safe = False.

    :param content: array of bytes with the contents of the file (the code)
    :return: guessed result
    """
    res = lib.GetLanguageByEmacsModeline(content, len(content))
    return go_guess_to_py(lib, res)

def get_language_by_vim_modeline(content: bytes) -> Guess:
    """
    Return detected language by its vim modeline.
    If there are more than one possible language return the first language
    in alphabetical order and safe = False.

    :param content: array of bytes with the contents of the file (the code)
    :return: guessed result
    """
    res = lib.GetLanguageByVimModeline(content, len(content))
    return go_guess_to_py(lib, res)

def get_mime_type(path: str, language: str) -> str:
    """
    Return mime type of the file.

    :param path: path of the file
    :param language: language to get mime type from
    :return: mime type
    """
    res = lib.GetMimeType(path.encode(), language.encode())
    return go_str_to_py(lib, res)

def get_color(language: str) -> str:
    """
    Return color code for given language

    :param language:
    :return: color in hex format
    """
    res = lib.GetColor(language.encode())
    return go_str_to_py(lib, res)

def get_language_type(language: str) -> str:
    """
    Docstring for get_language_type
    
    :param language:
    :type language: str
    :return: type of the language
    :rtype: str
    """
    res = lib.GetLanguageType(language.encode())
    return go_str_to_py(lib, res)

# --- Boolean API ---

def is_binary(content: bytes) -> bool:
    """
    Docstring for is_binary
    
    :param content: array of bytes with the contents of the file (the code)
    :type content: bytes
    :return: whether the file is binary or not
    :rtype: bool
    """
    return bool(lib.IsBinary(content, len(content)))

def is_vendor(filename: str) -> bool:
    """
    Docstring for is_vendor
    
    :param filename: Description
    :type filename: str
    :return: Description
    :rtype: bool
    """
    return bool(lib.IsVendor(filename.encode()))

def is_generated(filename: str, content: bytes) -> bool:
    """
    Docstring for is_generated
    
    :param filename: 
    :type filename: str
    :param content: array of bytes with the contents of the file (the code)
    :type content: bytes
    :return: whether the file is generated or not
    :rtype: bool
    """
    return bool(lib.IsGenerated(filename.encode(), content, len(content)))

def is_configuration(path: str) -> bool:
    """
    Docstring for is_configuration
    
    :param path: path of the file
    :type path: str
    :return: whether the file is a configuration file or not
    :rtype: bool
    """
    return bool(lib.IsConfiguration(path.encode()))

def is_documentation(path: str) -> bool:
    """
    Docstring for is_documentation
    
    :param path: path of the file
    :type path: str
    :return: whether the file is a documentation file or not
    :rtype: bool
    """
    return bool(lib.IsDocumentation(path.encode()))

def is_dot_file(path: str) -> bool:
    """
    Docstring for is_dot_file
    
    :param path: path of the file
    :type path: str
    :return: whether the file is a dot file or not
    :rtype: bool
    """
    return bool(lib.IsDotFile(path.encode()))

def is_image(path: str) -> bool:
    """
    Docstring for is_image
    
    :param path: path of the file
    :type path: str
    :return: whether the file is an image or not
    :rtype: bool
    """
    return bool(lib.IsImage(path.encode()))

def is_test(path: str) -> bool:
    """
    Docstring for is_test
    
    :param path: path of the file
    :type path: str
    :return: whether the file is a test file or not
    :rtype: bool
    """
    return bool(lib.IsTest(path.encode()))

# --- Plural API ---

def get_languages(filename: str, content: bytes) -> List[str]:
    """
    Docstring for get_languages
    
    :param filename: name of the file with the extension
    :type filename: str
    :param content: array of bytes with the contents of the file (the code)
    :type content: bytes
    :return: list of languages that are detected in the file
    :rtype: List[str]
    """
    res = lib.GetLanguages(filename.encode(), content, len(content))
    return go_str_slice_to_py(lib, res)

def get_language_extensions(language: str) -> List[str]:
    """
    Docstring for get_language_extensions
    
    :param language: name of the programming language
    :type language: str
    :return: list of extensions for the given language
    :rtype: List[str]
    """
    res = lib.GetLanguageExtensions(language.encode())
    return go_str_slice_to_py(lib, res)

def get_languages_by_filename(filename: str, content: bytes = b"", candidates: List[str] = None) -> List[str]:
    """
    Docstring for get_languages_by_filename
    
    :param filename: name of the file with the extension
    :type filename: str
    :param content: array of bytes with the contents of the file (the code)
    :type content: bytes
    :param candidates: list of candidate languages to consider
    :type candidates: List[str]
    :return: list of languages that are detected in the file
    :rtype: List[str]
    """
    # If candidates is explicitly an empty list, short-circuit and return no languages.
    # This mirrors the intended semantics where None means "no filter" but [] means "no candidates".
    if candidates == []:
        return []

    c_cand, _keep_alive = prepare_candidates(candidates)
    res = lib.GetLanguagesByFilename(filename.encode(), content, len(content), c_cand)
    return go_str_slice_to_py(lib, res)

def get_languages_by_shebang(filename: str, content: bytes = b"", candidates: List[str] = None) -> List[str]:
    """
    Docstring for get_languages_by_shebang
    
    :param filename: name of the file with the extension
    :type filename: str
    :param content: array of bytes with the contents of the file (the code)
    :type content: bytes
    :param candidates: list of candidate languages to consider
    :type candidates: List[str]
    :return: list of languages that are detected in the file by shebang
    :rtype: List[str]
    """
    # If candidates is explicitly an empty list, short-circuit and return no languages.
    # (The underlying Go strategy ignores the candidates parameter, so handle this at
    #  the Python layer to provide the intended semantics.)
    if candidates == []:
        return []

    c_cand, _keep_alive = prepare_candidates(candidates)
    res = lib.GetLanguagesByShebang(filename.encode(), content, len(content), c_cand)
    return go_str_slice_to_py(lib, res)