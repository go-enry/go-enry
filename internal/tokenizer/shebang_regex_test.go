//go:build !flex
// +build !flex

package tokenizer

import "testing"

func TestShebangDoesNotConsumeNextLine(t *testing.T) {
	for name, src := range map[string]string{
		"LF":   "#!/usr/bin/perl\n-foo\nprint value\n",
	} {
		tokens := Tokenize([]byte(src))
		has := func(want string) bool {
			for _, tok := range tokens {
				if tok == want {
					return true
				}
			}
			return false
		}
		if !has("SHEBANG#!perl") {
			t.Errorf("%s: missing shebang token in %q", name, tokens)
		}
		if !has("foo") {
			t.Errorf("%s: body token \"foo\" was swallowed by the shebang match: %q", name, tokens)
		}
	}
}