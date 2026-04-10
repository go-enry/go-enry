package enry

import (
	"bufio"
	"bytes"
	"path/filepath"
	"strings"
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

// GetLanguage checks the linguist-language attribute for path.
// If set, the language name is resolved via GetLanguageByAlias for canonical naming.
// Returns the language and true if an override was found.
func (ga GitAttributes) GetLanguage(path string) (string, bool) {
	if val, ok := ga.getAttr(path, "linguist-language"); ok {
		if lang, ok := GetLanguageByAlias(val.textValue()); ok {
			return lang, true
		}
		// Return as-is if alias resolution fails (could be a valid language name)
		return val.textValue(), true
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

// globMatch performs recursive glob matching with support for *, **, and ?.
func globMatch(pattern, str string) bool {
	for len(pattern) > 0 {
		switch pattern[0] {
		case '*':
			if len(pattern) > 1 && pattern[1] == '*' {
				// "**" matches everything including "/"
				rest := pattern[2:]
				// Consume optional trailing slash after **
				if len(rest) > 0 && rest[0] == '/' {
					rest = rest[1:]
				}
				if len(rest) == 0 {
					return true
				}
				// Try matching rest at path-component boundaries only
				for i := 0; i <= len(str); i++ {
					if i == 0 || str[i-1] == '/' {
						if globMatch(rest, str[i:]) {
							return true
						}
					}
				}
				return false
			}
			// "*" matches anything except "/"
			rest := pattern[1:]
			for i := 0; i <= len(str); i++ {
				if i > 0 && str[i-1] == '/' {
					break
				}
				if globMatch(rest, str[i:]) {
					return true
				}
			}
			return false
		case '?':
			if len(str) == 0 || str[0] == '/' {
				return false
			}
			pattern = pattern[1:]
			str = str[1:]
		case '[':
			// Character class matching
			if len(str) == 0 || str[0] == '/' {
				return false
			}
			// Find closing bracket
			end := strings.Index(pattern, "]")
			if end == -1 {
				// No closing bracket, treat as literal
				if str[0] != pattern[0] {
					return false
				}
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
			matched := matchCharClass(class, str[0])
			if negate {
				matched = !matched
			}
			if !matched {
				return false
			}
			pattern = pattern[end+1:]
			str = str[1:]
		default:
			if len(str) == 0 || pattern[0] != str[0] {
				return false
			}
			pattern = pattern[1:]
			str = str[1:]
		}
	}
	return len(str) == 0
}

func matchCharClass(class string, ch byte) bool {
	for i := 0; i < len(class); i++ {
		if i+2 < len(class) && class[i+1] == '-' {
			if ch >= class[i] && ch <= class[i+2] {
				return true
			}
			i += 2
			continue
		}
		if class[i] == ch {
			return true
		}
	}
	return false
}
