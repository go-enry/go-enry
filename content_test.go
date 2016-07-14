package slinguist

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"text/tabwriter"

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
	s.testGetLanguageByContent(c, "Hack")
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

func (s *TSuite) TestGetLanguageByContentLinguistCorpus(c *C) {
	var total, failed, ok, other, unsafe int

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	filepath.Walk(".linguist/samples", func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			if f.Name() == "filenames" {
				return filepath.SkipDir
			}

			return nil
		}

		expected := filepath.Base(filepath.Dir(path))
		filename := filepath.Base(path)
		extension := filepath.Ext(path)
		content, _ := ioutil.ReadFile(path)

		if extension == "" {
			return nil
		}

		total++
		obtained, safe := GetLanguageByContent(filename, content)
		if obtained == OtherLanguage {
			other++
		}

		var status string
		if expected == obtained {
			status = "ok"
			ok++
		} else {
			status = "failed"
			failed++
			if !safe {
				unsafe++
			}
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%s\n", filename, expected, obtained, safe, status)

		return nil
	})

	fmt.Fprintln(w)
	w.Flush()

	fmt.Printf("total files: %d, ok: %d, failed: %d, unsafe: %d, other: %d\n", total, ok, failed, unsafe, other)

}
