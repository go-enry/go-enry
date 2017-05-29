package generator

import (
	"bytes"
	"io"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

// Aliases reads from buf and builds source file from aliasesTmplPath.
func Aliases(data []byte, aliasesTmplPath, aliasesTmplName, commit string) ([]byte, error) {
	languages := make(map[string]*languageInfo)
	if err := yaml.Unmarshal(data, &languages); err != nil {
		return nil, err
	}

	orderedLangList := getAlphabeticalOrderedKeys(languages)
	languagesByAlias := buildAliasLanguageMap(languages, orderedLangList)

	buf := &bytes.Buffer{}
	if err := executeAliasesTemplate(buf, languagesByAlias, aliasesTmplPath, aliasesTmplName, commit); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func buildAliasLanguageMap(languages map[string]*languageInfo, orderedLangList []string) map[string]string {
	aliasLangsMap := make(map[string]string)
	for _, lang := range orderedLangList {
		langInfo := languages[lang]
		key := convertToAliasKey(lang)
		aliasLangsMap[key] = lang
		for _, alias := range langInfo.Aliases {
			key := convertToAliasKey(alias)
			aliasLangsMap[key] = lang
		}
	}

	return aliasLangsMap
}

func convertToAliasKey(s string) (key string) {
	key = strings.Replace(s, ` `, `_`, -1)
	key = strings.ToLower(key)
	return
}

func executeAliasesTemplate(out io.Writer, languagesByAlias map[string]string, aliasesTmplPath, aliasesTmpl, commit string) error {
	fmap := template.FuncMap{
		"getCommit": func() string { return commit },
	}

	t := template.Must(template.New(aliasesTmpl).Funcs(fmap).ParseFiles(aliasesTmplPath))
	if err := t.Execute(out, languagesByAlias); err != nil {
		return err
	}

	return nil
}
