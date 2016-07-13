package slinguist

import (
	"path/filepath"

	"gopkg.in/toqueteos/substring.v1"
)

func GetLanguageByContent(filename string, content []byte) (lang string, safe bool) {
	if fnMatcher, ok := matchers[filepath.Ext(filename)]; ok {
		lang, safe = fnMatcher(content)
		return
	}

	return GetLanguageByExtension(filename)
}

type languageMatcher func([]byte) (string, bool)

var matchers = map[string]languageMatcher{
	".cl":  clExtLanguage,
	".cls": clsExtLanguage,
	".m":   mExtLanguage,
	".h":   hExtLanguage,
	".pl":  plExtLanguage,
}

var (
	cPlusPlusMatcher = substring.BytesOr(
		substring.BytesRegexp(`^\s*template\s*<`),
		substring.BytesRegexp(`^\s*#\s*include <(cstdint|string|vector|map|list|array|bitset|queue|stack|forward_list|unordered_map|unordered_set|(i|o|io)stream)>`),
		substring.BytesRegexp(`^[ \t]*try`),
		substring.BytesRegexp(`^[ \t]*(class|(using[ \t]+)?namespace)\s+\w+`),
		substring.BytesRegexp(`^[ \t]*(private|public|protected):$`),
		substring.BytesRegexp(`std::\w+`),
		substring.BytesRegexp(`^[ \t]*catch\s*`),
	)
)

func hExtLanguage(input []byte) (string, bool) {
	if objectiveCMatcher.Match(input) {
		return "Objective-C", true
	} else if cPlusPlusMatcher.Match(input) {
		return "C++", true
	}

	return "C", true
}

var (
	commonLispMatcher = substring.BytesRegexp("(?i)(defpackage|defun|in-package)")
	coolMatcher       = substring.BytesRegexp("(?i)class")
	openCLMatcher     = substring.BytesOr(
		substring.BytesHas("\n}"),
		substring.BytesHas("}\n"),
		substring.BytesHas(`/*`),
		substring.BytesHas(`//`),
	)
)

func clExtLanguage(input []byte) (string, bool) {
	if commonLispMatcher.Match(input) {
		return "Common Lisp", true
	} else if coolMatcher.Match(input) {
		return "Cool", true
	} else if openCLMatcher.Match(input) {
		return "OpenCL", true
	}

	return OtherLanguage, false
}

var (
	apexMatcher = substring.BytesOr(
		substring.BytesHas("{\n"),
		substring.BytesHas("}\n"),
	)
	texMatcher = substring.BytesOr(
		substring.BytesHas(`%`),
		substring.BytesHas(`\`),
	)
	openEdgeABLMatcher = substring.BytesRegexp(`(?i)(class|define|interface|method|using)\b`)
	visualBasicMatcher = substring.BytesOr(
		substring.BytesHas("'*"),
		substring.BytesRegexp(`(?i)(attribute|option|sub|private|protected|public|friend)\b`),
	)
)

func clsExtLanguage(input []byte) (string, bool) {
	if texMatcher.Match(input) {
		return "TeX", true
	} else if visualBasicMatcher.Match(input) {
		return "Visual Basic", true
	} else if openEdgeABLMatcher.Match(input) {
		return "OpenEdge ABL", true
	} else if apexMatcher.Match(input) {
		return "Apex", true
	}

	return OtherLanguage, false
}

var (
	mathematicaMatcher = substring.BytesHas(`(*`)
	matlabMatcher      = substring.BytesRegexp(`\b(function\s*[\[a-zA-Z]+|classdef|figure|end|elseif)\b`)
	objectiveCMatcher  = substring.BytesRegexp(
		`@(class|end|implementation|interface|property|protocol|selector|synchronised)`)
)

func mExtLanguage(input []byte) (string, bool) {
	if objectiveCMatcher.Match(input) {
		return "Objective-C", true
	} else if matlabMatcher.Match(input) {
		return "Matlab", true
	} else if mathematicaMatcher.Match(input) {
		return "Mathematica", true
	}

	return OtherLanguage, false
}

var (
	prologMatcher = substring.BytesRegexp(`^[^#]+:-`)
	perl6Matcher  = substring.BytesRegexp(`^(use v6|(my )?class|module)`)
)

func plExtLanguage(input []byte) (string, bool) {
	if prologMatcher.Match(input) {
		return "Prolog", true
	} else if perl6Matcher.Match(input) {
		return "Perl6", true
	}

	return "Perl", false
}
