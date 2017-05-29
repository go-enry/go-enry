package generator

import (
	"bytes"
	"io"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

var typeToTypeConst = map[string]string{
	"data":        "Data",
	"programming": "Programming",
	"markup":      "Markup",
	"prose":       "Prose",
}

// Types reads from buf and builds source file from typeTmplPath.
func Types(data []byte, typeTmplPath, typeTmplName, commit string) ([]byte, error) {
	languages := make(map[string]*languageInfo)
	if err := yaml.Unmarshal(data, &languages); err != nil {
		return nil, err
	}

	langTypeMap := buildLanguageTypeMap(languages)

	buf := &bytes.Buffer{}
	if err := executeTypesTemplate(buf, langTypeMap, typeTmplPath, typeTmplName, commit); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func buildLanguageTypeMap(languages map[string]*languageInfo) map[string]string {
	langTypeMap := make(map[string]string)
	for lang, info := range languages {
		langTypeMap[lang] = typeToTypeConst[info.Type]
	}

	return langTypeMap
}

func executeTypesTemplate(out io.Writer, langTypeMap map[string]string, typeTmplPath, typeTmpl, commit string) error {
	fmap := template.FuncMap{
		"getCommit": func() string { return commit },
	}

	t := template.Must(template.New(typeTmpl).Funcs(fmap).ParseFiles(typeTmplPath))
	if err := t.Execute(out, langTypeMap); err != nil {
		return err
	}

	return nil
}
