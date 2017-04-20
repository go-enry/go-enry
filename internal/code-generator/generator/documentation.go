package generator

import (
	"bytes"
	"html/template"
	"io"

	yaml "gopkg.in/yaml.v2"
)

// Documentation reads from buf and builds documentation_matchers.go file from documentationTmplPath.
func Documentation(data []byte, documentationTmplPath, documentationTmplName, commit string) ([]byte, error) {
	var regexpList []string
	if err := yaml.Unmarshal(data, &regexpList); err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	if err := executeDocumentationTemplate(buf, regexpList, documentationTmplPath, documentationTmplName, commit); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func executeDocumentationTemplate(out io.Writer, regexpList []string, documentationTmplPath, documentationTmpl, commit string) error {
	fmap := template.FuncMap{
		"getCommit": func() string { return commit },
	}

	t := template.Must(template.New(documentationTmpl).Funcs(fmap).ParseFiles(documentationTmplPath))
	if err := t.Execute(out, regexpList); err != nil {
		return err
	}

	return nil
}
