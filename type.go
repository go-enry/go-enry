package slinguist

type Type int

const (
	// Language's type. Either data, programming, markup, prose, or unknown.
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
