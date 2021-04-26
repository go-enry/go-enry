package generator

import (
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-enry/go-enry/v2/internal/tokenizer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	linguistURL          = "https://github.com/github/linguist.git"
	linguistClonedEnvVar = "ENRY_TEST_REPO"
	commit               = "40992ba7f86889f80dfed3ba95e11e1082200bad"
	samplesDir           = "samples"
	languagesFile        = filepath.Join("lib", "linguist", "languages.yml")

	testDir   = "test_files"
	assetsDir = filepath.Join("..", "assets")

	// Extensions test
	extensionGold         = filepath.Join(testDir, "extension.gold")
	extensionTestTmplPath = filepath.Join(assetsDir, "extension.go.tmpl")
	extensionTestTmplName = "extension.go.tmpl"

	// Heuristics test
	heuristicsTestFile  = filepath.Join("lib", "linguist", "heuristics.yml")
	contentGold         = filepath.Join(testDir, "content.gold")
	contentTestTmplPath = filepath.Join(assetsDir, "content.go.tmpl")
	contentTestTmplName = "content.go.tmpl"

	// Vendor test
	vendorTestFile     = filepath.Join("lib", "linguist", "vendor.yml")
	vendorGold         = filepath.Join(testDir, "vendor.gold")
	vendorTestTmplPath = filepath.Join(assetsDir, "vendor.go.tmpl")
	vendorTestTmplName = "vendor.go.tmpl"

	// Documentation test
	documentationTestFile     = filepath.Join("lib", "linguist", "documentation.yml")
	documentationGold         = filepath.Join(testDir, "documentation.gold")
	documentationTestTmplPath = filepath.Join(assetsDir, "documentation.go.tmpl")
	documentationTestTmplName = "documentation.go.tmpl"

	// Types test
	typeGold         = filepath.Join(testDir, "type.gold")
	typeTestTmplPath = filepath.Join(assetsDir, "type.go.tmpl")
	typeTestTmplName = "type.go.tmpl"

	// Interpreters test
	interpreterGold         = filepath.Join(testDir, "interpreter.gold")
	interpreterTestTmplPath = filepath.Join(assetsDir, "interpreter.go.tmpl")
	interpreterTestTmplName = "interpreter.go.tmpl"

	// Filenames test
	filenameGold         = filepath.Join(testDir, "filename.gold")
	filenameTestTmplPath = filepath.Join(assetsDir, "filename.go.tmpl")
	filenameTestTmplName = "filename.go.tmpl"

	// Aliases test
	aliasGold         = filepath.Join(testDir, "alias.gold")
	aliasTestTmplPath = filepath.Join(assetsDir, "alias.go.tmpl")
	aliasTestTmplName = "alias.go.tmpl"

	// Frequencies test
	frequenciesGold         = filepath.Join(testDir, "frequencies.gold")
	frequenciesTestTmplPath = filepath.Join(assetsDir, "frequencies.go.tmpl")
	frequenciesTestTmplName = "frequencies.go.tmpl"

	// commit test
	commitGold         = filepath.Join(testDir, "commit.gold")
	commitTestTmplPath = filepath.Join(assetsDir, "commit.go.tmpl")
	commitTestTmplName = "commit.go.tmpl"

	// mime test
	mimeTypeGold         = filepath.Join(testDir, "mimeType.gold")
	mimeTypeTestTmplPath = filepath.Join(assetsDir, "mimeType.go.tmpl")
	mimeTypeTestTmplName = "mimeType.go.tmpl"

	// colors test
	colorsGold         = filepath.Join(testDir, "colors.gold")
	colorsTestTmplPath = filepath.Join(assetsDir, "colors.go.tmpl")
	colorsTestTmplName = "colors.go.tmpl"

	// colors test
	groupsGold         = filepath.Join(testDir, "groups.gold")
	groupsTestTmplPath = filepath.Join(assetsDir, "groups.go.tmpl")
	groupsTestTmplName = "groups.go.tmpl"
)

type GeneratorTestSuite struct {
	suite.Suite
	tmpLinguist string
	cloned      bool
	testCases   []testCase
}

type testCase struct {
	name        string
	fileToParse string
	samplesDir  string
	tmplPath    string
	tmplName    string
	commit      string
	generate    File
	wantOut     string
}

var updateGold = flag.Bool("update_gold", false, "Update golden test files")

func Test_GeneratorTestSuite(t *testing.T) {
	suite.Run(t, new(GeneratorTestSuite))
}

func (s *GeneratorTestSuite) maybeCloneLinguist() {
	var err error
	s.tmpLinguist = os.Getenv(linguistClonedEnvVar)
	s.cloned = s.tmpLinguist == ""
	if s.cloned {
		s.tmpLinguist, err = ioutil.TempDir("", "linguist-")
		assert.NoError(s.T(), err)
		cmd := exec.Command("git", "clone", linguistURL, s.tmpLinguist)
		err = cmd.Run()
		assert.NoError(s.T(), err)

		cwd, err := os.Getwd()
		assert.NoError(s.T(), err)

		err = os.Chdir(s.tmpLinguist)
		assert.NoError(s.T(), err)

		cmd = exec.Command("git", "checkout", commit)
		err = cmd.Run()
		assert.NoError(s.T(), err)

		err = os.Chdir(cwd)
		assert.NoError(s.T(), err)
	}
}

func (s *GeneratorTestSuite) SetupSuite() {
	s.maybeCloneLinguist()
	s.testCases = []testCase{
		{
			name:        "Extensions()",
			fileToParse: filepath.Join(s.tmpLinguist, languagesFile),
			samplesDir:  "",
			tmplPath:    extensionTestTmplPath,
			tmplName:    extensionTestTmplName,
			commit:      commit,
			generate:    Extensions,
			wantOut:     extensionGold,
		},
		{
			name:        "Heuristics()",
			fileToParse: filepath.Join(s.tmpLinguist, heuristicsTestFile),
			samplesDir:  "",
			tmplPath:    contentTestTmplPath,
			tmplName:    contentTestTmplName,
			commit:      commit,
			generate:    GenHeuristics,
			wantOut:     contentGold,
		},
		{
			name:        "Vendor()",
			fileToParse: filepath.Join(s.tmpLinguist, vendorTestFile),
			samplesDir:  "",
			tmplPath:    vendorTestTmplPath,
			tmplName:    vendorTestTmplName,
			commit:      commit,
			generate:    Vendor,
			wantOut:     vendorGold,
		},
		{
			name:        "Documentation()",
			fileToParse: filepath.Join(s.tmpLinguist, documentationTestFile),
			samplesDir:  "",
			tmplPath:    documentationTestTmplPath,
			tmplName:    documentationTestTmplName,
			commit:      commit,
			generate:    Documentation,
			wantOut:     documentationGold,
		},
		{
			name:        "Types()",
			fileToParse: filepath.Join(s.tmpLinguist, languagesFile),
			samplesDir:  "",
			tmplPath:    typeTestTmplPath,
			tmplName:    typeTestTmplName,
			commit:      commit,
			generate:    Types,
			wantOut:     typeGold,
		},
		{
			name:        "Interpreters()",
			fileToParse: filepath.Join(s.tmpLinguist, languagesFile),
			samplesDir:  "",
			tmplPath:    interpreterTestTmplPath,
			tmplName:    interpreterTestTmplName,
			commit:      commit,
			generate:    Interpreters,
			wantOut:     interpreterGold,
		},
		{
			name:        "Filenames()",
			fileToParse: filepath.Join(s.tmpLinguist, languagesFile),
			samplesDir:  filepath.Join(s.tmpLinguist, samplesDir),
			tmplPath:    filenameTestTmplPath,
			tmplName:    filenameTestTmplName,
			commit:      commit,
			generate:    Filenames,
			wantOut:     filenameGold,
		},
		{
			name:        "Aliases()",
			fileToParse: filepath.Join(s.tmpLinguist, languagesFile),
			samplesDir:  "",
			tmplPath:    aliasTestTmplPath,
			tmplName:    aliasTestTmplName,
			commit:      commit,
			generate:    Aliases,
			wantOut:     aliasGold,
		},
		{
			name:       "Frequencies()",
			samplesDir: filepath.Join(s.tmpLinguist, samplesDir),
			tmplPath:   frequenciesTestTmplPath,
			tmplName:   frequenciesTestTmplName,
			commit:     commit,
			generate:   Frequencies,
			wantOut:    frequenciesGold,
		},
		{
			name:       "Commit()",
			samplesDir: "",
			tmplPath:   commitTestTmplPath,
			tmplName:   commitTestTmplName,
			commit:     commit,
			generate:   Commit,
			wantOut:    commitGold,
		},
		{
			name:        "MimeType()",
			fileToParse: filepath.Join(s.tmpLinguist, languagesFile),
			samplesDir:  "",
			tmplPath:    mimeTypeTestTmplPath,
			tmplName:    mimeTypeTestTmplName,
			commit:      commit,
			generate:    MimeType,
			wantOut:     mimeTypeGold,
		},
		{
			name:        "Colors()",
			fileToParse: filepath.Join(s.tmpLinguist, languagesFile),
			samplesDir:  "",
			tmplPath:    colorsTestTmplPath,
			tmplName:    colorsTestTmplName,
			commit:      commit,
			generate:    Colors,
			wantOut:     colorsGold,
		},
		{
			name:        "Groups()",
			fileToParse: filepath.Join(s.tmpLinguist, languagesFile),
			samplesDir:  "",
			tmplPath:    groupsTestTmplPath,
			tmplName:    groupsTestTmplName,
			commit:      commit,
			generate:    Groups,
			wantOut:     groupsGold,
		},
	}
}

func (s *GeneratorTestSuite) TearDownSuite() {
	if s.cloned {
		err := os.RemoveAll(s.tmpLinguist)
		if err != nil {
			s.T().Logf("Failed to clean up %s after the test.\n", s.tmpLinguist)
		}
	}
}

// TestUpdateGeneratorTestSuiteGold is a Gold results generation automation.
// It should only be enabled&run manually on every new Linguist version
// to update *.gold files.
func (s *GeneratorTestSuite) TestUpdateGeneratorTestSuiteGold() {
	if !*updateGold {
		s.T().Skip()
	}
	s.T().Logf("Generating new *.gold test files")
	for _, test := range s.testCases {
		dst := test.wantOut
		s.T().Logf("Generating %s from %s\n", dst, test.fileToParse)
		err := test.generate(test.fileToParse, test.samplesDir, dst, test.tmplPath, test.tmplName, test.commit)
		assert.NoError(s.T(), err)
	}
}

func (s *GeneratorTestSuite) TestGenerationFiles() {
	for _, test := range s.testCases {
		gold, err := ioutil.ReadFile(test.wantOut)
		assert.NoError(s.T(), err)

		outPath, err := ioutil.TempFile("", "generator-test-")
		assert.NoError(s.T(), err)
		defer os.Remove(outPath.Name())
		err = test.generate(test.fileToParse, test.samplesDir, outPath.Name(), test.tmplPath, test.tmplName, test.commit)
		assert.NoError(s.T(), err)
		out, err := ioutil.ReadFile(outPath.Name())
		assert.NoError(s.T(), err)

		expected := normalizeSpaces(string(gold))
		actual := normalizeSpaces(string(out))
		assert.Equal(s.T(), expected, actual, "Test %s", test.name)
		if expected != actual {
			s.T().Logf("%s generated is different from %q", test.name, test.wantOut)
			s.T().Logf("Expected %q", expected[:400])
			s.T().Logf("Actual %q", actual[:400])
		}

	}
}

func (s *GeneratorTestSuite) TestTokenizerOnATS() {
	const suspiciousSample = "samples/ATS/csv_parse.hats"
	sFile := filepath.Join(s.tmpLinguist, suspiciousSample)
	content, err := ioutil.ReadFile(sFile)
	require.NoError(s.T(), err)

	tokens := tokenizer.Tokenize(content)
	assert.Equal(s.T(), 381, len(tokens), "Number of tokens using LF as line endings")
}

// normalizeSpaces returns a copy of str with whitespaces normalized.
// We use this to compare generated source as gofmt format may change.
// E.g for changes between Go 1.10 and 1.11 see
// https://go-review.googlesource.com/c/go/+/122295/
func normalizeSpaces(str string) string {
	return strings.Join(strings.Fields(str), " ")
}
