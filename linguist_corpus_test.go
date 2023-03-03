package enry

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-enry/go-enry/v2/data"
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
		"drop_stuff.sql":        true, // https://github.com/src-d/enry/issues/194
		"textobj-rubyblock.vba": true, // Because of unsupported negative lookahead RE syntax (https://github.com/github/linguist/blob/8083cb5a89cee2d99f5a988f165994d0243f0d1e/lib/linguist/heuristics.yml#L521)
		// .es and .ice fail heuristics parsing, but do not fail any tests
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
			s.T().Logf(fmt.Sprintf("\t\t[corner case] %s", errMsg))
		} else {
			s.Equal(expected, got, errMsg)
		}
		return nil
	})
	s.T().Logf("\t\ttotal files: %d, ok: %d, failed: %d, other: %d\n", total, ok, failed, other)
}
