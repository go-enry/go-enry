package generator

import (
	"bytes"
	"io"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

type languageInfo struct {
	Type         string   `yaml:"type,omitempty"`
	Aliases      []string `yaml:"aliases,omitempty,flow"`
	Extensions   []string `yaml:"extensions,omitempty,flow"`
	Interpreters []string `yaml:"interpreters,omitempty,flow"`
	Group        string   `yaml:"group,omitempty"`
}

// Languages reads from buf and builds languages.go file from languagesTmplPath.
func Languages(data []byte, languagesTmplPath, languagesTmplName, commit string) ([]byte, error) {
	languages := make(map[string]*languageInfo)
	if err := yaml.Unmarshal(data, &languages); err != nil {
		return nil, err
	}

	orderedKeyList, err := getAlphabeticalOrderedKeys(data)
	if err != nil {
		return nil, err
	}

	languagesByExtension := buildExtensionLanguageMap(languages, orderedKeyList)

	buf := &bytes.Buffer{}
	if err := executeLanguagesTemplate(buf, languagesByExtension, languagesTmplPath, languagesTmplName, commit); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func getAlphabeticalOrderedKeys(data []byte) ([]string, error) {
	var yamlSlice yaml.MapSlice
	if err := yaml.Unmarshal(data, &yamlSlice); err != nil {
		return nil, err
	}

	orderedKeyList := make([]string, 0)
	for _, lang := range yamlSlice {
		orderedKeyList = append(orderedKeyList, lang.Key.(string))
	}

	return orderedKeyList, nil
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
