package generator

import (
	"bytes"
	"io"
	"io/ioutil"
	"sort"

	"gopkg.in/yaml.v2"
)

type languageInfo struct {
	FSName         string   `yaml:"fs_name"`
	Type           string   `yaml:"type,omitempty"`
	Color          string   `yaml:"color,omitempty"`
	Group          string   `yaml:"group,omitempty"`
	Aliases        []string `yaml:"aliases,omitempty"`
	Extensions     []string `yaml:"extensions,omitempty,flow"`
	Interpreters   []string `yaml:"interpreters,omitempty,flow"`
	Filenames      []string `yaml:"filenames,omitempty,flow"`
	MimeType       string   `yaml:"codemirror_mime_type,omitempty,flow"`
	TMScope        string   `yaml:"tm_scope"`
	AceMode        string   `yaml:"ace_mode"`
	CodeMirrorMode string   `yaml:"codemirror_mode"`
	Wrap           bool     `yaml:"wrap"`
	LanguageID     *int     `yaml:"language_id,omitempty"`
}

func getAlphabeticalOrderedKeys(languages map[string]*languageInfo) []string {
	keyList := make([]string, 0)
	for lang := range languages {
		keyList = append(keyList, lang)
	}

	sort.Strings(keyList)
	return keyList
}

// LanguageInfo generates maps in Go with language name -> LanguageInfo and language ID -> LanguageInfo.
// It is of generator.File type.
func LanguageInfo(fileToParse, samplesDir, outPath, tmplPath, tmplName, commit string) error {
	data, err := ioutil.ReadFile(fileToParse)
	if err != nil {
		return err
	}

	languages := make(map[string]*languageInfo)
	if err := yaml.Unmarshal(data, &languages); err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	if err := executeLanguageInfoTemplate(buf, languages, tmplPath, tmplName, commit); err != nil {
		return err
	}

	return formatedWrite(outPath, buf.Bytes())
}

func executeLanguageInfoTemplate(out io.Writer, languages map[string]*languageInfo, tmplPath, tmplName, commit string) error {
	return executeTemplate(out, tmplName, tmplPath, commit, nil, languages)
}
