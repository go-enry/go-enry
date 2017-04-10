package slinguist

import (
	"bytes"
	"path/filepath"
	"strings"

	"gopkg.in/toqueteos/substring.v1"
)

func IsAuxiliaryLanguage(lang string) bool {
	_, ok := auxiliaryLanguages[lang]
	return ok
}

func IsConfiguration(path string) bool {
	lang, _ := GetLanguageByExtension(path)
	_, is := configurationLanguages[lang]

	return is
}

func IsDotFile(path string) bool {
	return strings.HasPrefix(filepath.Base(path), ".")
}

func IsVendor(path string) bool {
	return findIndex(path, vendorMatchers) >= 0
}

func IsDocumentation(path string) bool {
	return findIndex(path, documentationMatchers) >= 0
}

func findIndex(path string, matchers substring.StringsMatcher) int {
	return matchers.MatchIndex(path)
}

const sniffLen = 8000

//IsBinary detects if data is a binary value based on:
//http://git.kernel.org/cgit/git/git.git/tree/xdiff-interface.c?id=HEAD#n198
func IsBinary(data []byte) bool {
	if len(data) > sniffLen {
		data = data[:sniffLen]
	}

	if bytes.IndexByte(data, byte(0)) == -1 {
		return false
	}

	return true
}

var configurationLanguages = map[string]bool{
	"XML": true, "JSON": true, "TOML": true, "YAML": true, "INI": true, "SQL": true,
}
