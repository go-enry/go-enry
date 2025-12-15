from _c_enry import ffi, lib
from enry.types import Guess
from functools import wraps
import gc
from typing import Hashable, List, Sequence


def py_bytes_to_go(py_bytes: bytes):
    """Create a C buffer and Go slice pointing to `py_bytes`.

    Returns a tuple `(GoSlice, cdata_buffer)` where the first element
    is a `GoSlice` suitable for passing to the Go FFI and the second
    is the underlying `cdata` buffer (must be kept alive while the
    Go call is in progress).
    """
    c_bytes = ffi.new("char[]", py_bytes)
    go_slice = ffi.new("GoSlice *", [c_bytes, len(py_bytes), len(py_bytes)])
    return (go_slice[0], c_bytes)


def py_str_to_go(py_str: str):
    """Create a GoString and C buffer for `py_str`.

    Returns `(GoString, cdata_buffer)`. Keep the `cdata_buffer`
    alive until the FFI call completes to avoid use-after-free.
    """
    str_bytes = py_str.encode()
    c_str = ffi.new("char[]", str_bytes)
    go_str = ffi.new("_GoString_ *", [c_str, len(str_bytes)])
    return (go_str[0], c_str)


def go_str_to_py(go_str: str):
    """Convert a Go string to a Python `str`.

    Copies bytes from the Go string (`go_str.p`, `go_str.n`) into
    Python memory and returns a decoded `str`. Do NOT free the
    pointer here — the pointer is owned by the Go/cgo layer and
    must not be freed from Python. Only C-allocated buffers created
    with `C.CString` are freed on the Python side (see
    `go_str_slice_to_py`).
    """
    str_len = go_str.n
    if str_len <= 0 or go_str.p == ffi.NULL:
        return ""

    # Copy bytes out of the C buffer. Do NOT free the pointer here —
    # single-string returns are owned by Go/FFI wrapper and must not
    # be freed with C.free. Freeing them can cause invalid pointer
    # errors. Only free pointers that were allocated via C.CString
    # (handled in slice conversion).
    return ffi.unpack(go_str.p, str_len).decode()


def init_go_slice():
    """Allocate and return a pointer to an empty `GoSlice`.

    Used as an out-parameter for FFI functions that populate a slice
    (e.g. functions returning arrays of C strings).
    """
    return ffi.new("GoSlice *")


def go_str_slice_to_py(str_slice) -> List[str]:
    """Convert a C array of C strings (GoSlice) to Python `List[str]`.

    This function copies each C string into Python memory and frees
    the original C buffer using `lib.FreeCString` because those
    buffers are allocated by helpers using `C.CString` on the Go side.
    """
    slice_len = str_slice.len
    char_arr = ffi.cast("char **", str_slice.data)
    result = []
    for i in range(slice_len):
        p = char_arr[i]
        if p == ffi.NULL:
            result.append("")
            continue
        s = ffi.string(p).decode()
        result.append(s)
        # Free memory allocated by C.CString in Go helper
        lib.FreeCString(p)
    return result


def go_bool_to_py(go_bool: bool):
    return go_bool == 1


def go_guess_to_py(guess) -> Guess:
    return Guess(go_str_to_py(guess.r0), go_bool_to_py(guess.r1))


py_to_go = {
    str: py_str_to_go,
    bytes: py_bytes_to_go,
}


go_to_py = {
    str: go_str_to_py,
    bool: go_bool_to_py,
    Guess: go_guess_to_py,
}


def transform_types(in_types: Sequence[Hashable], out_type: Hashable):
    def decorator(fn):
        @wraps(fn)
        def inner(*args):
            args_transformed = [py_to_go[type_](arg) for type_, arg in zip(in_types, args)]
            c_args = [arg[0] for arg in args_transformed]
            # release Python-side references to C buffers so they can be GC'd
            del args_transformed
            res = fn(*c_args)
            # force immediate collection of cdata buffers allocated by cffi
            gc.collect()
            return go_to_py[out_type](res)
        return inner
    return decorator

    # Note: the decorator converts Python inputs to C-compatible
    # representations (keeping underlying cdata alive for the duration
    # of the call) and converts the return value back to Python. It
    # triggers a `gc.collect()` after the FFI call to speed up
    # reclamation of temporary cffi allocations during tight loops.

def transform_types_ret_str_slice(in_types: Sequence[Hashable]):
    def decorator(fn):
        @wraps(fn)
        def inner(*args):
            ret_slice = init_go_slice()
            args_transformed = [py_to_go[type_](arg) for type_, arg in zip(in_types, args)]
            c_args = [arg[0] for arg in args_transformed]
            del args_transformed
            fn(*c_args, ret_slice)
            # free temporary cffi allocations
            gc.collect()
            return go_str_slice_to_py(ret_slice)
        return inner
    return decorator
