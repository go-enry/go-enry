package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/src-d/simple-linguist.v1/internal/code-generator/generator"
)

const (
	// languages info file
	languagesYAML = ".linguist/lib/linguist/languages.yml"

	// extensions_map.go generation
	extensionsFile     = "extensions_map.go"
	extensionsTmplPath = "internal/code-generator/assets/extensions.go.tmpl"
	extensionsTmpl     = "extensions.go.tmpl"

	// content.go generation
	heuristicsRuby  = ".linguist/lib/linguist/heuristics.rb"
	contentFile     = "content.go"
	contentTmplPath = "internal/code-generator/assets/content.go.tmpl"
	contentTmpl     = "content.go.tmpl"

	// vendor_matchers.go generation
	vendorYAML     = ".linguist/lib/linguist/vendor.yml"
	vendorFile     = "vendor_matchers.go"
	vendorTmplPath = "internal/code-generator/assets/vendor.go.tmpl"
	vendorTmpl     = "vendor.go.tmpl"

	// documentation_matchers.go generation
	documentationYAML     = ".linguist/lib/linguist/documentation.yml"
	documentationFile     = "documentation_matchers.go"
	documentationTmplPath = "internal/code-generator/assets/documentation.go.tmpl"
	documentationTmpl     = "documentation.go.tmpl"

	// type.go generation
	typeFile     = "types_map.go"
	typeTmplPath = "internal/code-generator/assets/types.go.tmpl"
	typeTmpl     = "types.go.tmpl"

	// interpreters_map.go generation
	interpretersFile     = "interpreters_map.go"
	interpretersTmplPath = "internal/code-generator/assets/interpreters.go.tmpl"
	interpretersTmpl     = "interpreters.go.tmpl"

	// filenames_map.go generation
	filenamesFile     = "filenames_map.go"
	filenamesTmplPath = "internal/code-generator/assets/filenames.go.tmpl"
	filenamesTmpl     = "filenames.go.tmpl"

	// aliases_map.go generation
	aliasesFile     = "aliases_map.go"
	aliasesTmplPath = "internal/code-generator/assets/aliases.go.tmpl"
	aliasesTmpl     = "aliases.go.tmpl"

	commitPath = ".git/refs/heads/master"
)

type generatorArgs struct {
	fileToParse string
	outPath     string
	tmplPath    string
	tmplName    string
	commit      string
	generate    generator.Func
}

func main() {
	commit, err := getCommit(commitPath)
	if err != nil {
		log.Printf("couldn't find commit: %v", err)
	}

	argsList := []*generatorArgs{
		&generatorArgs{languagesYAML, extensionsFile, extensionsTmplPath, extensionsTmpl, commit, generator.Extensions},
		&generatorArgs{heuristicsRuby, contentFile, contentTmplPath, contentTmpl, commit, generator.Heuristics},
		&generatorArgs{vendorYAML, vendorFile, vendorTmplPath, vendorTmpl, commit, generator.Vendor},
		&generatorArgs{documentationYAML, documentationFile, documentationTmplPath, documentationTmpl, commit, generator.Documentation},
		&generatorArgs{languagesYAML, typeFile, typeTmplPath, typeTmpl, commit, generator.Types},
		&generatorArgs{languagesYAML, interpretersFile, interpretersTmplPath, interpretersTmpl, commit, generator.Interpreters},
		&generatorArgs{languagesYAML, filenamesFile, filenamesTmplPath, filenamesTmpl, commit, generator.Filenames},
		&generatorArgs{languagesYAML, aliasesFile, aliasesTmplPath, aliasesTmpl, commit, generator.Aliases},
	}

	for _, args := range argsList {
		if err := generator.FromFile(args.fileToParse, args.outPath, args.tmplPath, args.tmplName, args.commit, args.generate); err != nil {
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
