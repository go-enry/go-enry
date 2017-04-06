package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	// FromFile test
	formatedLangGold    = "test_files/formated_languages.gold"
	formatedContentGold = "test_files/formated_content.gold"
	formatedUtilsGold   = "test_files/formated_utils.gold"

	// Languages test
	ymlTestFile           = "test_files/languages.test.yml"
	langGold              = "test_files/languages.gold"
	languagesTestTmplPath = "test_files/languages.test.tmpl"
	languagesTestTmplName = "languages.test.tmpl"
	commitLangTest        = "fe8b44ab8a225b1ffa75b983b916ea22fee5b6f7"

	// Heuristics test
	heuristicsTestFile   = "test_files/heuristics.test.rb"
	contentGold          = "test_files/content.gold"
	contentTestTmplPath  = "test_files/content.test.go.tmpl"
	contentTestTmplName  = "content.test.go.tmpl"
	commitHeuristicsTest = "fe8b44ab8a225b1ffa75b983b916ea22fee5b6f7"

	// Vendor test
	vendorTestFile    = "test_files/vendor.test.yml"
	utilsGold         = "test_files/utils.gold"
	utilsTestTmplPath = "test_files/utils.test.go.tmpl"
	utilsTestTmplName = "utils.test.go.tmpl"
	commitVendorTest  = "fe8b44ab8a225b1ffa75b983b916ea22fee5b6f7"
)

func TestFromFile(t *testing.T) {
	goldLang, err := ioutil.ReadFile(formatedLangGold)
	assert.NoError(t, err)

	goldContent, err := ioutil.ReadFile(formatedContentGold)
	assert.NoError(t, err)

	goldUtils, err := ioutil.ReadFile(formatedUtilsGold)
	assert.NoError(t, err)

	outPathLang, err := ioutil.TempFile("/tmp", "generator-test-")
	assert.NoError(t, err)
	defer os.Remove(outPathLang.Name())

	outPathContent, err := ioutil.TempFile("/tmp", "generator-test-")
	assert.NoError(t, err)
	defer os.Remove(outPathContent.Name())

	outPathUtils, err := ioutil.TempFile("/tmp", "generator-test-")
	assert.NoError(t, err)
	defer os.Remove(outPathContent.Name())

	tests := []struct {
		name        string
		fileToParse string
		outPath     string
		tmplPath    string
		tmplName    string
		commit      string
		generate    Func
		wantOut     []byte
	}{
		{
			name:        "TestFromFile_Language",
			fileToParse: ymlTestFile,
			outPath:     outPathLang.Name(),
			tmplPath:    languagesTestTmplPath,
			tmplName:    languagesTestTmplName,
			commit:      commitLangTest,
			generate:    Languages,
			wantOut:     goldLang,
		},
		{
			name:        "TestFromFile_Heuristics",
			fileToParse: heuristicsTestFile,
			outPath:     outPathContent.Name(),
			tmplPath:    contentTestTmplPath,
			tmplName:    contentTestTmplName,
			commit:      commitHeuristicsTest,
			generate:    Heuristics,
			wantOut:     goldContent,
		},
		{
			name:        "TestFromFile_Vendor",
			fileToParse: vendorTestFile,
			outPath:     outPathUtils.Name(),
			tmplPath:    utilsTestTmplPath,
			tmplName:    utilsTestTmplName,
			commit:      commitVendorTest,
			generate:    Vendor,
			wantOut:     goldUtils,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FromFile(tt.fileToParse, tt.outPath, tt.tmplPath, tt.tmplName, tt.commit, tt.generate)
			assert.NoError(t, err)
			out, err := ioutil.ReadFile(tt.outPath)
			assert.NoError(t, err)
			assert.EqualValues(t, tt.wantOut, out, fmt.Sprintf("FromFile() = %v, want %v", string(out), string(tt.wantOut)))
		})
	}
}

func TestLanguages(t *testing.T) {
	gold, err := ioutil.ReadFile(langGold)
	assert.NoError(t, err)

	input, err := ioutil.ReadFile(ymlTestFile)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		input    []byte
		tmplPath string
		tmplName string
		commit   string
		wantOut  []byte
	}{
		{
			name:     "TestLanguages",
			input:    input,
			tmplPath: languagesTestTmplPath,
			tmplName: languagesTestTmplName,
			commit:   commitLangTest,
			wantOut:  gold,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := Languages(tt.input, tt.tmplPath, tt.tmplName, tt.commit)
			assert.NoError(t, err)
			assert.EqualValues(t, tt.wantOut, out, fmt.Sprintf("Languages() = %v, want %v", string(out), string(tt.wantOut)))
		})
	}
}

func TestHeuristics(t *testing.T) {
	gold, err := ioutil.ReadFile(contentGold)
	assert.NoError(t, err)

	input, err := ioutil.ReadFile(heuristicsTestFile)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		input    []byte
		tmplPath string
		tmplName string
		commit   string
		wantOut  []byte
	}{
		{
			name:     "TestHeuristics",
			input:    input,
			tmplPath: contentTestTmplPath,
			tmplName: contentTestTmplName,
			commit:   commitHeuristicsTest,
			wantOut:  gold,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := Heuristics(tt.input, tt.tmplPath, tt.tmplName, tt.commit)
			assert.NoError(t, err)
			assert.EqualValues(t, tt.wantOut, out, fmt.Sprintf("Heuristics() = %v, want %v", string(out), string(tt.wantOut)))
		})
	}
}

func TestVendor(t *testing.T) {
	gold, err := ioutil.ReadFile(utilsGold)
	assert.NoError(t, err)

	input, err := ioutil.ReadFile(vendorTestFile)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		input    []byte
		tmplPath string
		tmplName string
		commit   string
		wantOut  []byte
	}{
		{
			name:     "TestVendor",
			input:    input,
			tmplPath: utilsTestTmplPath,
			tmplName: utilsTestTmplName,
			commit:   commitVendorTest,
			wantOut:  gold,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := Vendor(tt.input, tt.tmplPath, tt.tmplName, tt.commit)
			assert.NoError(t, err)
			assert.EqualValues(t, tt.wantOut, out, fmt.Sprintf("Vendor() = %v, want %v", string(out), string(tt.wantOut)))
		})
	}
}
