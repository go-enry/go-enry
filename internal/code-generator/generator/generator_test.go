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

	// Filenames test
	filenamesTestFile     = "test_files/filenames.test.yml"
	filenamesGold         = "test_files/filenames.gold"
	filenamesTestTmplPath = "test_files/filenames.test.go.tmpl"
	filenamesTestTmplName = "filenames.test.go.tmpl"
)

func TestFromFile(t *testing.T) {
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
			fileToParse: ymlTestFile,
			tmplPath:    languagesTestTmplPath,
			tmplName:    languagesTestTmplName,
			commit:      commitTest,
			generate:    Languages,
			wantOut:     langGold,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gold, err := ioutil.ReadFile(tt.wantOut)
			assert.NoError(t, err)

			outPath, err := ioutil.TempFile("/tmp", "generator-test-")
			assert.NoError(t, err)
			defer os.Remove(outPath.Name())

			err = FromFile(tt.fileToParse, outPath.Name(), tt.tmplPath, tt.tmplName, tt.commit, tt.generate)
			assert.NoError(t, err)
			out, err := ioutil.ReadFile(outPath.Name())
			assert.NoError(t, err)
			assert.EqualValues(t, gold, out, fmt.Sprintf("FromFile() = %v, want %v", string(out), string(tt.wantOut)))
		})
	}
}
