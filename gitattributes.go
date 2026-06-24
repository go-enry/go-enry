package enry

import (
	"bufio"
	"bytes"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-enry/go-enry/v2/data"
)

// GitAttributes holds parsed .gitattributes rules for overriding language detection.
type GitAttributes struct {
	rules []gitAttributeRule
}

type gitAttributeRule struct {
	pattern string
	attrs   map[string]gitAttributeValue
}

type gitAttributeValueKind uint8

const (
	gitAttributeValueSet gitAttributeValueKind = iota
	gitAttributeValueUnset
	gitAttributeValueUnspecified
	gitAttributeValueString
)

type gitAttributeValue struct {
	kind  gitAttributeValueKind
	value string
}

type gitAttributeAssignment struct {
	name        string
	value       gitAttributeValue
	expandMacro bool
}

type gitAttributeRawRule struct {
	pattern     string
	assignments []gitAttributeAssignment
}

var (
	gitAttributeLinguistAliasOnce sync.Once
	gitAttributeLinguistAliasMap  map[string]string
)

// ParseGitAttributes parses the content of a .gitattributes file.
// Each non-comment, non-empty line has the form:
//   - pattern attr1 attr2=value ...
//   - [attr]macroName attr1 attr2=value ...
//
// Attributes can be bare (set), prefixed with "-" (unset), prefixed with "!"
// (reset to unspecified), or assigned via "=" (string value).
// Bare macro references are expanded recursively after parsing, matching Git's
// behavior even when the macro is defined later in the file.
func ParseGitAttributes(content []byte) (GitAttributes, error) {
	var rawRules []gitAttributeRawRule
	macros := make(map[string][]gitAttributeAssignment)
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		if strings.HasPrefix(fields[0], "[attr]") {
			name := strings.TrimPrefix(fields[0], "[attr]")
			if name == "" {
				continue
			}
			macros[name] = parseGitAttributeAssignments(fields[1:])
			continue
		}

		// Negative patterns are ignored in gitattributes (Git warns and skips).
		if strings.HasPrefix(fields[0], "!") {
			continue
		}

		if len(fields) < 2 {
			continue
		}

		rawRules = append(rawRules, gitAttributeRawRule{
			pattern:     fields[0],
			assignments: parseGitAttributeAssignments(fields[1:]),
		})
	}

	if err := scanner.Err(); err != nil {
		return GitAttributes{}, err
	}

	rules := make([]gitAttributeRule, 0, len(rawRules))
	memo := make(map[string]map[string]gitAttributeValue)
	for _, rule := range rawRules {
		rules = append(rules, gitAttributeRule{
			pattern: rule.pattern,
			attrs:   resolveGitAttributeAssignments(rule.assignments, macros, memo, map[string]bool{}),
		})
	}

	return GitAttributes{rules: rules}, nil
}

func parseGitAttributeAssignments(fields []string) []gitAttributeAssignment {
	assignments := make([]gitAttributeAssignment, 0, len(fields))
	for _, field := range fields {
		switch {
		case strings.HasPrefix(field, "!"):
			assignments = append(assignments, gitAttributeAssignment{
				name:  field[1:],
				value: gitAttributeValue{kind: gitAttributeValueUnspecified},
			})
		case strings.HasPrefix(field, "-"):
			assignments = append(assignments, gitAttributeAssignment{
				name:  field[1:],
				value: gitAttributeValue{kind: gitAttributeValueUnset},
			})
		case strings.Contains(field, "="):
			idx := strings.Index(field, "=")
			assignments = append(assignments, gitAttributeAssignment{
				name:  field[:idx],
				value: gitAttributeValue{kind: gitAttributeValueString, value: field[idx+1:]},
			})
		default:
			assignments = append(assignments, gitAttributeAssignment{
				name:        field,
				value:       gitAttributeValue{kind: gitAttributeValueSet},
				expandMacro: true,
			})
		}
	}
	return assignments
}

func resolveGitAttributeAssignments(assignments []gitAttributeAssignment, macros map[string][]gitAttributeAssignment, memo map[string]map[string]gitAttributeValue, visiting map[string]bool) map[string]gitAttributeValue {
	resolved := make(map[string]gitAttributeValue)
	for _, assignment := range assignments {
		if assignment.expandMacro {
			if expanded, ok := resolveGitAttributeMacro(assignment.name, macros, memo, visiting); ok {
				for name, value := range expanded {
					resolved[name] = value
				}
				continue
			}
		}
		resolved[assignment.name] = assignment.value
	}
	return resolved
}

func resolveGitAttributeMacro(name string, macros map[string][]gitAttributeAssignment, memo map[string]map[string]gitAttributeValue, visiting map[string]bool) (map[string]gitAttributeValue, bool) {
	if expanded, ok := memo[name]; ok {
		return cloneGitAttributeValues(expanded), true
	}

	assignments, ok := macros[name]
	if !ok {
		return nil, false
	}

	if visiting[name] {
		return map[string]gitAttributeValue{}, true
	}

	visiting[name] = true
	expanded := resolveGitAttributeAssignments(assignments, macros, memo, visiting)
	delete(visiting, name)

	memo[name] = cloneGitAttributeValues(expanded)
	return expanded, true
}

func cloneGitAttributeValues(attrs map[string]gitAttributeValue) map[string]gitAttributeValue {
	cloned := make(map[string]gitAttributeValue, len(attrs))
	for name, value := range attrs {
		cloned[name] = value
	}
	return cloned
}

func (v gitAttributeValue) boolValue() bool {
	switch v.kind {
	case gitAttributeValueUnset:
		return false
	case gitAttributeValueString:
		return v.value != "false"
	default:
		return true
	}
}

func (v gitAttributeValue) textValue() string {
	switch v.kind {
	case gitAttributeValueSet:
		return "true"
	case gitAttributeValueUnset:
		return "false"
	default:
		return v.value
	}
}

// getAttr returns the value of a named attribute for the given path.
// It returns the value from the last matching rule (git precedence: later rules win).
// Iterates in reverse so we can return immediately on the first (i.e. last-defined) match.
// If the attribute is unspecified (set via !attr or via a macro that resolves to
// !attr), it is treated as not set, causing the caller to fall back to default
// detection.
func (ga GitAttributes) getAttr(path string, attr string) (gitAttributeValue, bool) {
	for i := len(ga.rules) - 1; i >= 0; i-- {
		rule := ga.rules[i]
		if matchGitPattern(rule.pattern, path) {
			if val, ok := rule.attrs[attr]; ok {
				if val.kind == gitAttributeValueUnspecified {
					return gitAttributeValue{}, false
				}
				return val, true
			}
		}
	}
	return gitAttributeValue{}, false
}

// IsVendor checks the linguist-vendored attribute for path.
// If no rule matches, it falls back to the default IsVendor detection.
func (ga GitAttributes) IsVendor(path string) bool {
	if val, ok := ga.getAttr(path, "linguist-vendored"); ok {
		return val.boolValue()
	}
	return IsVendor(path)
}

// IsDocumentation checks the linguist-documentation attribute for path.
// If no rule matches, it falls back to the default IsDocumentation detection.
func (ga GitAttributes) IsDocumentation(path string) bool {
	if val, ok := ga.getAttr(path, "linguist-documentation"); ok {
		return val.boolValue()
	}
	return IsDocumentation(path)
}

// IsGenerated checks the linguist-generated attribute for path.
// If no rule matches, it falls back to the default IsGenerated detection.
func (ga GitAttributes) IsGenerated(path string, content []byte) bool {
	if val, ok := ga.getAttr(path, "linguist-generated"); ok {
		return val.boolValue()
	}
	return IsGenerated(path, content)
}

// IsDetectable checks the linguist-detectable attribute for path.
// Returns (value, hasOverride). When hasOverride is true, value indicates
// whether the file should be included in language statistics regardless of
// its language type. When hasOverride is false, the caller should fall back
// to default behavior (include programming/markup, exclude data/prose).
//
// Per Linguist semantics:
//   - "linguist-detectable" forces inclusion (e.g., data/prose languages in stats)
//   - "-linguist-detectable" forces exclusion (e.g., hide a programming language)
//   - No rule means use default language-type-based detection
func (ga GitAttributes) IsDetectable(path string) (bool, bool) {
	if val, ok := ga.getAttr(path, "linguist-detectable"); ok {
		return val.boolValue(), true
	}
	return false, false
}

// getLanguageByGitAttributeAlias resolves linguist-language values using the
// current enry alias normalization first, and then Linguist's hyphen alias convention.
// TODO(#200): remove the second lookup once GetLanguageByAlias accepts
// Linguist-style aliases directly.
func getLanguageByGitAttributeAlias(text string) (string, bool) {
	if lang, ok := GetLanguageByAlias(text); ok {
		return lang, true
	}

	gitAttributeLinguistAliasOnce.Do(buildGitAttributeLinguistAliasMap)
	lang, ok := gitAttributeLinguistAliasMap[gitAttributeLinguistAlias(text)]
	return lang, ok
}

// buildGitAttributeLinguistAliasMap builds the second alias lookup table used only
// for Linguist's alias normalization with hyphen https://github.com/github-linguist/linguist/blob/2409807814a3ff386294b1f217b886a1594642cd/docs/overrides.md#using-gitattributes .
func buildGitAttributeLinguistAliasMap() {
	aliases := make(map[string]string, len(data.LanguageInfoByID)*2)
	for _, info := range data.LanguageInfoByID {
		addGitAttributeLinguistAlias(aliases, info.Name, info.Name)
		for _, alias := range info.Aliases {
			addGitAttributeLinguistAlias(aliases, alias, info.Name)
		}
	}
	gitAttributeLinguistAliasMap = aliases
}

// addGitAttributeLinguistAlias adds an alias unless it already belongs to a different language.
func addGitAttributeLinguistAlias(aliases map[string]string, alias string, lang string) {
	key := gitAttributeLinguistAlias(alias)
	if existing, exists := aliases[key]; !exists || existing == lang {
		aliases[key] = lang
	}
}

func gitAttributeLinguistAlias(text string) string {
	text = strings.ReplaceAll(text, " ", "-")
	return strings.ToLower(text)
}

// GetLanguage checks the linguist-language attribute for path.
// If set, the language name is resolved via current enry aliases first, then
// a scoped Linguist-style .gitattributes alias fallback for canonical naming.
// Returns the language and true if an override was found.
func (ga GitAttributes) GetLanguage(path string) (string, bool) {
	if val, ok := ga.getAttr(path, "linguist-language"); ok {
		text := val.textValue()
		// TODO(#200): replace with GetLanguageByAlias once alias normalization is changed
		if lang, ok := getLanguageByGitAttributeAlias(text); ok {
			return lang, true
		}
		// Return as-is if alias resolution fails (could be a valid language name)
		return text, true
	}
	return "", false
}

// HasPotentialOverride reports whether attr is ever explicitly unset or reset
// in .gitattributes. Directory-level pruning in callers must be conservative
// when such rules exist, because a later descendant rule may need to override
// the parent directory classification.
func (ga GitAttributes) HasPotentialOverride(attr string) bool {
	for _, rule := range ga.rules {
		if val, ok := rule.attrs[attr]; ok && (val.kind == gitAttributeValueUnset || val.kind == gitAttributeValueUnspecified) {
			return true
		}
	}
	return false
}

// matchGitPattern implements git-style glob matching compatible with
// linguist's fnmatch(FNM_PATHNAME):
//   - "*" matches anything except "/"
//   - "**" matches everything including "/"
//   - "?" matches single char except "/"
//   - Leading "/" anchors to the repo root; otherwise matches basename
//     or anywhere if the pattern contains "/"
//   - Trailing "/" matches the directory itself and all files inside it
func matchGitPattern(pattern, path string) bool {
	// Trailing slash means directory-only pattern. Unlike .gitignore,
	// .gitattributes does NOT recurse into directories — the pattern only
	// matches paths that are themselves directories (ending with "/").
	// The Git docs call trailing "/" in .gitattributes "pointless";
	// users should use "vendor/**" instead of "vendor/".
	if strings.HasSuffix(pattern, "/") {
		if !strings.HasSuffix(path, "/") {
			return false
		}
		pattern = strings.TrimSuffix(pattern, "/")
		path = strings.TrimSuffix(path, "/")
	}

	anchored := strings.HasPrefix(pattern, "/")
	if anchored {
		pattern = pattern[1:]
	}

	// If pattern contains no slash (and wasn't anchored), match against basename
	if !anchored && !strings.Contains(pattern, "/") {
		path = filepath.Base(path)
	}

	return globMatch(pattern, path)
}

// globMatch performs recursive glob matching with support for *, **, ?, and \.
func globMatch(pattern, str string) bool {
	return globMatchInternal(pattern, str, 0)
}

// prev tracks the raw pattern byte immediately before the current pattern
// position. It is 0 at the start of a pattern context. This is used for
// the ** boundary check (Git wildmatch.c:95 "prev_p - pattern < 2").
func globMatchInternal(pattern, str string, prev byte) bool {
	for len(pattern) > 0 {
		switch pattern[0] {
		case '*':
			if len(pattern) > 1 && pattern[1] == '*' {
				// Check if ** is at a path boundary:
				// preceded by start-of-pattern or '/'
				// AND followed by end-of-pattern, '/', or '\/'
				atBoundary := (prev == 0 || prev == '/') &&
					(len(pattern) == 2 || pattern[2] == '/' ||
						(len(pattern) > 3 && pattern[2] == '\\' && pattern[3] == '/'))
				if atBoundary {
					rest := pattern[2:]
					restPrev := byte('*')
					if len(rest) > 0 && rest[0] == '/' {
						restPrev = '/'
						rest = rest[1:]
					}
					if len(rest) == 0 {
						return true
					}
					for i := 0; i <= len(str); i++ {
						if i == 0 || str[i-1] == '/' {
							if globMatchInternal(rest, str[i:], restPrev) {
								return true
							}
						}
					}
					return false
				}
				// ** not at boundary: skip extra * and treat as single *
				pattern = pattern[1:]
			}
			// "*" matches anything except "/"
			rest := pattern[1:]
			for i := 0; i <= len(str); i++ {
				if i > 0 && str[i-1] == '/' {
					break
				}
				if globMatchInternal(rest, str[i:], '*') {
					return true
				}
			}
			return false
		case '?':
			if len(str) == 0 || str[0] == '/' {
				return false
			}
			prev = '?'
			pattern = pattern[1:]
			str = str[1:]
		case '[':
			// Character class matching
			if len(str) == 0 || str[0] == '/' {
				return false
			}
			// Find closing bracket, skipping POSIX class [:..:]
			end := findClosingBracket(pattern)
			if end == -1 {
				// No closing bracket, treat as literal
				if str[0] != pattern[0] {
					return false
				}
				prev = pattern[0]
				pattern = pattern[1:]
				str = str[1:]
				continue
			}
			class := pattern[1:end]
			negate := false
			if len(class) > 0 && (class[0] == '!' || class[0] == '^') {
				negate = true
				class = class[1:]
			}
			matched, valid := matchCharClass(class, str[0])
			if !valid {
				return false
			}
			if negate {
				matched = !matched
			}
			if !matched {
				return false
			}
			prev = ']'
			pattern = pattern[end+1:]
			str = str[1:]
		case '\\':
			pattern = pattern[1:]
			if len(pattern) == 0 {
				return false
			}
			if len(str) == 0 || pattern[0] != str[0] {
				return false
			}
			prev = pattern[0]
			pattern = pattern[1:]
			str = str[1:]
		default:
			if len(str) == 0 || pattern[0] != str[0] {
				return false
			}
			prev = pattern[0]
			pattern = pattern[1:]
			str = str[1:]
		}
	}
	return len(str) == 0
}

// findClosingBracket returns the index of the closing ']' in a bracket
// expression starting at pattern[0] == '[', skipping POSIX classes [:name:].
// Per POSIX/Git wildmatch: ']' immediately after '[' (or '[!' / '[^') is a
// literal class member, not the closing bracket. Backslash escapes are also
// handled: '\]' inside a bracket is a literal ']'.
func findClosingBracket(pattern string) int {
	i := 1
	if i < len(pattern) && (pattern[i] == '!' || pattern[i] == '^') {
		i++
	}
	// ] immediately after [ (or [!/[^) is a literal class member
	if i < len(pattern) && pattern[i] == ']' {
		i++
	}
	for i < len(pattern) {
		if pattern[i] == '\\' && i+1 < len(pattern) {
			i += 2 // skip escaped character
			continue
		}
		if pattern[i] == '[' && i+1 < len(pattern) && pattern[i+1] == ':' {
			if end := strings.Index(pattern[i+2:], ":]"); end != -1 {
				i += end + 4
				continue
			}
		}
		if pattern[i] == ']' {
			return i
		}
		i++
	}
	return -1
}

// matchCharClass returns (matched, valid). valid is false if the class
// contains an unknown POSIX class name, in which case the caller should
// treat the entire pattern as non-matching. Backslash escapes are handled:
// '\x' inside a class is treated as literal 'x'.
func matchCharClass(class string, ch byte) (bool, bool) {
	if !hasValidPOSIXClasses(class) {
		return false, false
	}
	for len(class) > 0 {
		if strings.HasPrefix(class, "[:") {
			if end := strings.Index(class[2:], ":]"); end != -1 {
				m, _ := matchPOSIXClass(class[2:2+end], ch)
				if m {
					return true, true
				}
				class = class[2+end+2:]
				continue
			}
		}
		// Get current character, handling escape
		cur := class[0]
		advance := 1
		if cur == '\\' && len(class) > 1 {
			cur = class[1]
			advance = 2
		}
		// Check for range: cur-high
		rest := class[advance:]
		if len(rest) >= 2 && rest[0] == '-' && rest[1] != ']' {
			high := rest[1]
			highAdv := 2
			if high == '\\' && len(rest) > 2 {
				high = rest[2]
				highAdv = 3
			}
			if ch >= cur && ch <= high {
				return true, true
			}
			class = rest[highAdv:]
			continue
		}
		if cur == ch {
			return true, true
		}
		class = rest
	}
	return false, true
}

// hasValidPOSIXClasses checks that every [:name:] in the class string
// refers to a known POSIX class.
func hasValidPOSIXClasses(class string) bool {
	for {
		i := strings.Index(class, "[:")
		if i == -1 {
			return true
		}
		end := strings.Index(class[i+2:], ":]")
		if end == -1 {
			return true
		}
		if _, valid := matchPOSIXClass(class[i+2:i+2+end], 0); !valid {
			return false
		}
		class = class[i+2+end+2:]
	}
}

// matchPOSIXClass returns (matched, valid). valid is false for unknown class names.
func matchPOSIXClass(name string, ch byte) (bool, bool) {
	var m bool
	switch name {
	case "alnum":
		m = (ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')
	case "alpha":
		m = (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')
	case "blank":
		m = ch == ' ' || ch == '\t'
	case "cntrl":
		m = ch < 0x20 || ch == 0x7f
	case "digit":
		m = ch >= '0' && ch <= '9'
	case "graph":
		m = ch >= 0x21 && ch <= 0x7e
	case "lower":
		m = ch >= 'a' && ch <= 'z'
	case "print":
		m = ch >= 0x20 && ch <= 0x7e
	case "punct":
		m = ch >= 0x21 && ch <= 0x7e &&
			!((ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z'))
	case "space":
		m = ch == '\t' || ch == '\n' || ch == '\v' || ch == '\f' || ch == '\r' || ch == ' '
	case "upper":
		m = ch >= 'A' && ch <= 'Z'
	case "xdigit":
		m = (ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'F') || (ch >= 'a' && ch <= 'f')
	default:
		return false, false
	}
	return m, true
}
