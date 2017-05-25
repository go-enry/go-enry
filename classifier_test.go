package slinguist

import (
	"io/ioutil"
	"path/filepath"

	. "gopkg.in/check.v1"
)

func (s *TSuite) TestGetLanguageByClassifier(c *C) {
	const samples = `.linguist/samples/`
	test := []struct {
		filename     string
		candidates   []string
		expectedLang string
	}{
		{filename: filepath.Join(samples, "C/blob.c"), candidates: []string{"python", "ruby", "c", "c++"}, expectedLang: "C"},
		{filename: filepath.Join(samples, "C/blob.c"), candidates: nil, expectedLang: "C"},
		{filename: filepath.Join(samples, "C/main.c"), candidates: nil, expectedLang: "C"},
		{filename: filepath.Join(samples, "C/blob.c"), candidates: []string{"python", "ruby", "c++"}, expectedLang: "C++"},
		{filename: filepath.Join(samples, "C/blob.c"), candidates: []string{"ruby"}, expectedLang: "Ruby"},
		{filename: filepath.Join(samples, "Python/django-models-base.py"), candidates: []string{"python", "ruby", "c", "c++"}, expectedLang: "Python"},
		{filename: filepath.Join(samples, "Python/django-models-base.py"), candidates: nil, expectedLang: "Python"},
	}

	for _, test := range test {
		content, err := ioutil.ReadFile(test.filename)
		c.Assert(err, Equals, nil)
		lang := GetLanguageByClassifier(content, test.candidates, nil)
		c.Assert(lang, Equals, test.expectedLang)
	}
}
