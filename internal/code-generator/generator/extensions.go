package generator

import (
	"bytes"
	"io"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

type extensionsInfo struct {
	LanguagesByExtension map[string][]string
	ExtensionsByLanguage map[string][]string
}

// Extensions reads from buf and builds source file from extensionsTmplPath.
func Extensions(data []byte, extensionsTmplPath, extensionsTmplName, commit string) ([]byte, error) {
	languages := make(map[string]*languageInfo)
	if err := yaml.Unmarshal(data, &languages); err != nil {
		return nil, err
	}

	extInfo := &extensionsInfo{}
	orderedKeyList := getAlphabeticalOrderedKeys(languages)
	extInfo.LanguagesByExtension = buildExtensionLanguageMap(languages, orderedKeyList)
	extInfo.ExtensionsByLanguage = buildLanguageExtensionsMap(languages)

	buf := &bytes.Buffer{}
	if err := executeExtensionsTemplate(buf, extInfo, extensionsTmplPath, extensionsTmplName, commit); err != nil {
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

func buildLanguageExtensionsMap(languages map[string]*languageInfo) map[string][]string {
	langExtensionMap := make(map[string][]string, len(languages))
	for lang, info := range languages {
		if len(info.Extensions) > 0 {
			langExtensionMap[lang] = info.Extensions
		}
	}

	return langExtensionMap
}

func executeExtensionsTemplate(out io.Writer, extInfo *extensionsInfo, extensionsTmplPath, extensionsTmpl, commit string) error {
	fmap := template.FuncMap{
		"getCommit":         func() string { return commit },
		"formatStringSlice": func(slice []string) string { return `"` + strings.Join(slice, `","`) + `"` },
	}

	t := template.Must(template.New(extensionsTmpl).Funcs(fmap).ParseFiles(extensionsTmplPath))
	if err := t.Execute(out, extInfo); err != nil {
		return err
	}

	return nil
}
