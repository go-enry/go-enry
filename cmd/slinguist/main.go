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
	root := flag.Arg(0)

	o := make(map[string][]string, 0)
	filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
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

		o[l] = append(o[l], path)
		return nil
	})

	js, _ := json.MarshalIndent(o, "", "  ")
	fmt.Printf("%s\n", js)
}
