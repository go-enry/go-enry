package generator

import (
	"bytes"
	"html/template"
	"io"

	yaml "gopkg.in/yaml.v2"
)

// Vendor reads from buf and builds vendor_matchers.go file from VendorTmplPath.
func Vendor(data []byte, vendorTmplPath, vendorTmplName, commit string) ([]byte, error) {
	var regexpList []string
	if err := yaml.Unmarshal(data, &regexpList); err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	if err := executeVendorTemplate(buf, regexpList, vendorTmplPath, vendorTmplName, commit); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func executeVendorTemplate(out io.Writer, regexpList []string, vendorTmplPath, vendorTmpl, commit string) error {
	fmap := template.FuncMap{
		"getCommit": func() string { return commit },
	}

	t := template.Must(template.New(vendorTmpl).Funcs(fmap).ParseFiles(vendorTmplPath))
	if err := t.Execute(out, regexpList); err != nil {
		return err
	}

	return nil
}
