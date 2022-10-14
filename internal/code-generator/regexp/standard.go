// +build !oniguruma

package regexp

import (
	"regexp"
)

type ChosenRegexp = *regexp.Regexp

func Compile(str string) (ChosenRegexp, error) {
	return regexp.Compile(str)
}
