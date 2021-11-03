//go:build !oniguruma
// +build !oniguruma

package regex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustCompileRuby(t *testing.T) {
	re := MustCompileRuby(`a`)
	assert.Equal(t, "(?m)a", re.String())
}

func TestMustCompileRuby_NotSupported(t *testing.T) {
	re := MustCompileRuby(`\A(?=\w{6,10}\z)`)
	assert.Equal(t, EnryRegexp(nil), re)
}
