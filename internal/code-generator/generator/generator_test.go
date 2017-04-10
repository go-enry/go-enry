package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	commitTest = "fe8b44ab8a225b1ffa75b983b916ea22fee5b6f7"

	// FromFile test
	formatedLangGold          = "test_files/formated_languages.gold"
	formatedContentGold       = "test_files/formated_content.gold"
	formatedVendorGold        = "test_files/formated_vendor.gold"
	formatedDocumentationGold = "test_files/formated_documentation.gold"

	// Languages test
	ymlTestFile           = "test_files/languages.test.yml"
	langGold              = "test_files/languages.gold"
	languagesTestTmplPath = "test_files/languages.test.tmpl"
	languagesTestTmplName = "languages.test.tmpl"

	// Heuristics test
	heuristicsTestFile  = "test_files/heuristics.test.rb"
	contentGold         = "test_files/content.gold"
	contentTestTmplPath = "test_files/content.test.go.tmpl"
	contentTestTmplName = "content.test.go.tmpl"

	// Vendor test
	vendorTestFile     = "test_files/vendor.test.yml"
	vendorGold         = "test_files/vendor.gold"
	vendorTestTmplPath = "test_files/vendor.test.go.tmpl"
	vendorTestTmplName = "vendor.test.go.tmpl"

	// Documentation test
	documentationTestFile     = "test_files/documentation.test.yml"
	documentationGold         = "test_files/documentation.gold"
	documentationTestTmplPath = "test_files/documentation.test.go.tmpl"
	documentationTestTmplName = "documentation.test.go.tmpl"
)

func TestFromFile(t *testing.T) {
	goldLang, err := ioutil.ReadFile(formatedLangGold)
	assert.NoError(t, err)

	goldContent, err := ioutil.ReadFile(formatedContentGold)
	assert.NoError(t, err)

	goldVendor, err := ioutil.ReadFile(formatedVendorGold)
	assert.NoError(t, err)

	goldDocumentation, err := ioutil.ReadFile(formatedDocumentationGold)
	assert.NoError(t, err)

	outPathLang, err := ioutil.TempFile("/tmp", "generator-test-")
	assert.NoError(t, err)
	defer os.Remove(outPathLang.Name())

	outPathContent, err := ioutil.TempFile("/tmp", "generator-test-")
	assert.NoError(t, err)
	defer os.Remove(outPathContent.Name())

	outPathVendor, err := ioutil.TempFile("/tmp", "generator-test-")
	assert.NoError(t, err)
	defer os.Remove(outPathVendor.Name())

	outPathDocumentation, err := ioutil.TempFile("/tmp", "generator-test-")
	assert.NoError(t, err)
	defer os.Remove(outPathDocumentation.Name())

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
			commit:      commitTest,
			generate:    Languages,
			wantOut:     goldLang,
		},
		{
			name:        "TestFromFile_Heuristics",
			fileToParse: heuristicsTestFile,
			outPath:     outPathContent.Name(),
			tmplPath:    contentTestTmplPath,
			tmplName:    contentTestTmplName,
			commit:      commitTest,
			generate:    Heuristics,
			wantOut:     goldContent,
		},
		{
			name:        "TestFromFile_Vendor",
			fileToParse: vendorTestFile,
			outPath:     outPathVendor.Name(),
			tmplPath:    vendorTestTmplPath,
			tmplName:    vendorTestTmplName,
			commit:      commitTest,
			generate:    Vendor,
			wantOut:     goldVendor,
		},
		{
			name:        "TestFromFile_Documentation",
			fileToParse: documentationTestFile,
			outPath:     outPathDocumentation.Name(),
			tmplPath:    documentationTestTmplPath,
			tmplName:    documentationTestTmplName,
			commit:      commitTest,
			generate:    Documentation,
			wantOut:     goldDocumentation,
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
			commit:   commitTest,
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
			commit:   commitTest,
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
	gold, err := ioutil.ReadFile(vendorGold)
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
			tmplPath: vendorTestTmplPath,
			tmplName: vendorTestTmplName,
			commit:   commitTest,
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

func TestDocumentation(t *testing.T) {
	gold, err := ioutil.ReadFile(documentationGold)
	assert.NoError(t, err)

	input, err := ioutil.ReadFile(documentationTestFile)
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
			name:     "TestDocumentation",
			input:    input,
			tmplPath: documentationTestTmplPath,
			tmplName: documentationTestTmplName,
			commit:   commitTest,
			wantOut:  gold,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := Documentation(tt.input, tt.tmplPath, tt.tmplName, tt.commit)
			assert.NoError(t, err)
			assert.EqualValues(t, tt.wantOut, out, fmt.Sprintf("Documentation() = %v, want %v", string(out), string(tt.wantOut)))
		})
	}
}
