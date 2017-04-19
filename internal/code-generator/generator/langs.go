package generator

import (
	"bytes"
	"io"
	"sort"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

type languageInfo struct {
	Type         string   `yaml:"type,omitempty"`
	Extensions   []string `yaml:"extensions,omitempty,flow"`
	Interpreters []string `yaml:"interpreters,omitempty,flow"`
	Filenames    []string `yaml:"filenames,omitempty,flow"`
}

// Languages reads from buf and builds languages.go file from languagesTmplPath.
func Languages(data []byte, languagesTmplPath, languagesTmplName, commit string) ([]byte, error) {
	languages := make(map[string]*languageInfo)
	if err := yaml.Unmarshal(data, &languages); err != nil {
		return nil, err
	}

	orderedKeyList := getAlphabeticalOrderedKeys(languages)
	languagesByExtension := buildExtensionLanguageMap(languages, orderedKeyList)

	buf := &bytes.Buffer{}
	if err := executeLanguagesTemplate(buf, languagesByExtension, languagesTmplPath, languagesTmplName, commit); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func getAlphabeticalOrderedKeys(languages map[string]*languageInfo) []string {
	keyList := make([]string, 0)
	for lang := range languages {
		keyList = append(keyList, lang)
	}

	sort.Strings(keyList)
	return keyList
}

func buildExtensionLanguageMap(languages map[string]*languageInfo, orderedKeyList []string) map[string][]string {
	extensionLangsMap := make(map[string][]string)
	for _, lang := range orderedKeyList {
		langInfo := languages[lang]
		for _, extension := range langInfo.Extensions {
			extensionLangsMap[extension] = append(extensionLangsMap[extension], lang)
		}
	}

	return extensionLangsMap
}

func executeLanguagesTemplate(out io.Writer, languagesByExtension map[string][]string, languagesTmplPath, languagesTmpl, commit string) error {
	fmap := template.FuncMap{
		"getCommit":         func() string { return commit },
		"formatStringSlice": func(slice []string) string { return `"` + strings.Join(slice, `","`) + `"` },
	}

	t := template.Must(template.New(languagesTmpl).Funcs(fmap).ParseFiles(languagesTmplPath))
	if err := t.Execute(out, languagesByExtension); err != nil {
		return err
	}

	return nil
}
