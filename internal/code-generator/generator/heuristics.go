package generator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const (
	multilinePrefix = "(?m)"
	orPipe          = "|"
)

// GenHeuristics generates language identification heuristics in Go.
// It is of generator.File type.
func GenHeuristics(fileToParse, _, outPath, tmplPath, tmplName, commit string) error {
	heuristicsYaml, err := parseYaml(fileToParse)
	if err != nil {
		return err
	}

	langPatterns, err := loadHeuristics(heuristicsYaml)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	err = executeTemplate(buf, tmplName, tmplPath, commit, nil, langPatterns)
	if err != nil {
		return err
	}

	return formatedWrite(outPath, buf.Bytes())
}

// loadHeuristics transforms parsed YAML to map[".ext"]->IR for code generation.
func loadHeuristics(yaml *Heuristics) (map[string][]*LanguagePattern, error) {
	var patterns = make(map[string][]*LanguagePattern)
	for _, disambiguation := range yaml.Disambiguations {
		var rules []*LanguagePattern
		for _, rule := range disambiguation.Rules {
			langPattern := loadRule(yaml.NamedPatterns, rule)
			if langPattern != nil {
				rules = append(rules, langPattern)
			}
		}
		// unroll to a single map
		for _, ext := range disambiguation.Extensions {
			if _, ok := patterns[ext]; ok {
				return nil, fmt.Errorf("cannot add extension '%s', it already exists for %+v", ext, patterns[ext])
			}
			patterns[ext] = rules
		}

	}
	return patterns, nil
}

// loadRule transforms single rule from parsed YAML to IR for code generation.
// For OrPattern case, it always combines multiple patterns into a single one.
func loadRule(namedPatterns map[string]StringArray, rule *Rule) *LanguagePattern {
	var result *LanguagePattern
	if len(rule.And) != 0 { // AndPattern
		var subPatterns []*LanguagePattern
		for _, r := range rule.And {
			subp := loadRule(namedPatterns, r)
			subPatterns = append(subPatterns, subp)
		}
		result = &LanguagePattern{"And", rule.Languages, "", subPatterns}
	} else if len(rule.Pattern) != 0 { // OrPattern
		conjunction := strings.Join(rule.Pattern, orPipe)
		pattern := convertToValidRegexp(conjunction)
		result = &LanguagePattern{"Or", rule.Languages, pattern, nil}
	} else if rule.NegativePattern != "" { // NotPattern
		pattern := convertToValidRegexp(rule.NegativePattern)
		result = &LanguagePattern{"Not", rule.Languages, pattern, nil}
	} else if rule.NamedPattern != "" { // Named OrPattern
		conjunction := strings.Join(namedPatterns[rule.NamedPattern], orPipe)
		pattern := convertToValidRegexp(conjunction)
		result = &LanguagePattern{"Or", rule.Languages, pattern, nil}
	} else { // AlwaysPattern
		result = &LanguagePattern{"Always", rule.Languages, "", nil}
	}

	if isUnsupportedRegexpSyntax(result.Pattern) {
		log.Printf("skipping rule: language:'%q', rule:'%q'\n", rule.Languages, result.Pattern)
		return nil
	}
	return result
}

// LanguagePattern is an IR of parsed Rule suitable for code generations.
// Strings are used as this is to be be consumed by text/template.
type LanguagePattern struct {
	Op      string
	Langs   []string
	Pattern string
	Rules   []*LanguagePattern
}

type Heuristics struct {
	Disambiguations []*Disambiguation
	NamedPatterns   map[string]StringArray `yaml:"named_patterns"`
}

type Disambiguation struct {
	Extensions []string `yaml:"extensions,flow"`
	Rules      []*Rule  `yaml:"rules"`
}

type Rule struct {
	Patterns  `yaml:",inline"`
	Languages StringArray `yaml:"language"`
	And       []*Rule
}

type Patterns struct {
	Pattern         StringArray `yaml:"pattern,omitempty"`
	NamedPattern    string      `yaml:"named_pattern,omitempty"`
	NegativePattern string      `yaml:"negative_pattern,omitempty"`
}

// StringArray is workaround for parsing named_pattern,
// wich is sometimes arry and sometimes not.
// See https://github.com/go-yaml/yaml/issues/100
type StringArray []string

// UnmarshalYAML allowes to parse element always as a []string
func (sa *StringArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	if err := unmarshal(&multi); err != nil {
		var single string
		if err := unmarshal(&single); err != nil {
			return err
		}
		*sa = []string{single}
	} else {
		*sa = multi
	}
	return nil
}

func parseYaml(file string) (*Heuristics, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	h := &Heuristics{}
	if err := yaml.Unmarshal(data, &h); err != nil {
		return nil, err
	}

	return h, nil
}

// isUnsupportedRegexpSyntax filters regexp syntax that is not supported by RE2.
// In particular, we stumbled up on usage of next cases:
// - lookbehind & lookahead
// - named & numbered capturing group/after text matching
// - backreference
// - possessive quantifier
// For referece on supported syntax see https://github.com/google/re2/wiki/Syntax
func isUnsupportedRegexpSyntax(reg string) bool {
	return strings.Contains(reg, `(?<`) || strings.Contains(reg, `(?=`) ||
		strings.Contains(reg, `\1`) || strings.Contains(reg, `*+`) ||
		// See https://github.com/github/linguist/pull/4243#discussion_r246105067
		(strings.HasPrefix(reg, multilinePrefix+`/`) && strings.HasSuffix(reg, `/`))
}

// convertToValidRegexp converts Ruby regexp syntaxt to RE2 equivalent.
// Does not work with Ruby regexp literals.
func convertToValidRegexp(rubyRegexp string) string {
	return multilinePrefix + rubyRegexp
}
