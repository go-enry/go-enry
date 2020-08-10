from cffi import FFI
import os
from pathlib import Path

ffibuilder = FFI()

# cdef() expects a single string declaring the C types, functions and
# globals needed to use the shared object. It must be in valid C syntax.
# Taken from java/shared/libenry.h
ffibuilder.cdef(
    """
typedef unsigned char GoUint8;
typedef long long GoInt64;
typedef GoInt64 GoInt;

typedef struct { const char *p; ptrdiff_t n; } _GoString_;
typedef _GoString_ GoString;

typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;


extern GoString GetLanguage(GoString p0, GoSlice p1);

/* Return type for GetLanguageByContent */
struct GetLanguageByContent_return {
    GoString r0; /* language */
    GoUint8 r1; /* safe */
};

extern struct GetLanguageByContent_return GetLanguageByContent(GoString p0, GoSlice p1);

/* Return type for GetLanguageByEmacsModeline */
struct GetLanguageByEmacsModeline_return {
    GoString r0; /* language */
    GoUint8 r1; /* safe */
};

extern struct GetLanguageByEmacsModeline_return GetLanguageByEmacsModeline(GoSlice p0);

/* Return type for GetLanguageByExtension */
struct GetLanguageByExtension_return {
    GoString r0; /* language */
    GoUint8 r1; /* safe */
};

extern struct GetLanguageByExtension_return GetLanguageByExtension(GoString p0);

/* Return type for GetLanguageByFilename */
struct GetLanguageByFilename_return {
    GoString r0; /* language */
    GoUint8 r1; /* safe */
};

extern struct GetLanguageByFilename_return GetLanguageByFilename(GoString p0);

/* Return type for GetLanguageByModeline */
struct GetLanguageByModeline_return {
    GoString r0; /* language */
    GoUint8 r1; /* safe */
};

extern struct GetLanguageByModeline_return GetLanguageByModeline(GoSlice p0);

/* Return type for GetLanguageByShebang */
struct GetLanguageByShebang_return {
    GoString r0; /* language */
    GoUint8 r1; /* safe */
};

extern struct GetLanguageByShebang_return GetLanguageByShebang(GoSlice p0);

/* Return type for GetLanguageByVimModeline */
struct GetLanguageByVimModeline_return {
    GoString r0; /* language */
    GoUint8 r1; /* safe */
};

extern struct GetLanguageByVimModeline_return GetLanguageByVimModeline(GoSlice p0);

extern void GetLanguageExtensions(GoString p0, GoSlice* p1);

extern void GetLanguages(GoString p0, GoSlice p1, GoSlice* p2);

extern void GetLanguagesByContent(GoString p0, GoSlice p1, GoSlice p2, GoSlice* p3);

extern void GetLanguagesByEmacsModeline(GoString p0, GoSlice p1, GoSlice p2, GoSlice* p3);

extern void GetLanguagesByExtension(GoString p0, GoSlice p1, GoSlice p2, GoSlice* p3);

extern void GetLanguagesByFilename(GoString p0, GoSlice p1, GoSlice p2, GoSlice* p3);

extern void GetLanguagesByModeline(GoString p0, GoSlice p1, GoSlice p2, GoSlice* p3);

extern void GetLanguagesByShebang(GoString p0, GoSlice p1, GoSlice p2, GoSlice* p3);

extern void GetLanguagesByVimModeline(GoString p0, GoSlice p1, GoSlice p2, GoSlice* p3);

extern GoString GetMimeType(GoString p0, GoString p1);

extern GoUint8 IsBinary(GoSlice p0);

extern GoUint8 IsConfiguration(GoString p0);

extern GoUint8 IsDocumentation(GoString p0);

extern GoUint8 IsDotFile(GoString p0);

extern GoUint8 IsImage(GoString p0);

extern GoUint8 IsVendor(GoString p0);

extern GoUint8 IsGenerated(GoString p0, GoSlice p1);

extern GoString GetColor(GoString p0);
"""
)

# set_source() gives the name of the python extension module to
# produce, and some C source code as a string.  This C code needs
# to make the declarated functions, types and globals available,
# so it is often just the "#include".
lib_dir = Path(__file__).resolve().parent.parent / ".shared"
lib_header = lib_dir / "libenry.h"

ffibuilder.set_source(
    "_c_enry",
    f'#include "{lib_header.absolute()}"',
    libraries=["enry"],
    library_dirs=[str(lib_dir.absolute())],
)  # library name, for the linker


if __name__ == "__main__":
    ffibuilder.compile(verbose=True)
