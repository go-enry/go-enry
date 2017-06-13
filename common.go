package enry

import (
	"math"
	"path/filepath"
	"strings"
)

// OtherLanguage is used as a zero value when a function can not return a specific language.
const OtherLanguage = "Other"

// Strategy type fix the signature for the functions that can be used as a strategy.
type Strategy func(filename string, content []byte) (languages []string)

var strategies = []Strategy{
	GetLanguagesByModeline,
	GetLanguagesByFilename,
	GetLanguagesByShebang,
	GetLanguagesByExtension,
	GetLanguagesByContent,
}

// GetLanguage applies a sequence of strategies based on the given filename and content
// to find out the most probably language to return.
func GetLanguage(filename string, content []byte) string {
	candidates := map[string]float64{}
	for _, strategy := range strategies {
		languages := strategy(filename, content)
		if len(languages) == 1 {
			return languages[0]
		}

		if len(languages) > 0 {
			for _, language := range languages {
				candidates[language]++
			}
		}
	}

	if len(candidates) == 0 {
		return OtherLanguage
	}

	lang := GetLanguageByClassifier(content, candidates, nil)
	return lang
}

// GetLanguageByModeline returns the language of the given content looking for the modeline,
// and safe to indicate the sureness of returned language.
func GetLanguageByModeline(content []byte) (lang string, safe bool) {
	return getLangAndSafe("", content, GetLanguagesByModeline)
}

// GetLanguageByFilename returns a language based on the given filename, and safe to indicate
// the sureness of returned language.
func GetLanguageByFilename(filename string) (lang string, safe bool) {
	return getLangAndSafe(filename, nil, GetLanguagesByFilename)
}

// GetLanguagesByFilename returns a slice of possible languages for the given filename, content will be ignored.
// It accomplish the signature to be a Strategy type.
func GetLanguagesByFilename(filename string, content []byte) []string {
	return languagesByFilename[filename]
}

// GetLanguageByShebang returns the language of the given content looking for the shebang line,
// and safe to indicate the sureness of returned language.
func GetLanguageByShebang(content []byte) (lang string, safe bool) {
	return getLangAndSafe("", content, GetLanguagesByShebang)
}

// GetLanguageByExtension returns a language based on the given filename, and safe to indicate
// the sureness of returned language.
func GetLanguageByExtension(filename string) (lang string, safe bool) {
	return getLangAndSafe(filename, nil, GetLanguagesByExtension)
}

// GetLanguagesByExtension returns a slice of possible languages for the given filename, content will be ignored.
// It accomplish the signature to be a Strategy type.
func GetLanguagesByExtension(filename string, content []byte) []string {
	if !strings.Contains(filename, ".") {
		return nil
	}

	filename = strings.ToLower(filename)
	dots := getDotIndexes(filename)
	for _, dot := range dots {
		ext := filename[dot:]
		languages, ok := languagesByExtension[ext]
		if ok {
			return languages
		}
	}

	return nil
}

func getDotIndexes(filename string) []int {
	dots := make([]int, 0, 2)
	for i, letter := range filename {
		if letter == rune('.') {
			dots = append(dots, i)
		}
	}

	return dots
}

// GetLanguageByContent returns a language based on the filename and heuristics applies to the content,
// and safe to indicate the sureness of returned language.
func GetLanguageByContent(filename string, content []byte) (lang string, safe bool) {
	return getLangAndSafe(filename, content, GetLanguagesByContent)
}

// GetLanguagesByContent returns a slice of possible languages for the given content, filename will be ignored.
// It accomplish the signature to be a Strategy type.
func GetLanguagesByContent(filename string, content []byte) []string {
	ext := strings.ToLower(filepath.Ext(filename))
	fnMatcher, ok := contentMatchers[ext]
	if !ok {
		return nil
	}

	return fnMatcher(content)
}

func getLangAndSafe(filename string, content []byte, getLanguageByStrategy Strategy) (lang string, safe bool) {
	languages := getLanguageByStrategy(filename, content)
	if len(languages) == 0 {
		lang = OtherLanguage
		return
	}

	lang = languages[0]
	safe = len(languages) == 1
	return
}

// GetLanguageByClassifier takes in a content and a list of candidates, and apply the classifier's Classify method to
// get the most probably language. If classifier is null then DefaultClassfier will be used. If there aren't candidates
// OtherLanguage is returned.
func GetLanguageByClassifier(content []byte, candidates map[string]float64, classifier Classifier) string {
	scores := GetLanguagesByClassifier(content, candidates, classifier)
	if len(scores) == 0 {
		return OtherLanguage
	}

	return getLangugeHigherScore(scores)
}

func getLangugeHigherScore(scores map[string]float64) string {
	var language string
	higher := -math.MaxFloat64
	for lang, score := range scores {
		if higher < score {
			language = lang
			higher = score
		}
	}

	return language
}

// GetLanguagesByClassifier returns a map of possible languages as keys and a score as value based on content and candidates. The values can be ordered
// with the highest value as the most probably language. If classifier is null then DefaultClassfier will be used.
func GetLanguagesByClassifier(content []byte, candidates map[string]float64, classifier Classifier) map[string]float64 {
	if classifier == nil {
		classifier = DefaultClassifier
	}

	return classifier.Classify(content, candidates)
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
