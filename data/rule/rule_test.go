package rule

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

const lang = "ActionScript"

var fixtures = []struct {
	name     string
	rule     Heuristic
	numLangs int
	matching string
	noMatch  string
}{
	{"Always", Always(MatchingLanguages(lang)), 1, "a", ""},
	{"Not", Not(MatchingLanguages(lang), regexp.MustCompile(`a`)), 1, "b", "a"},
	{"And", And(MatchingLanguages(lang), regexp.MustCompile(`a`), regexp.MustCompile(`b`)), 1, "ab", "a"},
	{"Or", Or(MatchingLanguages(lang), regexp.MustCompile(`a|b`)), 1, "ab", "c"},
}

func TestRules(t *testing.T) {
	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			assert.NotNil(t, f.rule)
			assert.NotNil(t, f.rule.Languages())
			assert.Equal(t, f.numLangs, len(f.rule.Languages()))
			assert.Truef(t, f.rule.Match([]byte(f.matching)),
				"'%s' is expected to .Match() by rule %s%v", f.matching, f.name, f.rule)
			if f.noMatch != "" {
				assert.Falsef(t, f.rule.Match([]byte(f.noMatch)),
					"'%s' is expected NOT to .Match() by rule %s%v", f.noMatch, f.name, f.rule)
			}
		})
	}
}
