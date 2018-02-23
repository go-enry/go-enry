package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/src-d/enry.v1"
	"gopkg.in/src-d/enry.v1/data"
)

var (
	version = "undefined"
	build   = "undefined"
	commit  = "undefined"
)

func main() {
	flag.Usage = usage
	breakdownFlag := flag.Bool("breakdown", false, "")
	jsonFlag := flag.Bool("json", false, "")
	showVersion := flag.Bool("version", false, "Show the enry version information")
	onlyProg := flag.Bool("prog", false, "Only show programming file types in output")
	countMode := flag.String("mode", "file", "the method used to count file size. Available options are: file, line and byte")

	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	root, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	fileInfo, err := os.Stat(root)
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.Mode().IsRegular() {
		err = printFileAnalysis(root)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	out := make(map[string][]string, 0)
	err = filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return filepath.SkipDir
		}

		if !f.Mode().IsDir() && !f.Mode().IsRegular() {
			return nil
		}

		relativePath, err := filepath.Rel(root, path)
		if err != nil {
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

		language, ok := enry.GetLanguageByExtension(path)
		if !ok {
			if language, ok = enry.GetLanguageByFilename(path); !ok {
				content, err := ioutil.ReadFile(path)
				if err != nil {
					log.Println(err)
					return nil
				}

				language = enry.GetLanguage(filepath.Base(path), content)
				if language == enry.OtherLanguage {
					return nil
				}
			}
		}

		// If we are displaying only prog. and language is not prog. skip it.
		if *onlyProg && enry.GetLanguageType(language) != enry.Programming {
			return nil
		}

		out[language] = append(out[language], relativePath)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	var buff bytes.Buffer
	switch {
	case *jsonFlag && !*breakdownFlag:
		printJson(out, &buff)
	case *jsonFlag && *breakdownFlag:
		printBreakDown(out, &buff)
	case *breakdownFlag:
		printPercents(out, &buff, *countMode)
		buff.WriteByte('\n')
		printBreakDown(out, &buff)
	default:
		printPercents(out, &buff, *countMode)
	}

	fmt.Print(buff.String())
}

func usage() {
	fmt.Fprintf(
		os.Stderr,
		`  %[1]s %[2]s build: %[3]s commit: %[4]s, based on linguist commit: %[5]s
  %[1]s, A simple (and faster) implementation of github/linguist
  usage: %[1]s [-mode=(file|line|byte)] [-prog] <path>
         %[1]s [-mode=(file|line|byte)] [-prog] [-json] [-breakdown] <path>
         %[1]s [-mode=(file|line|byte)] [-prog] [-json] [-breakdown]
         %[1]s [-version]
`,
		os.Args[0], version, build, commit, data.LinguistCommit[:7],
	)
}

func printBreakDown(out map[string][]string, buff *bytes.Buffer) {
	for name, language := range out {
		writeStringLn(name, buff)
		for _, file := range language {
			writeStringLn(file, buff)
		}

		writeStringLn("", buff)
	}
}

func printJson(out map[string][]string, buff *bytes.Buffer) {
	data, _ := json.Marshal(out)
	buff.Write(data)
	buff.WriteByte('\n')
}

// filelistError represents a failed operation that took place across multiple files.
type filelistError []string

func (e filelistError) Error() string {
	return fmt.Sprintf("Could not process the following files:\n%s", strings.Join(e, "\n"))
}

func printPercents(fSummary map[string][]string, buff *bytes.Buffer, mode string) {
	// Select the way we quantify 'amount' of code.
	var reducer func([]string) (float64, filelistError)
	switch mode {
	case "file":
		reducer = fileCountValues
	case "line":
		reducer = lineCountValues
	case "byte":
		reducer = byteCountValues
	default:
		reducer = fileCountValues
	}

	// Reduce the list of files to a quantity of file type.
	var total float64
	fileValues := make(map[string]float64)
	keys := []string{}
	var unreadableFiles filelistError
	for fType, files := range fSummary {
		val, err := reducer(files)
		if err != nil {
			unreadableFiles = append(unreadableFiles, err...)
		}
		fileValues[fType] = val
		keys = append(keys, fType)
		total += val
	}

	// Slice the keys by their quantity (file count, line count, byte size, etc.).
	sort.Slice(keys, func(i, j int) bool {
		return fileValues[keys[i]] > fileValues[keys[j]]
	})

	// Calculate and write percentages of each file type.
	for _, fType := range keys {
		val := fileValues[fType]
		percent := val / total * 100.0
		buff.WriteString(fmt.Sprintf("%.2f%%\t%s\n", percent, fType))
		if unreadableFiles != nil {
			buff.WriteString(fmt.Sprintf("\n%s", unreadableFiles.Error()))
		}
	}
}

func fileCountValues(files []string) (float64, filelistError) {
	return float64(len(files)), nil
}

func lineCountValues(files []string) (float64, filelistError) {
	var filesErr filelistError
	var t float64
	for _, fName := range files {
		content, err := ioutil.ReadFile(fName)
		if err != nil {
			filesErr = append(filesErr, fName)
			continue
		}
		l, _ := getLines(content)
		t += float64(l)
	}
	return t, filesErr
}

func byteCountValues(files []string) (float64, filelistError) {
	var filesErr filelistError
	var t float64
	for _, fName := range files {
		f, err := os.Open(fName)
		if err != nil {
			filesErr = append(filesErr, fName)
			continue
		}
		fi, err := f.Stat()
		f.Close()
		if err != nil {
			filesErr = append(filesErr, fName)
			continue
		}
		t += float64(fi.Size())
	}
	return t, filesErr
}

func printFileAnalysis(fName string) error {
	content, err := ioutil.ReadFile(fName)
	if err != nil {
		return err
	}

	totalLines, nonBlank := getLines(content)
	fileType := getFileType(fName, content)
	language := enry.GetLanguage(fName, content)
	mimeType := enry.GetMimeType(fName, language)

	fmt.Printf(
		`%s: %d lines (%d sloc)
  type:      %s
  mime_type: %s
  language:  %s
`,
		filepath.Base(fName), totalLines, nonBlank, fileType, mimeType, language,
	)
	return nil
}

func getLines(b []byte) (total int, nonBlank int) {
	scanner := bufio.NewScanner(bytes.NewReader(b))
	lineCt := 0
	blankCt := 0

	for scanner.Scan() {
		lineCt++
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			blankCt++
		}
	}
	// Scanner doesn't catch the case of last byte newline.
	if len(b) > 0 && b[len(b)-1] == '\n' {
		lineCt++
		blankCt++
	}

	return lineCt, lineCt - blankCt
}

func getFileType(file string, content []byte) string {
	switch {
	case enry.IsImage(file):
		return "Image"
	case enry.IsBinary(content):
		return "Binary"
	default:
		return "Text"
	}
}

func writeStringLn(s string, buff *bytes.Buffer) {
	buff.WriteString(s)
	buff.WriteByte('\n')
}
