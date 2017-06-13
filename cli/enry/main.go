package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/src-d/enry.v1"
)

func main() {
	flag.Usage = usage
	flag.Parse()
	root, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	errors := false
	out := make(map[string][]string, 0)
	err = filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			errors = true
			log.Println(err)
			return filepath.SkipDir
		}

		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			errors = true
			log.Println(err)
			return nil
		}

		if relativePath == "." {
			return nil
		}

		if f.IsDir() {
			relativePath = relativePath + "/"
		}

		if enry.IsVendor(relativePath) || enry.IsDotFile(relativePath) ||
			enry.IsDocumentation(relativePath) || enry.IsConfiguration(relativePath) {
			if f.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		if f.IsDir() {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			errors = true
			log.Println(err)
			return nil
		}

		language := enry.GetLanguage(filepath.Base(path), content)
		if language == enry.OtherLanguage {
			return nil
		}

		out[language] = append(out[language], relativePath)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	data, _ := json.MarshalIndent(out, "", "  ")
	fmt.Printf("%s\n", data)

	if errors {
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(
		os.Stderr, "simple-linguist, A simple (and faster) implementation of linguist \nusage: %s <path>\n",
		os.Args[0],
	)

	flag.PrintDefaults()
}
