// +build darwin,cgo linux,cgo
// +build amd64

package main

import "C"
import "github.com/go-enry/go-enry/v2"
import "github.com/go-enry/go-enry/v2/data"

//export GetLanguage
func GetLanguage(filename string, content []byte) string {
	return enry.GetLanguage(filename, content)
}

//export GetLanguageByContent
func GetLanguageByContent(filename string, content []byte) (language string, safe bool) {
	return enry.GetLanguageByContent(filename, content)
}

//export GetLanguageByEmacsModeline
func GetLanguageByEmacsModeline(content []byte) (language string, safe bool) {
	return enry.GetLanguageByModeline(content)
}

//export GetLanguageByExtension
func GetLanguageByExtension(filename string) (language string, safe bool) {
	return enry.GetLanguageByExtension(filename)
}

//export GetLanguageByFilename
func GetLanguageByFilename(filename string) (language string, safe bool) {
	return enry.GetLanguageByFilename(filename)
}

//export GetLanguageByModeline
func GetLanguageByModeline(content []byte) (language string, safe bool) {
	return enry.GetLanguageByModeline(content)
}

//export GetLanguageByShebang
func GetLanguageByShebang(content []byte) (language string, safe bool) {
	return enry.GetLanguageByShebang(content)
}

//export GetLanguageByVimModeline
func GetLanguageByVimModeline(content []byte) (language string, safe bool) {
	return enry.GetLanguageByVimModeline(content)
}

//export GetLanguageExtensions
func GetLanguageExtensions(language string, result *[]*C.char) {
	strSliceCopy(result, enry.GetLanguageExtensions(language))
}

//export GetLanguages
func GetLanguages(filename string, content []byte, result *[]*C.char) {
	strSliceCopy(result, enry.GetLanguages(filename, content))
}

//export GetLanguagesByContent
func GetLanguagesByContent(filename string, content []byte, candidates []string, result *[]*C.char) {
	strSliceCopy(result, enry.GetLanguagesByContent(filename, content, candidates))
}

//export GetLanguagesByEmacsModeline
func GetLanguagesByEmacsModeline(filename string, content []byte, candidates []string, result *[]*C.char) {
	strSliceCopy(result, enry.GetLanguagesByEmacsModeline(filename, content, candidates))
}

//export GetLanguagesByExtension
func GetLanguagesByExtension(filename string, content []byte, candidates []string, result *[]*C.char) {
	strSliceCopy(result, enry.GetLanguagesByExtension(filename, content, candidates))
}

//export GetLanguagesByFilename
func GetLanguagesByFilename(filename string, content []byte, candidates []string, result *[]*C.char) {
	strSliceCopy(result, enry.GetLanguagesByFilename(filename, content, candidates))
}

//export GetLanguagesByModeline
func GetLanguagesByModeline(filename string, content []byte, candidates []string, result *[]*C.char) {
	strSliceCopy(result, enry.GetLanguagesByModeline(filename, content, candidates))
}

//export GetLanguagesByShebang
func GetLanguagesByShebang(filename string, content []byte, candidates []string, result *[]*C.char) {
	strSliceCopy(result, enry.GetLanguagesByShebang(filename, content, candidates))
}

//export GetLanguagesByVimModeline
func GetLanguagesByVimModeline(filename string, content []byte, candidates []string, result *[]*C.char) {
	strSliceCopy(result, enry.GetLanguagesByVimModeline(filename, content, candidates))
}

//export GetMimeType
func GetMimeType(path string, language string) string {
	return enry.GetMIMEType(path, language)
}

//export IsBinary
func IsBinary(data []byte) bool {
	return enry.IsBinary(data)
}

//export IsConfiguration
func IsConfiguration(path string) bool {
	return enry.IsConfiguration(path)
}

//export IsDocumentation
func IsDocumentation(path string) bool {
	return enry.IsDocumentation(path)
}

//export IsDotFile
func IsDotFile(path string) bool {
	return enry.IsDotFile(path)
}

//export IsImage
func IsImage(path string) bool {
	return enry.IsImage(path)
}

//export IsVendor
func IsVendor(path string) bool {
	return enry.IsVendor(path)
}

//export IsGenerated
func IsGenerated(path string, content []byte) bool {
	return enry.IsGenerated(path, content)
}

//export GetColor
func GetColor(language string) string {
	return enry.GetColor(language)
}

//export IsTest
func IsTest(path string) bool {
	return enry.IsTest(path)
}

//export GetLanguageType
func GetLanguageType(language string) string {
	return data.Type(enry.GetLanguageType(language)).String()
}

func strSliceCopy(result *[]*C.char, slice []string) {
	for _, str := range slice {
		*result = append(*result, C.CString(str))
	}
}

func main() {}