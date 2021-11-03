//go:build oniguruma
// +build oniguruma

package regex

import (
	rubex "github.com/go-enry/go-oniguruma"
)

type EnryRegexp = *rubex.Regexp

var Name = "oniguruma"

// MustCompileRuby same as MustCompile.
func MustCompileRuby(str string) EnryRegexp {
	return MustCompile(str)
}

func MustCompile(str string) EnryRegexp {
	return rubex.MustCompileASCII(str)
}

func QuoteMeta(s string) string {
	return rubex.QuoteMeta(s)
}
