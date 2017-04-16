package slinguist

// CODE GENERATED AUTOMATICALLY WITH github.com/src-d/simple-linguist/cli/slinguist-generate
// THIS FILE SHOULD NOT BE EDITED BY HAND
// Extracted from github/linguist commit: dae33dc2b20cddc85d1300435c3be7118a7115a9

import (
	"path/filepath"
	"regexp"
	"strings"
)

func GetLanguageByContent(filename string, content []byte) (lang string, safe bool) {
	ext := strings.ToLower(filepath.Ext(filename))
	if fnMatcher, ok := matchers[ext]; ok {
		lang, safe = fnMatcher(content)
		return
	}

	return GetLanguageByExtension(filename)
}

type languageMatcher func([]byte) (string, bool)

var matchers = map[string]languageMatcher{
	".asc": func(i []byte) (string, bool) {
		if asc_PublicKey_Matcher_0.Match(i) {
			return "Public Key", true
		} else if asc_AsciiDoc_Matcher_0.Match(i) {
			return "AsciiDoc", true
		} else if asc_AGSScript_Matcher_0.Match(i) {
			return "AGS Script", true
		}

		return OtherLanguage, false
	},
	".bb": func(i []byte) (string, bool) {
		if bb_BlitzBasic_Matcher_0.Match(i) || bb_BlitzBasic_Matcher_1.Match(i) {
			return "BlitzBasic", true
		} else if bb_BitBake_Matcher_0.Match(i) {
			return "BitBake", true
		}

		return OtherLanguage, false
	},
	".builds": func(i []byte) (string, bool) {
		if builds_XML_Matcher_0.Match(i) {
			return "XML", true
		}

		return "Text", true
	},
	".ch": func(i []byte) (string, bool) {
		if ch_xBase_Matcher_0.Match(i) {
			return "xBase", true
		}

		return OtherLanguage, false
	},
	".cl": func(i []byte) (string, bool) {
		if cl_CommonLisp_Matcher_0.Match(i) {
			return "Common Lisp", true
		} else if cl_Cool_Matcher_0.Match(i) {
			return "Cool", true
		} else if cl_OpenCL_Matcher_0.Match(i) {
			return "OpenCL", true
		}

		return OtherLanguage, false
	},
	".cls": func(i []byte) (string, bool) {
		if cls_TeX_Matcher_0.Match(i) {
			return "TeX", true
		}

		return OtherLanguage, false
	},
	".cs": func(i []byte) (string, bool) {
		if cs_Smalltalk_Matcher_0.Match(i) {
			return "Smalltalk", true
		} else if cs_CSharp_Matcher_0.Match(i) || cs_CSharp_Matcher_1.Match(i) {
			return "C#", true
		}

		return OtherLanguage, false
	},
	".d": func(i []byte) (string, bool) {
		if d_D_Matcher_0.Match(i) {
			return "D", true
		} else if d_DTrace_Matcher_0.Match(i) {
			return "DTrace", true
		} else if d_Makefile_Matcher_0.Match(i) {
			return "Makefile", true
		}

		return OtherLanguage, false
	},
	".ecl": func(i []byte) (string, bool) {
		if ecl_ECLiPSe_Matcher_0.Match(i) {
			return "ECLiPSe", true
		} else if ecl_ECL_Matcher_0.Match(i) {
			return "ECL", true
		}

		return OtherLanguage, false
	},
	".es": func(i []byte) (string, bool) {
		if es_Erlang_Matcher_0.Match(i) {
			return "Erlang", true
		}

		return OtherLanguage, false
	},
	".f": func(i []byte) (string, bool) {
		if f_Forth_Matcher_0.Match(i) {
			return "Forth", true
		} else if f_FilebenchWML_Matcher_0.Match(i) {
			return "Filebench WML", true
		} else if f_FORTRAN_Matcher_0.Match(i) {
			return "FORTRAN", true
		}

		return OtherLanguage, false
	},
	".for": func(i []byte) (string, bool) {
		if for_Forth_Matcher_0.Match(i) {
			return "Forth", true
		} else if for_FORTRAN_Matcher_0.Match(i) {
			return "FORTRAN", true
		}

		return OtherLanguage, false
	},
	".fr": func(i []byte) (string, bool) {
		if fr_Forth_Matcher_0.Match(i) {
			return "Forth", true
		} else if fr_Frege_Matcher_0.Match(i) {
			return "Frege", true
		}

		return "Text", true
	},
	".fs": func(i []byte) (string, bool) {
		if fs_Forth_Matcher_0.Match(i) {
			return "Forth", true
		} else if fs_FSharp_Matcher_0.Match(i) {
			return "F#", true
		} else if fs_GLSL_Matcher_0.Match(i) {
			return "GLSL", true
		} else if fs_Filterscript_Matcher_0.Match(i) {
			return "Filterscript", true
		}

		return OtherLanguage, false
	},
	".gs": func(i []byte) (string, bool) {
		if gs_Gosu_Matcher_0.Match(i) {
			return "Gosu", true
		}

		return OtherLanguage, false
	},
	".h": func(i []byte) (string, bool) {
		if h_ObjectiveDashC_Matcher_0.Match(i) {
			return "Objective-C", true
		} else if h_CPlusPlus_Matcher_0.Match(i) || h_CPlusPlus_Matcher_1.Match(i) || h_CPlusPlus_Matcher_2.Match(i) || h_CPlusPlus_Matcher_3.Match(i) || h_CPlusPlus_Matcher_4.Match(i) || h_CPlusPlus_Matcher_5.Match(i) || h_CPlusPlus_Matcher_6.Match(i) {
			return "C++", true
		}

		return OtherLanguage, false
	},
	".inc": func(i []byte) (string, bool) {
		if inc_PHP_Matcher_0.Match(i) {
			return "PHP", true
		} else if inc_POVDashRaySDL_Matcher_0.Match(i) {
			return "POV-Ray SDL", true
		}

		return OtherLanguage, false
	},
	".l": func(i []byte) (string, bool) {
		if l_CommonLisp_Matcher_0.Match(i) {
			return "Common Lisp", true
		} else if l_Lex_Matcher_0.Match(i) {
			return "Lex", true
		} else if l_Groff_Matcher_0.Match(i) {
			return "Groff", true
		} else if l_PicoLisp_Matcher_0.Match(i) {
			return "PicoLisp", true
		}

		return OtherLanguage, false
	},
	".ls": func(i []byte) (string, bool) {
		if ls_LoomScript_Matcher_0.Match(i) {
			return "LoomScript", true
		}

		return "LiveScript", true
	},
	".lsp": func(i []byte) (string, bool) {
		if lsp_CommonLisp_Matcher_0.Match(i) {
			return "Common Lisp", true
		} else if lsp_NewLisp_Matcher_0.Match(i) {
			return "NewLisp", true
		}

		return OtherLanguage, false
	},
	".lisp": func(i []byte) (string, bool) {
		if lisp_CommonLisp_Matcher_0.Match(i) {
			return "Common Lisp", true
		} else if lisp_NewLisp_Matcher_0.Match(i) {
			return "NewLisp", true
		}

		return OtherLanguage, false
	},
	".m": func(i []byte) (string, bool) {
		if m_ObjectiveDashC_Matcher_0.Match(i) {
			return "Objective-C", true
		} else if m_Mercury_Matcher_0.Match(i) {
			return "Mercury", true
		} else if m_MUF_Matcher_0.Match(i) {
			return "MUF", true
		} else if m_M_Matcher_0.Match(i) {
			return "M", true
		} else if m_Mathematica_Matcher_0.Match(i) {
			return "Mathematica", true
		} else if m_Matlab_Matcher_0.Match(i) {
			return "Matlab", true
		} else if m_Limbo_Matcher_0.Match(i) {
			return "Limbo", true
		}

		return OtherLanguage, false
	},
	".md": func(i []byte) (string, bool) {
		if md_Markdown_Matcher_0.Match(i) || md_Markdown_Matcher_1.Match(i) {
			return "Markdown", true
		} else if md_GCCmachinedescription_Matcher_0.Match(i) {
			return "GCC machine description", true
		}

		return "Markdown", true
	},
	".ml": func(i []byte) (string, bool) {
		if ml_OCaml_Matcher_0.Match(i) {
			return "OCaml", true
		} else if ml_StandardML_Matcher_0.Match(i) {
			return "Standard ML", true
		}

		return OtherLanguage, false
	},
	".mod": func(i []byte) (string, bool) {
		if mod_XML_Matcher_0.Match(i) {
			return "XML", true
		} else if mod_ModulaDash2_Matcher_0.Match(i) || mod_ModulaDash2_Matcher_1.Match(i) {
			return "Modula-2", true
		}

		return "Linux Kernel Module", false
	},
	".ms": func(i []byte) (string, bool) {
		if ms_Groff_Matcher_0.Match(i) {
			return "Groff", true
		}

		return "MAXScript", true
	},
	".n": func(i []byte) (string, bool) {
		if n_Groff_Matcher_0.Match(i) {
			return "Groff", true
		} else if n_Nemerle_Matcher_0.Match(i) {
			return "Nemerle", true
		}

		return OtherLanguage, false
	},
	".ncl": func(i []byte) (string, bool) {
		if ncl_Text_Matcher_0.Match(i) {
			return "Text", true
		}

		return OtherLanguage, false
	},
	".nl": func(i []byte) (string, bool) {
		if nl_NL_Matcher_0.Match(i) {
			return "NL", true
		}

		return "NewLisp", true
	},
	".php": func(i []byte) (string, bool) {
		if php_Hack_Matcher_0.Match(i) {
			return "Hack", true
		} else if php_PHP_Matcher_0.Match(i) {
			return "PHP", true
		}

		return OtherLanguage, false
	},
	".pl": func(i []byte) (string, bool) {
		if pl_Prolog_Matcher_0.Match(i) {
			return "Prolog", true
		} else if pl_Perl_Matcher_0.Match(i) {
			return "Perl", true
		} else if pl_Perl6_Matcher_0.Match(i) {
			return "Perl6", true
		}

		return OtherLanguage, false
	},
	".pm": func(i []byte) (string, bool) {
		if pm_Perl_Matcher_0.Match(i) {
			return "Perl", true
		} else if pm_Perl6_Matcher_0.Match(i) {
			return "Perl6", true
		}

		return OtherLanguage, false
	},
	".t": func(i []byte) (string, bool) {
		if t_Perl_Matcher_0.Match(i) {
			return "Perl", true
		} else if t_Perl6_Matcher_0.Match(i) {
			return "Perl6", true
		}

		return OtherLanguage, false
	},
	".pod": func(i []byte) (string, bool) {
		if pod_Pod_Matcher_0.Match(i) {
			return "Pod", true
		}

		return "Perl", true
	},
	".pro": func(i []byte) (string, bool) {
		if pro_Prolog_Matcher_0.Match(i) {
			return "Prolog", true
		} else if pro_INI_Matcher_0.Match(i) {
			return "INI", true
		} else if pro_QMake_Matcher_0.Match(i) && pro_QMake_Matcher_1.Match(i) {
			return "QMake", true
		} else if pro_IDL_Matcher_0.Match(i) {
			return "IDL", true
		}

		return OtherLanguage, false
	},
	".props": func(i []byte) (string, bool) {
		if props_XML_Matcher_0.Match(i) {
			return "XML", true
		} else if props_INI_Matcher_0.Match(i) {
			return "INI", true
		}

		return OtherLanguage, false
	},
	".r": func(i []byte) (string, bool) {
		if r_Rebol_Matcher_0.Match(i) {
			return "Rebol", true
		} else if r_R_Matcher_0.Match(i) {
			return "R", true
		}

		return OtherLanguage, false
	},
	".rno": func(i []byte) (string, bool) {
		if rno_RUNOFF_Matcher_0.Match(i) {
			return "RUNOFF", true
		} else if rno_Groff_Matcher_0.Match(i) {
			return "Groff", true
		}

		return OtherLanguage, false
	},
	".rpy": func(i []byte) (string, bool) {
		if rpy_Python_Matcher_0.Match(i) {
			return "Python", true
		}

		return "Ren'Py", true
	},
	".rs": func(i []byte) (string, bool) {
		if rs_Rust_Matcher_0.Match(i) {
			return "Rust", true
		} else if rs_RenderScript_Matcher_0.Match(i) {
			return "RenderScript", true
		}

		return OtherLanguage, false
	},
	".sc": func(i []byte) (string, bool) {
		if sc_SuperCollider_Matcher_0.Match(i) || sc_SuperCollider_Matcher_1.Match(i) || sc_SuperCollider_Matcher_2.Match(i) {
			return "SuperCollider", true
		} else if sc_Scala_Matcher_0.Match(i) || sc_Scala_Matcher_1.Match(i) || sc_Scala_Matcher_2.Match(i) {
			return "Scala", true
		}

		return OtherLanguage, false
	},
	".sql": func(i []byte) (string, bool) {
		if sql_PLpgSQL_Matcher_0.Match(i) || sql_PLpgSQL_Matcher_1.Match(i) || sql_PLpgSQL_Matcher_2.Match(i) {
			return "PLpgSQL", true
		} else if sql_SQLPL_Matcher_0.Match(i) || sql_SQLPL_Matcher_1.Match(i) {
			return "SQLPL", true
		} else if sql_PLSQL_Matcher_0.Match(i) || sql_PLSQL_Matcher_1.Match(i) {
			return "PLSQL", true
		} else if sql_SQL_Matcher_0.Match(i) {
			return "SQL", true
		}

		return OtherLanguage, false
	},
	".srt": func(i []byte) (string, bool) {
		if srt_SubRipText_Matcher_0.Match(i) {
			return "SubRip Text", true
		}

		return OtherLanguage, false
	},
	".toc": func(i []byte) (string, bool) {
		if toc_WorldofWarcraftAddonData_Matcher_0.Match(i) {
			return "World of Warcraft Addon Data", true
		} else if toc_TeX_Matcher_0.Match(i) {
			return "TeX", true
		}

		return OtherLanguage, false
	},
	".ts": func(i []byte) (string, bool) {
		if ts_XML_Matcher_0.Match(i) {
			return "XML", true
		}

		return "TypeScript", true
	},
	".tst": func(i []byte) (string, bool) {
		if tst_GAP_Matcher_0.Match(i) {
			return "GAP", true
		}

		return "Scilab", true
	},
	".tsx": func(i []byte) (string, bool) {
		if tsx_TypeScript_Matcher_0.Match(i) {
			return "TypeScript", true
		} else if tsx_XML_Matcher_0.Match(i) {
			return "XML", true
		}

		return OtherLanguage, false
	},
}

var (
	asc_PublicKey_Matcher_0                = regexp.MustCompile(`(?m)^(----[- ]BEGIN|ssh-(rsa|dss)) `)
	asc_AsciiDoc_Matcher_0                 = regexp.MustCompile(`(?m)^[=-]+(\s|\n)|{{[A-Za-z]`)
	asc_AGSScript_Matcher_0                = regexp.MustCompile(`(?m)^(\/\/.+|((import|export)\s+)?(function|int|float|char)\s+((room|repeatedly|on|game)_)?([A-Za-z]+[A-Za-z_0-9]+)\s*[;\(])`)
	bb_BlitzBasic_Matcher_0                = regexp.MustCompile(`(?m)^\s*; `)
	bb_BlitzBasic_Matcher_1                = regexp.MustCompile(`(?m)End Function`)
	bb_BitBake_Matcher_0                   = regexp.MustCompile(`(?m)^\s*(# |include|require)\b`)
	builds_XML_Matcher_0                   = regexp.MustCompile(`(?mi)^(\s*)(<Project|<Import|<Property|<?xml|xmlns)`)
	ch_xBase_Matcher_0                     = regexp.MustCompile(`(?mi)^\s*#\s*(if|ifdef|ifndef|define|command|xcommand|translate|xtranslate|include|pragma|undef)\b`)
	cl_CommonLisp_Matcher_0                = regexp.MustCompile(`(?mi)^\s*\((defun|in-package|defpackage) `)
	cl_Cool_Matcher_0                      = regexp.MustCompile(`(?m)^class`)
	cl_OpenCL_Matcher_0                    = regexp.MustCompile(`(?m)\/\* |\/\/ |^\}`)
	cls_TeX_Matcher_0                      = regexp.MustCompile(`(?m)\\\w+{`)
	cs_Smalltalk_Matcher_0                 = regexp.MustCompile(`(?m)![\w\s]+methodsFor: `)
	cs_CSharp_Matcher_0                    = regexp.MustCompile(`(?m)^\s*namespace\s*[\w\.]+\s*{`)
	cs_CSharp_Matcher_1                    = regexp.MustCompile(`(?m)^\s*\/\/`)
	d_D_Matcher_0                          = regexp.MustCompile(`(?m)^module\s+[\w.]*\s*;|import\s+[\w\s,.:]*;|\w+\s+\w+\s*\(.*\)(?:\(.*\))?\s*{[^}]*}|unittest\s*(?:\(.*\))?\s*{[^}]*}`)
	d_DTrace_Matcher_0                     = regexp.MustCompile(`(?m)^(\w+:\w*:\w*:\w*|BEGIN|END|provider\s+|(tick|profile)-\w+\s+{[^}]*}|#pragma\s+D\s+(option|attributes|depends_on)\s|#pragma\s+ident\s)`)
	d_Makefile_Matcher_0                   = regexp.MustCompile(`(?m)([\/\\].*:\s+.*\s\\$|: \\$|^ : |^[\w\s\/\\.]+\w+\.\w+\s*:\s+[\w\s\/\\.]+\w+\.\w+)`)
	ecl_ECLiPSe_Matcher_0                  = regexp.MustCompile(`(?m)^[^#]+:-`)
	ecl_ECL_Matcher_0                      = regexp.MustCompile(`(?m):=`)
	es_Erlang_Matcher_0                    = regexp.MustCompile(`(?m)^\s*(?:%%|main\s*\(.*?\)\s*->)`)
	f_Forth_Matcher_0                      = regexp.MustCompile(`(?m)^: `)
	f_FilebenchWML_Matcher_0               = regexp.MustCompile(`(?m)flowop`)
	f_FORTRAN_Matcher_0                    = regexp.MustCompile(`(?mi)^([c*][^abd-z]|      (subroutine|program|end|data)\s|\s*!)`)
	for_Forth_Matcher_0                    = regexp.MustCompile(`(?m)^: `)
	for_FORTRAN_Matcher_0                  = regexp.MustCompile(`(?mi)^([c*][^abd-z]|      (subroutine|program|end|data)\s|\s*!)`)
	fr_Forth_Matcher_0                     = regexp.MustCompile(`(?m)^(: |also |new-device|previous )`)
	fr_Frege_Matcher_0                     = regexp.MustCompile(`(?m)^\s*(import|module|package|data|type) `)
	fs_Forth_Matcher_0                     = regexp.MustCompile(`(?m)^(: |new-device)`)
	fs_FSharp_Matcher_0                    = regexp.MustCompile(`(?m)^\s*(#light|import|let|module|namespace|open|type)`)
	fs_GLSL_Matcher_0                      = regexp.MustCompile(`(?m)^\s*(#version|precision|uniform|varying|vec[234])`)
	fs_Filterscript_Matcher_0              = regexp.MustCompile(`(?m)#include|#pragma\s+(rs|version)|__attribute__`)
	gs_Gosu_Matcher_0                      = regexp.MustCompile(`(?m)^uses java\.`)
	h_ObjectiveDashC_Matcher_0             = regexp.MustCompile(`(?m)^\s*(@(interface|class|protocol|property|end|synchronised|selector|implementation)\b|#import\s+.+\.h[">])`)
	h_CPlusPlus_Matcher_0                  = regexp.MustCompile(`(?m)^\s*#\s*include <(cstdint|string|vector|map|list|array|bitset|queue|stack|forward_list|unordered_map|unordered_set|(i|o|io)stream)>`)
	h_CPlusPlus_Matcher_1                  = regexp.MustCompile(`(?m)^\s*template\s*<`)
	h_CPlusPlus_Matcher_2                  = regexp.MustCompile(`(?m)^[ \t]*try`)
	h_CPlusPlus_Matcher_3                  = regexp.MustCompile(`(?m)^[ \t]*catch\s*\(`)
	h_CPlusPlus_Matcher_4                  = regexp.MustCompile(`(?m)^[ \t]*(class|(using[ \t]+)?namespace)\s+\w+`)
	h_CPlusPlus_Matcher_5                  = regexp.MustCompile(`(?m)^[ \t]*(private|public|protected):$`)
	h_CPlusPlus_Matcher_6                  = regexp.MustCompile(`(?m)std::\w+`)
	inc_PHP_Matcher_0                      = regexp.MustCompile(`(?m)^<\?(?:php)?`)
	inc_POVDashRaySDL_Matcher_0            = regexp.MustCompile(`(?m)^\s*#(declare|local|macro|while)\s`)
	l_CommonLisp_Matcher_0                 = regexp.MustCompile(`(?m)\(def(un|macro)\s`)
	l_Lex_Matcher_0                        = regexp.MustCompile(`(?m)^(%[%{}]xs|<.*>)`)
	l_Groff_Matcher_0                      = regexp.MustCompile(`(?mi)^\.[a-z][a-z](\s|$)`)
	l_PicoLisp_Matcher_0                   = regexp.MustCompile(`(?m)^\((de|class|rel|code|data|must)\s`)
	ls_LoomScript_Matcher_0                = regexp.MustCompile(`(?m)^\s*package\s*[\w\.\/\*\s]*\s*{`)
	lsp_CommonLisp_Matcher_0               = regexp.MustCompile(`(?mi)^\s*\((defun|in-package|defpackage) `)
	lsp_NewLisp_Matcher_0                  = regexp.MustCompile(`(?m)^\s*\(define `)
	lisp_CommonLisp_Matcher_0              = regexp.MustCompile(`(?mi)^\s*\((defun|in-package|defpackage) `)
	lisp_NewLisp_Matcher_0                 = regexp.MustCompile(`(?m)^\s*\(define `)
	m_ObjectiveDashC_Matcher_0             = regexp.MustCompile(`(?m)^\s*(@(interface|class|protocol|property|end|synchronised|selector|implementation)\b|#import\s+.+\.h[">])`)
	m_Mercury_Matcher_0                    = regexp.MustCompile(`(?m):- module`)
	m_MUF_Matcher_0                        = regexp.MustCompile(`(?m)^: `)
	m_M_Matcher_0                          = regexp.MustCompile(`(?m)^\s*;`)
	m_Mathematica_Matcher_0                = regexp.MustCompile(`(?m)\*\)$`)
	m_Matlab_Matcher_0                     = regexp.MustCompile(`(?m)^\s*%`)
	m_Limbo_Matcher_0                      = regexp.MustCompile(`(?m)^\w+\s*:\s*module\s*{`)
	md_Markdown_Matcher_0                  = regexp.MustCompile(`(?mi)(^[-a-z0-9=#!\*\[|>])|<\/`)
	md_Markdown_Matcher_1                  = regexp.MustCompile(`(?m)^$`)
	md_GCCmachinedescription_Matcher_0     = regexp.MustCompile(`(?m)^(;;|\(define_)`)
	ml_OCaml_Matcher_0                     = regexp.MustCompile(`(?m)(^\s*module)|let rec |match\s+(\S+\s)+with`)
	ml_StandardML_Matcher_0                = regexp.MustCompile(`(?m)=> |case\s+(\S+\s)+of`)
	mod_XML_Matcher_0                      = regexp.MustCompile(`(?m)<!ENTITY `)
	mod_ModulaDash2_Matcher_0              = regexp.MustCompile(`(?mi)^\s*MODULE [\w\.]+;`)
	mod_ModulaDash2_Matcher_1              = regexp.MustCompile(`(?mi)^\s*END [\w\.]+;`)
	ms_Groff_Matcher_0                     = regexp.MustCompile(`(?mi)^[.'][a-z][a-z](\s|$)`)
	n_Groff_Matcher_0                      = regexp.MustCompile(`(?m)^[.']`)
	n_Nemerle_Matcher_0                    = regexp.MustCompile(`(?m)^(module|namespace|using)\s`)
	ncl_Text_Matcher_0                     = regexp.MustCompile(`(?m)THE_TITLE`)
	nl_NL_Matcher_0                        = regexp.MustCompile(`(?m)^(b|g)[0-9]+ `)
	php_Hack_Matcher_0                     = regexp.MustCompile(`(?m)<\?hh`)
	php_PHP_Matcher_0                      = regexp.MustCompile(`(?m)<?[^h]`)
	pl_Prolog_Matcher_0                    = regexp.MustCompile(`(?m)^[^#]*:-`)
	pl_Perl_Matcher_0                      = regexp.MustCompile(`(?m)use strict|use\s+v?5\.`)
	pl_Perl6_Matcher_0                     = regexp.MustCompile(`(?m)^(use v6|(my )?class|module)`)
	pm_Perl_Matcher_0                      = regexp.MustCompile(`(?m)use strict|use\s+v?5\.`)
	pm_Perl6_Matcher_0                     = regexp.MustCompile(`(?m)^(use v6|(my )?class|module)`)
	t_Perl_Matcher_0                       = regexp.MustCompile(`(?m)use strict|use\s+v?5\.`)
	t_Perl6_Matcher_0                      = regexp.MustCompile(`(?m)^(use v6|(my )?class|module)`)
	pod_Pod_Matcher_0                      = regexp.MustCompile(`(?m)^=\w+\b`)
	pro_Prolog_Matcher_0                   = regexp.MustCompile(`(?m)^[^#]+:-`)
	pro_INI_Matcher_0                      = regexp.MustCompile(`(?m)last_client=`)
	pro_QMake_Matcher_0                    = regexp.MustCompile(`(?m)HEADERS`)
	pro_QMake_Matcher_1                    = regexp.MustCompile(`(?m)SOURCES`)
	pro_IDL_Matcher_0                      = regexp.MustCompile(`(?m)^\s*function[ \w,]+$`)
	props_XML_Matcher_0                    = regexp.MustCompile(`(?mi)^(\s*)(<Project|<Import|<Property|<?xml|xmlns)`)
	props_INI_Matcher_0                    = regexp.MustCompile(`(?mi)\w+\s*=\s*`)
	r_Rebol_Matcher_0                      = regexp.MustCompile(`(?mi)\bRebol\b`)
	r_R_Matcher_0                          = regexp.MustCompile(`(?m)<-|^\s*#`)
	rno_RUNOFF_Matcher_0                   = regexp.MustCompile(`(?mi)^\.!|^\.end lit(?:eral)?\b`)
	rno_Groff_Matcher_0                    = regexp.MustCompile(`(?m)^\.\\" `)
	rpy_Python_Matcher_0                   = regexp.MustCompile(`(?ms)(^(import|from|class|def)\s)`)
	rs_Rust_Matcher_0                      = regexp.MustCompile(`(?m)^(use |fn |mod |pub |macro_rules|impl|#!?\[)`)
	rs_RenderScript_Matcher_0              = regexp.MustCompile(`(?m)#include|#pragma\s+(rs|version)|__attribute__`)
	sc_SuperCollider_Matcher_0             = regexp.MustCompile(`(?m)\^(this|super)\.`)
	sc_SuperCollider_Matcher_1             = regexp.MustCompile(`(?m)^\s*(\+|\*)\s*\w+\s*{`)
	sc_SuperCollider_Matcher_2             = regexp.MustCompile(`(?m)^\s*~\w+\s*=\.`)
	sc_Scala_Matcher_0                     = regexp.MustCompile(`(?m)^\s*import (scala|java)\.`)
	sc_Scala_Matcher_1                     = regexp.MustCompile(`(?m)^\s*val\s+\w+\s*=`)
	sc_Scala_Matcher_2                     = regexp.MustCompile(`(?m)^\s*class\b`)
	sql_PLpgSQL_Matcher_0                  = regexp.MustCompile(`(?mi)^\\i\b|AS \$\$|LANGUAGE '?plpgsql'?`)
	sql_PLpgSQL_Matcher_1                  = regexp.MustCompile(`(?mi)SECURITY (DEFINER|INVOKER)`)
	sql_PLpgSQL_Matcher_2                  = regexp.MustCompile(`(?mi)BEGIN( WORK| TRANSACTION)?;`)
	sql_SQLPL_Matcher_0                    = regexp.MustCompile(`(?mi)(alter module)|(language sql)|(begin( NOT)+ atomic)`)
	sql_SQLPL_Matcher_1                    = regexp.MustCompile(`(?mi)signal SQLSTATE '[0-9]+'`)
	sql_PLSQL_Matcher_0                    = regexp.MustCompile(`(?mi)\$\$PLSQL_|XMLTYPE|sysdate|systimestamp|\.nextval|connect by|AUTHID (DEFINER|CURRENT_USER)`)
	sql_PLSQL_Matcher_1                    = regexp.MustCompile(`(?mi)constructor\W+function`)
	sql_SQL_Matcher_0                      = regexp.MustCompile(`(?mi)! /begin|boolean|package|exception`)
	srt_SubRipText_Matcher_0               = regexp.MustCompile(`(?m)^(\d{2}:\d{2}:\d{2},\d{3})\s*(-->)\s*(\d{2}:\d{2}:\d{2},\d{3})$`)
	toc_WorldofWarcraftAddonData_Matcher_0 = regexp.MustCompile(`(?m)^## |@no-lib-strip@`)
	toc_TeX_Matcher_0                      = regexp.MustCompile(`(?m)^\\(contentsline|defcounter|beamer|boolfalse)`)
	ts_XML_Matcher_0                       = regexp.MustCompile(`(?m)<TS`)
	tst_GAP_Matcher_0                      = regexp.MustCompile(`(?m)gap> `)
	tsx_TypeScript_Matcher_0               = regexp.MustCompile(`(?m)^\s*(import.+(from\s+|require\()['"]react|\/\/\/\s*<reference\s)`)
	tsx_XML_Matcher_0                      = regexp.MustCompile(`(?mi)^\s*<\?xml\s+version`)
)
