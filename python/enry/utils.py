from _c_enry import ffi
from enry.types import Guess
from functools import wraps
from typing import Hashable, List, Sequence


def py_bytes_to_go(py_bytes: bytes):
    c_bytes = ffi.new("char[]", py_bytes)
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


def init_go_slice():
    return ffi.new("GoSlice *")


def go_str_slice_to_py(str_slice) -> List[str]:
    slice_len = str_slice.len
    char_arr = ffi.cast("char **", str_slice.data)
    return [ffi.string(char_arr[i]).decode() for i in range(slice_len)]


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
            return go_to_py[out_type](fn(*(arg[0] for arg in args_transformed)))
        return inner
    return decorator


def transform_types_ret_str_slice(in_types: Sequence[Hashable]):
    def decorator(fn):
        @wraps(fn)
        def inner(*args):
            ret_slice = init_go_slice()
            args_transformed = [py_to_go[type_](arg) for type_, arg in zip(in_types, args)]
            fn(*(arg[0] for arg in args_transformed), ret_slice)
            return go_str_slice_to_py(ret_slice)
        return inner
    return decorator
