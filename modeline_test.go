package slinguist

import (
	"io/ioutil"
	"path/filepath"

	. "gopkg.in/check.v1"
)

const (
	modelinesDir = ".linguist/test/fixtures/Data/Modelines"
)

func (s *TSuite) TestGetLanguageByModeline(c *C) {
	linguistTests := []struct {
		filename     string
		expectedLang string
		expectedSafe bool
	}{
		// Emacs
		{filename: "example_smalltalk.md", expectedLang: "Smalltalk", expectedSafe: true},
		{filename: "fundamentalEmacs.c", expectedLang: "Text", expectedSafe: true},
		{filename: "iamphp.inc", expectedLang: "PHP", expectedSafe: true},
		{filename: "seeplusplusEmacs1", expectedLang: "C++", expectedSafe: true},
		{filename: "seeplusplusEmacs2", expectedLang: "C++", expectedSafe: true},
		{filename: "seeplusplusEmacs3", expectedLang: "C++", expectedSafe: true},
		{filename: "seeplusplusEmacs4", expectedLang: "C++", expectedSafe: true},
		{filename: "seeplusplusEmacs5", expectedLang: "C++", expectedSafe: true},
		{filename: "seeplusplusEmacs6", expectedLang: "C++", expectedSafe: true},
		{filename: "seeplusplusEmacs7", expectedLang: "C++", expectedSafe: true},
		{filename: "seeplusplusEmacs9", expectedLang: "C++", expectedSafe: true},
		{filename: "seeplusplusEmacs10", expectedLang: "C++", expectedSafe: true},
		{filename: "seeplusplusEmacs11", expectedLang: "C++", expectedSafe: true},
		{filename: "seeplusplusEmacs12", expectedLang: "C++", expectedSafe: true},

		// Vim
		{filename: "seeplusplus", expectedLang: "C++", expectedSafe: true},
		{filename: "iamjs.pl", expectedLang: "JavaScript", expectedSafe: true},
		{filename: "iamjs2.pl", expectedLang: "JavaScript", expectedSafe: true},
		{filename: "not_perl.pl", expectedLang: "Prolog", expectedSafe: true},
		{filename: "ruby", expectedLang: "Ruby", expectedSafe: true},
		{filename: "ruby2", expectedLang: "Ruby", expectedSafe: true},
		{filename: "ruby3", expectedLang: "Ruby", expectedSafe: true},
		{filename: "ruby4", expectedLang: "Ruby", expectedSafe: true},
		{filename: "ruby5", expectedLang: "Ruby", expectedSafe: true},
		{filename: "ruby6", expectedLang: "Ruby", expectedSafe: true},
		{filename: "ruby7", expectedLang: "Ruby", expectedSafe: true},
		{filename: "ruby8", expectedLang: "Ruby", expectedSafe: true},
		{filename: "ruby9", expectedLang: "Ruby", expectedSafe: true},
		{filename: "ruby10", expectedLang: "Ruby", expectedSafe: true},
		{filename: "ruby11", expectedLang: "Ruby", expectedSafe: true},
		{filename: "ruby12", expectedLang: "Ruby", expectedSafe: true},
	}

	for _, test := range linguistTests {
		content, err := ioutil.ReadFile(filepath.Join(modelinesDir, test.filename))
		c.Assert(err, Equals, nil)

		lang, safe := GetLanguageByModeline(content)
		c.Assert(lang, Equals, test.expectedLang)
		c.Assert(safe, Equals, test.expectedSafe)
	}

	const (
		wrongVim = `# vim: set syntax=ruby ft  =python filetype=perl :`
		rightVim = `/* vim: set syntax=python ft   =python filetype=python */`
	)

	tests := []struct {
		content      []byte
		expectedLang string
		expectedSafe bool
	}{
		{content: []byte(wrongVim), expectedLang: OtherLanguage, expectedSafe: false},
		{content: []byte(rightVim), expectedLang: "Python", expectedSafe: true},
	}

	for _, test := range tests {
		lang, safe := GetLanguageByModeline(test.content)
		c.Assert(lang, Equals, test.expectedLang)
		c.Assert(safe, Equals, test.expectedSafe)
	}
}
