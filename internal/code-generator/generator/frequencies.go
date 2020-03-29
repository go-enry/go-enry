package generator

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/go-enry/go-enry/v2/internal/tokenizer"
)

type samplesFrequencies struct {
	LanguageTotal  int                       `json:"language_total,omitempty"`
	Languages      map[string]int            `json:"languages,omitempty"`
	TokensTotal    int                       `json:"tokens_total,omitempty"`
	Tokens         map[string]map[string]int `json:"tokens,omitempty"`
	LanguageTokens map[string]int            `json:"language_tokens,omitempty"`
}

// Frequencies reads directories in samplesDir, retrieves information about frequencies of languages and tokens, and write
// the file outPath using tmplName as a template. It complies with type File signature.
func Frequencies(fileToParse, samplesDir, outPath, tmplPath, tmplName, commit string) error {
	freqs, err := getFrequencies(samplesDir)
	if err != nil {
		return err
	}

	if _, ok := os.LookupEnv("ENRY_DEBUG"); ok {
		log.Printf("Total samples: %d\n", freqs.LanguageTotal)
		log.Printf("Total tokens: %d\n", freqs.TokensTotal)

		keys := make([]string, 0, len(freqs.Languages))
		for k := range freqs.Languages {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fmt.Printf(" %s: %d\n", k, freqs.Languages[k])
		}
	}

	buf := &bytes.Buffer{}
	if err := executeFrequenciesTemplate(buf, freqs, tmplPath, tmplName, commit); err != nil {
		return err
	}

	return formatedWrite(outPath, buf.Bytes())
}

func getFrequencies(samplesDir string) (*samplesFrequencies, error) {
	langDirs, err := ioutil.ReadDir(samplesDir)
	if err != nil {
		return nil, err
	}

	var languageTotal int
	var languages = make(map[string]int)
	var tokensTotal int
	var tokens = make(map[string]map[string]int)
	var languageTokens = make(map[string]int)

	for _, langDir := range langDirs {
		if !langDir.IsDir() {
			continue
		}

		lang := langDir.Name()
		samples, err := readSamples(filepath.Join(samplesDir, lang))
		if err != nil {
			log.Println(err)
		}

		if len(samples) == 0 {
			continue
		}

		samplesTokens, err := getTokens(samples)
		if err != nil {
			log.Println(err)
			continue
		}

		languageTotal += len(samples)
		languages[lang] = len(samples)
		tokensTotal += len(samplesTokens)
		languageTokens[lang] = len(samplesTokens)
		tokens[lang] = make(map[string]int)
		for _, token := range samplesTokens {
			tokens[lang][token]++
		}
	}

	return &samplesFrequencies{
		TokensTotal:    tokensTotal,
		LanguageTotal:  languageTotal,
		Tokens:         tokens,
		LanguageTokens: languageTokens,
		Languages:      languages,
	}, nil
}

// readSamples collects ./samples/ filenames from the Linguist codebase, skiping symlinks.
func readSamples(samplesLangDir string) ([]string, error) {
	const specialSubDir = "filenames"
	var samples []string

	err := filepath.Walk(samplesLangDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			switch info.Name() {
			case filepath.Base(samplesLangDir):
				return nil
			case specialSubDir:
				return nil
			default:
				return filepath.SkipDir
			}
		}
		// skip git file symlinks on win and *nix
		if isKnownSymlinkInLinguist(path) || !info.Mode().IsRegular() {
			return nil
		}
		samples = append(samples, path)
		return nil
	})

	return samples, err
}

// isKnownSymlinkInLinguist checks if the file name is on the list of known symlinks.
// On Windows, there is no symlink support in Git [1] and those become regular text files,
// so we have to skip these files manually, maintaing a list here :/
//  1. https://github.com/git-for-windows/git/wiki/Symbolic-Links
//
// $ find -L .linguist/samples -xtype l
func isKnownSymlinkInLinguist(path string) bool {
	return strings.HasSuffix(path, filepath.Join("Ant Build System", "filenames", "build.xml")) ||
		strings.HasSuffix(path, filepath.Join("Markdown", "symlink.md"))
}

func getTokens(samples []string) ([]string, error) {
	tokens := make([]string, 0, 20)
	var anyError error
	for _, sample := range samples {
		content, err := ioutil.ReadFile(sample)
		if err != nil {
			anyError = err
			continue
		}

		t := tokenizer.Tokenize(content)
		tokens = append(tokens, t...)
	}

	return tokens, anyError
}

func executeFrequenciesTemplate(out io.Writer, freqs *samplesFrequencies, tmplPath, tmplName, commit string) error {
	fmap := template.FuncMap{
		"toFloat64": func(num int) string { return fmt.Sprintf("%f", float64(num)) },
		"orderKeys": func(m map[string]int) []string {
			keys := make([]string, 0, len(m))
			for key := range m {
				keys = append(keys, key)
			}

			sort.Strings(keys)
			return keys
		},
		"languageLogProbability": func(language string) string {
			num := math.Log(float64(freqs.Languages[language]) / float64(freqs.LanguageTotal))
			return fmt.Sprintf("%f", num)
		},
		"orderMapMapKeys": func(mm map[string]map[string]int) []string {
			keys := make([]string, 0, len(mm))
			for key := range mm {
				keys = append(keys, key)
			}

			sort.Strings(keys)
			return keys
		},
		"tokenLogProbability": func(language, token string) string {
			num := math.Log(float64(freqs.Tokens[language][token]) / float64(freqs.LanguageTokens[language]))
			return fmt.Sprintf("%f", num)
		},
		"quote": strconv.Quote,
	}
	return executeTemplate(out, tmplName, tmplPath, commit, fmap, freqs)
}
