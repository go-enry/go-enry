package rule

import (
	"testing"

	"github.com/go-enry/go-enry/v2/regex"
	"github.com/stretchr/testify/assert"
)

const lang = "ActionScript"

var fixtures = []struct {
	name     string
	rule     Heuristic
	numLangs int
	match    string
	noMatch  string
}{
	{"Always", Always(MatchingLanguages(lang)), 1, "a", ""},
	{"Not", Not(MatchingLanguages(lang), regex.MustCompile(`a`)), 1, "b", "a"},
	{"And", And(MatchingLanguages(lang), regex.MustCompile(`a`), regex.MustCompile(`b`)), 1, "ab", "a"},
	{"Or", Or(MatchingLanguages(lang), regex.MustCompile(`a|b`)), 1, "ab", "c"},

	{"OrNil", Or(MatchingLanguages(lang), regex.EnryRegexp(nil)), 1, "", "c"},
	{"NotNil", Not(MatchingLanguages(lang), regex.EnryRegexp(nil)), 1, "b", ""},
}

func TestRules(t *testing.T) {
	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			assert.NotNil(t, f.rule)
			assert.NotNil(t, f.rule.Languages())
			assert.Equal(t, f.numLangs, len(f.rule.Languages()))

			if f.match != "" {
				assert.Truef(t, f.rule.Match([]byte(f.match)),
					"'%s' is expected to .Match() by rule %s%v", f.match, f.name, f.rule)
			}

			if f.noMatch != "" {
				assert.Falsef(t, f.rule.Match([]byte(f.noMatch)),
					"'%s' is expected NOT to .Match() by rule %s%v", f.noMatch, f.name, f.rule)
			}
		})
	}
}
