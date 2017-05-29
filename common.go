package slinguist

import (
	"path/filepath"
	"strings"
)

// OtherLanguage is used as a zero value when a function can not return a specific language.
const OtherLanguage = "Other"

// GetLanguage applies a sequence of strategies based on the given filename and content
// to find out the most probably language to return.
func GetLanguage(filename string, content []byte) string {
	if lang, safe := GetLanguageByModeline(content); safe {
		return lang
	}

	if lang, safe := GetLanguageByFilename(filename); safe {
		return lang
	}

	if lang, safe := GetLanguageByShebang(content); safe {
		return lang
	}

	if lang, safe := GetLanguageByExtension(filename); safe {
		return lang
	}

	if lang, safe := GetLanguageByContent(filename, content); safe {
		return lang
	}

	lang := GetLanguageByClassifier(content, nil, nil)
	return lang
}

// GetLanguageByModeline returns the language of the given content looking for the modeline,
// and safe to indicate the sureness of returned language.
func GetLanguageByModeline(content []byte) (lang string, safe bool) {
	return getLanguageByModeline(content)
}

// GetLanguageByFilename returns a language based on the given filename, and safe to indicate
// the sureness of returned language.
func GetLanguageByFilename(filename string) (lang string, safe bool) {
	return getLanguageByFilename(filename)
}

func getLanguageByFilename(filename string) (lang string, safe bool) {
	lang, safe = languagesByFilename[filename]
	if lang == "" {
		lang = OtherLanguage
	}

	return
}

// GetLanguageByShebang returns the language of the given content looking for the shebang line,
// and safe to indicate the sureness of returned language.
func GetLanguageByShebang(content []byte) (lang string, safe bool) {
	return getLanguageByShebang(content)
}

// GetLanguageByExtension returns a language based on the given filename, and safe to indicate
// the sureness of returned language.
func GetLanguageByExtension(filename string) (lang string, safe bool) {
	return getLanguageByExtension(filename)
}

func getLanguageByExtension(filename string) (lang string, safe bool) {
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

// GetLanguageByContent returns a language based on the filename and heuristics applies to the content,
// and safe to indicate the sureness of returned language.
func GetLanguageByContent(filename string, content []byte) (lang string, safe bool) {
	return getLanguageByContent(filename, content)
}

func getLanguageByContent(filename string, content []byte) (lang string, safe bool) {
	ext := strings.ToLower(filepath.Ext(filename))
	if fnMatcher, ok := contentMatchers[ext]; ok {
		lang, safe = fnMatcher(content)
	} else {
		lang = OtherLanguage
	}

	return
}

// GetLanguageByClassifier takes in a content and a list of candidates, and apply the classifier's Classify method to
// get the most probably language. If classifier is null then DefaultClassfier will be used.
func GetLanguageByClassifier(content []byte, candidates []string, classifier Classifier) string {
	return getLanguageByClassifier(content, candidates, classifier)
}

// GetLanguageExtensions returns the different extensions being used by the language.
func GetLanguageExtensions(language string) []string {
	return extensionsByLanguage[language]
}

// Type represent language's type. Either data, programming, markup, prose, or unknown.
type Type int

// Type's values.
const (
	Unknown Type = iota
	Data
	Programming
	Markup
	Prose
)

// GetLanguageType returns the given language's type.
func GetLanguageType(language string) (langType Type) {
	langType, ok := languagesType[language]
	if !ok {
		langType = Unknown
	}

	return langType
}

// GetLanguageByAlias returns either the language related to the given alias and ok set to true
// or Otherlanguage and ok set to false if the alias is not recognized.
func GetLanguageByAlias(alias string) (lang string, ok bool) {
	a := strings.Split(alias, `,`)[0]
	a = strings.ToLower(a)
	lang, ok = languagesByAlias[a]
	if !ok {
		lang = OtherLanguage
	}

	return
}
