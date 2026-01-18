//go:build (darwin && cgo) || (linux && cgo)
// +build darwin,cgo linux,cgo

package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/go-enry/go-enry/v2"
	"github.com/go-enry/go-enry/v2/data"
)

//export FreeCString
func FreeCString(str *C.char) {
	C.free(unsafe.Pointer(str))
}

//export FreeStringArray
func FreeStringArray(ptr **C.char) {
	if ptr == nil {
		return
	}
	cArray := (*[1 << 28]*C.char)(unsafe.Pointer(ptr))
	for i := 0; cArray[i] != nil; i++ {
		C.free(unsafe.Pointer(cArray[i]))
	}
	C.free(unsafe.Pointer(ptr))
}

// --- Internal Helpers ---

func toGoBytes(content *C.char, length C.int) []byte {
	if content == nil || length <= 0 {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(content), length)
}

// strArrayToC converts a Go []string into a C-style char** array.
// It allocates memory on the C heap that must be freed by FreeStringArray.
func strArrayToC(langs []string) **C.char {
	size := unsafe.Sizeof((*C.char)(nil))
	// Allocate space for N pointers + 1 NULL terminator
	ptr := C.malloc(C.size_t(len(langs)+1) * C.size_t(size))
	cArray := (*[1 << 28]*C.char)(ptr)

	for i, s := range langs {
		cArray[i] = C.CString(s)
	}
	cArray[len(langs)] = nil // Set the NULL terminator
	return (**C.char)(ptr)
}

func cArrayToGoSlice(cArray **C.char) []string {
	if cArray == nil {
		return nil
	}
	var slice []string
	// Iterate until we hit the NULL terminator we'll set in Python
	p := unsafe.Pointer(cArray)
	for *(*unsafe.Pointer)(p) != nil {
		slice = append(slice, C.GoString(*(**C.char)(p)))
		p = unsafe.Pointer(uintptr(p) + unsafe.Sizeof(uintptr(0)))
	}
	return slice
}

// --- Singular API ---

//export GetLanguage
func GetLanguage(filename *C.char, content *C.char, length C.int) *C.char {
	return C.CString(enry.GetLanguage(C.GoString(filename), toGoBytes(content, length)))
}

//export GetLanguageByContent
func GetLanguageByContent(filename *C.char, content *C.char, length C.int) *C.char {
	lang, _ := enry.GetLanguageByContent(C.GoString(filename), toGoBytes(content, length))
	return C.CString(lang)
}

//export GetLanguageByEmacsModeline
func GetLanguageByEmacsModeline(content *C.char, length C.int) *C.char {
	lang, _ := enry.GetLanguageByModeline(toGoBytes(content, length))
	return C.CString(lang)
}

//export GetLanguageByExtension
func GetLanguageByExtension(filename *C.char) *C.char {
	lang, _ := enry.GetLanguageByExtension(C.GoString(filename))
	return C.CString(lang)
}

//export GetLanguageByFilename
func GetLanguageByFilename(filename *C.char) *C.char {
	lang, _ := enry.GetLanguageByFilename(C.GoString(filename))
	return C.CString(lang)
}

//export GetLanguageByModeline
func GetLanguageByModeline(content *C.char, length C.int) *C.char {
	lang, _ := enry.GetLanguageByModeline(toGoBytes(content, length))
	return C.CString(lang)
}

//export GetLanguageByShebang
func GetLanguageByShebang(content *C.char, length C.int) *C.char {
	lang, _ := enry.GetLanguageByShebang(toGoBytes(content, length))
	return C.CString(lang)
}

//export GetLanguageByVimModeline
func GetLanguageByVimModeline(content *C.char, length C.int) *C.char {
	lang, _ := enry.GetLanguageByVimModeline(toGoBytes(content, length))
	return C.CString(lang)
}

//export GetMimeType
func GetMimeType(path *C.char, language *C.char) *C.char {
	return C.CString(enry.GetMIMEType(C.GoString(path), C.GoString(language)))
}

//export GetColor
func GetColor(language *C.char) *C.char {
	return C.CString(enry.GetColor(C.GoString(language)))
}

//export GetLanguageType
func GetLanguageType(language *C.char) *C.char {
	return C.CString(data.Type(enry.GetLanguageType(C.GoString(language))).String())
}

// --- Boolean API (Bools return C.int: 1 for true, 0 for false) ---

//export IsBinary
func IsBinary(content *C.char, length C.int) C.int {
	if enry.IsBinary(toGoBytes(content, length)) {
		return 1
	}
	return 0
}

//export IsConfiguration
func IsConfiguration(path *C.char) C.int {
	if enry.IsConfiguration(C.GoString(path)) {
		return 1
	}
	return 0
}

//export IsDocumentation
func IsDocumentation(path *C.char) C.int {
	if enry.IsDocumentation(C.GoString(path)) {
		return 1
	}
	return 0
}

//export IsDotFile
func IsDotFile(path *C.char) C.int {
	if enry.IsDotFile(C.GoString(path)) {
		return 1
	}
	return 0
}

//export IsImage
func IsImage(path *C.char) C.int {
	if enry.IsImage(C.GoString(path)) {
		return 1
	}
	return 0
}

//export IsVendor
func IsVendor(path *C.char) C.int {
	if enry.IsVendor(C.GoString(path)) {
		return 1
	}
	return 0
}

//export IsGenerated
func IsGenerated(path *C.char, content *C.char, length C.int) C.int {
	if enry.IsGenerated(C.GoString(path), toGoBytes(content, length)) {
		return 1
	}
	return 0
}

//export IsTest
func IsTest(path *C.char) C.int {
	if enry.IsTest(C.GoString(path)) {
		return 1
	}
	return 0
}

// --- Plural API (Returning arrays) ---

//export GetLanguages
func GetLanguages(filename *C.char, content *C.char, length C.int) **C.char {
	langs := enry.GetLanguages(C.GoString(filename), toGoBytes(content, length))
	return strArrayToC(langs)
}

//export GetLanguageExtensions
func GetLanguageExtensions(language *C.char) **C.char {
	langs := enry.GetLanguageExtensions(C.GoString(language))
	return strArrayToC(langs)
}

//export GetLanguagesByContent
func GetLanguagesByContent(filename *C.char, content *C.char, length C.int, candidates **C.char) **C.char {
	langs := enry.GetLanguagesByContent(C.GoString(filename), toGoBytes(content, length), cArrayToGoSlice(candidates))
	return strArrayToC(langs)
}

//export GetLanguagesByEmacsModeline
func GetLanguagesByEmacsModeline(filename *C.char, content *C.char, length C.int, candidates **C.char) **C.char {
	langs := enry.GetLanguagesByEmacsModeline(C.GoString(filename), toGoBytes(content, length), cArrayToGoSlice(candidates))
	return strArrayToC(langs)
}

//export GetLanguagesByExtension
func GetLanguagesByExtension(filename *C.char, content *C.char, length C.int, candidates **C.char) **C.char {
	langs := enry.GetLanguagesByExtension(C.GoString(filename), toGoBytes(content, length), cArrayToGoSlice(candidates))
	return strArrayToC(langs)
}

//export GetLanguagesByFilename
func GetLanguagesByFilename(filename *C.char, content *C.char, length C.int, candidates **C.char) **C.char {
	langs := enry.GetLanguagesByFilename(C.GoString(filename), toGoBytes(content, length), cArrayToGoSlice(candidates))
	return strArrayToC(langs)
}

//export GetLanguagesByModeline
func GetLanguagesByModeline(filename *C.char, content *C.char, length C.int, candidates **C.char) **C.char {
	langs := enry.GetLanguagesByModeline(C.GoString(filename), toGoBytes(content, length), cArrayToGoSlice(candidates))
	return strArrayToC(langs)
}

//export GetLanguagesByShebang
func GetLanguagesByShebang(filename *C.char, content *C.char, length C.int, candidates **C.char) **C.char {
	langs := enry.GetLanguagesByShebang(C.GoString(filename), toGoBytes(content, length), cArrayToGoSlice(candidates))
	return strArrayToC(langs)
}

//export GetLanguagesByVimModeline
func GetLanguagesByVimModeline(filename *C.char, content *C.char, length C.int, candidates **C.char) **C.char {
	langs := enry.GetLanguagesByVimModeline(C.GoString(filename), toGoBytes(content, length), cArrayToGoSlice(candidates))
	return strArrayToC(langs)
}

func main() {}
