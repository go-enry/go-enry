package generator

import "sort"

type languageInfo struct {
	Type         string   `yaml:"type,omitempty"`
	Color        string   `yaml:"color,omitempty"`
	Group        string   `yaml:"group,omitempty"`
	Aliases      []string `yaml:"aliases,omitempty"`
	Extensions   []string `yaml:"extensions,omitempty,flow"`
	Interpreters []string `yaml:"interpreters,omitempty,flow"`
	Filenames    []string `yaml:"filenames,omitempty,flow"`
	MimeType     string   `yaml:"codemirror_mime_type,omitempty,flow"`
	LanguageID   *int     `yaml:"language_id,omitempty"`
}

func getAlphabeticalOrderedKeys(languages map[string]*languageInfo) []string {
	keyList := make([]string, 0)
	for lang := range languages {
		keyList = append(keyList, lang)
	}

	sort.Strings(keyList)
	return keyList
}
