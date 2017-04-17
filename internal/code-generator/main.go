package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/src-d/simple-linguist.v1/internal/code-generator/generator"
)

const (
	// languages.go generation
	languagesYAML     = ".linguist/lib/linguist/languages.yml"
	langFile          = "languages.go"
	languagesTmplPath = "internal/code-generator/assets/languages.go.tmpl"
	languagesTmpl     = "languages.go.tmpl"

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
	typeFile     = "type.go"
	typeTmplPath = "internal/code-generator/assets/type.go.tmpl"
	typeTmpl     = "type.go.tmpl"

	// interpreters_map.go generation
	interpretersFile     = "interpreters_map.go"
	interpretersTmplPath = "internal/code-generator/assets/interpreters.go.tmpl"
	interpretersTmpl     = "interpreters.go.tmpl"

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
		&generatorArgs{languagesYAML, langFile, languagesTmplPath, languagesTmpl, commit, generator.Languages},
		&generatorArgs{heuristicsRuby, contentFile, contentTmplPath, contentTmpl, commit, generator.Heuristics},
		&generatorArgs{vendorYAML, vendorFile, vendorTmplPath, vendorTmpl, commit, generator.Vendor},
		&generatorArgs{documentationYAML, documentationFile, documentationTmplPath, documentationTmpl, commit, generator.Documentation},
		&generatorArgs{languagesYAML, typeFile, typeTmplPath, typeTmpl, commit, generator.Types},
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
