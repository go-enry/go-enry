package slinguist

const OtherLanguage = "Other"

var (
	ExtensionsByLanguage map[string][]string
	ignoredExtensions    = map[string]bool{
		".asc": true, ".cgi": true, ".fcgi": true, ".gml": true, ".fx": true,
		".vhost": true,
	}
	auxiliaryLanguages = map[string]bool{
		"Other": true, "XML": true, "YAML": true, "TOML": true, "INI": true,
		"JSON": true, "TeX": true, "Public Key": true, "AsciiDoc": true,
		"AGS Script": true, "VimL": true, "Diff": true, "CMake": true, "fish": true,
		"Awk": true, "Graphviz (DOT)": true, "Markdown": true, "desktop": true,
		"XSLT": true, "SQL": true, "RMarkdown": true, "IRC log": true,
		"reStructuredText": true, "Twig": true, "CSS": true, "Batchfile": true,
		"Text": true, "HTML+ERB": true, "HTML": true, "Gettext Catalog": true,
		"Smarty": true, "Raw token data": true,
	}
)

func init() {
	for l, _ := range ignoredExtensions {
		languagesByExtension[l] = []string{OtherLanguage}
	}

	ExtensionsByLanguage = reverseStringListMap(languagesByExtension)
}

// GetLanguageExtensions returns the different extensions being used by the
// language.
func GetLanguageExtensions(language string) []string {
	return ExtensionsByLanguage[language]
}

// GetLanguage return the Language for a given filename and file content.
func GetLanguage(filename string, content []byte) string {
	if lang, safe := GetLanguageByShebang(content); safe {
		return lang
	}

	if lang, safe := GetLanguageByExtension(filename); safe {
		return lang
	}

	lang, _ := GetLanguageByContent(filename, content)
	return lang
}

func reverseStringListMap(i map[string][]string) (o map[string][]string) {
	o = map[string][]string{}
	for key, set := range i {
		for _, value := range set {
			o[value] = append(o[value], key)
		}
	}

	return
}
