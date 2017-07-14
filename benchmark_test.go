package enry

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const samplesDir = ".linguist/samples"

type sample struct {
	filename string
	content  []byte
}

var (
	slow              bool
	overcomeLanguage  string
	overcomeLanguages []string
	samples           []*sample
)

func TestMain(m *testing.M) {
	flag.BoolVar(&slow, "slow", false, "run benchmarks per sample for strategies too")
	flag.Parse()
	var err error
	samples, err = getSamples(samplesDir)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
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
				o = DefaultClassifier.Classify(sample.content, nil)
			}

			overcomeLanguages = o
		}
	})
}

func BenchmarkStrategiesTotal(b *testing.B) {
	if slow {
		b.SkipNow()
	}

	benchmarks := []struct {
		name       string
		strategy   Strategy
		candidates []string
	}{
		{name: "GetLanguagesByModeline()_TOTAL", strategy: GetLanguagesByModeline},
		{name: "GetLanguagesByFilename()_TOTAL", strategy: GetLanguagesByFilename},
		{name: "GetLanguagesByShebang()_TOTAL", strategy: GetLanguagesByShebang},
		{name: "GetLanguagesByExtension()_TOTAL", strategy: GetLanguagesByExtension},
		{name: "GetLanguagesByContent()_TOTAL", strategy: GetLanguagesByContent},
	}

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
				o = DefaultClassifier.Classify(sample.content, nil)
			}

			overcomeLanguages = o
		})
	}
}

func BenchmarkStrategiesPerSample(b *testing.B) {
	if !slow {
		b.SkipNow()
	}

	benchmarks := []struct {
		name       string
		strategy   Strategy
		candidates []string
	}{
		{name: "GetLanguagesByModeline()_SAMPLE_", strategy: GetLanguagesByModeline},
		{name: "GetLanguagesByFilename()_SAMPLE_", strategy: GetLanguagesByFilename},
		{name: "GetLanguagesByShebang()_SAMPLE_", strategy: GetLanguagesByShebang},
		{name: "GetLanguagesByExtension()_SAMPLE_", strategy: GetLanguagesByExtension},
		{name: "GetLanguagesByContent()_SAMPLE_", strategy: GetLanguagesByContent},
	}

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
