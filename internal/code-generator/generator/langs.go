package generator

import (
	"bytes"
	"io"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

type languageInfo struct {
	Type         string   `yaml:"type,omitempty" json:"type,omitempty"`
	Aliases      []string `yaml:"aliases,omitempty,flow" json:"aliases,omitempty"`
	Extensions   []string `yaml:"extensions,omitempty,flow" json:"extensions,omitempty"`
	Interpreters []string `yaml:"interpreters,omitempty,flow" json:"interpreters,omitempty"`
	Group        string   `yaml:"group,omitempty" json:"group,omitempty"`
}

// Languages reads from buf and builds languages.go file from languagesTmplPath.
func Languages(data []byte, languagesTmplPath, languagesTmplName, commit string) ([]byte, error) {
	languages := make(map[string]*languageInfo)
	if err := yaml.Unmarshal(data, &languages); err != nil {
		return nil, err
	}

	languagesByExtension, err := buildExtensionLanguageMap(languages)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	if err := executeLanguagesTemplate(buf, languagesByExtension, languagesTmplPath, languagesTmplName, commit); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func buildExtensionLanguageMap(languages map[string]*languageInfo) (map[string][]string, error) {
	extensionLangsMap := make(map[string][]string)
	for lang, info := range languages {
		for _, extension := range info.Extensions {
			extensionLangsMap[extension] = append(extensionLangsMap[extension], lang)
		}
	}

	return extensionLangsMap, nil
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
