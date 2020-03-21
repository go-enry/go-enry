package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYAMLParsing(t *testing.T) {
	heuristics, err := parseYaml("test_files/heuristics.yml")

	require.NoError(t, err)
	assert.NotNil(t, heuristics)

	// extensions
	require.NotNil(t, heuristics.Disambiguations)
	assert.Equal(t, 4, len(heuristics.Disambiguations))
	assert.Equal(t, 2, len(heuristics.Disambiguations[0].Extensions))

	rules := heuristics.Disambiguations[0].Rules
	assert.Equal(t, 2, len(rules))
	require.Equal(t, "Objective-C", rules[0].Languages[0])
	assert.Equal(t, 1, len(rules[0].Pattern))

	rules = heuristics.Disambiguations[1].Rules
	assert.Equal(t, 3, len(rules))
	require.Equal(t, "Forth", rules[0].Languages[0])
	require.Equal(t, 2, len(rules[0].Pattern))

	rules = heuristics.Disambiguations[2].Rules
	assert.Equal(t, 3, len(rules))
	require.Equal(t, "Unix Assembly", rules[1].Languages[0])
	require.NotNil(t, rules[1].And)
	assert.Equal(t, 2, len(rules[1].And))
	require.NotNil(t, rules[1].And[0].NegativePattern)
	assert.Equal(t, "np", rules[1].And[0].NegativePattern)

	rules = heuristics.Disambiguations[3].Rules
	assert.Equal(t, 1, len(rules))
	assert.Equal(t, "Linux Kernel Module", rules[0].Languages[0])
	assert.Equal(t, "AMPL", rules[0].Languages[1])

	// named_patterns
	require.NotNil(t, heuristics.NamedPatterns)
	assert.Equal(t, 2, len(heuristics.NamedPatterns))
	assert.Equal(t, 1, len(heuristics.NamedPatterns["fortran"]))
	assert.Equal(t, 2, len(heuristics.NamedPatterns["cpp"]))
}

func TestSingleRuleLoading(t *testing.T) {
	namedPatterns := map[string]StringArray{"cpp": []string{"cpp_ptrn1", "cpp_ptrn2"}}
	rules := []*Rule{
		&Rule{Languages: []string{"a"}, Patterns: Patterns{NamedPattern: "cpp"}},
		&Rule{Languages: []string{"b"}, And: []*Rule{}},
	}

	// named_pattern case
	langPattern := loadRule(namedPatterns, rules[0])
	require.Equal(t, "a", langPattern.Langs[0])
	assert.NotEmpty(t, langPattern.Pattern)

	// and case
	langPattern = loadRule(namedPatterns, rules[1])
	require.Equal(t, "b", langPattern.Langs[0])
}

func TestLoadingAllHeuristics(t *testing.T) {
	parsedYaml, err := parseYaml("test_files/heuristics.yml")
	require.NoError(t, err)

	hs, err := loadHeuristics(parsedYaml)

	// grep -Eo "extensions:\ (.*)" internal/code-generator/generator/test_files/heuristics.yml
	assert.Equal(t, 5, len(hs))
}

func TestLoadingHeuristicsForSameExt(t *testing.T) {
	parsedYaml := &Heuristics{
		Disambiguations: []*Disambiguation{
			&Disambiguation{
				Extensions: []string{".a", ".b"},
				Rules:      []*Rule{&Rule{Languages: []string{"A"}}},
			},
			&Disambiguation{
				Extensions: []string{".b"},
				Rules:      []*Rule{&Rule{Languages: []string{"B"}}},
			},
		},
	}

	_, err := loadHeuristics(parsedYaml)
	require.Error(t, err)
}

func TestTemplateMatcherVars(t *testing.T) {
	parsed, err := parseYaml("test_files/heuristics.yml")
	require.NoError(t, err)

	heuristics, err := loadHeuristics(parsed)
	require.NoError(t, err)

	// render a tmpl
	const contentTmpl = "../assets/content.go.tmpl"
	tmpl, err := template.New("content.go.tmpl").Funcs(template.FuncMap{
		"stringVal": func(val string) string {
			return fmt.Sprintf("`%s`", val)
		},
	}).ParseFiles(contentTmpl)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, heuristics)
	require.NoError(t, err, fmt.Sprintf("%+v", tmpl))
	require.NotEmpty(t, buf)

	// TODO(bzz) add more advanced test using go/ast package, to verify the
	// strucutre of generated code:
	//  - check key literal exists in map for each extension:

	src, err := format.Source(buf.Bytes())
	require.NoError(t, err, "\n%s\n", string(src))
}
