package generator

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"path/filepath"
	"sort"
	"strconv"
	"text/template"

	"github.com/bzz/enry/v2/internal/tokenizer"
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

func readSamples(samplesLangDir string) ([]string, error) {
	const samplesLangFilesDir = "filenames"
	sampleFiles, err := ioutil.ReadDir(samplesLangDir)
	if err != nil {
		return nil, err
	}

	var samples []string
	for _, sampleFile := range sampleFiles {
		filename := filepath.Join(samplesLangDir, sampleFile.Name())
		if sampleFile.Mode().IsRegular() {
			samples = append(samples, filename)
			continue
		}

		if sampleFile.IsDir() && sampleFile.Name() == samplesLangFilesDir {
			subSamples, err := readSubSamples(filename)
			if err != nil {
				return nil, err
			}

			samples = append(samples, subSamples...)
		}

	}

	return samples, nil
}

func readSubSamples(path string) ([]string, error) {
	subSamples := []string{}
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.Mode().IsRegular() {
			subSamples = append(subSamples, filepath.Join(path, entry.Name()))
		}
	}

	return subSamples, nil
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
