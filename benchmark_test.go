package enry

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/bzz/enry/v2/data"
)

type sample struct {
	filename string
	content  []byte
}

var (
	slow              bool
	overcomeLanguage  string
	overcomeLanguages []string
	samples           []*sample
	samplesDir        string
	cloned            bool
)

func TestMain(m *testing.M) {
	flag.BoolVar(&slow, "slow", false, "run benchmarks per sample for strategies too")
	flag.Parse()

	if err := cloneLinguist(linguistURL); err != nil {
		log.Fatal(err)
	}
	if cloned {
		defer os.RemoveAll(filepath.Dir(samplesDir))
	}

	var err error
	samples, err = getSamples(samplesDir)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func cloneLinguist(linguistURL string) error {
	repoLinguist := os.Getenv(linguistClonedEnvVar)
	cloned = repoLinguist == ""
	if cloned {
		var err error
		repoLinguist, err = ioutil.TempDir("", "linguist-")
		if err != nil {
			return err
		}
	}

	samplesDir = filepath.Join(repoLinguist, "samples")

	if cloned {
		cmd := exec.Command("git", "clone", linguistURL, repoLinguist)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err = os.Chdir(repoLinguist); err != nil {
		return err
	}

	cmd := exec.Command("git", "checkout", data.LinguistCommit)
	if err := cmd.Run(); err != nil {
		return err
	}

	if err = os.Chdir(cwd); err != nil {
		return err
	}

	return nil
}

func getSamples(dir string) ([]*sample, error) {
	samples := make([]*sample, 0, 2000)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		s := &sample{
			filename: path,
			content:  content,
		}
		samples = append(samples, s)
		return nil
	})
	return samples, err
}

func BenchmarkGetLanguageTotal(b *testing.B) {
	if slow {
		b.SkipNow()
	}

	var o string
	b.Run("GetLanguage()_TOTAL", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			for _, sample := range samples {
				o = GetLanguage(sample.filename, sample.content)
			}
		}

		overcomeLanguage = o
	})
}

func BenchmarkClassifyTotal(b *testing.B) {
	if slow {
		b.SkipNow()
	}

	var o []string
	b.Run("Classify()_TOTAL", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			for _, sample := range samples {
				o = defaultClassifier.classify(sample.content, nil)
			}

			overcomeLanguages = o
		}
	})
}

func BenchmarkStrategiesTotal(b *testing.B) {
	if slow {
		b.SkipNow()
	}

	benchmarks := benchmarkForAllStrategies("TOTAL")

	var o []string
	for _, benchmark := range benchmarks {
		b.Run(benchmark.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				for _, sample := range samples {
					o = benchmark.strategy(sample.filename, sample.content, benchmark.candidates)
				}

				overcomeLanguages = o
			}
		})
	}
}

func BenchmarkGetLanguagePerSample(b *testing.B) {
	if !slow {
		b.SkipNow()
	}

	var o string
	for _, sample := range samples {
		b.Run("GetLanguage()_SAMPLE_"+sample.filename, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				o = GetLanguage(sample.filename, sample.content)
			}

			overcomeLanguage = o
		})
	}
}

func BenchmarkClassifyPerSample(b *testing.B) {
	if !slow {
		b.SkipNow()
	}

	var o []string
	for _, sample := range samples {
		b.Run("Classify()_SAMPLE_"+sample.filename, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				o = defaultClassifier.classify(sample.content, nil)
			}

			overcomeLanguages = o
		})
	}
}

func BenchmarkStrategiesPerSample(b *testing.B) {
	if !slow {
		b.SkipNow()
	}

	benchmarks := benchmarkForAllStrategies("SAMPLE")

	var o []string
	for _, benchmark := range benchmarks {
		for _, sample := range samples {
			b.Run(benchmark.name+sample.filename, func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					o = benchmark.strategy(sample.filename, sample.content, benchmark.candidates)
				}

				overcomeLanguages = o
			})
		}
	}
}

type strategyName struct {
	name       string
	strategy   Strategy
	candidates []string
}

func benchmarkForAllStrategies(class string) []strategyName {
	return []strategyName{
		{name: fmt.Sprintf("GetLanguagesByModeline()_%s_", class), strategy: GetLanguagesByModeline},
		{name: fmt.Sprintf("GetLanguagesByFilename()_%s_", class), strategy: GetLanguagesByFilename},
		{name: fmt.Sprintf("GetLanguagesByShebang()_%s_", class), strategy: GetLanguagesByShebang},
		{name: fmt.Sprintf("GetLanguagesByExtension()_%s_", class), strategy: GetLanguagesByExtension},
		{name: fmt.Sprintf("GetLanguagesByContent()_%s_", class), strategy: GetLanguagesByContent},
	}
}
