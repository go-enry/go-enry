package generator

import (
	"bytes"
	"html/template"
	"io"

	yaml "gopkg.in/yaml.v2"
)

// Vendor reads from buf and builds utils.go file from utilsTmplPath.
func Vendor(data []byte, uitlsTmplPath, utilsTmplName, commit string) ([]byte, error) {
	var regexpList []string
	if err := yaml.Unmarshal(data, &regexpList); err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	if err := executeVendorTemplate(buf, regexpList, uitlsTmplPath, utilsTmplName, commit); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func executeVendorTemplate(out io.Writer, regexpList []string, languagesTmplPath, languagesTmpl, commit string) error {
	fmap := template.FuncMap{
		"getCommit": func() string { return commit },
	}

	t := template.Must(template.New(languagesTmpl).Funcs(fmap).ParseFiles(languagesTmplPath))
	if err := t.Execute(out, regexpList); err != nil {
		return err
	}

	return nil
}
