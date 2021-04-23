// Package generator provides facilities to generate Go code for the
// package data in enry from YAML files describing supported languages in Linguist.
package generator

import (
	"bytes"
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
		return err
	}
	if err := ioutil.WriteFile(outPath, formatedSource, 0666); err != nil {
		return err
	}
	return nil
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

	const headerTmpl = "header.go.tmpl"
	headerPath := filepath.Join(filepath.Dir(path), headerTmpl)

	h := template.Must(template.New(headerTmpl).Funcs(template.FuncMap{
		"getCommit":     getCommit,
		"stringVal":     stringVal,
		"languageConst": languageConst,
	}).ParseFiles(headerPath))

	buf := bytes.NewBuffer(nil)
	if err := h.Execute(buf, data); err != nil {
		return err
	}

	if fmap == nil {
		fmap = make(template.FuncMap)
	}
	fmap["getCommit"] = getCommit
	fmap["stringVal"] = stringVal
	fmap["languageConst"] = languageConst

	t := template.Must(template.New(name).Funcs(fmap).ParseFiles(path))
	if err := t.Execute(buf, data); err != nil {
		return err
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	_, err = w.Write(src)
	return err
}

// replaceConstTable defines the chars replaced on the constant generation
// the `^` means that should be at the begining of the string.
var replaceConstTable = map[string]string{
	".":  "",
	" ":  "",
	"-":  "",
	"'":  "",
	"+":  "Plus",
	"#":  "Sharp",
	"*":  "Star",
	"^1": "One",
	"^2": "Two",
	"^3": "Three",
	"^4": "Four",
	"^5": "Five",
	"^6": "Six",
	"^7": "Seven",
	"^8": "Eight",
	"^9": "Nine",
	"^0": "Zero",
}

// languageConst translate an language as string on in a constant.
func languageConst(lang string) string {
	for i, o := range replaceConstTable {
		if i[0] != '^' {
			lang = strings.ReplaceAll(lang, i, o)
			continue
		}

		if len(lang) <= 1 || i[1] != lang[0] {
			continue
		}

		fmt.Println(o, lang[1:], lang)
		lang = o + lang[1:]
	}

	lang = strings.Title(lang)
	return lang
}
