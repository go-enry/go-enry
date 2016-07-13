package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/src-d/simple-linguist"
)

func main() {
	flag.Parse()
	root, err := filepath.Abs(flag.Arg(0))
	ifError(err)

	if root == "" {
		usage()
	}

	o := make(map[string][]string, 0)
	err = filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
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

		l, safe := slinguist.GetLanguageByExtension(path)
		if !safe {
			content, _ := ioutil.ReadFile(path)
			l, safe = slinguist.GetLanguageByContent(path, content)

		}

		r, _ := filepath.Rel(root, path)
		o[l] = append(o[l], r)
		return nil
	})

	ifError(err)

	js, _ := json.MarshalIndent(o, "", "  ")
	fmt.Printf("%s\n", js)
}

func usage() {
	fmt.Fprintf(
		os.Stderr, "simple-linguist, A simple (and faster) implementation of linguist \nusage: %s <path>\n",
		os.Args[0],
	)

	flag.PrintDefaults()
	os.Exit(2)
}

func ifError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(2)
	}
}
