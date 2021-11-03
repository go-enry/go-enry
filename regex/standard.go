//go:build !oniguruma
// +build !oniguruma

package regex

import (
	"regexp"
	"strings"
)

type EnryRegexp = *regexp.Regexp

var Name = "standard"

const multilinePrefix = "(?m)"

// MustCompileRuby returns a new EnryRegexp based on the given Ruby regexp, if
// the regexp is not supported by the standard library regexp engine, returns a
// EnryRegexp(nil). If supported it compiles the regexp converting it in a
// multiline regexp.
func MustCompileRuby(str string) EnryRegexp {
	if isUnsupportedRegexpSyntax(str) {
		return nil
	}

	return regexp.MustCompile(convertToValidRegexp(str))
}

func MustCompile(str string) EnryRegexp {
	return regexp.MustCompile(str)
}

func QuoteMeta(s string) string {
	return regexp.QuoteMeta(s)
}

// isUnsupportedRegexpSyntax filters regexp syntax that is not supported by RE2.
// In particular, we stumbled up on usage of next cases:
// - lookbehind & lookahead
// - named & numbered capturing group/after text matching
// - backreference
// - possessive quantifier
// For referece on supported syntax see https://github.com/google/re2/wiki/Syntax
func isUnsupportedRegexpSyntax(reg string) bool {
	return strings.Contains(reg, `(?<`) || strings.Contains(reg, `(?=`) || strings.Contains(reg, `(?!`) ||
		strings.Contains(reg, `\1`) || strings.Contains(reg, `*+`) ||
		// See https://github.com/github/linguist/pull/4243#discussion_r246105067
		(strings.HasPrefix(reg, multilinePrefix+`/`) && strings.HasSuffix(reg, `/`))
}

// convertToValidRegexp converts Ruby regexp syntaxt to RE2 equivalent.
// Does not work with Ruby regexp literals.
func convertToValidRegexp(rubyRegexp string) string {
	return multilinePrefix + rubyRegexp
}
