package generator

import (
	"go/format"
	"io/ioutil"
)

// Func is the function's type that generate source file from a data to be parsed and a template.
type Func func(dataToParse []byte, templatePath string, template string, commit string) ([]byte, error)

// FromFile read data to parse from a file named fileToParse and write the generated source code to a file named outPath. The generated
// source code is formated with gofmt and tagged with commit.
func FromFile(fileToParse, outPath, tmplPath, tmplName, commit string, generate Func) error {
	buf, err := ioutil.ReadFile(fileToParse)
	if err != nil {
		return err
	}

	source, err := generate(buf, tmplPath, tmplName, commit)
	if err != nil {
		return err
	}

	return formatedWrite(outPath, source)
}

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
