package slinguist

import (
	"bytes"
	"regexp"
)

// GetLanguageByModeline returns the language of the given content looking for the modeline,
// and safe to indicate the sureness of returned language.
func GetLanguageByModeline(content []byte) (lang string, safe bool) {
	headFoot := getHeaderAndFooter(content)
	for _, getLang := range modelinesFunc {
		lang = getLang(headFoot)
		safe = lang != OtherLanguage
		if safe {
			return
		}
	}

	return
}

func getHeaderAndFooter(content []byte) []byte {
	const (
		searchScope = 5
		eol         = `\n`
	)

	if bytes.Count(content, []byte(eol)) < 2*searchScope {
		return content
	}

	splitted := bytes.Split(content, []byte(eol))
	header := splitted[:searchScope]
	footer := splitted[len(splitted)-searchScope:]
	headerAndFooter := append(header, footer...)
	return bytes.Join(headerAndFooter, []byte(eol))
}

var modelinesFunc = []func(content []byte) string{
	GetLanguageByEmacsModeline,
	GetLanguageByVimModeline,
}

var (
	reEmacsModeline = regexp.MustCompile(`.*-\*-\s*(.+?)\s*-\*-.*(?m:$)`)
	reEmacsLang     = regexp.MustCompile(`.*(?i:mode)\s*:\s*([^\s;]+)\s*;*.*`)
	reVimModeline   = regexp.MustCompile(`(?:(?m:\s|^)vi(?:m[<=>]?\d+|m)?|[\t\x20]*ex)\s*[:]\s*(.*)(?m:$)`)
	reVimLang       = regexp.MustCompile(`(?i:filetype|ft|syntax)\s*=(\w+)(?:\s|:|$)`)
)

// GetLanguageByEmacsModeline detecs if the content has a emacs modeline and try to get a
// language basing on alias. If couldn't retrieve a valid language, it returns OtherLanguage.
func GetLanguageByEmacsModeline(content []byte) (lang string) {
	matched := reEmacsModeline.FindAllSubmatch(content, -1)
	if matched == nil {
		return OtherLanguage
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

	lang = GetLanguageByAlias(alias)
	return
}

// GetLanguageByVimModeline detecs if the content has a vim modeline and try to get a
// language basing on alias. If couldn't retrieve a valid language, it returns OtherLanguage.
func GetLanguageByVimModeline(content []byte) (lang string) {
	matched := reVimModeline.FindAllSubmatch(content, -1)
	if matched == nil {
		return OtherLanguage
	}

	// only take the last matched line, discard previous lines
	lastLineMatched := matched[len(matched)-1][1]
	matchedAlias := reVimLang.FindAllSubmatch(lastLineMatched, -1)
	if matchedAlias == nil {
		return OtherLanguage
	}

	alias := string(matchedAlias[0][1])
	if len(matchedAlias) > 1 {
		// cases:
		// matchedAlias = [["syntax=ruby " "ruby"] ["ft=python " "python"] ["filetype=perl " "perl"]] returns OtherLanguage;
		// matchedAlias = [["syntax=python " "python"] ["ft=python " "python"] ["filetype=python " "python"]] returns "Python";
		for _, match := range matchedAlias {
			otherAlias := string(match[1])
			if otherAlias != alias {
				alias = OtherLanguage
				break
			}
		}
	}

	lang = GetLanguageByAlias(alias)
	return
}
