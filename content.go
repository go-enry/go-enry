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
	".bf":   bfExtLanguage,
	".b":    bExtLanguage,
	".cl":   clExtLanguage,
	".inc":  incExtLanguage,
	".cls":  clsExtLanguage,
	".m":    mExtLanguage,
	".ms":   msExtLanguage,
	".md":   mdExtLanguage,
	".fs":   fsExtLanguage,
	".h":    hExtLanguage,
	".hh":   hhExtLanguage,
	".l":    lExtLanguage,
	".n":    nExtLanguage,
	".lisp": lispExtLanguage,
	".lsp":  lispExtLanguage,
	".pm":   pmExtLanguage,
	".t":    tExtLanguage,
	".ts":   tsExtLanguage,
	".tsx":  tsxExtLanguage,
	".rs":   rsExtLanguage,
	".pl":   plExtLanguage,
	".pro":  proExtLanguage,
	".toc":  tocExtLanguage,
	".sls":  slsExtLanguage,
	".sql":  sqlExtLanguage,
}

var (
	cPlusPlusMatcher = substring.BytesOr(
		substring.BytesRegexp(`\s*template\s*<`),
		substring.BytesRegexp(`\s*#\s*include <(cstdint|string|vector|map|list|array|bitset|queue|stack|forward_list|unordered_map|unordered_set|(i|o|io)stream)>`),
		substring.BytesRegexp(`\n[ \t]*try`),
		substring.BytesRegexp(`\n[ \t]*(class|(using[ \t]+)?namespace)\s+\w+`),
		substring.BytesRegexp(`\n[ \t]*(private|public|protected):$`),
		substring.BytesRegexp(`std::\w+`),
		substring.BytesRegexp(`[ \t]*catch\s*`),
	)
)

func bExtLanguage(input []byte) (string, bool) {
	if substring.BytesRegexp(`(include|modules)`).Match(input) {
		return "Limbo", true
	}

	return "Brainfuck", false
}

func bfExtLanguage(input []byte) (string, bool) {
	if substring.BytesRegexp(`(fprintf|function|return)`).Match(input) {
		return "HyPhy", true
	}

	return "Brainfuck", false
}

func incExtLanguage(input []byte) (string, bool) {
	if substring.BytesRegexp(`^<\?(?:php)?`).Match(input) {
		return "PHP", true
	}

	return OtherLanguage, true
}

func fsExtLanguage(input []byte) (string, bool) {
	if substring.BytesRegexp(`\n(: |new-device)`).Match(input) {
		return "Forth", true
	} else if substring.BytesRegexp(`\s*(#light|import|let|module|namespace|open|type)`).Match(input) {
		return "F#", true
	} else if substring.BytesRegexp(`(#version|precision|uniform|varying|vec[234])`).Match(input) {
		return "GLSL", true
	} else if substring.BytesRegexp(`#include|#pragma\s+(rs|version)|__attribute__`).Match(input) {
		return "Filterscript", true
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

func hhExtLanguage(input []byte) (string, bool) {
	if substring.BytesRegexp(`^<\?(?:hh)?`).Match(input) {
		return "Hack", true
	} else if cPlusPlusMatcher.Match(input) {
		return "C++", true
	}

	return OtherLanguage, false
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
	mathematicaMatcher = substring.BytesHas(`\n\s*\(\*`)
	matlabMatcher      = substring.BytesRegexp(`\b(function\s*[\[a-zA-Z]+|pcolor|classdef|figure|end|elseif)\b`)
	objectiveCMatcher  = substring.BytesRegexp(
		`@(interface|class|protocol|property|end|synchronised|selector|implementation)\b|#import\s+.+\.h[">]`)
)

func mExtLanguage(input []byte) (string, bool) {
	if objectiveCMatcher.Match(input) {
		return "Objective-C", true
	} else if substring.BytesHas(`:- module`).Match(input) {
		return "Mercury", true
	} else if substring.BytesRegexp(`\n: `).Match(input) {
		return "MUF", true
	} else if substring.BytesRegexp(`^\s*;`).Match(input) {
		return "M", true
	} else if mathematicaMatcher.Match(input) {
		return "Mathematica", true
	} else if matlabMatcher.Match(input) {
		return "Matlab", true
	} else if substring.BytesRegexp(`^\w+\s*:\s*module\s*{`).Match(input) {
		return "Matlab", true
	}

	return OtherLanguage, false
}

func mdExtLanguage(input []byte) (string, bool) {
	if substring.BytesRegexp(`\n[-a-z0-9=#!\*\[|]`).Match(input) {
		return "Markdown", true
	} else if substring.BytesRegexp(`\n(;;|\(define_)`).Match(input) {
		return "GCC Machine Description", true
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

func tExtLanguage(input []byte) (string, bool) {
	if perlMatcher.Match(input) {
		return "Perl", true
	} else if perl6Matcher.Match(input) {
		return "Perl6", true
	} else if substring.BytesRegexp(`^\s*%|^\s*var\s+\w+\s*:\s*\w+`).Match(input) {
		return "RenderScript", true
	} else if substring.BytesRegexp(`^\s*use\s+v6\s*;`).Match(input) {
		return "RenderScript", true
	}

	return "Perl", false
}

func rsExtLanguage(input []byte) (string, bool) {
	if substring.BytesRegexp(`(use |fn |mod |pub |macro_rules|impl|#!?\[)`).Match(input) {
		return "Rust", true
	} else if substring.BytesRegexp(`#include|#pragma\s+(rs|version)|__attribute__`).Match(input) {
		return "RenderScript", true
	}

	return OtherLanguage, false
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

func slsExtLanguage(input []byte) (string, bool) {
	if substring.BytesRegexp("## |@no-lib-strip@").Match(input) {
		return "World of Warcraft Addon Data", true
	} else if substring.BytesRegexp("(contentsline|defcounter|beamer|boolfalse)").Match(input) {
		return "TeX", true
	}

	return OtherLanguage, false
}

var (
	pgSQLMatcher = substring.BytesOr(
		substring.BytesRegexp(`(?i)\\i\b|AS \$\$|LANGUAGE '?plpgsql'?`),
		substring.BytesRegexp(`(?i)SECURITY (DEFINER|INVOKER)`),
		substring.BytesRegexp(`BEGIN( WORK| TRANSACTION)?;`),
	)
	db2SQLMatcher = substring.BytesOr(
		substring.BytesRegexp(`(?i)(alter module)|(language sql)|(begin( NOT)+ atomic)`),
		substring.BytesRegexp(`(?i)signal SQLSTATE '[0-9]+'`),
	)
	oracleSQLMatcher = substring.BytesOr(
		substring.BytesRegexp(`(?i)\$\$PLSQL_|XMLTYPE|sysdate|systimestamp|\.nextval|connect by|AUTHID (DEFINER|CURRENT_USER)`),
		substring.BytesRegexp(`(?i)constructor\W+function`),
	)
)

func sqlExtLanguage(input []byte) (string, bool) {
	if pgSQLMatcher.Match(input) {
		return "PLpgSQL", true
	} else if db2SQLMatcher.Match(input) {
		return "SQLPL", true
	} else if oracleSQLMatcher.Match(input) {
		return "PLSQL", true
	}

	return "SQL", false
}

func tsExtLanguage(input []byte) (string, bool) {
	if substring.BytesHas("</TS>").Match(input) {
		return "XML", true
	}

	return "TypeScript", true
}

func tsxExtLanguage(input []byte) (string, bool) {
	if substring.BytesHas("</tileset>").Match(input) {
		return "XML", true
	}

	return "TypeScript", true
}
