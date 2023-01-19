package rule

import (
	"testing"

	"github.com/go-enry/go-enry/v2/regex"
	"github.com/stretchr/testify/assert"
)

const lang = "ActionScript"

type fixture struct {
	name     string
	rule     Heuristic
	numLangs int
	match    string
	noMatch  string
}

var specificFixtures = map[string][]fixture{
	"": { // cases that don't vary between the engines
		{"Always", Always(MatchingLanguages(lang)), 1, "a", ""},
		{"Not", Not(MatchingLanguages(lang), regex.MustCompile(`a`)), 1, "b", "a"},
		{"And", And(MatchingLanguages(lang), regex.MustCompile(`a`), regex.MustCompile(`b`)), 1, "ab", "a"},
		{"Or", Or(MatchingLanguages(lang), regex.MustCompile(`a|b`)), 1, "ab", "c"},
		// the results of these depend on the regex engine
		// {"NilOr", Or(noLanguages(), regex.MustCompileRuby(``)), 0, "", "a"},
		// {"NilNot", Not(noLanguages(), regex.MustCompileRuby(`a`)), 0, "", "a"},
	},
	regex.RE2: {
		{"NilAnd", And(noLanguages(), regex.MustCompileRuby(`a`), regex.MustCompile(`b`)), 0, "b", "a"},
		{"NilNot", Not(noLanguages(), regex.MustCompileRuby(`a`), regex.MustCompile(`b`)), 0, "c", "b"},
	},
	regex.Oniguruma: {
		{"NilAnd", And(noLanguages(), regex.MustCompileRuby(`a`), regex.MustCompile(`b`)), 0, "ab", "c"},
		{"NilNot", Not(noLanguages(), regex.MustCompileRuby(`a`), regex.MustCompile(`b`)), 0, "c", "a"},
		{"NilOr", Or(noLanguages(), regex.MustCompileRuby(`a`) /*, regexp.MustCompile(`b`)*/), 0, "a", "b"},
	},
}

func testRulesForEngine(t *testing.T, engine string) {
	if engine != "" && regex.Name != engine {
		return
	}
	for _, f := range specificFixtures[engine] {
		t.Run(engine+f.name, func(t *testing.T) {
			check(t, f)
		})
	}
}

func TestRules(t *testing.T) {
	//TODO(bzz): can all be run in parallel
	testRulesForEngine(t, "")
	testRulesForEngine(t, regex.RE2)
	testRulesForEngine(t, regex.Oniguruma)
}

func check(t *testing.T, f fixture) {
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
}
