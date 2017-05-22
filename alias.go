package slinguist

import "strings"

// GetLanguageByAlias returns the language related to the given alias or Otherlanguage otherwise.
func GetLanguageByAlias(alias string) (lang string) {
	a := strings.Split(alias, `,`)[0]
	a = strings.ToLower(a)
	lang, ok := languagesByAlias[a]
	if !ok {
		lang = OtherLanguage
	}

	return
}
