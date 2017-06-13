package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/src-d/enry.v1/internal/code-generator/generator"
)

const (
	// languages info file
	languagesYAML = ".linguist/lib/linguist/languages.yml"

	// linguist's samples directory
	samplesDir = ".linguist/samples"

	// extension.go generation
	extensionsFile     = "extension.go"
	extensionsTmplPath = "internal/code-generator/assets/extensions.go.tmpl"
	extensionsTmpl     = "extensions.go.tmpl"

	// content.go generation
	heuristicsRuby  = ".linguist/lib/linguist/heuristics.rb"
	contentFile     = "content.go"
	contentTmplPath = "internal/code-generator/assets/content.go.tmpl"
	contentTmpl     = "content.go.tmpl"

	// vendor.go generation
	vendorYAML     = ".linguist/lib/linguist/vendor.yml"
	vendorFile     = "vendor.go"
	vendorTmplPath = "internal/code-generator/assets/vendor.go.tmpl"
	vendorTmpl     = "vendor.go.tmpl"

	// documentation.go generation
	documentationYAML     = ".linguist/lib/linguist/documentation.yml"
	documentationFile     = "documentation.go"
	documentationTmplPath = "internal/code-generator/assets/documentation.go.tmpl"
	documentationTmpl     = "documentation.go.tmpl"

	// type.go generation
	typeFile     = "type.go"
	typeTmplPath = "internal/code-generator/assets/types.go.tmpl"
	typeTmpl     = "types.go.tmpl"

	// interpreter.go generation
	interpretersFile     = "interpreter.go"
	interpretersTmplPath = "internal/code-generator/assets/interpreters.go.tmpl"
	interpretersTmpl     = "interpreters.go.tmpl"

	// filename.go generation
	filenamesFile     = "filename.go"
	filenamesTmplPath = "internal/code-generator/assets/filenames.go.tmpl"
	filenamesTmpl     = "filenames.go.tmpl"

	// alias.go generation
	aliasesFile     = "alias.go"
	aliasesTmplPath = "internal/code-generator/assets/aliases.go.tmpl"
	aliasesTmpl     = "aliases.go.tmpl"

	// frequencies.go generation
	frequenciesFile     = "frequencies.go"
	frequenciesTmplPath = "internal/code-generator/assets/frequencies.go.tmpl"
	frequenciesTmpl     = "frequencies.go.tmpl"

	commitPath = ".linguist/.git/refs/heads/master"
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
		&generatorFiles{generator.Extensions, languagesYAML, "", extensionsFile, extensionsTmplPath, extensionsTmpl, commit},
		&generatorFiles{generator.Heuristics, heuristicsRuby, "", contentFile, contentTmplPath, contentTmpl, commit},
		&generatorFiles{generator.Vendor, vendorYAML, "", vendorFile, vendorTmplPath, vendorTmpl, commit},
		&generatorFiles{generator.Documentation, documentationYAML, "", documentationFile, documentationTmplPath, documentationTmpl, commit},
		&generatorFiles{generator.Types, languagesYAML, "", typeFile, typeTmplPath, typeTmpl, commit},
		&generatorFiles{generator.Interpreters, languagesYAML, "", interpretersFile, interpretersTmplPath, interpretersTmpl, commit},
		&generatorFiles{generator.Filenames, languagesYAML, samplesDir, filenamesFile, filenamesTmplPath, filenamesTmpl, commit},
		&generatorFiles{generator.Aliases, languagesYAML, "", aliasesFile, aliasesTmplPath, aliasesTmpl, commit},
		&generatorFiles{generator.Frequencies, "", samplesDir, frequenciesFile, frequenciesTmplPath, frequenciesTmpl, commit},
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

	return string(commit), nil
}
