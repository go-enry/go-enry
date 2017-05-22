package slinguist

import (
	"path/filepath"
	"strings"
)

func GetLanguageByExtension(filename string) (lang string, safe bool) {
	ext := strings.ToLower(filepath.Ext(filename))
	lang = OtherLanguage
	langs, ok := languagesByExtension[ext]
	if !ok {
		return
	}

	lang = langs[0]
	safe = len(langs) == 1
	return
}
