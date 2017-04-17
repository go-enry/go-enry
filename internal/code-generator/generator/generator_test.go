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

	// Languages test
	ymlTestFile           = "test_files/languages.test.yml"
	langGold              = "test_files/languages.gold"
	languagesTestTmplPath = "test_files/languages.test.go.tmpl"
	languagesTestTmplName = "languages.test.go.tmpl"

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

	// Types test
	typesTestFile     = "test_files/type.test.yml"
	typesGold         = "test_files/type.gold"
	typesTestTmplPath = "test_files/type.test.go.tmpl"
	typesTestTmplName = "type.test.go.tmpl"

	// Interpreters test
	interpretersTestFile     = "test_files/interpreters.test.yml"
	interpretersGold         = "test_files/interpreters.gold"
	interpretersTestTmplPath = "test_files/interpreters.test.go.tmpl"
	interpretersTestTmplName = "interpreters.test.go.tmpl"
)

func TestFromFile(t *testing.T) {
	goldLang, err := ioutil.ReadFile(langGold)
	assert.NoError(t, err)

	goldContent, err := ioutil.ReadFile(contentGold)
	assert.NoError(t, err)

	goldVendor, err := ioutil.ReadFile(vendorGold)
	assert.NoError(t, err)

	goldDocumentation, err := ioutil.ReadFile(documentationGold)
	assert.NoError(t, err)

	goldTypes, err := ioutil.ReadFile(typesGold)
	assert.NoError(t, err)

	goldInterpreters, err := ioutil.ReadFile(interpretersGold)
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

	outPathTypes, err := ioutil.TempFile("/tmp", "generator-test-")
	assert.NoError(t, err)
	defer os.Remove(outPathTypes.Name())

	outPathInterpreters, err := ioutil.TempFile("/tmp", "generator-test-")
	assert.NoError(t, err)
	defer os.Remove(outPathInterpreters.Name())

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
		{
			name:        "TestFromFile_Types",
			fileToParse: typesTestFile,
			outPath:     outPathTypes.Name(),
			tmplPath:    typesTestTmplPath,
			tmplName:    typesTestTmplName,
			commit:      commitTest,
			generate:    Types,
			wantOut:     goldTypes,
		},
		{
			name:        "TestFromFile_Interpreters",
			fileToParse: interpretersTestFile,
			outPath:     outPathInterpreters.Name(),
			tmplPath:    interpretersTestTmplPath,
			tmplName:    interpretersTestTmplName,
			commit:      commitTest,
			generate:    Interpreters,
			wantOut:     goldInterpreters,
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
