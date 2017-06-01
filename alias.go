package slinguist

import "strings"

// GetLanguageByAlias returns the language related to the given alias and ok set to true,
// or Otherlanguage and ok set to false otherwise.
func GetLanguageByAlias(alias string) (lang string, ok bool) {
	a := strings.Split(alias, `,`)[0]
	a = strings.ToLower(a)
	lang, ok = languagesByAlias[a]
	if !ok {
		lang = OtherLanguage
	}

	return
}
