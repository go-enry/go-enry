package data

import (
	"regexp"
	"testing"

	"github.com/bzz/enry/v2/data/rule"
	"github.com/stretchr/testify/assert"
)

var testContentHeuristics = map[string]*Heuristics{
	".md": &Heuristics{ // final pattern for parsed YAML rule
		rule.Or(
			rule.MatchingLanguages("Markdown"),
			regexp.MustCompile(`(^[-A-Za-z0-9=#!\*\[|>])|<\/ | \A\z`),
		),
		rule.Or(
			rule.MatchingLanguages("GCC Machine Description"),
			regexp.MustCompile(`^(;;|\(define_)`),
		),
		rule.Always(
			rule.MatchingLanguages("Markdown"),
		),
	},
	".ms": &Heuristics{
		// Order defines precedence: And, Or, Not, Named, Always
		rule.And(
			rule.MatchingLanguages("Unix Assembly"),
			rule.Not(rule.MatchingLanguages(""), regexp.MustCompile(`/\*`)),
			rule.Or(
				rule.MatchingLanguages(""),
				regexp.MustCompile(`^\s*\.(?:include\s|globa?l\s|[A-Za-z][_A-Za-z0-9]*:)`),
			),
		),
		rule.Or(
			rule.MatchingLanguages("Roff"),
			regexp.MustCompile(`^[.''][A-Za-z]{2}(\s|$)`),
		),
		rule.Always(
			rule.MatchingLanguages("MAXScript"),
		),
	},
}

func TestContentHeuristic_MatchingAlways(t *testing.T) {
	lang := testContentHeuristics[".md"].matchString("")
	assert.Equal(t, []string{"Markdown"}, lang)

	lang = testContentHeuristics[".ms"].matchString("")
	assert.Equal(t, []string{"MAXScript"}, lang)
}

func TestContentHeuristic_MatchingAnd(t *testing.T) {
	lang := testContentHeuristics[".md"].matchString(";;")
	assert.Equal(t, []string{"GCC Machine Description"}, lang)
}

func TestContentHeuristic_MatchingOr(t *testing.T) {
	lang := testContentHeuristics[".ms"].matchString("	.include \"math.s\"")
	assert.Equal(t, []string{"Unix Assembly"}, lang)
}
