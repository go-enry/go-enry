package slinguist

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/tabwriter"

	. "gopkg.in/check.v1"
)

func (s *TSuite) TestGetLanguageByContentLinguistCorpus(c *C) {
	c.Skip("report")

	var total, failed, ok, other, unsafe int

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	filepath.Walk(".linguist/samples", func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			if f.Name() == "filenames" {
				return filepath.SkipDir
			}

			return nil
		}

		expected := filepath.Base(filepath.Dir(path))
		filename := filepath.Base(path)
		extension := filepath.Ext(path)
		content, _ := ioutil.ReadFile(path)

		if extension == "" {
			return nil
		}

		total++
		obtained, safe := GetLanguageByContent(filename, content)
		if obtained == OtherLanguage {
			other++
		}

		var status string
		if expected == obtained {
			status = "ok"
			ok++
		} else {
			status = "failed"
			failed++
			if !safe {
				unsafe++
			}
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%s\n", filename, expected, obtained, safe, status)

		return nil
	})

	fmt.Fprintln(w)
	w.Flush()

	fmt.Printf("total files: %d, ok: %d, failed: %d, unsafe: %d, other: %d\n", total, ok, failed, unsafe, other)

}
