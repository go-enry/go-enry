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
	attrs   map[string]string
}

// ParseGitAttributes parses the content of a .gitattributes file.
// Each non-comment, non-empty line has the form: pattern attr1 attr2=value ...
// Attributes can be:
//   - "linguist-vendored" (set/true), "-linguist-vendored" (unset/false)
//   - "linguist-language=Go"
//   - etc.
func ParseGitAttributes(content []byte) (GitAttributes, error) {
	var rules []gitAttributeRule
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		pattern := fields[0]
		attrs := make(map[string]string)
		for _, field := range fields[1:] {
			if strings.HasPrefix(field, "!") {
				// !attr means unspecified (reset to default)
				attrs[field[1:]] = "unspecified"
			} else if strings.HasPrefix(field, "-") {
				// -attr means unset (false)
				attrs[field[1:]] = "false"
			} else if idx := strings.Index(field, "="); idx != -1 {
				// attr=value
				attrs[field[:idx]] = field[idx+1:]
			} else {
				// attr alone means set (true)
				attrs[field] = "true"
			}
		}

		rules = append(rules, gitAttributeRule{pattern: pattern, attrs: attrs})
	}

	if err := scanner.Err(); err != nil {
		return GitAttributes{}, err
	}

	return GitAttributes{rules: rules}, nil
}

// getAttr returns the value of a named attribute for the given path.
// It returns the value from the last matching rule (git precedence: later rules win).
// Iterates in reverse so we can return immediately on the first (i.e. last-defined) match.
// If the attribute is "unspecified" (set via !attr), it is treated as not set,
// causing the caller to fall back to default detection.
func (ga GitAttributes) getAttr(path string, attr string) (string, bool) {
	for i := len(ga.rules) - 1; i >= 0; i-- {
		rule := ga.rules[i]
		if matchGitPattern(rule.pattern, path) {
			if val, ok := rule.attrs[attr]; ok {
				if val == "unspecified" {
					return "", false
				}
				return val, true
			}
		}
	}
	return "", false
}

// IsVendor checks the linguist-vendored attribute for path.
// If no rule matches, it falls back to the default IsVendor detection.
func (ga GitAttributes) IsVendor(path string) bool {
	if val, ok := ga.getAttr(path, "linguist-vendored"); ok {
		return val != "false"
	}
	return IsVendor(path)
}

// IsDocumentation checks the linguist-documentation attribute for path.
// If no rule matches, it falls back to the default IsDocumentation detection.
func (ga GitAttributes) IsDocumentation(path string) bool {
	if val, ok := ga.getAttr(path, "linguist-documentation"); ok {
		return val != "false"
	}
	return IsDocumentation(path)
}

// IsGenerated checks the linguist-generated attribute for path.
// If no rule matches, it falls back to the default IsGenerated detection.
func (ga GitAttributes) IsGenerated(path string, content []byte) bool {
	if val, ok := ga.getAttr(path, "linguist-generated"); ok {
		return val != "false"
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
		return val != "false", true
	}
	return false, false
}

// GetLanguage checks the linguist-language attribute for path.
// If set, the language name is resolved via GetLanguageByAlias for canonical naming.
// Returns the language and true if an override was found.
func (ga GitAttributes) GetLanguage(path string) (string, bool) {
	if val, ok := ga.getAttr(path, "linguist-language"); ok {
		if lang, ok := GetLanguageByAlias(val); ok {
			return lang, true
		}
		// Return as-is if alias resolution fails (could be a valid language name)
		return val, true
	}
	return "", false
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
	// Trailing slash means directory pattern: matches the directory path
	// (ending with "/") or any file inside the directory.
	if strings.HasSuffix(pattern, "/") {
		dir := strings.TrimSuffix(pattern, "/")
		if strings.HasSuffix(path, "/") {
			// Directory path itself — match against dir name
			path = strings.TrimSuffix(path, "/")
			pattern = dir
		} else if strings.HasPrefix(path, dir+"/") {
			// File inside the directory — matches
			return true
		} else {
			return false
		}
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
