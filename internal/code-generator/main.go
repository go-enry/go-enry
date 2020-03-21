package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/go-enry/go-enry/v2/internal/code-generator/generator"
)

var (
	// directories
	samplesDir = filepath.Join(".linguist", "samples")
	libDir     = filepath.Join(".linguist", "lib", "linguist")
	assetsDir  = filepath.Join("internal", "code-generator", "assets")

	// languages info file
	languagesYAML = filepath.Join(libDir, "languages.yml")

	// extension.go generation
	extensionsFile     = filepath.Join("data", "extension.go")
	extensionsTmplPath = filepath.Join(assetsDir, "extension.go.tmpl")
	extensionsTmpl     = "extension.go.tmpl"

	// content.go generation
	heuristicsYAML  = filepath.Join(libDir, "heuristics.yml")
	contentFile     = filepath.Join("data", "content.go")
	contentTmplPath = filepath.Join(assetsDir, "content.go.tmpl")
	contentTmpl     = "content.go.tmpl"

	// vendor.go generation
	vendorYAML     = filepath.Join(libDir, "vendor.yml")
	vendorFile     = filepath.Join("data", "vendor.go")
	vendorTmplPath = filepath.Join(assetsDir, "vendor.go.tmpl")
	vendorTmpl     = "vendor.go.tmpl"

	// documentation.go generation
	documentationYAML     = filepath.Join(libDir, "documentation.yml")
	documentationFile     = filepath.Join("data", "documentation.go")
	documentationTmplPath = filepath.Join(assetsDir, "documentation.go.tmpl")
	documentationTmpl     = "documentation.go.tmpl"

	// type.go generation
	typeFile     = filepath.Join("data", "type.go")
	typeTmplPath = filepath.Join(assetsDir, "type.go.tmpl")
	typeTmpl     = "type.go.tmpl"

	// interpreter.go generation
	interpretersFile     = filepath.Join("data", "interpreter.go")
	interpretersTmplPath = filepath.Join(assetsDir, "interpreter.go.tmpl")
	interpretersTmpl     = "interpreter.go.tmpl"

	// filename.go generation
	filenamesFile     = filepath.Join("data", "filename.go")
	filenamesTmplPath = filepath.Join(assetsDir, "filename.go.tmpl")
	filenamesTmpl     = "filename.go.tmpl"

	// alias.go generation
	aliasesFile     = filepath.Join("data", "alias.go")
	aliasesTmplPath = filepath.Join(assetsDir, "alias.go.tmpl")
	aliasesTmpl     = "alias.go.tmpl"

	// frequencies.go generation
	frequenciesFile     = filepath.Join("data", "frequencies.go")
	frequenciesTmplPath = filepath.Join(assetsDir, "frequencies.go.tmpl")
	frequenciesTmpl     = "frequencies.go.tmpl"

	// commit.go generation
	commitFile     = filepath.Join("data", "commit.go")
	commitTmplPath = filepath.Join(assetsDir, "commit.go.tmpl")
	commitTmpl     = "commit.go.tmpl"

	// mimeType.go generation
	mimeTypeFile     = filepath.Join("data", "mimeType.go")
	mimeTypeTmplPath = filepath.Join(assetsDir, "mimeType.go.tmpl")
	mimeTypeTmpl     = "mimeType.go.tmpl"

	// colors.go generation
	colorsFile     = filepath.Join("data", "colors.go")
	colorsTmplPath = filepath.Join(assetsDir, "colors.go.tmpl")
	colorsTmpl     = "colors.go.tmpl"

	// groups.go generation
	groupsFile     = filepath.Join("data", "groups.go")
	groupsTmplPath = filepath.Join(assetsDir, "groups.go.tmpl")
	groupsTmpl     = "groups.go.tmpl"

	commitPath = filepath.Join(".linguist", ".git", "HEAD")
)

type generatorFiles struct {
	generate    generator.File
	fileToParse string
	samplesDir  string
	outPath     string
	tmplPath    string
	tmplName    string
	commit      string
}

func main() {
	commit, err := getCommit(commitPath)
	if err != nil {
		log.Printf("couldn't find commit: %v", err)
	}

	fileList := []*generatorFiles{
		{generator.Extensions, languagesYAML, "", extensionsFile, extensionsTmplPath, extensionsTmpl, commit},
		{generator.GenHeuristics, heuristicsYAML, "", contentFile, contentTmplPath, contentTmpl, commit},
		{generator.Vendor, vendorYAML, "", vendorFile, vendorTmplPath, vendorTmpl, commit},
		{generator.Documentation, documentationYAML, "", documentationFile, documentationTmplPath, documentationTmpl, commit},
		{generator.Types, languagesYAML, "", typeFile, typeTmplPath, typeTmpl, commit},
		{generator.Interpreters, languagesYAML, "", interpretersFile, interpretersTmplPath, interpretersTmpl, commit},
		{generator.Filenames, languagesYAML, samplesDir, filenamesFile, filenamesTmplPath, filenamesTmpl, commit},
		{generator.Aliases, languagesYAML, "", aliasesFile, aliasesTmplPath, aliasesTmpl, commit},
		{generator.Frequencies, "", samplesDir, frequenciesFile, frequenciesTmplPath, frequenciesTmpl, commit},
		{generator.Commit, "", "", commitFile, commitTmplPath, commitTmpl, commit},
		{generator.MimeType, languagesYAML, "", mimeTypeFile, mimeTypeTmplPath, mimeTypeTmpl, commit},
		{generator.Colors, languagesYAML, "", colorsFile, colorsTmplPath, colorsTmpl, commit},
		{generator.Groups, languagesYAML, "", groupsFile, groupsTmplPath, groupsTmpl, commit},
	}

	for _, file := range fileList {
		if err := file.generate(file.fileToParse, file.samplesDir, file.outPath, file.tmplPath, file.tmplName, file.commit); err != nil {
			log.Println(err)
		}
	}
}

func getCommit(path string) (string, error) {
	commit, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	if string(commit) == "ref: refs/heads/master\n" {
		path = filepath.Join(".linguist", ".git", string(commit[5:len(commit)-1]))
		commit, err = ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}
	}

	return string(commit[:len(commit)-1]), nil
}
