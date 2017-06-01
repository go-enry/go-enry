package slinguist

import (
	"io/ioutil"
	"path/filepath"

	. "gopkg.in/check.v1"
)

const (
	modelinesDir = ".linguist/test/fixtures/Data/Modelines"
	samplesDir   = ".linguist/samples"
)

func (s *TSuite) TestGetLanguageByModeline(c *C) {
	linguistTests := []struct {
		filename     string
		expectedLang string
		expectedSafe bool
	}{
		// Emacs
		{filename: filepath.Join(modelinesDir, "example_smalltalk.md"), expectedLang: "Smalltalk", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "fundamentalEmacs.c"), expectedLang: "Text", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "iamphp.inc"), expectedLang: "PHP", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "seeplusplusEmacs1"), expectedLang: "C++", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "seeplusplusEmacs2"), expectedLang: "C++", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "seeplusplusEmacs3"), expectedLang: "C++", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "seeplusplusEmacs4"), expectedLang: "C++", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "seeplusplusEmacs5"), expectedLang: "C++", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "seeplusplusEmacs6"), expectedLang: "C++", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "seeplusplusEmacs7"), expectedLang: "C++", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "seeplusplusEmacs9"), expectedLang: "C++", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "seeplusplusEmacs10"), expectedLang: "C++", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "seeplusplusEmacs11"), expectedLang: "C++", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "seeplusplusEmacs12"), expectedLang: "C++", expectedSafe: true},

		// Vim
		{filename: filepath.Join(modelinesDir, "seeplusplus"), expectedLang: "C++", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "iamjs.pl"), expectedLang: "JavaScript", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "iamjs2.pl"), expectedLang: "JavaScript", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "not_perl.pl"), expectedLang: "Prolog", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "ruby"), expectedLang: "Ruby", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "ruby2"), expectedLang: "Ruby", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "ruby3"), expectedLang: "Ruby", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "ruby4"), expectedLang: "Ruby", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "ruby5"), expectedLang: "Ruby", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "ruby6"), expectedLang: "Ruby", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "ruby7"), expectedLang: "Ruby", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "ruby8"), expectedLang: "Ruby", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "ruby9"), expectedLang: "Ruby", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "ruby10"), expectedLang: "Ruby", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "ruby11"), expectedLang: "Ruby", expectedSafe: true},
		{filename: filepath.Join(modelinesDir, "ruby12"), expectedLang: "Ruby", expectedSafe: true},
		{filename: filepath.Join(samplesDir, "C/main.c"), expectedLang: OtherLanguage, expectedSafe: false},
	}

	for _, test := range linguistTests {
		content, err := ioutil.ReadFile(test.filename)
		c.Assert(err, Equals, nil)

		lang, safe := GetLanguageByModeline(content)
		c.Assert(lang, Equals, test.expectedLang)
		c.Assert(safe, Equals, test.expectedSafe)
	}

	const (
		wrongVim  = `# vim: set syntax=ruby ft  =python filetype=perl :`
		rightVim  = `/* vim: set syntax=python ft   =python filetype=python */`
		noLangVim = `/* vim: set shiftwidth=4 softtabstop=0 cindent cinoptions={1s: */`
	)

	tests := []struct {
		content      []byte
		expectedLang string
		expectedSafe bool
	}{
		{content: []byte(wrongVim), expectedLang: OtherLanguage, expectedSafe: false},
		{content: []byte(rightVim), expectedLang: "Python", expectedSafe: true},
		{content: []byte(noLangVim), expectedLang: OtherLanguage, expectedSafe: false},
	}

	for _, test := range tests {
		lang, safe := GetLanguageByModeline(test.content)
		c.Assert(lang, Equals, test.expectedLang)
		c.Assert(safe, Equals, test.expectedSafe)
	}
}
