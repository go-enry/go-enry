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
	"text/template"

	"gopkg.in/src-d/simple-linguist.v1/internal/tokenizer"
)

const samplesSubDir = "filenames"

type samplesFrequencies struct {
	LanguageTotal  int                       `json:"language_total,omitempty"`
	Languages      map[string]int            `json:"languages,omitempty"`
	TokensTotal    int                       `json:"tokens_total,omitempty"`
	Tokens         map[string]map[string]int `json:"tokens,omitempty"`
	LanguageTokens map[string]int            `json:"language_tokens,omitempty"`
}

// Frequencies reads directories in samplesDir, retrieves information about frequencies of languages and tokens, and write
// the file outPath using frequenciesTmplName as a template.
func Frequencies(samplesDir, frequenciesTmplPath, frequenciesTmplName, commit, outPath string) error {
	freqs, err := getFrequencies(samplesDir)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	if err := executeFrequenciesTemplate(buf, freqs, frequenciesTmplPath, frequenciesTmplName, commit); err != nil {
		return err
	}

	return formatedWrite(outPath, buf.Bytes())
}

func getFrequencies(samplesDir string) (*samplesFrequencies, error) {
	entries, err := ioutil.ReadDir(samplesDir)
	if err != nil {
		return nil, err
	}

	var languageTotal int
	var languages = make(map[string]int)
	var tokensTotal int
	var tokens = make(map[string]map[string]int)
	var languageTokens = make(map[string]int)

	for _, entry := range entries {
		if !entry.IsDir() {
			log.Println(err)
			continue
		}

		samples, err := getSamples(samplesDir, entry)
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

		lang := entry.Name()
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

func getSamples(samplesDir string, langDir os.FileInfo) ([]string, error) {
	samples := []string{}
	path := filepath.Join(samplesDir, langDir.Name())
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.Mode().IsRegular() {
			samples = append(samples, filepath.Join(path, entry.Name()))
		}

		if entry.IsDir() && entry.Name() == samplesSubDir {
			subSamples, err := getSubSamples(samplesDir, langDir.Name(), entry)
			if err != nil {
				return nil, err
			}

			samples = append(samples, subSamples...)
		}

	}

	return samples, nil
}

func getSubSamples(samplesDir, langDir string, subLangDir os.FileInfo) ([]string, error) {
	subSamples := []string{}
	path := filepath.Join(samplesDir, langDir, subLangDir.Name())
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

func executeFrequenciesTemplate(out io.Writer, freqs *samplesFrequencies, frequenciesTmplPath, frequenciesTmpl, commit string) error {
	fmap := template.FuncMap{
		"getCommit": func() string { return commit },
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

	t := template.Must(template.New(frequenciesTmpl).Funcs(fmap).ParseFiles(frequenciesTmplPath))
	if err := t.Execute(out, freqs); err != nil {
		return err
	}

	return nil
}
