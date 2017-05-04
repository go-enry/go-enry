package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/src-d/simple-linguist.v1"
)

func main() {
	flag.Usage = usage
	flag.Parse()
	root, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	errors := false
	o := make(map[string][]string, 0)
	err = filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			errors = true
			log.Println(err)
			return filepath.SkipDir
		}

		if slinguist.IsVendor(f.Name()) || slinguist.IsDotFile(f.Name()) {
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

		l := slinguist.GetLanguage(path, content)

		r, err := filepath.Rel(root, path)
		if err != nil {
			errors = true
			log.Println(err)
			return nil
		}

		o[l] = append(o[l], r)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	js, _ := json.MarshalIndent(o, "", "  ")
	fmt.Printf("%s\n", js)

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
