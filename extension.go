package slinguist

import (
	"path/filepath"
)

func GetLanguageByExtension(filename string) (lang string, safe bool) {
	lang = OtherLanguage
	langs, ok := LanguagesByExtension[filepath.Ext(filename)]
	if !ok {
		return
	}

	lang = langs[0]
	if len(langs) == 1 {
		safe = true
	}

	return
}
