//go:build !oniguruma
// +build !oniguruma

package regex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustCompileMultiline(t *testing.T) {
	const re = `^\.(.*)!$`
	want := MustCompileMultiline(re)
	assert.Equal(t, "(?m)"+re, want.String())

	const s = `.one
.two!
thre!`
	if !want.MatchString(s) {
		t.Fatalf("MustCompileMultiline(`%s`) must match multiline %q\n", re, s)
	}
}

func TestMustCompileRuby(t *testing.T) {
	assert.Nil(t, MustCompileRuby(``))
}
