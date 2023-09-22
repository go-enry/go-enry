package enry

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-enry/go-enry/v2/data"
	"github.com/go-enry/go-enry/v2/regex"
	"github.com/stretchr/testify/suite"
)

type linguistCorpusSuite struct {
	enryBaseTestSuite
}

func Test_EnryOnLinguistCorpus(t *testing.T) {
	suite.Run(t, new(linguistCorpusSuite))
}

// First part of the test_blob.rb#test_language
// https://github.com/github/linguist/blob/59b2d88b2242e6062384e5fb876668cc30ead951/test/test_blob.rb#L258
func (s *linguistCorpusSuite) TestLinguistSamples() {
	const filenamesDir = "filenames"
	var cornerCases = map[string]bool{
		"textobj-rubyblock.vba": true, // unsupported negative lookahead RE syntax (https://github.com/github/linguist/blob/8083cb5a89cee2d99f5a988f165994d0243f0d1e/lib/linguist/heuristics.yml#L521)
		// .es and .ice fail heuristics parsing, but do not fail any tests
		// 'Adblock Filter List' hack https://github.com/github/linguist/blob/bf853f1c663903e3ee35935189760191f1c45e1c/lib/linguist/heuristics.yml#L680-L702
		"Imperial Units Remover.txt": true,
		"abp-filters-anti-cv.txt":    true,
		"anti-facebook.txt":          true,
		"fake-news.txt":              true,
		"test_rules.txt":             true,
	}

	var total, failed, ok, other int
	var expected string
	filepath.Walk(s.samplesDir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			if f.Name() != filenamesDir {
				expected, _ = data.LanguageByAlias(f.Name())
			}

			return nil
		}

		filename := filepath.Base(path)
		content, _ := ioutil.ReadFile(path)

		total++
		got := GetLanguage(filename, content)
		if got == OtherLanguage {
			got = "Other"
			other++
		}

		if expected == got {
			ok++
		} else {
			failed++
		}

		errMsg := fmt.Sprintf("file: %q\texpected: %q\tgot: %q\n", path, expected, got)
		if _, ok := cornerCases[filename]; ok {
			s.T().Logf(fmt.Sprintf("\t[corner case] %s", errMsg))
		} else {
			s.Equal(expected, got, errMsg)
		}
		return nil
	})
	s.T().Logf("\ttotal files: %d, ok: %d, failed: %d, other: %d\n", total, ok, failed, other)
}

// Second part of the test_blob.rb#test_language
// https://github.com/github/linguist/blob/59b2d88b2242e6062384e5fb876668cc30ead951/test/test_blob.rb#L275
func (s *linguistCorpusSuite) TestLinguistGenericFixtures() {
	if regex.Name != regex.Oniguruma {
		s.T().Skip("requires Ruby regexp support")
	}
	const filenamesDir = "filenames"

	var total, failed, ok, multiple int
	var expected string
	filepath.Walk(filepath.Join(s.testFixturesDir, "Generic"), func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			if f.Name() != filenamesDir {
				expected, _ = data.LanguageByAlias(f.Name())
			}
			return nil
		}

		filename := filepath.Base(path)
		content, _ := ioutil.ReadFile(path)

		total++
		result := GetLanguagesByContent(filename, content, nil)

		var obtained, status string
		switch len(result) {
		case 0:
			obtained = ""
		case 1:
			obtained = result[0]
		default:
			obtained = ""
			multiple++
		}

		if expected == obtained {
			status = "ok"
			ok++
		} else {
			status = "failed"
			failed++
		}
		s.Equal(expected, obtained, fmt.Sprintf("%s\texpected: %s\tobtained: %s\tstatus: %s\n", filename, expected, obtained, status))
		return nil
	})
	s.T().Logf("\t\ttotal files: %d, ok: %d, failed: %d, multiple: %d\n", total, ok, failed, multiple)
}
