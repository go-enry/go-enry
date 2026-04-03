package enry

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// EC-1: matchGitPattern — path-component boundary for **
// ---------------------------------------------------------------------------

// TestEdgeDoubleStarPathBoundary verifies that ** only matches at path
// component boundaries, not in the middle of a filename.
func TestEdgeDoubleStarPathBoundary(t *testing.T) {
	tests := []struct {
		pattern string
		path    string
		want    bool
		note    string
	}{
		// ** should NOT match inside a filename component
		{"**/foo.go", "barfoo.go", false, "** must not mid-component match: barfoo.go contains 'foo.go' at offset 3"},
		{"**/foo.go", "xfoo.go", false, "** must not mid-component match: xfoo.go"},
		{"**/test", "atest", false, "** must not match 'atest' for pattern **/test"},
		{"**/vendor/**", "bvendor/lib.go", false, "** must not match 'bvendor' as 'vendor'"},
		{"**/vendor/**", "notvendor/x.go", false, "** must not match 'notvendor' prefix"},

		// ** SHOULD match at component boundaries
		{"**/foo.go", "foo.go", true, "** matches at root"},
		{"**/foo.go", "src/foo.go", true, "** matches one level deep"},
		{"**/foo.go", "src/lib/foo.go", true, "** matches two levels deep"},
		{"**/vendor/**", "vendor/foo.go", true, "** matches vendor at root"},
		{"**/vendor/**", "a/vendor/foo.go", true, "** matches vendor one level deep"},
		{"**/vendor/**", "a/b/vendor/foo.go", true, "** matches vendor two levels deep"},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"__"+tt.path, func(t *testing.T) {
			got := matchGitPattern(tt.pattern, tt.path)
			assert.Equal(t, tt.want, got, "%s: matchGitPattern(%q, %q)", tt.note, tt.pattern, tt.path)
		})
	}
}

// ---------------------------------------------------------------------------
// EC-2: Trailing-slash patterns — directory contents
// ---------------------------------------------------------------------------

// TestEdgeTrailingSlashDirectory verifies that a trailing-slash pattern
// (e.g. "vendor/") matches files INSIDE that directory, not just the literal
// path ending in "/".
func TestEdgeTrailingSlashDirectory(t *testing.T) {
	tests := []struct {
		pattern string
		path    string
		want    bool
		note    string
	}{
		// Files inside vendor/ should match vendor/ pattern
		{"vendor/", "vendor/foo.go", true, "files in vendor/ should be matched"},
		{"vendor/", "vendor/lib/bar.go", true, "nested files in vendor/ should be matched"},
		{"vendor/", "vendor/", true, "the directory itself should match"},

		// Files outside should not match
		{"vendor/", "src/vendor/foo.go", false, "files in src/vendor/ should not match anchored vendor/"},
		{"vendor/", "notvendor/foo.go", false, "files in notvendor/ should not match vendor/"},
		{"vendor/", "vendor", false, "bare 'vendor' without trailing slash should not match"},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"__"+tt.path, func(t *testing.T) {
			got := matchGitPattern(tt.pattern, tt.path)
			assert.Equal(t, tt.want, got, "%s: matchGitPattern(%q, %q)", tt.note, tt.pattern, tt.path)
		})
	}
}

// ---------------------------------------------------------------------------
// EC-3: Double-star in middle of pattern
// ---------------------------------------------------------------------------

func TestEdgeDoubleStarMiddle(t *testing.T) {
	tests := []struct {
		pattern string
		path    string
		want    bool
	}{
		{"src/**/test.go", "src/test.go", true},       // zero dirs
		{"src/**/test.go", "src/a/test.go", true},     // one dir
		{"src/**/test.go", "src/a/b/test.go", true},   // two dirs
		{"src/**/test.go", "src/a/b/c/test.go", true}, // three dirs
		{"src/**/test.go", "test.go", false},          // missing src prefix
		{"src/**/test.go", "other/a/test.go", false},  // wrong top dir
		{"src/**/*.go", "src/foo.go", true},           // zero dirs
		{"src/**/*.go", "src/pkg/foo.go", true},       // one dir
		{"src/**/*.go", "src/pkg/sub/foo.go", true},   // nested
		{"src/**/*.go", "pkg/foo.go", false},          // missing src
		{"a/**/b/**/c.go", "a/b/c.go", true},          // two double-stars, both match zero
		{"a/**/b/**/c.go", "a/x/b/y/c.go", true},      // both expanding
		{"a/**/b/**/c.go", "a/b/x/c.go", true},        // first ** matches zero, second ** matches x/
		{"a/**/b/**/c.go", "a/x/c.go", false},         // no b/ component anywhere
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"__"+tt.path, func(t *testing.T) {
			got := matchGitPattern(tt.pattern, tt.path)
			assert.Equal(t, tt.want, got, "matchGitPattern(%q, %q)", tt.pattern, tt.path)
		})
	}
}

// ---------------------------------------------------------------------------
// EC-4: Single-star must not cross path separators
// ---------------------------------------------------------------------------

func TestEdgeSingleStarNoCrossSlash(t *testing.T) {
	tests := []struct {
		pattern string
		path    string
		want    bool
	}{
		{"*/*.go", "src/main.go", true},  // explicit one-level
		{"*/*.go", "a/b/main.go", false}, // single * can't cross two slashes
		{"*/*.go", "main.go", false},     // needs a directory segment
		{"*.go", "a/b/c.go", true},       // basename-only match (no slash in pattern)
		{"*.go", "a/b/c.rb", false},      // extension mismatch
		{"src/*.go", "src/main.go", true},
		{"src/*.go", "src/sub/main.go", false}, // * can't cross slash
		{"src/*.go", "src/main.rb", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"__"+tt.path, func(t *testing.T) {
			got := matchGitPattern(tt.pattern, tt.path)
			assert.Equal(t, tt.want, got, "matchGitPattern(%q, %q)", tt.pattern, tt.path)
		})
	}
}

// ---------------------------------------------------------------------------
// EC-5: Character classes
// ---------------------------------------------------------------------------

func TestEdgeCharClasses(t *testing.T) {
	tests := []struct {
		pattern string
		path    string
		want    bool
		note    string
	}{
		{"*.[ch]", "foo.c", true, "c in class"},
		{"*.[ch]", "foo.h", true, "h in class"},
		{"*.[ch]", "foo.go", false, "go not in [ch]"},
		{"*.[ch]", "foo./", false, "slash should not match char class"},
		{"*.[^ch]", "foo.c", false, "negated [^ch]"},
		{"*.[^ch]", "foo.g", true, "g not in negated [^ch]"},
		{"*.[!ch]", "foo.c", false, "negated [!ch]"},
		{"*.[!ch]", "foo.g", true, "g not in negated [!ch]"},
		{"*.[a-z]", "foo.c", true, "c in range a-z"},
		{"*.[a-z]", "foo.C", false, "uppercase C not in a-z"},
		{"*.[A-Z]", "foo.C", true, "C in range A-Z"},
		{"*.[0-9]", "foo.3", true, "3 in range 0-9"},
		{"*.[0-9]", "foo.a", false, "a not in range 0-9"},
		// Char class should not match path separator
		{"[a-z]oo.go", "foo.go", true, "f matches [a-z]"},
		{"[a-z]oo.go", "/oo.go", false, "/ must not match char class"},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"__"+tt.path, func(t *testing.T) {
			got := matchGitPattern(tt.pattern, tt.path)
			assert.Equal(t, tt.want, got, "%s", tt.note)
		})
	}
}

// ---------------------------------------------------------------------------
// EC-6: Anchored patterns
// ---------------------------------------------------------------------------

func TestEdgeAnchoredPatterns(t *testing.T) {
	tests := []struct {
		pattern string
		path    string
		want    bool
	}{
		{"/Makefile", "Makefile", true},
		{"/Makefile", "src/Makefile", false},
		{"/Makefile", "a/b/Makefile", false},
		{"/vendor/**", "vendor/foo.go", true},
		{"/vendor/**", "vendor/a/b/c.go", true},
		{"/vendor/**", "src/vendor/foo.go", false},
		{"/src/*.go", "src/main.go", true},
		{"/src/*.go", "src/sub/main.go", false},
		{"/src/*.go", "other/src/main.go", false},
		{"/.hidden", ".hidden", true},
		{"/.hidden", "dir/.hidden", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"__"+tt.path, func(t *testing.T) {
			got := matchGitPattern(tt.pattern, tt.path)
			assert.Equal(t, tt.want, got, "matchGitPattern(%q, %q)", tt.pattern, tt.path)
		})
	}
}

// ---------------------------------------------------------------------------
// EC-7: Basename-only matching (pattern with no slash, not anchored)
// ---------------------------------------------------------------------------

func TestEdgeBasenameMatching(t *testing.T) {
	tests := []struct {
		pattern string
		path    string
		want    bool
	}{
		{"Makefile", "Makefile", true},
		{"Makefile", "src/Makefile", true},
		{"Makefile", "a/b/c/Makefile", true},
		{"Makefile", "MakefileExtra", false},
		{"*.go", "main.go", true},
		{"*.go", "src/main.go", true},
		{"*.go", "a/b/main.go", true},
		{"*.go", "main.rb", false},
		{"Dockerfile.*", "Dockerfile.dev", true},
		{"Dockerfile.*", "src/Dockerfile.prod", true},
		{"Dockerfile.*", "Dockerfile", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"__"+tt.path, func(t *testing.T) {
			got := matchGitPattern(tt.pattern, tt.path)
			assert.Equal(t, tt.want, got, "matchGitPattern(%q, %q)", tt.pattern, tt.path)
		})
	}
}

// ---------------------------------------------------------------------------
// EC-8: Edge cases for empty/special inputs to matchGitPattern
// ---------------------------------------------------------------------------

func TestEdgeMatchSpecialInputs(t *testing.T) {
	// Empty path
	assert.False(t, matchGitPattern("*.go", ""), "empty path should not match *.go")

	// Empty pattern — filepath.Base("") returns "." so empty pattern vs empty path is non-trivial
	assert.False(t, matchGitPattern("", "foo.go"), "empty pattern should not match foo.go")
	// Note: matchGitPattern("","") returns false because filepath.Base("") == "." in Go;
	// this is a known quirk — an empty pattern in .gitattributes is not meaningful in practice.

	// Just a star — should match any single filename (basename match)
	assert.True(t, matchGitPattern("*", "foo.go"), "* matches any file")
	assert.True(t, matchGitPattern("*", "main.rb"), "* matches main.rb")

	// ** alone should match everything
	assert.True(t, matchGitPattern("**", "foo.go"), "** matches any file")
	assert.True(t, matchGitPattern("**", "a/b/c.go"), "** matches deep path")
	assert.True(t, matchGitPattern("**", ""), "** matches empty path")

	// Literal "/" pattern: HasSuffix("/","/") fires first, so it becomes a trailing-slash dir pattern;
	// the path "" doesn't end with "/" so this returns false.
	assert.False(t, matchGitPattern("/", ""), "/ pattern is treated as trailing-slash dir pattern, not root anchor")

	// Pattern with special regex chars (should be literal)
	assert.False(t, matchGitPattern("foo.go", "fooXgo"), ". is literal, not regex wildcard")
	assert.True(t, matchGitPattern("foo.go", "foo.go"), "literal dot matches literal dot")

	// Pattern with + (should be literal)
	assert.True(t, matchGitPattern("*.c++", "main.c++"), "c++ extension literal match")
	assert.False(t, matchGitPattern("*.c++", "main.cpp"), "c++ does not match cpp")
}

// ---------------------------------------------------------------------------
// EC-9: Parsing edge cases
// ---------------------------------------------------------------------------

func TestEdgeParsing(t *testing.T) {
	t.Run("empty_value", func(t *testing.T) {
		// linguist-language= with empty value
		ga, err := ParseGitAttributes([]byte("*.go linguist-language=\n"))
		require.NoError(t, err)
		require.Len(t, ga.rules, 1)
		requireGitAttributeValue(t, ga.rules[0].attrs, "linguist-language", gitAttributeValue{kind: gitAttributeValueString, value: ""})
	})

	t.Run("value_with_multiple_equals", func(t *testing.T) {
		// Only split on first =
		ga, err := ParseGitAttributes([]byte("*.go linguist-language=a=b\n"))
		require.NoError(t, err)
		require.Len(t, ga.rules, 1)
		requireGitAttributeValue(t, ga.rules[0].attrs, "linguist-language", gitAttributeValue{kind: gitAttributeValueString, value: "a=b"})
	})

	t.Run("double_dash_prefix", func(t *testing.T) {
		// --linguist-vendored should NOT be treated as -linguist-vendored
		ga, err := ParseGitAttributes([]byte("*.go --linguist-vendored\n"))
		require.NoError(t, err)
		require.Len(t, ga.rules, 1)
		// The attr key should be "-linguist-vendored" (double dash → strip one dash, key="-linguist-vendored")
		_, hasVendored := ga.rules[0].attrs["linguist-vendored"]
		assert.False(t, hasVendored, "--linguist-vendored should NOT parse as linguist-vendored=false")
	})

	t.Run("crlf_line_endings", func(t *testing.T) {
		ga, err := ParseGitAttributes([]byte("*.go linguist-language=Go\r\n*.rb linguist-language=Ruby\r\n"))
		require.NoError(t, err)
		assert.Len(t, ga.rules, 2)
		requireGitAttributeValue(t, ga.rules[0].attrs, "linguist-language", gitAttributeValue{kind: gitAttributeValueString, value: "Go"})
		requireGitAttributeValue(t, ga.rules[1].attrs, "linguist-language", gitAttributeValue{kind: gitAttributeValueString, value: "Ruby"})
	})

	t.Run("inline_comment", func(t *testing.T) {
		// Inline comments are NOT part of the gitattributes spec — # within a line is not a comment
		// The parser should treat # as an attribute token
		ga, err := ParseGitAttributes([]byte("*.go linguist-vendored # comment\n"))
		require.NoError(t, err)
		require.Len(t, ga.rules, 1)
		// linguist-vendored should still be parsed
		requireGitAttributeValue(t, ga.rules[0].attrs, "linguist-vendored", gitAttributeValue{kind: gitAttributeValueSet})
		// "#" and "comment" are parsed as bare attrs (not ignored)
		_, hasHash := ga.rules[0].attrs["#"]
		assert.True(t, hasHash, "inline # is parsed as attribute name (not a comment)")
	})

	t.Run("tab_separated", func(t *testing.T) {
		ga, err := ParseGitAttributes([]byte("*.go\tlinguist-language=Go\n"))
		require.NoError(t, err)
		require.Len(t, ga.rules, 1)
		requireGitAttributeValue(t, ga.rules[0].attrs, "linguist-language", gitAttributeValue{kind: gitAttributeValueString, value: "Go"})
	})

	t.Run("language_with_special_chars", func(t *testing.T) {
		// C++ and C# are real language names with special chars
		ga, err := ParseGitAttributes([]byte("*.cpp linguist-language=C++\n"))
		require.NoError(t, err)
		require.Len(t, ga.rules, 1)
		requireGitAttributeValue(t, ga.rules[0].attrs, "linguist-language", gitAttributeValue{kind: gitAttributeValueString, value: "C++"})
	})

	t.Run("language_csharp", func(t *testing.T) {
		ga, err := ParseGitAttributes([]byte("*.cs linguist-language=C#\n"))
		require.NoError(t, err)
		require.Len(t, ga.rules, 1)
		requireGitAttributeValue(t, ga.rules[0].attrs, "linguist-language", gitAttributeValue{kind: gitAttributeValueString, value: "C#"})
	})

	t.Run("whitespace_only_line", func(t *testing.T) {
		ga, err := ParseGitAttributes([]byte("   \t   \n"))
		require.NoError(t, err)
		assert.Len(t, ga.rules, 0)
	})

	t.Run("no_newline_at_end", func(t *testing.T) {
		ga, err := ParseGitAttributes([]byte("*.go linguist-language=Go"))
		require.NoError(t, err)
		assert.Len(t, ga.rules, 1)
	})

	t.Run("nul_byte_in_content", func(t *testing.T) {
		// Should not panic
		content := []byte("*.go linguist-language=Go\x00\n*.rb linguist-language=Ruby\n")
		_, err := ParseGitAttributes(content)
		// Either succeeds or returns an error — must not panic
		_ = err
	})

	t.Run("spaces_in_pattern_not_escaped", func(t *testing.T) {
		// strings.Fields splits on whitespace, so a space in pattern breaks parsing
		// "path with spaces" → pattern="path", attrs include "with", "spaces", "linguist-vendored"
		ga, err := ParseGitAttributes([]byte("path with spaces linguist-vendored\n"))
		require.NoError(t, err)
		// Pattern is just "path" (first field), rest are parsed as attrs
		if len(ga.rules) > 0 {
			assert.Equal(t, "path", ga.rules[0].pattern, "spaces in path break parsing: only 'path' is the pattern")
		}
	})

	t.Run("very_many_rules", func(t *testing.T) {
		// Performance/stability: 1000 rules
		var sb strings.Builder
		for i := 0; i < 1000; i++ {
			sb.WriteString("*.go linguist-language=Go\n")
		}
		ga, err := ParseGitAttributes([]byte(sb.String()))
		require.NoError(t, err)
		assert.Len(t, ga.rules, 1000)
	})
}

// ---------------------------------------------------------------------------
// EC-10: Attribute semantics — precedence, negation, unknowns
// ---------------------------------------------------------------------------

func TestEdgeAttributeSemantics(t *testing.T) {
	t.Run("explicit_false_overrides_default_vendor", func(t *testing.T) {
		// vendor/ paths are detected as vendor by default
		// -linguist-vendored should override and force NOT vendor
		ga, err := ParseGitAttributes([]byte("vendor/** -linguist-vendored\n"))
		require.NoError(t, err)
		// Default IsVendor("vendor/foo.go") is true, but the override should flip it
		assert.False(t, ga.IsVendor("vendor/foo.go"), "-linguist-vendored should override default vendor detection")
	})

	t.Run("last_rule_wins_for_same_attr", func(t *testing.T) {
		content := "*.go linguist-language=Ruby\n*.go linguist-language=Go\n"
		ga, err := ParseGitAttributes([]byte(content))
		require.NoError(t, err)
		lang, ok := ga.GetLanguage("main.go")
		assert.True(t, ok)
		assert.Equal(t, "Go", lang, "last matching rule should win")
	})

	t.Run("last_rule_wins_vendor_override", func(t *testing.T) {
		content := "vendor/** linguist-vendored\nvendor/mine/** -linguist-vendored\n"
		ga, err := ParseGitAttributes([]byte(content))
		require.NoError(t, err)
		assert.False(t, ga.IsVendor("vendor/mine/foo.go"), "later -linguist-vendored should override")
		assert.True(t, ga.IsVendor("vendor/other/foo.go"), "earlier rule still applies to unoverridden paths")
	})

	t.Run("different_attrs_from_different_rules", func(t *testing.T) {
		// Rule 1 sets language, rule 2 sets vendored — both should apply to *.go
		content := "*.go linguist-language=Go\n*.go linguist-vendored\n"
		ga, err := ParseGitAttributes([]byte(content))
		require.NoError(t, err)
		lang, ok := ga.GetLanguage("main.go")
		assert.True(t, ok)
		assert.Equal(t, "Go", lang)
		assert.True(t, ga.IsVendor("main.go"), "linguist-vendored from rule 2 should also apply")
	})

	t.Run("rule2_matches_but_has_no_attr_falls_back_to_rule1", func(t *testing.T) {
		// Rule 1: *.go → language=Go; Rule 2: *.go → vendored (no language attr)
		// GetLanguage for *.go should find rule2 doesn't have language, then rule1 does
		content := "*.go linguist-language=Go\n*.go linguist-vendored\n"
		ga, err := ParseGitAttributes([]byte(content))
		require.NoError(t, err)
		lang, ok := ga.GetLanguage("main.go")
		assert.True(t, ok)
		assert.Equal(t, "Go", lang, "rule2 matches but has no language attr; rule1's language should be used")
	})

	t.Run("lowercase_language_alias", func(t *testing.T) {
		ga, err := ParseGitAttributes([]byte("*.go linguist-language=go\n"))
		require.NoError(t, err)
		lang, ok := ga.GetLanguage("main.go")
		assert.True(t, ok)
		assert.Equal(t, "Go", lang, "lowercase alias 'go' should resolve to 'Go'")
	})

	t.Run("python_alias", func(t *testing.T) {
		ga, err := ParseGitAttributes([]byte("*.py linguist-language=python\n"))
		require.NoError(t, err)
		lang, ok := ga.GetLanguage("main.py")
		assert.True(t, ok)
		assert.Equal(t, "Python", lang)
	})

	t.Run("unknown_language_returned_as_is", func(t *testing.T) {
		ga, err := ParseGitAttributes([]byte("*.xyz linguist-language=MyFakeLang9000\n"))
		require.NoError(t, err)
		lang, ok := ga.GetLanguage("test.xyz")
		assert.True(t, ok, "GetLanguage should still return true for unknown language")
		assert.Equal(t, "MyFakeLang9000", lang, "unknown lang returned as-is")
	})

	t.Run("empty_language_value", func(t *testing.T) {
		ga, err := ParseGitAttributes([]byte("*.go linguist-language=\n"))
		require.NoError(t, err)
		lang, ok := ga.GetLanguage("main.go")
		// Empty string: alias resolution will fail; should return "" as-is and true
		assert.True(t, ok, "even empty lang value returns true (override exists)")
		_ = lang // could be "" or resolved
	})

	t.Run("is_detectable_no_override", func(t *testing.T) {
		// No rules → IsDetectable returns (false, false) — no override
		ga, err := ParseGitAttributes([]byte(""))
		require.NoError(t, err)
		_, ok := ga.IsDetectable("main.go")
		assert.False(t, ok, "IsDetectable with no rules returns no override")
		_, ok = ga.IsDetectable("schema.sql")
		assert.False(t, ok, "IsDetectable with no rules returns no override for SQL")
	})

	t.Run("vendored_arbitrary_value", func(t *testing.T) {
		// linguist-vendored=anything (not "false") should be treated as true
		ga, err := ParseGitAttributes([]byte("*.go linguist-vendored=maybe\n"))
		require.NoError(t, err)
		assert.True(t, ga.IsVendor("main.go"), "any value other than 'false' is treated as vendored")
	})

	t.Run("empty_gitattributes_fallback", func(t *testing.T) {
		ga := GitAttributes{}
		// All methods should fall back to defaults
		assert.Equal(t, IsVendor("vendor/foo.go"), ga.IsVendor("vendor/foo.go"), "empty attrs falls back to IsVendor")
		assert.Equal(t, IsDocumentation("docs/guide.md"), ga.IsDocumentation("docs/guide.md"), "empty attrs falls back to IsDocumentation")
		assert.Equal(t, IsGenerated("foo_generated.go", nil), ga.IsGenerated("foo_generated.go", nil), "empty attrs falls back to IsGenerated")
		_, ok := ga.IsDetectable("main.go")
		assert.False(t, ok, "IsDetectable with no rules returns no override")
		_, ok = ga.GetLanguage("main.go")
		assert.False(t, ok, "GetLanguage with no rules returns false")
	})

	t.Run("multi_attr_same_line_both_apply", func(t *testing.T) {
		ga, err := ParseGitAttributes([]byte("*.pb.go linguist-generated linguist-vendored\n"))
		require.NoError(t, err)
		assert.True(t, ga.IsGenerated("foo.pb.go", nil))
		assert.True(t, ga.IsVendor("foo.pb.go"))
	})

	t.Run("negation_of_multi_attr", func(t *testing.T) {
		content := "*.pb.go linguist-generated linguist-vendored\nspecial.pb.go -linguist-generated -linguist-vendored\n"
		ga, err := ParseGitAttributes([]byte(content))
		require.NoError(t, err)
		assert.False(t, ga.IsGenerated("special.pb.go", nil), "last rule negates generated")
		assert.False(t, ga.IsVendor("special.pb.go"), "last rule negates vendored")
		assert.True(t, ga.IsGenerated("other.pb.go", nil), "other files still match first rule")
	})
}

// ---------------------------------------------------------------------------
// EC-11: matchGitPattern robustness — should not panic on pathological input
// ---------------------------------------------------------------------------

func TestEdgeMatchPatternNoPanic(t *testing.T) {
	pathological := []struct{ pattern, path string }{
		{"[", "foo"},
		{"[]", "foo"},
		{"[!]", "foo"},
		{"[^]", "foo"},
		{"[a-", "foo"},
		{"[-]", "-"},
		{"[---]", "-"},
		{"***", "foo.go"},
		{"****", "a/b/c"},
		{"**/**", "a/b/c"},
		{"**/**/**", "a/b/c"},
		{"[\\]", "\\"},
		{"[abc", "a"},
		{"a[b", "ab"},
	}

	for _, tt := range pathological {
		t.Run(tt.pattern+"__"+tt.path, func(t *testing.T) {
			// Should not panic; result is don't care
			assert.NotPanics(t, func() {
				_ = matchGitPattern(tt.pattern, tt.path)
			})
		})
	}
}

// ---------------------------------------------------------------------------
// EC-12: Real-world .gitattributes patterns from major open source projects
// ---------------------------------------------------------------------------

func TestEdgeRealWorldPatterns(t *testing.T) {
	// Patterns extracted from real repositories
	content := `
# Kubernetes-style
vendor/** linguist-vendored
*.generated.go linguist-generated
zz_generated*.go linguist-generated
staging/src/** linguist-vendored

# Tensorflow-style
tensorflow/core/api_def/base_api/** linguist-documentation
*.pb.go linguist-generated
third_party/** linguist-vendored

# Rails-style
app/assets/javascripts/cable.js linguist-generated
db/schema.rb linguist-generated
Gemfile.lock linguist-generated

# Custom language overrides
Podfile linguist-language=Ruby
*.podspec linguist-language=Ruby
Fastfile linguist-language=Ruby
Dangerfile linguist-language=Ruby
*.jbuilder linguist-language=Ruby
`

	ga, err := ParseGitAttributes([]byte(content))
	require.NoError(t, err)

	// Kubernetes patterns
	assert.True(t, ga.IsVendor("vendor/k8s.io/client-go/rest/client.go"))
	assert.True(t, ga.IsVendor("vendor/github.com/pkg/errors/errors.go"))
	assert.True(t, ga.IsGenerated("pkg/apis/core/v1/zz_generated.conversion.go", nil))
	assert.True(t, ga.IsGenerated("pkg/apis/core/v1/types.generated.go", nil))
	assert.True(t, ga.IsVendor("staging/src/k8s.io/apimachinery/pkg/api.go"))

	// Tensorflow patterns
	assert.True(t, ga.IsVendor("third_party/protobuf/src/google/protobuf/any.proto"))
	assert.True(t, ga.IsGenerated("tensorflow/go/op/wrappers.pb.go", nil))
	assert.True(t, ga.IsDocumentation("tensorflow/core/api_def/base_api/api_def_Abs.pbtxt"))

	// Language overrides
	lang, ok := ga.GetLanguage("Podfile")
	assert.True(t, ok)
	assert.Equal(t, "Ruby", lang)

	lang, ok = ga.GetLanguage("MyApp.podspec")
	assert.True(t, ok)
	assert.Equal(t, "Ruby", lang)

	lang, ok = ga.GetLanguage("Fastfile")
	assert.True(t, ok)
	assert.Equal(t, "Ruby", lang)

	lang, ok = ga.GetLanguage("app/views/users/index.jbuilder")
	assert.True(t, ok)
	assert.Equal(t, "Ruby", lang)
}

// ---------------------------------------------------------------------------
// EC-13: linguist-detectable semantics
// ---------------------------------------------------------------------------

func TestEdgeDetectable(t *testing.T) {
	t.Run("sql_made_detectable", func(t *testing.T) {
		ga, err := ParseGitAttributes([]byte("*.sql linguist-detectable\n"))
		require.NoError(t, err)
		val, ok := ga.IsDetectable("schema.sql")
		assert.True(t, ok, "should have override")
		assert.True(t, val, "should be detectable")
		val, ok = ga.IsDetectable("migrations/001_init.sql")
		assert.True(t, ok)
		assert.True(t, val)
		_, ok = ga.IsDetectable("main.go")
		assert.False(t, ok, "go files not in rule have no override")
	})

	t.Run("detectable_then_not", func(t *testing.T) {
		content := "*.sql linguist-detectable\nexception.sql -linguist-detectable\n"
		ga, err := ParseGitAttributes([]byte(content))
		require.NoError(t, err)
		val, ok := ga.IsDetectable("exception.sql")
		assert.True(t, ok, "should have override")
		assert.False(t, val, "explicitly not detectable")
		val, ok = ga.IsDetectable("other.sql")
		assert.True(t, ok)
		assert.True(t, val, "other sql files still detectable")
	})

	t.Run("html_made_detectable", func(t *testing.T) {
		ga, err := ParseGitAttributes([]byte("*.html linguist-detectable\n"))
		require.NoError(t, err)
		val, ok := ga.IsDetectable("index.html")
		assert.True(t, ok)
		assert.True(t, val)
	})

	t.Run("programming_language_made_undetectable", func(t *testing.T) {
		// -linguist-detectable can exclude a programming language from stats
		ga, err := ParseGitAttributes([]byte("tools/*.py -linguist-detectable\n"))
		require.NoError(t, err)
		val, ok := ga.IsDetectable("tools/export.py")
		assert.True(t, ok, "should have override")
		assert.False(t, val, "-linguist-detectable should exclude from stats")
		// Other py files have no override
		_, ok = ga.IsDetectable("src/main.py")
		assert.False(t, ok, "no override for src/main.py")
	})
}

// ---------------------------------------------------------------------------
// EC-14: Question-mark wildcard edge cases
// ---------------------------------------------------------------------------

func TestEdgeQuestionMark(t *testing.T) {
	tests := []struct {
		pattern string
		path    string
		want    bool
	}{
		{"?.go", "a.go", true},
		{"?.go", "ab.go", false}, // too long
		{"?.go", ".go", false},   // too short (0 chars before .)
		{"?.go", "/.go", false},  // ? doesn't match /
		{"??.go", "ab.go", true},
		{"??.go", "a.go", false},
		{"a?.go", "ab.go", true},
		{"a?.go", "a/.go", false},  // ? doesn't match /
		{"a?.go", "abc.go", false}, // pattern expects exactly 1 char after a
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"__"+tt.path, func(t *testing.T) {
			got := matchGitPattern(tt.pattern, tt.path)
			assert.Equal(t, tt.want, got, "matchGitPattern(%q, %q)", tt.pattern, tt.path)
		})
	}
}

// ---------------------------------------------------------------------------
// EC-15: Interaction between getAttr reverse-iteration and multi-attr rules
// ---------------------------------------------------------------------------

func TestEdgeGetAttrPrecedence(t *testing.T) {
	// Rule 1: *.go → language=Go, vendored=true
	// Rule 2: *.go → language=Ruby  (no vendored)
	// For main.go:
	//   language → rule2 wins (last match with this attr) → Ruby
	//   vendored → rule2 has no vendored, rule1 has vendored=true → true
	content := "*.go linguist-language=Go linguist-vendored\n*.go linguist-language=Ruby\n"
	ga, err := ParseGitAttributes([]byte(content))
	require.NoError(t, err)

	lang, ok := ga.GetLanguage("main.go")
	assert.True(t, ok)
	assert.Equal(t, "Ruby", lang, "last rule with language attr wins")

	assert.True(t, ga.IsVendor("main.go"), "rule1's vendored attr applies since rule2 doesn't have it")
}

// ---------------------------------------------------------------------------
// EC-16: Very deep paths and complex ** expansion
// ---------------------------------------------------------------------------

func TestEdgeDeepPaths(t *testing.T) {
	tests := []struct {
		pattern string
		path    string
		want    bool
	}{
		{"**/*.go", "a/b/c/d/e/f/g/h/i/j/main.go", true},
		{"vendor/**", "vendor/a/b/c/d/e/f/g/h/i/main.go", true},
		{"/vendor/**", "vendor/a/b/c/d/e/f.go", true},
		{"src/**/pkg/**/main.go", "src/a/b/c/pkg/d/e/f/main.go", true},
		{"src/**/pkg/**/main.go", "src/pkg/main.go", true},
		{"*.go", "a/b/c/d/e/f/g/h/i/j/main.go", true}, // basename match
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"__"+tt.path, func(t *testing.T) {
			got := matchGitPattern(tt.pattern, tt.path)
			assert.Equal(t, tt.want, got, "matchGitPattern(%q, %q)", tt.pattern, tt.path)
		})
	}
}

// ---------------------------------------------------------------------------
// EC-17: ParseGitAttributes returns zero rules for pattern-only lines
// ---------------------------------------------------------------------------

func TestEdgePatternOnlyLines(t *testing.T) {
	// Lines with only a pattern and no attributes should be silently ignored
	inputs := []string{
		"*.go\n",
		"/vendor/**\n",
		"Makefile\n",
		"**\n",
	}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			ga, err := ParseGitAttributes([]byte(input))
			require.NoError(t, err)
			assert.Len(t, ga.rules, 0, "pattern-only line '%s' should produce 0 rules", input)
		})
	}
}
