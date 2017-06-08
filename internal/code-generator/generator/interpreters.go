package generator

import (
	"bytes"
	"io"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

// Interpreters reads from buf and builds source file from interpretersTmplPath.
func Interpreters(data []byte, interpretersTmplPath, interpretersTmplName, commit string) ([]byte, error) {
	languages := make(map[string]*languageInfo)
	if err := yaml.Unmarshal(data, &languages); err != nil {
		return nil, err
	}

	orderedKeys := getAlphabeticalOrderedKeys(languages)
	languagesByInterpreter := buildInterpreterLanguagesMap(languages, orderedKeys)

	buf := &bytes.Buffer{}
	if err := executeInterpretersTemplate(buf, languagesByInterpreter, interpretersTmplPath, interpretersTmplName, commit); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func buildInterpreterLanguagesMap(languages map[string]*languageInfo, orderedKeys []string) map[string][]string {
	interpreterLangsMap := make(map[string][]string)
	for _, lang := range orderedKeys {
		langInfo := languages[lang]
		for _, interpreter := range langInfo.Interpreters {
			interpreterLangsMap[interpreter] = append(interpreterLangsMap[interpreter], lang)
		}
	}

	return interpreterLangsMap
}

func executeInterpretersTemplate(out io.Writer, languagesByInterpreter map[string][]string, interpretersTmplPath, interpretersTmpl, commit string) error {
	fmap := template.FuncMap{
		"getCommit":         func() string { return commit },
		"formatStringSlice": func(slice []string) string { return `"` + strings.Join(slice, `","`) + `"` },
	}

	t := template.Must(template.New(interpretersTmpl).Funcs(fmap).ParseFiles(interpretersTmplPath))
	if err := t.Execute(out, languagesByInterpreter); err != nil {
		return err
	}

	return nil
}
