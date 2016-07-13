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
	".cl":   clExtLanguage,
	".inc":  incExtLanguage,
	".cls":  clsExtLanguage,
	".m":    mExtLanguage,
	".ms":   msExtLanguage,
	".h":    hExtLanguage,
	".l":    lExtLanguage,
	".n":    nExtLanguage,
	".lisp": lispExtLanguage,
	".lsp":  lispExtLanguage,
	".pm":   pmExtLanguage,
	".t":    pmExtLanguage,
	".pl":   plExtLanguage,
	".pro":  proExtLanguage,
	".toc":  tocExtLanguage,
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

func incExtLanguage(input []byte) (string, bool) {
	if substring.BytesRegexp(`^<\?(?:php)?`).Match(input) {
		return "PHP", true
	}

	return OtherLanguage, true
}

func hExtLanguage(input []byte) (string, bool) {
	if objectiveCMatcher.Match(input) {
		return "Objective-C", true
	} else if cPlusPlusMatcher.Match(input) {
		return "C++", true
	}

	return "C", true
}

func msExtLanguage(input []byte) (string, bool) {
	if substring.BytesRegexp(`[.'][a-z][a-z](\s|$)`).Match(input) {
		return "Groff", true
	}

	return "MAXScript", true
}

func nExtLanguage(input []byte) (string, bool) {
	if substring.BytesRegexp(`^[.']`).Match(input) {
		return "Groff", true
	} else if substring.BytesRegexp(`(module|namespace|using)`).Match(input) {
		return "Nemerle", true
	}

	return OtherLanguage, false
}

func lExtLanguage(input []byte) (string, bool) {
	if substring.BytesRegexp(`\(def(un|macro)\s`).Match(input) {
		return "Common Lisp", true
	} else if substring.BytesRegexp(`(%[%{}]xs|<.*>)`).Match(input) {
		return "Lex", true
	} else if substring.BytesRegexp(`\.[a-z][a-z](\s|$)`).Match(input) {
		return "Groff", true
	} else if substring.BytesRegexp(`(de|class|rel|code|data|must)`).Match(input) {
		return "PicoLisp", true
	}

	return OtherLanguage, false
}

func lispExtLanguage(input []byte) (string, bool) {
	if commonLispMatcher.Match(input) {
		return "Common Lisp", true
	} else if substring.BytesRegexp(`\s*\(define `).Match(input) {
		return "NewLisp", true
	}

	return OtherLanguage, false
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
	mathematicaMatcher = substring.BytesHas(`\s*\(\*`)
	matlabMatcher      = substring.BytesRegexp(`\b(function\s*[\[a-zA-Z]+|pcolor|classdef|figure|end|elseif)\b`)
	objectiveCMatcher  = substring.BytesRegexp(
		`@(interface|class|protocol|property|end|synchronised|selector|implementation)\b|#import\s+.+\.h[">]`)
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
	perlMatcher   = substring.BytesRegexp(`use strict|use\s+v?5\.`)
	perl6Matcher  = substring.BytesRegexp(`(use v6|(my )?class|module)`)
)

func plExtLanguage(input []byte) (string, bool) {
	if prologMatcher.Match(input) {
		return "Prolog", true
	} else if perl6Matcher.Match(input) {
		return "Perl6", true
	}

	return "Perl", false
}

func pmExtLanguage(input []byte) (string, bool) {
	if perlMatcher.Match(input) {
		return "Perl", true
	} else if perl6Matcher.Match(input) {
		return "Perl6", true
	}

	return "Perl", false
}

func proExtLanguage(input []byte) (string, bool) {
	if prologMatcher.Match(input) {
		return "Prolog", true
	}

	return OtherLanguage, false
}

func tocExtLanguage(input []byte) (string, bool) {
	if substring.BytesRegexp("## |@no-lib-strip@").Match(input) {
		return "World of Warcraft Addon Data", true
	} else if substring.BytesRegexp("(contentsline|defcounter|beamer|boolfalse)").Match(input) {
		return "TeX", true
	}

	return OtherLanguage, false
}
