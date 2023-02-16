// Package generator provides facilities to generate Go code for the
// package data in enry from YAML files describing supported languages in Linguist.
package generator

import (
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"
)

// File is a common type for all generator functions.
// It generates Go source code file based on template in tmplPath,
// by parsing the data in fileToParse and linguist's samplesDir
// saving results to an outFile.
type File func(fileToParse, samplesDir, outPath, tmplPath, tmplName, commit string) error

func formatedWrite(outPath string, source []byte) error {
	formatedSource, err := format.Source(source)
	if err != nil {
		err = fmt.Errorf("'go fmt' fails on %v", err)
		// write un-formatter source to simplify debugging
		formatedSource = source
	}

	if err := ioutil.WriteFile(outPath, formatedSource, 0666); err != nil {
		return err
	}
	return err
}

func executeTemplate(w io.Writer, name, path, commit string, fmap template.FuncMap, data interface{}) error {
	getCommit := func() string {
		return commit
	}
	// stringVal returns escaped string that can be directly placed into go code.
	// for value test`s it would return `test`+"`"+`s`
	stringVal := func(val string) string {
		val = strings.ReplaceAll(val, "`", "`+\"`\"+`")
		return fmt.Sprintf("`%s`", val)
	}
	if fmap == nil {
		fmap = make(template.FuncMap)
	}
	fmap["getCommit"] = getCommit
	fmap["stringVal"] = stringVal
	fmap["isRE2"] = isRE2

	const headerTmpl = "header.go.tmpl"
	headerPath := filepath.Join(filepath.Dir(path), headerTmpl)

	h := template.Must(template.New(headerTmpl).Funcs(fmap).ParseFiles(headerPath))
	if err := h.Execute(w, data); err != nil {
		return err
	}

	t := template.Must(template.New(name).Funcs(fmap).ParseFiles(path))
	return t.Execute(w, data)
}
