package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	lingustURL = "https://github.com/github/linguist.git"
	commitTree = "60f864a138650dd17fafc94814be9ee2d3aaef8c"
	commitTest = "fe8b44ab8a225b1ffa75b983b916ea22fee5b6f7"

	// Extensions test
	extensionsTestFile     = "test_files/extensions.test.yml"
	extensionsGold         = "test_files/extensions.gold"
	extensionsTestTmplPath = "../assets/extensions.go.tmpl"
	extensionsTestTmplName = "extensions.go.tmpl"

	// Heuristics test
	heuristicsTestFile  = "test_files/heuristics.test.rb"
	contentGold         = "test_files/content.gold"
	contentTestTmplPath = "../assets/content.go.tmpl"
	contentTestTmplName = "content.go.tmpl"

	// Vendor test
	vendorTestFile     = "test_files/vendor.test.yml"
	vendorGold         = "test_files/vendor.gold"
	vendorTestTmplPath = "../assets/vendor.go.tmpl"
	vendorTestTmplName = "vendor.go.tmpl"

	// Documentation test
	documentationTestFile     = "test_files/documentation.test.yml"
	documentationGold         = "test_files/documentation.gold"
	documentationTestTmplPath = "../assets/documentation.go.tmpl"
	documentationTestTmplName = "documentation.go.tmpl"

	// Types test
	typesTestFile     = "test_files/types.test.yml"
	typesGold         = "test_files/types.gold"
	typesTestTmplPath = "../assets/types.go.tmpl"
	typesTestTmplName = "types.go.tmpl"

	// Interpreters test
	interpretersTestFile     = "test_files/interpreters.test.yml"
	interpretersGold         = "test_files/interpreters.gold"
	interpretersTestTmplPath = "../assets/interpreters.go.tmpl"
	interpretersTestTmplName = "interpreters.go.tmpl"

	// Filenames test
	filenamesTestFile     = "test_files/filenames.test.yml"
	filenamesGold         = "test_files/filenames.gold"
	filenamesTestTmplPath = "../assets/filenames.go.tmpl"
	filenamesTestTmplName = "filenames.go.tmpl"

	// Aliases test
	aliasesTestFile     = "test_files/aliases.test.yml"
	aliasesGold         = "test_files/aliases.gold"
	aliasesTestTmplPath = "../assets/aliases.go.tmpl"
	aliasesTestTmplName = "aliases.go.tmpl"

	// Frequencies test
	frequenciesTestDir      = "/samples"
	frequenciesGold         = "test_files/frequencies.gold"
	frequenciesTestTmplPath = "../assets/frequencies.go.tmpl"
	frequenciesTestTmplName = "frequencies.go.tmpl"
)

type GeneratorTestSuite struct {
	suite.Suite
	tmpLinguist string
}

func (g *GeneratorTestSuite) SetupSuite() {
	tmpLinguist, err := ioutil.TempDir("", "linguist-")
	assert.NoError(g.T(), err)
	g.tmpLinguist = tmpLinguist

	cmd := exec.Command("git", "clone", lingustURL, tmpLinguist)
	err = cmd.Run()
	assert.NoError(g.T(), err)

	cwd, err := os.Getwd()
	assert.NoError(g.T(), err)

	err = os.Chdir(tmpLinguist)
	assert.NoError(g.T(), err)

	cmd = exec.Command("git", "checkout", commitTree)
	err = cmd.Run()
	assert.NoError(g.T(), err)

	err = os.Chdir(cwd)
	assert.NoError(g.T(), err)
}

func (g *GeneratorTestSuite) TearDownSuite() {
	err := os.RemoveAll(g.tmpLinguist)
	assert.NoError(g.T(), err)
}

func (g *GeneratorTestSuite) TestFromFile() {
	tests := []struct {
		name        string
		fileToParse string
		tmplPath    string
		tmplName    string
		commit      string
		generate    Func
		wantOut     string
	}{
		{
			name:        "TestFromFile_Language",
			fileToParse: extensionsTestFile,
			tmplPath:    extensionsTestTmplPath,
			tmplName:    extensionsTestTmplName,
			commit:      commitTest,
			generate:    Extensions,
			wantOut:     extensionsGold,
		},
		{
			name:        "TestFromFile_Heuristics",
			fileToParse: heuristicsTestFile,
			tmplPath:    contentTestTmplPath,
			tmplName:    contentTestTmplName,
			commit:      commitTest,
			generate:    Heuristics,
			wantOut:     contentGold,
		},
		{
			name:        "TestFromFile_Vendor",
			fileToParse: vendorTestFile,
			tmplPath:    vendorTestTmplPath,
			tmplName:    vendorTestTmplName,
			commit:      commitTest,
			generate:    Vendor,
			wantOut:     vendorGold,
		},
		{
			name:        "TestFromFile_Documentation",
			fileToParse: documentationTestFile,
			tmplPath:    documentationTestTmplPath,
			tmplName:    documentationTestTmplName,
			commit:      commitTest,
			generate:    Documentation,
			wantOut:     documentationGold,
		},
		{
			name:        "TestFromFile_Types",
			fileToParse: typesTestFile,
			tmplPath:    typesTestTmplPath,
			tmplName:    typesTestTmplName,
			commit:      commitTest,
			generate:    Types,
			wantOut:     typesGold,
		},
		{
			name:        "TestFromFile_Interpreters",
			fileToParse: interpretersTestFile,
			tmplPath:    interpretersTestTmplPath,
			tmplName:    interpretersTestTmplName,
			commit:      commitTest,
			generate:    Interpreters,
			wantOut:     interpretersGold,
		},
		{
			name:        "TestFromFile_Filenames",
			fileToParse: filenamesTestFile,
			tmplPath:    filenamesTestTmplPath,
			tmplName:    filenamesTestTmplName,
			commit:      commitTest,
			generate:    Filenames,
			wantOut:     filenamesGold,
		},
		{
			name:        "TestFromFile_Aliases",
			fileToParse: aliasesTestFile,
			tmplPath:    aliasesTestTmplPath,
			tmplName:    aliasesTestTmplName,
			commit:      commitTest,
			generate:    Aliases,
			wantOut:     aliasesGold,
		},
	}

	for _, test := range tests {
		gold, err := ioutil.ReadFile(test.wantOut)
		assert.NoError(g.T(), err)

		outPath, err := ioutil.TempFile("/tmp", "generator-test-")
		assert.NoError(g.T(), err)
		defer os.Remove(outPath.Name())

		err = FromFile(test.fileToParse, outPath.Name(), test.tmplPath, test.tmplName, test.commit, test.generate)
		assert.NoError(g.T(), err)
		out, err := ioutil.ReadFile(outPath.Name())
		assert.NoError(g.T(), err)
		assert.EqualValues(g.T(), gold, out, fmt.Sprintf("FromFile() = %v, want %v", string(out), string(test.wantOut)))
	}
}

func (g *GeneratorTestSuite) TestFrequencies() {
	tests := []struct {
		name       string
		samplesDir string
		tmplPath   string
		tmplName   string
		commit     string
		wantOut    string
	}{
		{
			name:       "Frequencies_1",
			samplesDir: filepath.Join(g.tmpLinguist, frequenciesTestDir),
			tmplPath:   frequenciesTestTmplPath,
			tmplName:   frequenciesTestTmplName,
			commit:     commitTree,
			wantOut:    frequenciesGold,
		},
	}

	for _, test := range tests {
		gold, err := ioutil.ReadFile(test.wantOut)
		assert.NoError(g.T(), err)

		outPath, err := ioutil.TempFile("/tmp", "frequencies-test-")
		assert.NoError(g.T(), err)
		defer os.Remove(outPath.Name())

		err = Frequencies(test.samplesDir, test.tmplPath, test.tmplName, test.commit, outPath.Name())
		assert.NoError(g.T(), err)
		out, err := ioutil.ReadFile(outPath.Name())
		assert.NoError(g.T(), err)
		assert.EqualValues(g.T(), gold, out, fmt.Sprintf("Frequencies() = %v, want %v", string(out), string(test.wantOut)))
	}
}

func TestGeneratorTestSuite(t *testing.T) {
	suite.Run(t, new(GeneratorTestSuite))
}
