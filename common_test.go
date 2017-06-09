package enry

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SimpleLinguistTestSuite struct {
	suite.Suite
}

func TestSimpleLinguistTestSuite(t *testing.T) {
	suite.Run(t, new(SimpleLinguistTestSuite))
}

func (s *SimpleLinguistTestSuite) TestGetLanguage() {
	tests := []struct {
		name     string
		filename string
		content  []byte
		expected string
	}{
		{name: "TestGetLanguage_1", filename: "foo.py", content: []byte{}, expected: "Python"},
		{name: "TestGetLanguage_2", filename: "foo.m", content: []byte(":- module"), expected: "Mercury"},
		{name: "TestGetLanguage_3", filename: "foo.m", content: nil, expected: OtherLanguage},
	}

	for _, test := range tests {
		language := GetLanguage(test.filename, test.content)
		assert.Equal(s.T(), test.expected, language, fmt.Sprintf("%v: %v, expected: %v", test.name, language, test.expected))
	}
}

func (s *SimpleLinguistTestSuite) TestGetLanguageByModelineLinguist() {
	const (
		modelinesDir = ".linguist/test/fixtures/Data/Modelines"
		samplesDir   = ".linguist/samples"
	)

	tests := []struct {
		name         string
		filename     string
		expectedLang string
		expectedSafe bool
	}{
		// Emacs
		{name: "TestGetLanguageByModelineLinguist_1", filename: filepath.Join(modelinesDir, "example_smalltalk.md"), expectedLang: "Smalltalk", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_2", filename: filepath.Join(modelinesDir, "fundamentalEmacs.c"), expectedLang: "Text", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_3", filename: filepath.Join(modelinesDir, "iamphp.inc"), expectedLang: "PHP", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_4", filename: filepath.Join(modelinesDir, "seeplusplusEmacs1"), expectedLang: "C++", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_5", filename: filepath.Join(modelinesDir, "seeplusplusEmacs2"), expectedLang: "C++", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_6", filename: filepath.Join(modelinesDir, "seeplusplusEmacs3"), expectedLang: "C++", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_7", filename: filepath.Join(modelinesDir, "seeplusplusEmacs4"), expectedLang: "C++", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_8", filename: filepath.Join(modelinesDir, "seeplusplusEmacs5"), expectedLang: "C++", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_9", filename: filepath.Join(modelinesDir, "seeplusplusEmacs6"), expectedLang: "C++", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_10", filename: filepath.Join(modelinesDir, "seeplusplusEmacs7"), expectedLang: "C++", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_11", filename: filepath.Join(modelinesDir, "seeplusplusEmacs9"), expectedLang: "C++", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_12", filename: filepath.Join(modelinesDir, "seeplusplusEmacs10"), expectedLang: "C++", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_13", filename: filepath.Join(modelinesDir, "seeplusplusEmacs11"), expectedLang: "C++", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_14", filename: filepath.Join(modelinesDir, "seeplusplusEmacs12"), expectedLang: "C++", expectedSafe: true},

		// Vim
		{name: "TestGetLanguageByModelineLinguist_15", filename: filepath.Join(modelinesDir, "seeplusplus"), expectedLang: "C++", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_16", filename: filepath.Join(modelinesDir, "iamjs.pl"), expectedLang: "JavaScript", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_17", filename: filepath.Join(modelinesDir, "iamjs2.pl"), expectedLang: "JavaScript", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_18", filename: filepath.Join(modelinesDir, "not_perl.pl"), expectedLang: "Prolog", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_19", filename: filepath.Join(modelinesDir, "ruby"), expectedLang: "Ruby", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_20", filename: filepath.Join(modelinesDir, "ruby2"), expectedLang: "Ruby", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_21", filename: filepath.Join(modelinesDir, "ruby3"), expectedLang: "Ruby", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_22", filename: filepath.Join(modelinesDir, "ruby4"), expectedLang: "Ruby", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_23", filename: filepath.Join(modelinesDir, "ruby5"), expectedLang: "Ruby", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_24", filename: filepath.Join(modelinesDir, "ruby6"), expectedLang: "Ruby", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_25", filename: filepath.Join(modelinesDir, "ruby7"), expectedLang: "Ruby", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_26", filename: filepath.Join(modelinesDir, "ruby8"), expectedLang: "Ruby", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_27", filename: filepath.Join(modelinesDir, "ruby9"), expectedLang: "Ruby", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_28", filename: filepath.Join(modelinesDir, "ruby10"), expectedLang: "Ruby", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_29", filename: filepath.Join(modelinesDir, "ruby11"), expectedLang: "Ruby", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_30", filename: filepath.Join(modelinesDir, "ruby12"), expectedLang: "Ruby", expectedSafe: true},
		{name: "TestGetLanguageByModelineLinguist_31", filename: filepath.Join(samplesDir, "C/main.c"), expectedLang: OtherLanguage, expectedSafe: false},
	}

	for _, test := range tests {
		content, err := ioutil.ReadFile(test.filename)
		assert.NoError(s.T(), err)

		lang, safe := GetLanguageByModeline(content)
		assert.Equal(s.T(), test.expectedLang, lang, fmt.Sprintf("%v: lang = %v, expected: %v", test.name, lang, test.expectedLang))
		assert.Equal(s.T(), test.expectedSafe, safe, fmt.Sprintf("%v: safe = %v, expected: %v", test.name, safe, test.expectedSafe))
	}
}

func (s *SimpleLinguistTestSuite) TestGetLanguageByModeline() {
	const (
		wrongVim  = `# vim: set syntax=ruby ft  =python filetype=perl :`
		rightVim  = `/* vim: set syntax=python ft   =python filetype=python */`
		noLangVim = `/* vim: set shiftwidth=4 softtabstop=0 cindent cinoptions={1s: */`
	)

	tests := []struct {
		name         string
		content      []byte
		expectedLang string
		expectedSafe bool
	}{
		{name: "TestGetLanguageByModeline_1", content: []byte(wrongVim), expectedLang: OtherLanguage, expectedSafe: false},
		{name: "TestGetLanguageByModeline_2", content: []byte(rightVim), expectedLang: "Python", expectedSafe: true},
		{name: "TestGetLanguageByModeline_3", content: []byte(noLangVim), expectedLang: OtherLanguage, expectedSafe: false},
	}

	for _, test := range tests {
		lang, safe := GetLanguageByModeline(test.content)
		assert.Equal(s.T(), test.expectedLang, lang, fmt.Sprintf("%v: lang = %v, expected: %v", test.name, lang, test.expectedLang))
		assert.Equal(s.T(), test.expectedSafe, safe, fmt.Sprintf("%v: safe = %v, expected: %v", test.name, safe, test.expectedSafe))
	}
}

func (s *SimpleLinguistTestSuite) TestGetLanguageByFilename() {
	tests := []struct {
		name         string
		filename     string
		expectedLang string
		expectedSafe bool
	}{
		{name: "TestGetLanguageByFilename_1", filename: "unknown.interpreter", expectedLang: OtherLanguage, expectedSafe: false},
		{name: "TestGetLanguageByFilename_2", filename: ".bashrc", expectedLang: "Shell", expectedSafe: true},
		{name: "TestGetLanguageByFilename_3", filename: "Dockerfile", expectedLang: "Dockerfile", expectedSafe: true},
		{name: "TestGetLanguageByFilename_4", filename: "Makefile.frag", expectedLang: "Makefile", expectedSafe: true},
		{name: "TestGetLanguageByFilename_5", filename: "makefile", expectedLang: "Makefile", expectedSafe: true},
		{name: "TestGetLanguageByFilename_6", filename: "Vagrantfile", expectedLang: "Ruby", expectedSafe: true},
		{name: "TestGetLanguageByFilename_7", filename: "_vimrc", expectedLang: "Vim script", expectedSafe: true},
		{name: "TestGetLanguageByFilename_8", filename: "pom.xml", expectedLang: "Maven POM", expectedSafe: true},
	}

	for _, test := range tests {
		lang, safe := GetLanguageByFilename(test.filename)
		assert.Equal(s.T(), test.expectedLang, lang, fmt.Sprintf("%v: lang = %v, expected: %v", test.name, lang, test.expectedLang))
		assert.Equal(s.T(), test.expectedSafe, safe, fmt.Sprintf("%v: safe = %v, expected: %v", test.name, safe, test.expectedSafe))
	}
}

func (s *SimpleLinguistTestSuite) TestGetLanguageByShebang() {
	const (
		multilineExecHack = `#!/bin/sh
# Next line is comment in Tcl, but not in sh... \
exec tclsh "$0" ${1+"$@"}`

		multilineNoExecHack = `#!/bin/sh
#<<<#
echo "A shell script in a zkl program ($0)"
echo "Now run zkl <this file> with Hello World as args"
zkl $0 Hello World!
exit
#<<<#
println("The shell script says ",vm.arglist.concat(" "));`
	)

	tests := []struct {
		name         string
		content      []byte
		expectedLang string
		expectedSafe bool
	}{
		{name: "TestGetLanguageByShebang_1", content: []byte(`#!/unknown/interpreter`), expectedLang: OtherLanguage, expectedSafe: false},
		{name: "TestGetLanguageByShebang_2", content: []byte(`no shebang`), expectedLang: OtherLanguage, expectedSafe: false},
		{name: "TestGetLanguageByShebang_3", content: []byte(`#!/usr/bin/env`), expectedLang: OtherLanguage, expectedSafe: false},
		{name: "TestGetLanguageByShebang_4", content: []byte(`#!/usr/bin/python -tt`), expectedLang: "Python", expectedSafe: true},
		{name: "TestGetLanguageByShebang_5", content: []byte(`#!/usr/bin/env python2.6`), expectedLang: "Python", expectedSafe: true},
		{name: "TestGetLanguageByShebang_6", content: []byte(`#!/usr/bin/env perl`), expectedLang: "Perl", expectedSafe: true},
		{name: "TestGetLanguageByShebang_7", content: []byte(`#!	/bin/sh`), expectedLang: "Shell", expectedSafe: true},
		{name: "TestGetLanguageByShebang_8", content: []byte(`#!bash`), expectedLang: "Shell", expectedSafe: true},
		{name: "TestGetLanguageByShebang_9", content: []byte(multilineExecHack), expectedLang: "Tcl", expectedSafe: true},
		{name: "TestGetLanguageByShebang_10", content: []byte(multilineNoExecHack), expectedLang: "Shell", expectedSafe: true},
	}

	for _, test := range tests {
		lang, safe := GetLanguageByShebang(test.content)
		assert.Equal(s.T(), test.expectedLang, lang, fmt.Sprintf("%v: lang = %v, expected: %v", test.name, lang, test.expectedLang))
		assert.Equal(s.T(), test.expectedSafe, safe, fmt.Sprintf("%v: safe = %v, expected: %v", test.name, safe, test.expectedSafe))
	}
}

func (s *SimpleLinguistTestSuite) TestGetLanguageByExtension() {
	tests := []struct {
		name         string
		filename     string
		expectedLang string
		expectedSafe bool
	}{
		{name: "TestGetLanguageByExtension_1", filename: "foo.foo", expectedLang: OtherLanguage, expectedSafe: false},
		{name: "TestGetLanguageByExtension_2", filename: "foo.go", expectedLang: "Go", expectedSafe: true},
		{name: "TestGetLanguageByExtension_3", filename: "foo.go.php", expectedLang: "Hack", expectedSafe: false},
	}

	for _, test := range tests {
		lang, safe := GetLanguageByExtension(test.filename)
		assert.Equal(s.T(), test.expectedLang, lang, fmt.Sprintf("%v: lang = %v, expected: %v", test.name, lang, test.expectedLang))
		assert.Equal(s.T(), test.expectedSafe, safe, fmt.Sprintf("%v: safe = %v, expected: %v", test.name, safe, test.expectedSafe))
	}
}

func (s *SimpleLinguistTestSuite) TestGetLanguageByClassifier() {
	const samples = `.linguist/samples/`
	test := []struct {
		name       string
		filename   string
		candidates map[string]float64
		expected   string
	}{
		{name: "TestGetLanguageByClassifier_1", filename: filepath.Join(samples, "C/blob.c"), candidates: map[string]float64{"python": 1.00, "ruby": 1.00, "c": 1.00, "c++": 1.00}, expected: "C"},
		{name: "TestGetLanguageByClassifier_2", filename: filepath.Join(samples, "C/blob.c"), candidates: nil, expected: "C"},
		{name: "TestGetLanguageByClassifier_3", filename: filepath.Join(samples, "C/main.c"), candidates: nil, expected: "C"},
		{name: "TestGetLanguageByClassifier_4", filename: filepath.Join(samples, "C/blob.c"), candidates: map[string]float64{"python": 1.00, "ruby": 1.00, "c++": 1.00}, expected: "C++"},
		{name: "TestGetLanguageByClassifier_5", filename: filepath.Join(samples, "C/blob.c"), candidates: map[string]float64{"ruby": 1.00}, expected: "Ruby"},
		{name: "TestGetLanguageByClassifier_6", filename: filepath.Join(samples, "Python/django-models-base.py"), candidates: map[string]float64{"python": 1.00, "ruby": 1.00, "c": 1.00, "c++": 1.00}, expected: "Python"},
		{name: "TestGetLanguageByClassifier_7", filename: filepath.Join(samples, "Python/django-models-base.py"), candidates: nil, expected: "Python"},
	}

	for _, test := range test {
		content, err := ioutil.ReadFile(test.filename)
		assert.NoError(s.T(), err)

		lang := GetLanguageByClassifier(content, test.candidates, nil)
		assert.Equal(s.T(), test.expected, lang, fmt.Sprintf("%v: lang = %v, expected: %v", test.name, lang, test.expected))
	}
}

func (s *SimpleLinguistTestSuite) TestGetLanguageExtensions() {
	tests := []struct {
		name     string
		language string
		expected []string
	}{
		{name: "TestGetLanguageExtensions_1", language: "foo", expected: nil},
		{name: "TestGetLanguageExtensions_2", language: "COBOL", expected: []string{".cob", ".cbl", ".ccp", ".cobol", ".cpy"}},
		{name: "TestGetLanguageExtensions_3", language: "Maven POM", expected: nil},
	}

	for _, test := range tests {
		extensions := GetLanguageExtensions(test.language)
		assert.EqualValues(s.T(), test.expected, extensions, fmt.Sprintf("%v: extensions = %v, expected: %v", test.name, extensions, test.expected))
	}
}

func (s *SimpleLinguistTestSuite) TestGetLanguageType() {
	tests := []struct {
		name     string
		language string
		expected Type
	}{
		{name: "TestGetLanguageType_1", language: "BestLanguageEver", expected: Unknown},
		{name: "TestGetLanguageType_2", language: "JSON", expected: Data},
		{name: "TestGetLanguageType_3", language: "COLLADA", expected: Data},
		{name: "TestGetLanguageType_4", language: "Go", expected: Programming},
		{name: "TestGetLanguageType_5", language: "Brainfuck", expected: Programming},
		{name: "TestGetLanguageType_6", language: "HTML", expected: Markup},
		{name: "TestGetLanguageType_7", language: "Sass", expected: Markup},
		{name: "TestGetLanguageType_8", language: "AsciiDoc", expected: Prose},
		{name: "TestGetLanguageType_9", language: "Textile", expected: Prose},
	}

	for _, test := range tests {
		langType := GetLanguageType(test.language)
		assert.Equal(s.T(), test.expected, langType, fmt.Sprintf("%v: langType = %v, expected: %v", test.name, langType, test.expected))
	}
}

func (s *SimpleLinguistTestSuite) TestGetLanguageByAlias() {
	tests := []struct {
		name         string
		alias        string
		expectedLang string
		expectedOk   bool
	}{
		{name: "TestGetLanguageByAlias_1", alias: "BestLanguageEver", expectedLang: OtherLanguage, expectedOk: false},
		{name: "TestGetLanguageByAlias_2", alias: "aspx-vb", expectedLang: "ASP", expectedOk: true},
		{name: "TestGetLanguageByAlias_3", alias: "C++", expectedLang: "C++", expectedOk: true},
		{name: "TestGetLanguageByAlias_4", alias: "c++", expectedLang: "C++", expectedOk: true},
		{name: "TestGetLanguageByAlias_5", alias: "objc", expectedLang: "Objective-C", expectedOk: true},
		{name: "TestGetLanguageByAlias_6", alias: "golang", expectedLang: "Go", expectedOk: true},
		{name: "TestGetLanguageByAlias_7", alias: "GOLANG", expectedLang: "Go", expectedOk: true},
		{name: "TestGetLanguageByAlias_8", alias: "bsdmake", expectedLang: "Makefile", expectedOk: true},
		{name: "TestGetLanguageByAlias_9", alias: "xhTmL", expectedLang: "HTML", expectedOk: true},
		{name: "TestGetLanguageByAlias_10", alias: "python", expectedLang: "Python", expectedOk: true},
	}

	for _, test := range tests {
		lang, ok := GetLanguageByAlias(test.alias)
		assert.Equal(s.T(), test.expectedLang, lang, fmt.Sprintf("%v: lang = %v, expected: %v", test.name, lang, test.expectedLang))
		assert.Equal(s.T(), test.expectedOk, ok, fmt.Sprintf("%v: ok = %v, expected: %v", test.name, ok, test.expectedOk))
	}
}

func (s *SimpleLinguistTestSuite) TestLinguistCorpus() {
	const (
		samplesDir   = ".linguist/samples"
		filenamesDir = "filenames"
	)

	var cornerCases = map[string]bool{
		"hello.ms": true,
	}

	var total, failed, ok, other int
	var expected string
	filepath.Walk(samplesDir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			if f.Name() != filenamesDir {
				expected = f.Name()
			}

			return nil
		}

		filename := filepath.Base(path)
		content, _ := ioutil.ReadFile(path)

		total++
		obtained := GetLanguage(filename, content)
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

		}

		if _, ok := cornerCases[filename]; ok {
			fmt.Printf("\t\t[condidered corner case] %s\t%s\t%s\t%s\n", filename, expected, obtained, status)
		} else {
			assert.Equal(s.T(), expected, obtained, fmt.Sprintf("%s\t%s\t%s\t%s\n", filename, expected, obtained, status))
		}

		return nil
	})

	fmt.Printf("\t\ttotal files: %d, ok: %d, failed: %d, other: %d\n", total, ok, failed, other)
}
