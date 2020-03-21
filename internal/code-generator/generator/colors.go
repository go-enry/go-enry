package generator

import (
	"bytes"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Colors generates a map in Go with language name -> color string.
// It is of generator.File type.
func Colors(fileToParse, samplesDir, outPath, tmplPath, tmplName, commit string) error {
	data, err := ioutil.ReadFile(fileToParse)
	if err != nil {
		return err
	}

	languages := make(map[string]*languageInfo)
	if err := yaml.Unmarshal(data, &languages); err != nil {
		return err
	}

	langColorMap := buildLanguageColorMap(languages)

	buf := &bytes.Buffer{}
	if err := executeColorTemplate(buf, langColorMap, tmplPath, tmplName, commit); err != nil {
		return err
	}

	return formatedWrite(outPath, buf.Bytes())
}

func buildLanguageColorMap(languages map[string]*languageInfo) map[string]string {
	langColorMap := make(map[string]string)
	for lang, info := range languages {
		if len(info.Color) != 0 {
			langColorMap[lang] = info.Color
		}
	}

	return langColorMap
}

func executeColorTemplate(out io.Writer, langColorMap map[string]string, tmplPath, tmplName, commit string) error {
	return executeTemplate(out, tmplName, tmplPath, commit, nil, langColorMap)
}
