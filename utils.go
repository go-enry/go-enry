package enry

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-enry/go-enry/v2/data"
	"github.com/go-enry/go-enry/v2/regex"
)

const binSniffLen = 8000

var configurationLanguages = map[string]struct{}{
	"XML":  {},
	"JSON": {},
	"TOML": {},
	"YAML": {},
	"INI":  {},
	"SQL":  {},
}

// IsConfiguration tells if filename is in one of the configuration languages.
func IsConfiguration(path string) bool {
	language, _ := GetLanguageByExtension(path)
	_, is := configurationLanguages[language]
	return is
}

// IsImage tells if a given file is an image (PNG, JPEG or GIF format).
func IsImage(path string) bool {
	extension := filepath.Ext(path)
	if extension == ".png" || extension == ".jpg" || extension == ".jpeg" || extension == ".gif" {
		return true
	}

	return false
}

// GetMIMEType returns a MIME type of a given file based on its languages.
func GetMIMEType(path string, language string) string {
	if mime, ok := data.LanguagesMime[language]; ok {
		return mime
	}

	if IsImage(path) {
		return "image/" + filepath.Ext(path)[1:]
	}

	return "text/plain"
}

// IsDocumentation returns whether or not path is a documentation path.
func IsDocumentation(path string) bool {
	return matchRegexSlice(data.DocumentationMatchers, path)
}

// IsDotFile returns whether or not path has dot as a prefix.
func IsDotFile(path string) bool {
	base := filepath.Base(filepath.Clean(path))
	return strings.HasPrefix(base, ".") && base != "."
}

var allVendorRegExp regex.EnryRegexp

// IsVendor returns whether or not path is a vendor path.
func IsVendor(path string) bool {
	return allVendorRegExp.MatchString(path)
}

func init() {
	// We now collate all regexps from VendorMatchers to a single large regexp
	// which is at least twice as fast to test than simply iterating & matching.
	//
	// ---
	//
	// We could test each matcher from VendorMatchers in turn i.e.
	//
	//  	func IsVendor(filename string) bool {
	// 			for _, matcher := range data.VendorMatchers {
	// 				if matcher.MatchString(filename) {
	//					return true
	//				}
	//			}
	//			return false
	//		}
	//
	// Or naïvely concatentate all these regexps using groups i.e.
	//
	//		`(regexp1)|(regexp2)|(regexp3)|...`
	//
	// However, both of these are relatively slow and don't take advantage
	// of the inherent structure within our regexps.
	//
	// Imperical observation: by looking at the regexps, we only have 3 types.
	//  1. Those that start with `^`
	//  2. Those that start with `(^|/)`
	//  3. All the rest
	//
	// If we collate our regexps into these 3 groups - that will significantly
	// reduce the likelihood of backtracking within the regexp trie matcher.
	//
	// A further improvement is to use non-capturing groups (?:) as otherwise
	// the regexp parser, whilst matching, will have to allocate slices for
	// matching positions. (A future improvement left out could be to
	// enforce non-capturing groups within the sub-regexps.)

	matchers := data.VendorMatchers
	sort.SliceStable(matchers, func(i, j int) bool {
		return matchers[i].String() < matchers[j].String()
	})

	var caretPrefixed, caretOrSlashPrefixed, theRest []string
	// Check prefix, add to the respective group slices
	for _, matcher := range matchers {
		str := matcher.String()
		if strings.HasPrefix(str, "^") {
			caretPrefixed = append(caretPrefixed, str[1:])
		} else if strings.HasPrefix(str, "(^|/)") {
			caretOrSlashPrefixed = append(caretOrSlashPrefixed, str[5:])
		} else {
			theRest = append(theRest, str)
		}
	}
	var sb strings.Builder
	// group 1 - start with `^`
	appendGroupWithCommonPrefix(&sb, "^", caretPrefixed)
	sb.WriteString("|")
	// group 2 - start with `(^|/)`
	appendGroupWithCommonPrefix(&sb, "(?:^|/)", caretOrSlashPrefixed)
	sb.WriteString("|")
	// grou 3, all rest.
	appendGroupWithCommonPrefix(&sb, "", theRest)
	allVendorRegExp = regex.MustCompile(sb.String())
}

func appendGroupWithCommonPrefix(sb *strings.Builder, commonPrefix string, res []string) {
	sb.WriteString("(?:")
	if commonPrefix != "" {
		sb.WriteString(fmt.Sprintf("%s(?:(?:", commonPrefix))
	}
	sb.WriteString(strings.Join(res, ")|(?:"))
	if commonPrefix != "" {
		sb.WriteString("))")
	}
	sb.WriteString(")")
}

// IsTest returns whether or not path is a test path.
func IsTest(path string) bool {
	return matchRegexSlice(data.TestMatchers, path)
}

func matchRegexSlice(exprs []regex.EnryRegexp, str string) bool {
	for _, expr := range exprs {
		if expr.MatchString(str) {
			return true
		}
	}

	return false
}

// IsBinary detects if data is a binary value based on:
// http://git.kernel.org/cgit/git/git.git/tree/xdiff-interface.c?id=HEAD#n198
func IsBinary(data []byte) bool {
	if len(data) > binSniffLen {
		data = data[:binSniffLen]
	}

	if bytes.IndexByte(data, byte(0)) == -1 {
		return false
	}

	return true
}

// GetColor returns a HTML color code of a given language.
func GetColor(language string) string {
	if color, ok := data.LanguagesColor[language]; ok {
		return color
	}

	if color, ok := data.LanguagesColor[GetLanguageGroup(language)]; ok {
		return color
	}

	return "#cccccc"
}

// IsGenerated returns whether the file with the given path and content is a
// generated file.
func IsGenerated(path string, content []byte) bool {
	ext := strings.ToLower(filepath.Ext(path))
	if _, ok := data.GeneratedCodeExtensions[ext]; ok {
		return true
	}

	for _, m := range data.GeneratedCodeNameMatchers {
		if m(path) {
			return true
		}
	}

	path = strings.ToLower(path)
	for _, m := range data.GeneratedCodeMatchers {
		if m(path, ext, content) {
			return true
		}
	}

	return false
}
