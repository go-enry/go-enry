package slinguist

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	. "gopkg.in/check.v1"
)

func (s *TSuite) TestGetLanguageByContentH(c *C) {
	s.testGetLanguageByContent(c, "Objective-C")
	s.testGetLanguageByContent(c, "C++")
	s.testGetLanguageByContent(c, "C")
	s.testGetLanguageByContent(c, "Common Lisp")
	s.testGetLanguageByContent(c, "Cool")
	s.testGetLanguageByContent(c, "OpenCL")
	s.testGetLanguageByContent(c, "Groff")
	s.testGetLanguageByContent(c, "PicoLisp")
	s.testGetLanguageByContent(c, "PicoLisp")
	s.testGetLanguageByContent(c, "NewLisp")
	s.testGetLanguageByContent(c, "Lex")
	s.testGetLanguageByContent(c, "TeX")
	s.testGetLanguageByContent(c, "Visual Basic")
	s.testGetLanguageByContent(c, "Matlab")
	s.testGetLanguageByContent(c, "Mathematica")
	s.testGetLanguageByContent(c, "Prolog")
	s.testGetLanguageByContent(c, "Perl")
	s.testGetLanguageByContent(c, "Perl6")
}

func (s *TSuite) testGetLanguageByContent(c *C, expected string) {
	files, err := filepath.Glob(path.Join(".linguist/samples", expected, "*"))
	c.Assert(err, IsNil)

	for _, file := range files {
		s, _ := os.Stat(file)
		if s.IsDir() {
			continue
		}

		content, _ := ioutil.ReadFile(file)

		obtained, _ := GetLanguageByContent(path.Base(file), content)
		if obtained == OtherLanguage {
			continue
		}

		c.Check(obtained, Equals, expected, Commentf(file))
	}
}
