package generator

import (
	"bytes"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Documentation generates regex matchers in Go for documentation files/dirs.
// It is of generator.File type.
func Documentation(fileToParse, _, outFile, tmplPath, tmplName, commit string) error {
	data, err := ioutil.ReadFile(fileToParse)
	if err != nil {
		return err
	}

	var regexpList []string
	if err := yaml.Unmarshal(data, &regexpList); err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	err = executeTemplate(buf, tmplName, tmplPath, commit, nil, regexpList)
	if err != nil {
		return err
	}

	return formatedWrite(outFile, buf.Bytes())
}
