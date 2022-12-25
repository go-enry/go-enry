//go:build !oniguruma
// +build !oniguruma

package regex

import (
	"regexp"
)

type EnryRegexp = *regexp.Regexp

func MustCompile(str string) EnryRegexp {
	return regexp.MustCompile(str)
}

// MustCompileMultiline mimics Ruby defaults for regexp, where ^$ matches begin/end of line.
// I.e. it converts Ruby regexp syntaxt to RE2 equivalent
func MustCompileMultiline(s string) EnryRegexp {
	const multilineModeFlag = "(?m)"
	return regexp.MustCompile(multilineModeFlag + s)
}

// MustCompileRuby used for expressions with syntax not supported by RE2.
func MustCompileRuby(s string) EnryRegexp {
	// TODO(bzz): find a bettee way?
	// This will only trigger a panic on .Match() for the clients
	return nil
}

func QuoteMeta(s string) string {
	return regexp.QuoteMeta(s)
}
