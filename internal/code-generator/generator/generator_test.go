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
	typesTestFile     = "test_files/type.test.yml"
	typesGold         = "test_files/type.gold"
	typesTestTmplPath = "../assets/type.go.tmpl"
	typesTestTmplName = "type.go.tmpl"

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
