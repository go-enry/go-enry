from enry._c_enry import ffi
from enry.types import Guess
from typing import List

def go_str_to_py(lib, c_ptr) -> str:
    """Converts a char* to Python str and frees the C memory via FreeCString."""
    if c_ptr == ffi.NULL:
        return ""
    try:
        # ffi.string converts the C pointer to a Python bytes object then we decode
        return ffi.string(c_ptr).decode('utf-8')
    finally:
        # Every string returned as C.CString in Go MUST be freed
        # We call Free on the lib object passed in from definitions.py
        lib.FreeCString(c_ptr)

def go_str_slice_to_py(lib, c_ptr_array) -> List[str]:
    """Converts a char** array to Python List[str] and frees it via FreeStringArray."""
    if c_ptr_array == ffi.NULL:
        return []
    results = []
    try:
        i = 0
        # We iterate until we find the NULL terminator we added in Go
        while c_ptr_array[i] != ffi.NULL:
            results.append(ffi.string(c_ptr_array[i]).decode('utf-8'))
            i += 1
        return results
    finally:
        # This one call frees the array AND all strings inside it
        # We call Free on the lib object passed in from definitions.py
        lib.FreeStringArray(c_ptr_array)

def prepare_candidates(candidates: List[str]):
    """Converts a Python list to a NULL-terminated char** array."""
    if not candidates:
        return ffi.NULL
    
    # Create the array of pointers (size + 1 for the NULL terminator)
    c_list = ffi.new("char*[]", len(candidates) + 1)
    
    # CFFI is smart: as long as c_strings stays in scope during the 
    # lib call, the memory is safe.
    c_strings = [ffi.new("char[]", c.encode('utf-8')) for c in candidates]
    for i, c_str in enumerate(c_strings):
        c_list[i] = c_str
    
    c_list[len(candidates)] = ffi.NULL
    return c_list

def go_guess_to_py(lib, c_ptr) -> Guess:
    """Standardizes the return of single-string 'Guess' functions."""
    lang = go_str_to_py(lib, c_ptr)
    # Original enry logic: if a language is returned, it's considered safe
    return Guess(language=lang, safe=bool(lang))