package enry

import (
	"bytes"
	"regexp"
)

const (
	searchScope = 5
)

// GetLanguagesByModeline returns a slice of possible languages for the given content, filename will be ignored.
// It accomplish the signature to be a Strategy type.
func GetLanguagesByModeline(filename string, content []byte) []string {
	headFoot := getHeaderAndFooter(content)
	var languages []string
	for _, getLang := range modelinesFunc {
		languages = getLang("", headFoot)
		if len(languages) > 0 {
			break
		}
	}

	return languages
}

func getHeaderAndFooter(content []byte) []byte {
	if bytes.Count(content, []byte("\n")) < 2*searchScope {
		return content
	}

	header := headScope(content, searchScope)
	footer := footScope(content, searchScope)
	headerAndFooter := make([]byte, 0, len(content[:header])+len(content[footer:]))
	headerAndFooter = append(headerAndFooter, content[:header]...)
	headerAndFooter = append(headerAndFooter, content[footer:]...)
	return headerAndFooter
}

func headScope(content []byte, scope int) (index int) {
	for i := 0; i < scope; i++ {
		eol := bytes.IndexAny(content, "\n")
		content = content[eol+1:]
		index += eol
	}

	return index + scope - 1
}

func footScope(content []byte, scope int) (index int) {
	for i := 0; i < scope; i++ {
		index = bytes.LastIndexAny(content, "\n")
		content = content[:index]
	}

	return index + 1
}

var modelinesFunc = []func(filename string, content []byte) []string{
	GetLanguagesByEmacsModeline,
	GetLanguagesByVimModeline,
}

var (
	reEmacsModeline = regexp.MustCompile(`.*-\*-\s*(.+?)\s*-\*-.*(?m:$)`)
	reEmacsLang     = regexp.MustCompile(`.*(?i:mode)\s*:\s*([^\s;]+)\s*;*.*`)
	reVimModeline   = regexp.MustCompile(`(?:(?m:\s|^)vi(?:m[<=>]?\d+|m)?|[\t\x20]*ex)\s*[:]\s*(.*)(?m:$)`)
	reVimLang       = regexp.MustCompile(`(?i:filetype|ft|syntax)\s*=(\w+)(?:\s|:|$)`)
)

// GetLanguageByEmacsModeline detecs if the content has a emacs modeline and try to get a
// language basing on alias. If couldn't retrieve a valid language, it returns OtherLanguage and false.
func GetLanguageByEmacsModeline(content []byte) (string, bool) {
	languages := GetLanguagesByEmacsModeline("", content)
	if len(languages) == 0 {
		return OtherLanguage, false
	}

	return languages[0], true
}

// GetLanguagesByEmacsModeline returns a slice of possible languages for the given content, filename will be ignored.
// It accomplish the signature to be a Strategy type.
func GetLanguagesByEmacsModeline(filename string, content []byte) []string {
	matched := reEmacsModeline.FindAllSubmatch(content, -1)
	if matched == nil {
		return nil
	}

	// only take the last matched line, discard previous lines
	lastLineMatched := matched[len(matched)-1][1]
	matchedAlias := reEmacsLang.FindSubmatch(lastLineMatched)
	var alias string
	if matchedAlias != nil {
		alias = string(matchedAlias[1])
	} else {
		alias = string(lastLineMatched)
	}

	language, ok := GetLanguageByAlias(alias)
	if !ok {
		return nil
	}

	return []string{language}
}

// GetLanguageByVimModeline detecs if the content has a vim modeline and try to get a
// language basing on alias. If couldn't retrieve a valid language, it returns OtherLanguage and false.
func GetLanguageByVimModeline(content []byte) (string, bool) {
	languages := GetLanguagesByVimModeline("", content)
	if len(languages) == 0 {
		return OtherLanguage, false
	}

	return languages[0], true
}

// GetLanguagesByVimModeline returns a slice of possible languages for the given content, filename will be ignored.
// It accomplish the signature to be a Strategy type.
func GetLanguagesByVimModeline(filename string, content []byte) []string {
	matched := reVimModeline.FindAllSubmatch(content, -1)
	if matched == nil {
		return nil
	}

	// only take the last matched line, discard previous lines
	lastLineMatched := matched[len(matched)-1][1]
	matchedAlias := reVimLang.FindAllSubmatch(lastLineMatched, -1)
	if matchedAlias == nil {
		return nil
	}

	alias := string(matchedAlias[0][1])
	if len(matchedAlias) > 1 {
		// cases:
		// matchedAlias = [["syntax=ruby " "ruby"] ["ft=python " "python"] ["filetype=perl " "perl"]] returns OtherLanguage;
		// matchedAlias = [["syntax=python " "python"] ["ft=python " "python"] ["filetype=python " "python"]] returns "Python";
		for _, match := range matchedAlias {
			otherAlias := string(match[1])
			if otherAlias != alias {
				return nil
			}
		}
	}

	language, ok := GetLanguageByAlias(alias)
	if !ok {
		return nil
	}

	return []string{language}
}
