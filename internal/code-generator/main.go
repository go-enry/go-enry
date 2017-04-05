package main

import (
	"io/ioutil"
	"log"

	"srcd.works/simple-linguist.v1/internal/code-generator/generator"
)

const (
	languagesYAML     = ".linguist/lib/linguist/languages.yml"
	langFile          = "languages.go"
	languagesTmplPath = "internal/code-generator/assets/languages.go.tmpl"
	languagesTmpl     = "languages.go.tmpl"

	heuristicsRuby  = ".linguist/lib/linguist/heuristics.rb"
	contentFile     = "content.go"
	contentTmplPath = "internal/code-generator/assets/content.go.tmpl"
	contentTmpl     = "content.go.tmpl"

	commitPath = ".git/refs/heads/master"
)

func main() {
	commit, err := getCommit(commitPath)
	if err != nil {
		log.Printf("couldn't find commit: %v", err)
	}

	if err := generator.FromFile(languagesYAML, langFile, languagesTmplPath, languagesTmpl, commit, generator.Languages); err != nil {
		log.Println(err)
	}
}

func getCommit(path string) (string, error) {
	commit, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(commit), nil
}
