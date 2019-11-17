from cffi import FFI
ffibuilder = FFI()

# cdef() expects a single string declaring the C types, functions and
# globals needed to use the shared object. It must be in valid C syntax.
# Taken from java/shared/libenry.h
ffibuilder.cdef("""
typedef unsigned char GoUint8;
typedef long long GoInt64;
typedef GoInt64 GoInt;

typedef struct { const char *p; ptrdiff_t n; } _GoString_;
typedef _GoString_ GoString;

typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;


extern GoString GetLanguage(GoString p0, GoSlice p1);

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

extern GoUint8 IsVendor(GoString p0);
""")

# set_source() gives the name of the python extension module to
# produce, and some C source code as a string.  This C code needs
# to make the declarated functions, types and globals available,
# so it is often just the "#include".
ffibuilder.set_source("_c_enry",
                      """
     #include "../.shared/libenry.h"   // the C header of the library
""",
                      libraries=['enry'],
                      library_dirs=['../.shared'
                                    ])  # library name, for the linker

if __name__ == "__main__":
    ffibuilder.compile(verbose=True)
