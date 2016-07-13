package slinguist

var AuxiliaryLanguages = map[string]bool{
	"Other":            true,
	"XML":              true,
	"YAML":             true,
	"TOML":             true,
	"INI":              true,
	"JSON":             true,
	"TeX":              true,
	"Public Key":       true,
	"AsciiDoc":         true,
	"AGS Script":       true,
	"VimL":             true,
	"Diff":             true,
	"CMake":            true,
	"fish":             true,
	"Awk":              true,
	"Graphviz (DOT)":   true,
	"Markdown":         true,
	"desktop":          true,
	"XSLT":             true,
	"SQL":              true,
	"RMarkdown":        true,
	"IRC log":          true,
	"reStructuredText": true,
	"Twig":             true,
	"CSS":              true,
	"Batchfile":        true,
	"Text":             true,
	"HTML+ERB":         true,
	"HTML":             true,
	"Gettext Catalog":  true,
	"Smarty":           true,
	"Raw token data":   true,
}

func IsAuxiliaryLanguage(lang string) bool {
	_, ok := AuxiliaryLanguages[lang]
	return ok
}
