package generator

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"sort"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

// Vendor generates regex matchers in Go for vendoring files/dirs.
// It is of generator.File type.
func Vendor(fileToParse, samplesDir, outPath, tmplPath, tmplName, commit string) error {
	data, err := ioutil.ReadFile(fileToParse)
	if err != nil {
		return err
	}

	var regexps []string
	if err := yaml.Unmarshal(data, &regexps); err != nil {
		return fmt.Errorf("failed to parse YAML %s, %q", fileToParse, err)
	}

	for _, re := range regexps {
		if !isRE2(re) {
			log.Printf("RE2 incompatible syntax for vendor:'%s'\n", re)
		}
	}

	buf := &bytes.Buffer{}
	if err := executeVendorTemplate(buf, regexps, tmplPath, tmplName, commit); err != nil {
		return err
	}

	return formatedWrite(outPath, buf.Bytes())
}

func executeVendorTemplate(out io.Writer, regexps []string, tmplPath, tmplName, commit string) error {
	funcs := template.FuncMap{"collateAllRegexps": collateAllRegexps}
	return executeTemplate(out, tmplName, tmplPath, commit, funcs, regexps)
}

// collateAllRegexps all regexps to a single large regexp.
func collateAllRegexps(regexps []string) string {
	// which is at least twice as fast to test than simply iterating & matching.
	//
	// Imperical observation: by looking at the regexps, we only have 3 types.
	//  1. Those that start with `^`
	//  2. Those that start with `(^|/)`
	//  3. All the rest
	//
	// If we collate our regexps into these 3 groups - that will significantly
	// reduce the likelihood of backtracking within the regexp trie matcher.
	//
	// A further improvement is to use non-capturing groups (?:) as otherwise
	// the regexp parser, whilst matching, will have to allocate slices for
	// matching positions. (A future improvement left out could be to
	// enforce non-capturing groups within the sub-regexps.)
	const (
		caret        = "^"
		caretOrSlash = "(^|/)"
	)

	sort.Strings(regexps)

	// Check prefix, group expressions
	var caretPrefixed, caretOrSlashPrefixed, theRest []string
	for _, re := range regexps {
		if strings.HasPrefix(re, caret) {
			caretPrefixed = append(caretPrefixed, re[len(caret):])
		} else if strings.HasPrefix(re, caretOrSlash) {
			caretOrSlashPrefixed = append(caretOrSlashPrefixed, re[len(caretOrSlash):])
		} else {
			theRest = append(theRest, re)
		}
	}

	var sb strings.Builder
	appendGroupWithCommonPrefix(&sb, "^", caretPrefixed)
	sb.WriteString("|")

	appendGroupWithCommonPrefix(&sb, "(?:^|/)", caretOrSlashPrefixed)
	sb.WriteString("|")

	appendGroupWithCommonPrefix(&sb, "", theRest)
	return sb.String()
}

func appendGroupWithCommonPrefix(sb *strings.Builder, commonPrefix string, res []string) {
	sb.WriteString("(?:")
	if commonPrefix != "" {
		sb.WriteString(fmt.Sprintf("%s(?:(?:", commonPrefix))
	}
	sb.WriteString(strings.Join(res, ")|(?:"))
	if commonPrefix != "" {
		sb.WriteString("))")
	}
	sb.WriteString(")")
}
