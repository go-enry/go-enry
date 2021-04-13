package generator

import (
	"bytes"
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// ID generates a map in Go with language name -> language ID.
// It is of generator.File type.
func ID(fileToParse, samplesDir, outPath, tmplPath, tmplName, commit string) error {
	data, err := ioutil.ReadFile(fileToParse)
	if err != nil {
		return err
	}

	languages := make(map[string]*languageInfo)
	if err := yaml.Unmarshal(data, &languages); err != nil {
		return err
	}

	langMimeMap := buildLanguageIDMap(languages)

	buf := &bytes.Buffer{}
	if err := executeIDTemplate(buf, langMimeMap, tmplPath, tmplName, commit); err != nil {
		return err
	}

	return formatedWrite(outPath, buf.Bytes())
}

func buildLanguageIDMap(languages map[string]*languageInfo) map[string]int {
	langIDMap := make(map[string]int)
	for lang, info := range languages {
		// NOTE: 0 is a valid language ID so checking the zero value would skip one language
		if info.LanguageID != nil {
			langIDMap[lang] = *info.LanguageID
		}
	}

	return langIDMap
}

func executeIDTemplate(out io.Writer, langIDMap map[string]int, tmplPath, tmplName, commit string) error {
	return executeTemplate(out, tmplName, tmplPath, commit, nil, langIDMap)
}
