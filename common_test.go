package enry

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-enry/go-enry/v2/data"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const linguistURL = "https://github.com/github/linguist.git"
const linguistClonedEnvVar = "ENRY_TEST_REPO"

type EnryTestSuite struct {
	suite.Suite
	tmpLinguist     string
	needToClone     bool
	samplesDir      string
	testFixturesDir string
}

func (s *EnryTestSuite) TestRegexpEdgeCases() {
	var regexpEdgeCases = []struct {
		lang     string
		filename string
	}{
		{lang: "ActionScript", filename: "FooBar.as"},
		{lang: "Forth", filename: "asm.fr"},
		{lang: "X PixMap", filename: "cc-public_domain_mark_white.pm"},
		//{lang: "SQL", filename: "drop_stuff.sql"}, // https://github.com/src-d/enry/issues/194
		{lang: "Fstar", filename: "Hacl.Spec.Bignum.Fmul.fst"},
		{lang: "C++", filename: "Types.h"},
	}

	for _, r := range regexpEdgeCases {
		filename := filepath.Join(s.tmpLinguist, "samples", r.lang, r.filename)

		content, err := ioutil.ReadFile(filename)
		require.NoError(s.T(), err)

		lang := GetLanguage(r.filename, content)
		s.T().Logf("File:%s, lang:%s", filename, lang)

		expLang, _ := data.LanguageByAlias(r.lang)
		require.EqualValues(s.T(), expLang, lang)
	}
}

func Test_EnryTestSuite(t *testing.T) {
	suite.Run(t, new(EnryTestSuite))
}

func (s *EnryTestSuite) SetupSuite() {
	var err error
	s.tmpLinguist = os.Getenv(linguistClonedEnvVar)
	s.needToClone = s.tmpLinguist == ""
	if s.needToClone {
		s.tmpLinguist, err = ioutil.TempDir("", "linguist-")
		require.NoError(s.T(), err)
		s.T().Logf("Cloning Linguist repo to '%s' as %s was not set\n",
			s.tmpLinguist, linguistClonedEnvVar)
		cmd := exec.Command("git", "clone", linguistURL, s.tmpLinguist)
		err = cmd.Run()
		require.NoError(s.T(), err)
	}
	s.samplesDir = filepath.Join(s.tmpLinguist, "samples")
	s.T().Logf("using samples from %s", s.samplesDir)

	s.testFixturesDir = filepath.Join(s.tmpLinguist, "test", "fixtures")
	s.T().Logf("using test fixtures from %s", s.samplesDir)

	cwd, err := os.Getwd()
	assert.NoError(s.T(), err)

	err = os.Chdir(s.tmpLinguist)
	assert.NoError(s.T(), err)

	cmd := exec.Command("git", "checkout", data.LinguistCommit)
	err = cmd.Run()
	assert.NoError(s.T(), err)

	err = os.Chdir(cwd)
	assert.NoError(s.T(), err)
}

func (s *EnryTestSuite) TearDownSuite() {
	if s.needToClone {
		err := os.RemoveAll(s.tmpLinguist)
		assert.NoError(s.T(), err)
	}
}

func (s *EnryTestSuite) TestGetLanguage() {
	tests := []struct {
		name     string
		filename string
		content  []byte
		expected string
		safe     bool
	}{
		{name: "TestGetLanguage_0", filename: "foo.h", content: []byte{}, expected: "C"},
		{name: "TestGetLanguage_1", filename: "foo.py", content: []byte{}, expected: "Python"},
		{name: "TestGetLanguage_2", filename: "foo.m", content: []byte(":- module"), expected: "Mercury"},
		{name: "TestGetLanguage_3", filename: "foo.m", content: nil, expected: "MATLAB"},
		{name: "TestGetLanguage_4", filename: "foo.mo", content: []byte{0xDE, 0x12, 0x04, 0x95, 0x00, 0x00, 0x00, 0x00}, expected: OtherLanguage},
		{name: "TestGetLanguage_5", filename: "", content: nil, expected: OtherLanguage},
	}

	for _, test := range tests {
		language := GetLanguage(test.filename, test.content)
		assert.Equal(s.T(), test.expected, language, fmt.Sprintf("%v: %v, expected: %v", test.name, language, test.expected))
	}
}

func (s *EnryTestSuite) TestGetLanguages() {
	tests := []struct {
		name     string
		filename string
		content  []byte
		expected []string
	}{
		// With no content or filename, no language can be detected
		{name: "TestGetLanguages_0", filename: "", content: []byte{}, expected: nil},
		// The strategy that will match is GetLanguagesByExtension. Lacking content, it will return those results.
		{name: "TestGetLanguages_1", filename: "foo.h", content: []byte{}, expected: []string{"C"}},
		// GetLanguagesByExtension will return an unambiguous match when there is a single result.
		{name: "TestGetLanguages_2", filename: "foo.groovy", content: []byte{}, expected: []string{"Groovy"}},
		// GetLanguagesByExtension will return "Rust", "RenderScript" for .rs,
		// then GetLanguagesByContent will take the first rule that matches (in this case Rust)
		{name: "TestGetLanguages_3", filename: "foo.rs", content: []byte("use \n#include"), expected: []string{"Rust"}},
		// .. and in this case, RenderScript (no content that matches a Rust regex can be included, because it runs first.)
		{name: "TestGetLanguages_4", filename: "foo.rs", content: []byte("#include"), expected: []string{"RenderScript"}},
		// GetLanguagesByExtension will return "AMPL", "Linux Kernel Module", "Modula-2", "XML",
		// then GetLanguagesByContent will ALWAYS return Linux Kernel Module and AMPL when there is no content,
		// and no further classifier can do anything without content
		{name: "TestGetLanguages_5", filename: "foo.mod", content: []byte{}, expected: []string{"Linux Kernel Module", "AMPL"}},
		// ...with some AMPL tokens, the DefaultClassifier will pick AMPL as the most likely language.
		{name: "TestGetLanguages_6", filename: "foo.mod", content: []byte("BEAMS ROWS - TotalWeight"), expected: []string{"AMPL", "Linux Kernel Module"}},
	}

	for _, test := range tests {
		languages := GetLanguages(test.filename, test.content)
		assert.Equal(s.T(), test.expected, languages, fmt.Sprintf("%v: %v, expected: %v", test.name, languages, test.expected))
	}
}

func (s *EnryTestSuite) TestGetLanguagesByModelineLinguist() {
	var modelinesDir = filepath.Join(s.tmpLinguist, "test", "fixtures", "Data", "Modelines")

	tests := []struct {
		name       string
		filename   string
		candidates []string
		expected   []string
	}{
		// Emacs
		{name: "TestGetLanguagesByModelineLinguist_1", filename: filepath.Join(modelinesDir, "example_smalltalk.md"), expected: []string{"Smalltalk"}},
		{name: "TestGetLanguagesByModelineLinguist_2", filename: filepath.Join(modelinesDir, "fundamentalEmacs.c"), expected: []string{"Text"}},
		{name: "TestGetLanguagesByModelineLinguist_3", filename: filepath.Join(modelinesDir, "iamphp.inc"), expected: []string{"PHP"}},
		{name: "TestGetLanguagesByModelineLinguist_4", filename: filepath.Join(modelinesDir, "seeplusplusEmacs1"), expected: []string{"C++"}},
		{name: "TestGetLanguagesByModelineLinguist_5", filename: filepath.Join(modelinesDir, "seeplusplusEmacs2"), expected: []string{"C++"}},
		{name: "TestGetLanguagesByModelineLinguist_6", filename: filepath.Join(modelinesDir, "seeplusplusEmacs3"), expected: []string{"C++"}},
		{name: "TestGetLanguagesByModelineLinguist_7", filename: filepath.Join(modelinesDir, "seeplusplusEmacs4"), expected: []string{"C++"}},
		{name: "TestGetLanguagesByModelineLinguist_8", filename: filepath.Join(modelinesDir, "seeplusplusEmacs5"), expected: []string{"C++"}},
		{name: "TestGetLanguagesByModelineLinguist_9", filename: filepath.Join(modelinesDir, "seeplusplusEmacs6"), expected: []string{"C++"}},
		{name: "TestGetLanguagesByModelineLinguist_10", filename: filepath.Join(modelinesDir, "seeplusplusEmacs7"), expected: []string{"C++"}},
		{name: "TestGetLanguagesByModelineLinguist_11", filename: filepath.Join(modelinesDir, "seeplusplusEmacs9"), expected: []string{"C++"}},
		{name: "TestGetLanguagesByModelineLinguist_12", filename: filepath.Join(modelinesDir, "seeplusplusEmacs10"), expected: []string{"C++"}},
		{name: "TestGetLanguagesByModelineLinguist_13", filename: filepath.Join(modelinesDir, "seeplusplusEmacs11"), expected: []string{"C++"}},
		{name: "TestGetLanguagesByModelineLinguist_14", filename: filepath.Join(modelinesDir, "seeplusplusEmacs12"), expected: []string{"C++"}},

		// Vim
		{name: "TestGetLanguagesByModelineLinguist_15", filename: filepath.Join(modelinesDir, "seeplusplus"), expected: []string{"C++"}},
		{name: "TestGetLanguagesByModelineLinguist_16", filename: filepath.Join(modelinesDir, "iamjs.pl"), expected: []string{"JavaScript"}},
		{name: "TestGetLanguagesByModelineLinguist_17", filename: filepath.Join(modelinesDir, "iamjs2.pl"), expected: []string{"JavaScript"}},
		{name: "TestGetLanguagesByModelineLinguist_18", filename: filepath.Join(modelinesDir, "not_perl.pl"), expected: []string{"Prolog"}},
		{name: "TestGetLanguagesByModelineLinguist_19", filename: filepath.Join(modelinesDir, "ruby"), expected: []string{"Ruby"}},
		{name: "TestGetLanguagesByModelineLinguist_20", filename: filepath.Join(modelinesDir, "ruby2"), expected: []string{"Ruby"}},
		{name: "TestGetLanguagesByModelineLinguist_21", filename: filepath.Join(modelinesDir, "ruby3"), expected: []string{"Ruby"}},
		{name: "TestGetLanguagesByModelineLinguist_22", filename: filepath.Join(modelinesDir, "ruby4"), expected: []string{"Ruby"}},
		{name: "TestGetLanguagesByModelineLinguist_23", filename: filepath.Join(modelinesDir, "ruby5"), expected: []string{"Ruby"}},
		{name: "TestGetLanguagesByModelineLinguist_24", filename: filepath.Join(modelinesDir, "ruby6"), expected: []string{"Ruby"}},
		{name: "TestGetLanguagesByModelineLinguist_25", filename: filepath.Join(modelinesDir, "ruby7"), expected: []string{"Ruby"}},
		{name: "TestGetLanguagesByModelineLinguist_26", filename: filepath.Join(modelinesDir, "ruby8"), expected: []string{"Ruby"}},
		{name: "TestGetLanguagesByModelineLinguist_27", filename: filepath.Join(modelinesDir, "ruby9"), expected: []string{"Ruby"}},
		{name: "TestGetLanguagesByModelineLinguist_28", filename: filepath.Join(modelinesDir, "ruby10"), expected: []string{"Ruby"}},
		{name: "TestGetLanguagesByModelineLinguist_29", filename: filepath.Join(modelinesDir, "ruby11"), expected: []string{"Ruby"}},
		{name: "TestGetLanguagesByModelineLinguist_30", filename: filepath.Join(modelinesDir, "ruby12"), expected: []string{"Ruby"}},
		{name: "TestGetLanguagesByModelineLinguist_31", filename: filepath.Join(s.samplesDir, "C++/runtime-compiler.cc"), expected: nil},
		{name: "TestGetLanguagesByModelineLinguist_32", filename: "", expected: nil},
	}

	for _, test := range tests {
		var content []byte
		var err error

		if test.filename != "" {
			content, err = ioutil.ReadFile(test.filename)
			assert.NoError(s.T(), err)
		}

		languages := GetLanguagesByModeline(test.filename, content, test.candidates)
		assert.Equal(s.T(), test.expected, languages, fmt.Sprintf("%v: languages = %v, expected: %v", test.name, languages, test.expected))
	}
}

func (s *EnryTestSuite) TestGetLanguagesByModeline() {
	const (
		wrongVim  = `# vim: set syntax=ruby ft  =python filetype=perl :`
		rightVim  = `/* vim: set syntax=python ft   =python filetype=python */`
		noLangVim = `/* vim: set shiftwidth=4 softtabstop=0 cindent cinoptions={1s: */`
	)

	tests := []struct {
		name       string
		filename   string
		content    []byte
		candidates []string
		expected   []string
	}{
		{name: "TestGetLanguagesByModeline_1", content: []byte(wrongVim), expected: nil},
		{name: "TestGetLanguagesByModeline_2", content: []byte(rightVim), expected: []string{"Python"}},
		{name: "TestGetLanguagesByModeline_3", content: []byte(noLangVim), expected: nil},
		{name: "TestGetLanguagesByModeline_4", content: nil, expected: nil},
		{name: "TestGetLanguagesByModeline_5", content: []byte{}, expected: nil},
	}

	for _, test := range tests {
		languages := GetLanguagesByModeline(test.filename, test.content, test.candidates)
		assert.Equal(s.T(), test.expected, languages, fmt.Sprintf("%v: languages = %v, expected: %v", test.name, languages, test.expected))
	}
}

func (s *EnryTestSuite) TestGetLanguagesByFilename() {
	tests := []struct {
		name       string
		filename   string
		content    []byte
		candidates []string
		expected   []string
	}{
		{name: "TestGetLanguagesByFilename_1", filename: "unknown.interpreter", expected: nil},
		{name: "TestGetLanguagesByFilename_2", filename: ".bashrc", expected: []string{"Shell"}},
		{name: "TestGetLanguagesByFilename_3", filename: "Dockerfile", expected: []string{"Dockerfile"}},
		{name: "TestGetLanguagesByFilename_4", filename: "Makefile.frag", expected: []string{"Makefile"}},
		{name: "TestGetLanguagesByFilename_5", filename: "makefile", expected: []string{"Makefile"}},
		{name: "TestGetLanguagesByFilename_6", filename: "Vagrantfile", expected: []string{"Ruby"}},
		{name: "TestGetLanguagesByFilename_7", filename: "_vimrc", expected: []string{"Vim Script"}},
		{name: "TestGetLanguagesByFilename_8", filename: "pom.xml", expected: []string{"Maven POM"}},
		{name: "TestGetLanguagesByFilename_9", filename: "", expected: nil},
	}

	for _, test := range tests {
		languages := GetLanguagesByFilename(test.filename, test.content, test.candidates)
		assert.Equal(s.T(), len(test.expected), len(languages), fmt.Sprintf("%v: number of languages = %v, expected: %v", test.name, len(languages), len(test.expected)))
		for i := range languages { // case-insensitive name comparison
			assert.True(s.T(), strings.EqualFold(test.expected[i], languages[i]), fmt.Sprintf("%v: languages = %v, expected: %v", test.name, languages, test.expected))
		}
	}
}

func (s *EnryTestSuite) TestGetLanguagesByShebang() {
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
		name       string
		filename   string
		content    []byte
		candidates []string
		expected   []string
	}{
		{name: "TestGetLanguagesByShebang_1", content: []byte(`#!/unknown/interpreter`), expected: nil},
		{name: "TestGetLanguagesByShebang_2", content: []byte(`no shebang`), expected: nil},
		{name: "TestGetLanguagesByShebang_3", content: []byte(`#!/usr/bin/env`), expected: nil},
		{name: "TestGetLanguagesByShebang_4", content: []byte(`#!/usr/bin/python -tt`), expected: []string{"Python"}},
		{name: "TestGetLanguagesByShebang_5", content: []byte(`#!/usr/bin/env python2.6`), expected: []string{"Python"}},
		{name: "TestGetLanguagesByShebang_6", content: []byte(`#!/usr/bin/env perl`), expected: []string{"Perl", "Pod"}},
		{name: "TestGetLanguagesByShebang_7", content: []byte(`#!	/bin/sh`), expected: []string{"Shell"}},
		{name: "TestGetLanguagesByShebang_8", content: []byte(`#!bash`), expected: []string{"Shell"}},
		{name: "TestGetLanguagesByShebang_9", content: []byte(multilineExecHack), expected: []string{"Tcl"}},
		{name: "TestGetLanguagesByShebang_10", content: []byte(multilineNoExecHack), expected: []string{"Shell"}},
		{name: "TestGetLanguagesByShebang_11", content: []byte(`#!/envinpath/python`), expected: []string{"Python"}},

		{name: "TestGetLanguagesByShebang_12", content: []byte(""), expected: nil},
		{name: "TestGetLanguagesByShebang_13", content: []byte("foo"), expected: nil},
		{name: "TestGetLanguagesByShebang_14", content: []byte("#bar"), expected: nil},
		{name: "TestGetLanguagesByShebang_15", content: []byte("#baz"), expected: nil},
		{name: "TestGetLanguagesByShebang_16", content: []byte("///"), expected: nil},
		{name: "TestGetLanguagesByShebang_17", content: []byte("\n\n\n\n\n"), expected: nil},
		{name: "TestGetLanguagesByShebang_18", content: []byte(" #!/usr/sbin/ruby"), expected: nil},
		{name: "TestGetLanguagesByShebang_19", content: []byte("\n#!/usr/sbin/ruby"), expected: nil},
		{name: "TestGetLanguagesByShebang_20", content: []byte("#!"), expected: nil},
		{name: "TestGetLanguagesByShebang_21", content: []byte("#! "), expected: nil},
		{name: "TestGetLanguagesByShebang_22", content: []byte("#!/usr/bin/env"), expected: nil},
		{name: "TestGetLanguagesByShebang_23", content: []byte("#!/usr/bin/env osascript -l JavaScript"), expected: nil},
		{name: "TestGetLanguagesByShebang_24", content: []byte("#!/usr/bin/env osascript -l AppleScript"), expected: nil},
		{name: "TestGetLanguagesByShebang_25", content: []byte("#!/usr/bin/env osascript -l foobar"), expected: nil},
		{name: "TestGetLanguagesByShebang_26", content: []byte("#!/usr/bin/osascript -l JavaScript"), expected: nil},
		{name: "TestGetLanguagesByShebang_27", content: []byte("#!/usr/bin/osascript -l foobar"), expected: nil},

		{name: "TestGetLanguagesByShebang_28", content: []byte("#!/usr/sbin/ruby\n# bar"), expected: []string{"Ruby"}},
		{name: "TestGetLanguagesByShebang_29", content: []byte("#!/usr/bin/ruby\n# foo"), expected: []string{"Ruby"}},
		{name: "TestGetLanguagesByShebang_30", content: []byte("#!/usr/sbin/ruby"), expected: []string{"Ruby"}},
		{name: "TestGetLanguagesByShebang_31", content: []byte("#!/usr/sbin/ruby foo bar baz\n"), expected: []string{"Ruby"}},

		{name: "TestGetLanguagesByShebang_32", content: []byte("#!/usr/bin/env Rscript\n# example R script\n#\n"), expected: []string{"R"}},
		{name: "TestGetLanguagesByShebang_33", content: []byte("#!/usr/bin/env ruby\n# baz"), expected: []string{"Ruby"}},

		{name: "TestGetLanguagesByShebang_34", content: []byte("#!/usr/bin/bash\n"), expected: []string{"Shell"}},
		{name: "TestGetLanguagesByShebang_35", content: []byte("#!/bin/sh"), expected: []string{"Shell"}},
		{name: "TestGetLanguagesByShebang_36", content: []byte("#!/bin/python\n# foo\n# bar\n# baz"), expected: []string{"Python"}},
		{name: "TestGetLanguagesByShebang_37", content: []byte("#!/usr/bin/python2.7\n\n\n\n"), expected: []string{"Python"}},
		{name: "TestGetLanguagesByShebang_38", content: []byte("#!/usr/bin/python3\n\n\n\n"), expected: []string{"Python"}},
		{name: "TestGetLanguagesByShebang_39", content: []byte("#!/usr/bin/sbcl --script\n\n"), expected: []string{"Common Lisp"}},
		{name: "TestGetLanguagesByShebang_40", content: []byte("#! perl"), expected: []string{"Perl", "Pod"}},

		{name: "TestGetLanguagesByShebang_41", content: []byte("#!/bin/sh\n\n\nexec ruby $0 $@"), expected: []string{"Ruby"}},
		{name: "TestGetLanguagesByShebang_42", content: []byte("#! /usr/bin/env A=003 B=149 C=150 D=xzd E=base64 F=tar G=gz H=head I=tail sh"), expected: []string{"Shell"}},
		{name: "TestGetLanguagesByShebang_43", content: []byte("#!/usr/bin/env foo=bar bar=foo python -cos=__import__(\"os\");"), expected: []string{"Python"}},
		{name: "TestGetLanguagesByShebang_44", content: []byte("#!/usr/bin/env osascript"), expected: []string{"AppleScript"}},
		{name: "TestGetLanguagesByShebang_45", content: []byte("#!/usr/bin/osascript"), expected: []string{"AppleScript"}},

		{name: "TestGetLanguagesByShebang_46", content: []byte("#!/usr/bin/env -vS ruby -wKU\nputs ?t+?e+?s+?t"), expected: []string{"Ruby"}},
		{name: "TestGetLanguagesByShebang_47", content: []byte("#!/usr/bin/env --split-string sed -f\ny/a/A/"), expected: []string{"sed"}},
		{name: "TestGetLanguagesByShebang_48", content: []byte("#!/usr/bin/env -S GH_TOKEN=ghp_*** deno run --allow-net\nconsole.log(1);"), expected: []string{"TypeScript"}},
	}

	for _, test := range tests {
		languages := GetLanguagesByShebang(test.filename, test.content, test.candidates)
		assert.Equal(s.T(), test.expected, languages, fmt.Sprintf("%v: languages = %v, expected: %v", test.name, languages, test.expected))
	}
}

func (s *EnryTestSuite) TestGetLanguagesByExtension() {
	tests := []struct {
		name       string
		filename   string
		content    []byte
		candidates []string
		expected   []string
	}{
		{name: "TestGetLanguagesByExtension_0", filename: "foo.h", expected: []string{"C", "C++", "Objective-C"}},
		{name: "TestGetLanguagesByExtension_1", filename: "foo.foo", expected: nil},
		{name: "TestGetLanguagesByExtension_2", filename: "foo.go", expected: []string{"Go"}},
		{name: "TestGetLanguagesByExtension_3", filename: "foo.go.php", expected: []string{"Hack", "PHP"}},
		{name: "TestGetLanguagesByExtension_4", filename: "", expected: nil},
	}

	for _, test := range tests {
		languages := GetLanguagesByExtension(test.filename, test.content, test.candidates)
		assert.Equal(s.T(), test.expected, languages, fmt.Sprintf("%v: languages = %v, expected: %v", test.name, languages, test.expected))
	}
}

func (s *EnryTestSuite) TestGetLanguagesByManpage() {
	tests := []struct {
		name       string
		filename   string
		content    []byte
		candidates []string
		expected   []string
	}{
		{name: "TestGetLanguagesByManpage_1", filename: "bsdmalloc.3malloc", expected: []string{"Roff Manpage", "Roff"}},
		{name: "TestGetLanguagesByManpage_2", filename: "dirent.h.0p", expected: []string{"Roff Manpage", "Roff"}},
		{name: "TestGetLanguagesByManpage_3", filename: "linguist.1gh", expected: []string{"Roff Manpage", "Roff"}},
		{name: "TestGetLanguagesByManpage_4", filename: "test.1.in", expected: []string{"Roff Manpage", "Roff"}},
		{name: "TestGetLanguagesByManpage_5", filename: "test.man.in", expected: []string{"Roff Manpage", "Roff"}},
		{name: "TestGetLanguagesByManpage_6", filename: "test.mdoc.in", expected: []string{"Roff Manpage", "Roff"}},
		{name: "TestGetLanguagesByManpage_7", filename: "foo.h", expected: nil},
		{name: "TestGetLanguagesByManpage_8", filename: "", expected: nil},
	}

	for _, test := range tests {
		languages := GetLanguagesByManpage(test.filename, test.content, test.candidates)
		assert.Equal(s.T(), test.expected, languages, fmt.Sprintf("%v: languages = %v, expected: %v", test.name, languages, test.expected))
	}
}

func (s *EnryTestSuite) TestGetLanguagesByXML() {
	tests := []struct {
		name       string
		filename   string
		candidates []string
		expected   []string
	}{
		{name: "TestGetLanguagesByXML_1", filename: filepath.Join(s.testFixturesDir, "XML/app.config"), expected: []string{"XML"}},
		{name: "TestGetLanguagesByXML_2", filename: filepath.Join(s.testFixturesDir, "XML/AssertionIDRequestOptionalAttributes.xml.svn-base"), expected: []string{"XML"}},
		// no XML header so should not be identified by this strategy
		{name: "TestGetLanguagesByXML_3", filename: filepath.Join(s.samplesDir, "XML/libsomething.dll.config"), expected: nil},
		{name: "TestGetLanguagesByXML_4", filename: filepath.Join(s.samplesDir, "Eagle/Eagle.sch"), candidates: []string{"Eagle"}, expected: []string{"Eagle"}},
	}

	for _, test := range tests {
		content, err := ioutil.ReadFile(test.filename)
		assert.NoError(s.T(), err)

		languages := GetLanguagesByXML(test.filename, content, test.candidates)
		assert.Equal(s.T(), test.expected, languages, fmt.Sprintf("%v: languages = %v, expected: %v", test.name, languages, test.expected))
	}
}

func (s *EnryTestSuite) TestGetLanguagesByClassifier() {
	test := []struct {
		name       string
		filename   string
		candidates []string
		expected   string
	}{
		{name: "TestGetLanguagesByClassifier_1", filename: filepath.Join(s.samplesDir, "C/blob.c"), candidates: []string{"python", "ruby", "c", "c++"}, expected: "C"},
		{name: "TestGetLanguagesByClassifier_2", filename: filepath.Join(s.samplesDir, "C/blob.c"), candidates: nil, expected: OtherLanguage},
		{name: "TestGetLanguagesByClassifier_3", filename: filepath.Join(s.samplesDir, "C++/runtime-compiler.cc"), candidates: []string{}, expected: OtherLanguage},
		{name: "TestGetLanguagesByClassifier_4", filename: filepath.Join(s.samplesDir, "C/blob.c"), candidates: []string{"python", "ruby", "c++"}, expected: "C++"},
		{name: "TestGetLanguagesByClassifier_5", filename: filepath.Join(s.samplesDir, "C/blob.c"), candidates: []string{"ruby"}, expected: "Ruby"},
		{name: "TestGetLanguagesByClassifier_6", filename: filepath.Join(s.samplesDir, "Python/django-models-base.py"), candidates: []string{"python", "ruby", "c", "c++"}, expected: "Python"},
		{name: "TestGetLanguagesByClassifier_7", filename: "", candidates: []string{"python"}, expected: "Python"},
	}

	for _, test := range test {
		var content []byte
		var err error

		if test.filename != "" {
			content, err = ioutil.ReadFile(test.filename)
			assert.NoError(s.T(), err)
		}

		languages := GetLanguagesByClassifier(test.filename, content, test.candidates)
		var language string
		if len(languages) == 0 {
			language = OtherLanguage
		} else {
			language = languages[0]
		}

		assert.Equal(s.T(), test.expected, language, fmt.Sprintf("%v: language = %v, expected: %v", test.name, language, test.expected))
	}
}

func (s *EnryTestSuite) TestGetLanguagesBySpecificClassifier() {
	test := []struct {
		name       string
		filename   string
		candidates []string
		classifier classifier
		expected   string
	}{
		{name: "TestGetLanguagesByClassifier_1", filename: filepath.Join(s.samplesDir, "C/blob.c"), candidates: []string{"python", "ruby", "c", "c++"}, classifier: defaultClassifier, expected: "C"},
		{name: "TestGetLanguagesByClassifier_2", filename: filepath.Join(s.samplesDir, "C/blob.c"), candidates: nil, classifier: defaultClassifier, expected: "C"},
		{name: "TestGetLanguagesByClassifier_3", filename: filepath.Join(s.samplesDir, "C++/runtime-compiler.cc"), candidates: []string{}, classifier: defaultClassifier, expected: "C++"},
		{name: "TestGetLanguagesByClassifier_4", filename: filepath.Join(s.samplesDir, "C/blob.c"), candidates: []string{"python", "ruby", "c++"}, classifier: defaultClassifier, expected: "C++"},
		{name: "TestGetLanguagesByClassifier_5", filename: filepath.Join(s.samplesDir, "C/blob.c"), candidates: []string{"ruby"}, classifier: defaultClassifier, expected: "Ruby"},
		{name: "TestGetLanguagesByClassifier_6", filename: filepath.Join(s.samplesDir, "Python/django-models-base.py"), candidates: []string{"python", "ruby", "c", "c++"}, classifier: defaultClassifier, expected: "Python"},
		{name: "TestGetLanguagesByClassifier_7", filename: os.DevNull, candidates: nil, classifier: defaultClassifier, expected: "XML"},
	}

	for _, test := range test {
		content, err := ioutil.ReadFile(test.filename)
		assert.NoError(s.T(), err)

		languages := getLanguagesBySpecificClassifier(content, test.candidates, test.classifier)
		var language string
		if len(languages) == 0 {
			language = OtherLanguage
		} else {
			language = languages[0]
		}

		assert.Equal(s.T(), test.expected, language, fmt.Sprintf("%v: language = %v, expected: %v", test.name, language, test.expected))
	}
}

func (s *EnryTestSuite) TestGetLanguageExtensions() {
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

func (s *EnryTestSuite) TestGetLanguageType() {
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

func (s *EnryTestSuite) TestGetLanguageGroup() {
	tests := []struct {
		name     string
		language string
		expected string
	}{
		{name: "TestGetLanguageGroup_1", language: "BestLanguageEver", expected: ""},
		{name: "TestGetLanguageGroup_2", language: "Bison", expected: "Yacc"},
		{name: "TestGetLanguageGroup_3", language: "HTML+PHP", expected: "HTML"},
		{name: "TestGetLanguageGroup_4", language: "HTML", expected: ""},
	}

	for _, test := range tests {
		langGroup := GetLanguageGroup(test.language)
		assert.Equal(s.T(), test.expected, langGroup, fmt.Sprintf("%v: langGroup = %v, expected: %v", test.name, langGroup, test.expected))
	}
}

func (s *EnryTestSuite) TestGetLanguageByAlias() {
	tests := []struct {
		name         string
		alias        string
		expectedLang string
		expectedOk   bool
	}{
		{name: "TestGetLanguageByAlias_1", alias: "BestLanguageEver", expectedLang: OtherLanguage, expectedOk: false},
		{name: "TestGetLanguageByAlias_2", alias: "aspx-vb", expectedLang: "ASP.NET", expectedOk: true},
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

func (s *EnryTestSuite) TestLinguistCorpus() {
	const filenamesDir = "filenames"
	var cornerCases = map[string]bool{
		"drop_stuff.sql":        true, // https://github.com/src-d/enry/issues/194
		"textobj-rubyblock.vba": true, // Because of unsupported negative lookahead RE syntax (https://github.com/github/linguist/blob/8083cb5a89cee2d99f5a988f165994d0243f0d1e/lib/linguist/heuristics.yml#L521)
		// .es and .ice fail heuristics parsing, but do not fail any tests
	}

	var total, failed, ok, other int
	var expected string
	filepath.Walk(s.samplesDir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			if f.Name() != filenamesDir {
				expected, _ = data.LanguageByAlias(f.Name())
			}

			return nil
		}

		filename := filepath.Base(path)
		content, _ := ioutil.ReadFile(path)

		total++
		obtained := GetLanguage(filename, content)
		if obtained == OtherLanguage {
			obtained = "Other"
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
			s.T().Logf("\t\t[considered corner case] %s\texpected: %s\tobtained: %s\tstatus: %s\n", filename, expected, obtained, status)
		} else {
			assert.Equal(s.T(), expected, obtained, fmt.Sprintf("%s\texpected: %s\tobtained: %s\tstatus: %s\n", filename, expected, obtained, status))
		}

		return nil
	})

	s.T().Logf("\t\ttotal files: %d, ok: %d, failed: %d, other: %d\n", total, ok, failed, other)
}

func (s *EnryTestSuite) TestGetLanguageID() {
	tests := []struct {
		name       string
		language   string
		expectedID int
		found      bool
	}{
		{name: "TestGetLanguageID_1", language: "1C Enterprise", expectedID: 0, found: true},
		{name: "TestGetLanguageID_2", language: "BestLanguageEver", expectedID: 0, found: false},
		{name: "TestGetLanguageID_3", language: "C++", expectedID: 43, found: true},
		{name: "TestGetLanguageID_5", language: "Objective-C", expectedID: 257, found: true},
		{name: "TestGetLanguageID_6", language: "golang", expectedID: 0, found: false}, // Aliases are not supported
		{name: "TestGetLanguageID_7", language: "Go", expectedID: 132, found: true},
		{name: "TestGetLanguageID_8", language: "Makefile", expectedID: 220, found: true},
	}

	for _, test := range tests {
		id, found := GetLanguageID(test.language)
		assert.Equal(s.T(), test.expectedID, id, fmt.Sprintf("%v: id = %v, expected: %v", test.name, id, test.expectedID))
		assert.Equal(s.T(), test.found, found, fmt.Sprintf("%v: found = %t, expected: %t", test.name, found, test.found))
	}
}

func (s *EnryTestSuite) TestGetLanguageInfo() {
	tests := []struct {
		name       string
		language   string
		expectedID int
		error      bool
	}{
		{name: "TestGetLanguageID_1", language: "1C Enterprise", expectedID: 0},
		{name: "TestGetLanguageID_2", language: "BestLanguageEver", error: true},
		{name: "TestGetLanguageID_3", language: "C++", expectedID: 43},
		{name: "TestGetLanguageID_5", language: "Objective-C", expectedID: 257},
		{name: "TestGetLanguageID_6", language: "golang", error: true}, // Aliases are not supported
		{name: "TestGetLanguageID_7", language: "Go", expectedID: 132},
		{name: "TestGetLanguageID_8", language: "Makefile", expectedID: 220},
	}

	for _, test := range tests {
		info, err := GetLanguageInfo(test.language)
		if test.error {
			assert.Error(s.T(), err, "%v: expected error for %q", test.name, test.language)
		} else {
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), test.expectedID, info.LanguageID, fmt.Sprintf("%v: id = %v, expected: %v", test.name, info.LanguageID, test.expectedID))
		}
	}
}

func (s *EnryTestSuite) TestGetLanguageInfoByID() {
	tests := []struct {
		name         string
		id           int
		expectedName string
		error        bool
	}{
		{name: "TestGetLanguageID_1", id: 0, expectedName: "1C Enterprise"},
		{name: "TestGetLanguageID_2", id: -1, error: true},
		{name: "TestGetLanguageID_3", id: 43, expectedName: "C++"},
		{name: "TestGetLanguageID_5", id: 257, expectedName: "Objective-C"},
		{name: "TestGetLanguageID_7", id: 132, expectedName: "Go"},
		{name: "TestGetLanguageID_8", id: 220, expectedName: "Makefile"},
	}

	for _, test := range tests {
		info, err := GetLanguageInfoByID(test.id)
		if test.error {
			assert.Error(s.T(), err, "%v: expected error for %q", test.name, test.id)
		} else {
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), test.expectedName, info.Name, fmt.Sprintf("%v: id = %v, expected: %v", test.name, test.id, test.expectedName))
		}
	}
}
