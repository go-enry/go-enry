package generator

import (
	"bytes"
	"io"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

// Filenames reads from buf and builds filenames_map.go file from filenamesTmplPath.
func Filenames(data []byte, filenamesTmplPath, filenamesTmplName, commit string) ([]byte, error) {
	languages := make(map[string]*languageInfo)
	if err := yaml.Unmarshal(data, &languages); err != nil {
		return nil, err
	}

	languagesByFilename := buildFilenameLanguageMap(languages)

	buf := &bytes.Buffer{}
	if err := executeFilenamesTemplate(buf, languagesByFilename, filenamesTmplPath, filenamesTmplName, commit); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func buildFilenameLanguageMap(languages map[string]*languageInfo) map[string]string {
	filenameLangMap := make(map[string]string)
	for lang, langInfo := range languages {
		for _, filename := range langInfo.Filenames {
			filenameLangMap[filename] = lang
		}
	}

	return filenameLangMap
}

func executeFilenamesTemplate(out io.Writer, languagesByFilename map[string]string, filenamesTmplPath, filenamesTmpl, commit string) error {
	fmap := template.FuncMap{
		"getCommit": func() string { return commit },
	}

	t := template.Must(template.New(filenamesTmpl).Funcs(fmap).ParseFiles(filenamesTmplPath))
	if err := t.Execute(out, languagesByFilename); err != nil {
		return err
	}

	return nil
}
