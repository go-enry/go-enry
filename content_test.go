package slinguist

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	. "gopkg.in/check.v1"
)

func (s *TSuite) TestGetLanguageByContent(c *C) {
	s.testGetLanguageByContent(c, "C")
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
		c.Assert(obtained, Equals, expected, Commentf(file))
	}
}
