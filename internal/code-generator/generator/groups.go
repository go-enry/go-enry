package generator

import (
	"bytes"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Groups generates a map in Go with language name -> group name.
// It is of generator.File type.
func Groups(fileToParse, samplesDir, outPath, tmplPath, tmplName, commit string) error {
	data, err := ioutil.ReadFile(fileToParse)
	if err != nil {
		return err
	}

	languages := make(map[string]*languageInfo)
	if err := yaml.Unmarshal(data, &languages); err != nil {
		return err
	}

	langGroupMap := buildLanguageGroupMap(languages)

	buf := &bytes.Buffer{}
	if err := executeGroupTemplate(buf, langGroupMap, tmplPath, tmplName, commit); err != nil {
		return err
	}

	return formatedWrite(outPath, buf.Bytes())
}

func buildLanguageGroupMap(languages map[string]*languageInfo) map[string]string {
	langGroupMap := make(map[string]string)
	for lang, info := range languages {
		if len(info.Group) != 0 {
			langGroupMap[lang] = info.Group
		}
	}

	return langGroupMap
}

func executeGroupTemplate(out io.Writer, langColorMap map[string]string, tmplPath, tmplName, commit string) error {
	return executeTemplate(out, tmplName, tmplPath, commit, nil, langColorMap)
}
