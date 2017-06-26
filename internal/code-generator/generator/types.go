package generator

import (
	"bytes"
	"io"
	"io/ioutil"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

var typeToTypeConst = map[string]string{
	"data":        "Data",
	"programming": "Programming",
	"markup":      "Markup",
	"prose":       "Prose",
}

// Types reads from fileToParse and builds source file from tmplPath. It complies with type File signature.
func Types(fileToParse, samplesDir, outPath, tmplPath, tmplName, commit string) error {
	data, err := ioutil.ReadFile(fileToParse)
	if err != nil {
		return err
	}

	languages := make(map[string]*languageInfo)
	if err := yaml.Unmarshal(data, &languages); err != nil {
		return err
	}

	langTypeMap := buildLanguageTypeMap(languages)

	buf := &bytes.Buffer{}
	if err := executeTypesTemplate(buf, langTypeMap, tmplPath, tmplName, commit); err != nil {
		return err
	}

	return formatedWrite(outPath, buf.Bytes())
}

func buildLanguageTypeMap(languages map[string]*languageInfo) map[string]string {
	langTypeMap := make(map[string]string)
	for lang, info := range languages {
		langTypeMap[lang] = typeToTypeConst[info.Type]
	}

	return langTypeMap
}

func executeTypesTemplate(out io.Writer, langTypeMap map[string]string, tmplPath, tmplName, commit string) error {
	fmap := template.FuncMap{
		"getCommit": func() string { return commit },
	}

	t := template.Must(template.New(tmplName).Funcs(fmap).ParseFiles(tmplPath))
	if err := t.Execute(out, langTypeMap); err != nil {
		return err
	}

	return nil
}
