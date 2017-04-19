package generator

import (
	"bytes"
	"io"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

// Extensions reads from buf and builds extensions_map.go file from extensionsTmplPath.
func Extensions(data []byte, extensionsTmplPath, extensionsTmplName, commit string) ([]byte, error) {
	languages := make(map[string]*languageInfo)
	if err := yaml.Unmarshal(data, &languages); err != nil {
		return nil, err
	}

	orderedKeyList := getAlphabeticalOrderedKeys(languages)
	languagesByExtension := buildExtensionLanguageMap(languages, orderedKeyList)

	buf := &bytes.Buffer{}
	if err := executeExtensionsTemplate(buf, languagesByExtension, extensionsTmplPath, extensionsTmplName, commit); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
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

func executeExtensionsTemplate(out io.Writer, languagesByExtension map[string][]string, extensionsTmplPath, extensionsTmpl, commit string) error {
	fmap := template.FuncMap{
		"getCommit":         func() string { return commit },
		"formatStringSlice": func(slice []string) string { return `"` + strings.Join(slice, `","`) + `"` },
	}

	t := template.Must(template.New(extensionsTmpl).Funcs(fmap).ParseFiles(extensionsTmplPath))
	if err := t.Execute(out, languagesByExtension); err != nil {
		return err
	}

	return nil
}
