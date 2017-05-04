package slinguist

func GetLanguageByFilename(filename string) (lang string, safe bool) {
	lang, safe = languagesByFilename[filename]
	if lang == "" {
		lang = OtherLanguage
	}

	return
}
