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
	lingustURL    = "https://github.com/github/linguist.git"
	commit        = "b6460f8ed6b249281ada099ca28bd8f1230b8892"
	samplesDir    = "samples"
	languagesFile = "lib/linguist/languages.yml"

	// Extensions test
	extensionGold         = "test_files/extension.gold"
	extensionTestTmplPath = "../assets/extension.go.tmpl"
	extensionTestTmplName = "extension.go.tmpl"

	// Heuristics test
	heuristicsTestFile  = "lib/linguist/heuristics.rb"
	contentGold         = "test_files/content.gold"
	contentTestTmplPath = "../assets/content.go.tmpl"
	contentTestTmplName = "content.go.tmpl"

	// Vendor test
	vendorTestFile     = "lib/linguist/vendor.yml"
	vendorGold         = "test_files/vendor.gold"
	vendorTestTmplPath = "../assets/vendor.go.tmpl"
	vendorTestTmplName = "vendor.go.tmpl"

	// Documentation test
	documentationTestFile     = "lib/linguist/documentation.yml"
	documentationGold         = "test_files/documentation.gold"
	documentationTestTmplPath = "../assets/documentation.go.tmpl"
	documentationTestTmplName = "documentation.go.tmpl"

	// Types test
	typeGold         = "test_files/type.gold"
	typeTestTmplPath = "../assets/type.go.tmpl"
	typeTestTmplName = "type.go.tmpl"

	// Interpreters test
	interpreterGold         = "test_files/interpreter.gold"
	interpreterTestTmplPath = "../assets/interpreter.go.tmpl"
	interpreterTestTmplName = "interpreter.go.tmpl"

	// Filenames test
	filenameGold         = "test_files/filename.gold"
	filenameTestTmplPath = "../assets/filename.go.tmpl"
	filenameTestTmplName = "filename.go.tmpl"

	// Aliases test
	aliasGold         = "test_files/alias.gold"
	aliasTestTmplPath = "../assets/alias.go.tmpl"
	aliasTestTmplName = "alias.go.tmpl"

	// Frequencies test
	frequenciesGold         = "test_files/frequencies.gold"
	frequenciesTestTmplPath = "../assets/frequencies.go.tmpl"
	frequenciesTestTmplName = "frequencies.go.tmpl"
)

type GeneratorTestSuite struct {
	suite.Suite
	tmpLinguist string
}

func TestGeneratorTestSuite(t *testing.T) {
	suite.Run(t, new(GeneratorTestSuite))
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

	cmd = exec.Command("git", "checkout", commit)
	err = cmd.Run()
	assert.NoError(g.T(), err)

	err = os.Chdir(cwd)
	assert.NoError(g.T(), err)
}

func (g *GeneratorTestSuite) TearDownSuite() {
	err := os.RemoveAll(g.tmpLinguist)
	assert.NoError(g.T(), err)
}

func (g *GeneratorTestSuite) TestGenerationFiles() {
	tests := []struct {
		name        string
		fileToParse string
		samplesDir  string
		tmplPath    string
		tmplName    string
		commit      string
		generate    File
		wantOut     string
	}{
		{
			name:        "Extensions()",
			fileToParse: filepath.Join(g.tmpLinguist, languagesFile),
			samplesDir:  "",
			tmplPath:    extensionTestTmplPath,
			tmplName:    extensionTestTmplName,
			commit:      commit,
			generate:    Extensions,
			wantOut:     extensionGold,
		},
		{
			name:        "Heuristics()",
			fileToParse: filepath.Join(g.tmpLinguist, heuristicsTestFile),
			samplesDir:  "",
			tmplPath:    contentTestTmplPath,
			tmplName:    contentTestTmplName,
			commit:      commit,
			generate:    Heuristics,
			wantOut:     contentGold,
		},
		{
			name:        "Vendor()",
			fileToParse: filepath.Join(g.tmpLinguist, vendorTestFile),
			samplesDir:  "",
			tmplPath:    vendorTestTmplPath,
			tmplName:    vendorTestTmplName,
			commit:      commit,
			generate:    Vendor,
			wantOut:     vendorGold,
		},
		{
			name:        "Documentation()",
			fileToParse: filepath.Join(g.tmpLinguist, documentationTestFile),
			samplesDir:  "",
			tmplPath:    documentationTestTmplPath,
			tmplName:    documentationTestTmplName,
			commit:      commit,
			generate:    Documentation,
			wantOut:     documentationGold,
		},
		{
			name:        "Types()",
			fileToParse: filepath.Join(g.tmpLinguist, languagesFile),
			samplesDir:  "",
			tmplPath:    typeTestTmplPath,
			tmplName:    typeTestTmplName,
			commit:      commit,
			generate:    Types,
			wantOut:     typeGold,
		},
		{
			name:        "Interpreters()",
			fileToParse: filepath.Join(g.tmpLinguist, languagesFile),
			samplesDir:  "",
			tmplPath:    interpreterTestTmplPath,
			tmplName:    interpreterTestTmplName,
			commit:      commit,
			generate:    Interpreters,
			wantOut:     interpreterGold,
		},
		{
			name:        "Filenames()",
			fileToParse: filepath.Join(g.tmpLinguist, languagesFile),
			samplesDir:  filepath.Join(g.tmpLinguist, samplesDir),
			tmplPath:    filenameTestTmplPath,
			tmplName:    filenameTestTmplName,
			commit:      commit,
			generate:    Filenames,
			wantOut:     filenameGold,
		},
		{
			name:        "Aliases()",
			fileToParse: filepath.Join(g.tmpLinguist, languagesFile),
			samplesDir:  "",
			tmplPath:    aliasTestTmplPath,
			tmplName:    aliasTestTmplName,
			commit:      commit,
			generate:    Aliases,
			wantOut:     aliasGold,
		},
		{
			name:       "Frequencies()",
			samplesDir: filepath.Join(g.tmpLinguist, samplesDir),
			tmplPath:   frequenciesTestTmplPath,
			tmplName:   frequenciesTestTmplName,
			commit:     commit,
			generate:   Frequencies,
			wantOut:    frequenciesGold,
		},
	}

	for _, test := range tests {
		gold, err := ioutil.ReadFile(test.wantOut)
		assert.NoError(g.T(), err)

		outPath, err := ioutil.TempFile("/tmp", "generator-test-")
		assert.NoError(g.T(), err)
		defer os.Remove(outPath.Name())

		err = test.generate(test.fileToParse, test.samplesDir, outPath.Name(), test.tmplPath, test.tmplName, test.commit)
		assert.NoError(g.T(), err)
		out, err := ioutil.ReadFile(outPath.Name())
		assert.NoError(g.T(), err)
		assert.EqualValues(g.T(), gold, out, fmt.Sprintf("%v: %v, expected: %v", test.name, string(out), string(gold)))
	}
}
