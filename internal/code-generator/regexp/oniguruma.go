// +build oniguruma

package regexp

import (
	rubex "github.com/go-enry/go-oniguruma"
)

type ChosenRegexp = *rubex.Regexp

func Compile(str string) (ChosenRegexp, error) {
	return rubex.Compile(str)
}
