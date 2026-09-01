package enry

import (
	"testing"

	"github.com/go-enry/go-enry/v2/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func requireGitAttributeValue(t *testing.T, attrs map[string]gitAttributeValue, name string, want gitAttributeValue) {
	t.Helper()
	got, ok := attrs[name]
	require.True(t, ok, "expected attribute %q to be present", name)
	assert.Equal(t, want, got)
}

func TestParseGitAttributes(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantRules int
	}{
		{"empty", "", 0},
		{"comment only", "# this is a comment\n", 0},
		{"blank lines", "\n\n\n", 0},
		{"pattern without attrs", "*.go\n", 0},
		{"single attr set", "*.go linguist-language=Go\n", 1},
		{"single attr unset", "*.go -linguist-vendored\n", 1},
		{"single attr bare", "*.go linguist-vendored\n", 1},
		{"multiple attrs", "*.go linguist-vendored linguist-language=Go\n", 1},
		{"multiple rules", "*.go linguist-language=Go\n*.rb linguist-language=Ruby\n", 2},
		{"mixed with comments", "# header\n*.go linguist-language=Go\n# footer\n", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ga, err := ParseGitAttributes([]byte(tt.content))
			require.NoError(t, err)
			assert.Equal(t, tt.wantRules, len(ga.rules))
		})
	}
}

func TestParseGitAttributesValues(t *testing.T) {
	content := "*.go linguist-vendored linguist-language=Go -linguist-documentation\n"
	ga, err := ParseGitAttributes([]byte(content))
	require.NoError(t, err)
	require.Len(t, ga.rules, 1)

	rule := ga.rules[0]
	assert.Equal(t, "*.go", rule.pattern)
	requireGitAttributeValue(t, rule.attrs, "linguist-vendored", gitAttributeValue{kind: gitAttributeValueSet})
	requireGitAttributeValue(t, rule.attrs, "linguist-language", gitAttributeValue{kind: gitAttributeValueString, value: "Go"})
	requireGitAttributeValue(t, rule.attrs, "linguist-documentation", gitAttributeValue{kind: gitAttributeValueUnset})
}

func TestParseGitAttributesMacros(t *testing.T) {
	t.Run("macro_expands_linguist_attrs", func(t *testing.T) {
		ga, err := ParseGitAttributes([]byte("[attr]vendored linguist-vendored\n*.go vendored\n"))
		require.NoError(t, err)
		require.Len(t, ga.rules, 1)
		requireGitAttributeValue(t, ga.rules[0].attrs, "linguist-vendored", gitAttributeValue{kind: gitAttributeValueSet})
		assert.True(t, ga.IsVendor("main.go"))
	})

	t.Run("macro_can_be_defined_after_use", func(t *testing.T) {
		ga, err := ParseGitAttributes([]byte("*.go vendored\n[attr]vendored linguist-vendored\n"))
		require.NoError(t, err)
		assert.True(t, ga.IsVendor("main.go"))
	})

	t.Run("nested_macros_expand_recursively", func(t *testing.T) {
		content := "[attr]base linguist-vendored\n[attr]combo base linguist-language=Go\n*.go combo\n"
		ga, err := ParseGitAttributes([]byte(content))
		require.NoError(t, err)
		assert.True(t, ga.IsVendor("main.go"))
		lang, ok := ga.GetLanguage("main.go")
		assert.True(t, ok)
		assert.Equal(t, "Go", lang)
	})

	t.Run("special_macro_forms_do_not_expand", func(t *testing.T) {
		content := "[attr]lang linguist-language=Go\n*.rb !lang\n*.py -lang\n*.ts lang=value\n"
		ga, err := ParseGitAttributes([]byte(content))
		require.NoError(t, err)
		require.Len(t, ga.rules, 3)

		_, hasLang := ga.rules[0].attrs["linguist-language"]
		assert.False(t, hasLang)
		_, hasLang = ga.rules[1].attrs["linguist-language"]
		assert.False(t, hasLang)
		_, hasLang = ga.rules[2].attrs["linguist-language"]
		assert.False(t, hasLang)

		_, ok := ga.GetLanguage("main.rb")
		assert.False(t, ok)
		_, ok = ga.GetLanguage("main.py")
		assert.False(t, ok)
		_, ok = ga.GetLanguage("main.ts")
		assert.False(t, ok)
	})

	t.Run("macro_override_affects_has_potential_override", func(t *testing.T) {
		ga, err := ParseGitAttributes([]byte("[attr]clear !linguist-vendored\nvendor/** clear\n"))
		require.NoError(t, err)
		assert.True(t, ga.HasPotentialOverride("linguist-vendored"))
	})
}

func TestGitAttributesUnspecifiedVsStringValue(t *testing.T) {
	content := "[attr]clear !linguist-vendored\nclear.go clear\nstring.go linguist-vendored=unspecified\n"
	ga, err := ParseGitAttributes([]byte(content))
	require.NoError(t, err)

	assert.False(t, ga.IsVendor("clear.go"), "unspecified should fall back to default detection")
	assert.True(t, ga.IsVendor("string.go"), "string value 'unspecified' should not be treated as reset")
}

func TestGitAttributeLinguistAliasMapHasNoConflicts(t *testing.T) {
	seen := map[string]string{}
	for _, info := range data.LanguageInfoByID {
		assertNoGitAttributeAliasConflict(t, seen, info.Name, info.Name)
		for _, alias := range info.Aliases {
			assertNoGitAttributeAliasConflict(t, seen, alias, info.Name)
		}
	}
}

func assertNoGitAttributeAliasConflict(t *testing.T, seen map[string]string, alias string, lang string) {
	t.Helper()
	key := gitAttributeLinguistAlias(alias)
	if existing, exists := seen[key]; exists && existing != lang {
		t.Fatalf("git attribute alias %q maps to both %q and %q", key, existing, lang)
	}
	seen[key] = lang
}

func TestGitAttributesGetLanguageLinguistAliases(t *testing.T) {
	tests := []struct {
		attrValue string
		want      string
	}{
		{"OpenStep-Property-List", "OpenStep Property List"},
		{"Common-Lisp", "Common Lisp"},
		{"common-lisp", "Common Lisp"},
		{"common_lisp", "Common Lisp"},
		{"Emacs-Lisp", "Emacs Lisp"},
		{"POV-Ray-SDL", "POV-Ray SDL"},
		{"Tree-sitter-Query", "Tree-sitter Query"},
		{"coffee-script", "CoffeeScript"},
		{"Definitely-Not-A-Language", "Definitely-Not-A-Language"},
	}

	for _, tt := range tests {
		t.Run(tt.attrValue, func(t *testing.T) {
			ga, err := ParseGitAttributes([]byte("*.x linguist-language=" + tt.attrValue + "\n"))
			require.NoError(t, err)

			lang, ok := ga.GetLanguage("file.x")
			require.True(t, ok)
			assert.Equal(t, tt.want, lang)
		})
	}
}

func TestMatchGitPattern(t *testing.T) {
	tests := []struct {
		pattern string
		path    string
		want    bool
	}{
		// Basic extension matching
		{"*.go", "main.go", true},
		{"*.go", "src/main.go", true},
		{"*.go", "main.rb", false},

		// Question mark
		{"?.go", "a.go", true},
		{"?.go", "ab.go", false},
		{"?.go", "/.go", false},

		// Double star
		{"**/vendor/**", "a/vendor/b", true},
		{"**/vendor/**", "vendor/b", true},
		{"**/test", "a/b/test", true},
		{"**/*.go", "src/main.go", true},
		{"**/*.go", "main.go", true},

		// Anchored patterns (leading /)
		{"/vendor/**", "vendor/lib/foo.go", true},
		{"/vendor/**", "src/vendor/foo.go", false},
		{"/Makefile", "Makefile", true},
		{"/Makefile", "src/Makefile", false},

		// Basename matching (no slash in pattern, not anchored)
		{"Makefile", "Makefile", true},
		{"Makefile", "src/Makefile", true},

		// Slash in pattern (not anchored) - matches full path
		{"vendor/*.go", "vendor/foo.go", true},
		{"vendor/*.go", "src/vendor/foo.go", false},

		// Trailing slash (directory-only pattern) - only matches directory paths
		// Unlike .gitignore, .gitattributes trailing "/" does NOT recurse into files
		{"vendor/", "vendor/", true},
		{"vendor/", "vendor/foo.go", false},
		{"vendor/", "vendor/lib/bar.go", false},
		{"vendor/", "vendor", false},
		{"vendor/", "notvendor/foo.go", false},

		// Character classes
		{"*.[ch]", "foo.c", true},
		{"*.[ch]", "foo.h", true},
		{"*.[ch]", "foo.go", false},

		// Star does not match slash
		{"*.go", "a/b.go", true},   // basename match
		{"*/*.go", "a/b.go", true}, // explicit slash in pattern
		{"*/*.go", "b.go", false},  // needs a directory
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.path, func(t *testing.T) {
			got := matchGitPattern(tt.pattern, tt.path)
			assert.Equal(t, tt.want, got, "matchGitPattern(%q, %q)", tt.pattern, tt.path)
		})
	}
}

func TestGitAttributesIsVendor(t *testing.T) {
	ga, err := ParseGitAttributes([]byte("special/lib/* linguist-vendored\nmy/lib/* -linguist-vendored\n"))
	require.NoError(t, err)

	// Explicit vendored
	assert.True(t, ga.IsVendor("special/lib/foo.go"))

	// Explicit not vendored
	assert.False(t, ga.IsVendor("my/lib/bar.go"))

	// No rule - falls back to default
	assert.False(t, ga.IsVendor("src/main.go"))
}

func TestGitAttributesIsDocumentation(t *testing.T) {
	ga, err := ParseGitAttributes([]byte("docs/* linguist-documentation\napi-docs/* -linguist-documentation\n"))
	require.NoError(t, err)

	assert.True(t, ga.IsDocumentation("docs/guide.md"))
	assert.False(t, ga.IsDocumentation("api-docs/spec.md"))
}

func TestGitAttributesIsGenerated(t *testing.T) {
	ga, err := ParseGitAttributes([]byte("*.pb.go linguist-generated\nhand-written.pb.go -linguist-generated\n"))
	require.NoError(t, err)

	assert.True(t, ga.IsGenerated("foo.pb.go", nil))

	// Explicit override to not generated - last matching rule wins
	assert.False(t, ga.IsGenerated("hand-written.pb.go", nil))
}

func TestGitAttributesIsDetectable(t *testing.T) {
	ga, err := ParseGitAttributes([]byte("*.sql linguist-detectable\n"))
	require.NoError(t, err)

	val, ok := ga.IsDetectable("schema.sql")
	assert.True(t, ok, "should have override")
	assert.True(t, val, "should be detectable")

	_, ok = ga.IsDetectable("main.go")
	assert.False(t, ok, "no override for main.go")
}

func TestGitAttributesGetLanguage(t *testing.T) {
	ga, err := ParseGitAttributes([]byte("*.extension linguist-language=Go\n"))
	require.NoError(t, err)

	lang, ok := ga.GetLanguage("foo.extension")
	assert.True(t, ok)
	assert.Equal(t, "Go", lang)

	// No match
	_, ok = ga.GetLanguage("foo.other")
	assert.False(t, ok)
}

func TestGitAttributesGetLanguageAlias(t *testing.T) {
	// "py" and "python" should both resolve to "Python"
	ga, err := ParseGitAttributes([]byte("*.myext linguist-language=python\n"))
	require.NoError(t, err)

	lang, ok := ga.GetLanguage("test.myext")
	assert.True(t, ok)
	assert.Equal(t, "Python", lang)
}

func TestGitAttributesMultipleRules(t *testing.T) {
	// Last matching rule wins (git precedence)
	content := `*.go linguist-language=Go
special.go linguist-language=Ruby
`
	ga, err := ParseGitAttributes([]byte(content))
	require.NoError(t, err)

	// "special.go" matches both rules; last one wins
	lang, ok := ga.GetLanguage("special.go")
	assert.True(t, ok)
	assert.Equal(t, "Ruby", lang)

	// "main.go" only matches first rule
	lang, ok = ga.GetLanguage("main.go")
	assert.True(t, ok)
	assert.Equal(t, "Go", lang)
}

func TestGitAttributesMultipleAttrsPerLine(t *testing.T) {
	content := "vendor/** linguist-vendored linguist-generated\n"
	ga, err := ParseGitAttributes([]byte(content))
	require.NoError(t, err)

	assert.True(t, ga.IsVendor("vendor/lib.go"))
	assert.True(t, ga.IsGenerated("vendor/lib.go", nil))
}

func TestGitAttributesLinguistCompatibility(t *testing.T) {
	// Test patterns that match real linguist .gitattributes usage
	content := `
# Mark vendored dependencies
vendor/** linguist-vendored
third_party/** linguist-vendored

# Mark generated files
*.pb.go linguist-generated
*_generated.go linguist-generated

# Override language detection
Vagrantfile linguist-language=Ruby
*.podspec linguist-language=Ruby
Dockerfile.* linguist-language=Dockerfile

# Documentation
docs/** linguist-documentation
`
	ga, err := ParseGitAttributes([]byte(content))
	require.NoError(t, err)

	assert.True(t, ga.IsVendor("vendor/github.com/pkg/errors/errors.go"))
	assert.True(t, ga.IsVendor("third_party/protobuf/something.cc"))
	assert.True(t, ga.IsGenerated("api/service.pb.go", nil))
	assert.True(t, ga.IsGenerated("internal/types_generated.go", nil))
	assert.True(t, ga.IsDocumentation("docs/getting-started.md"))

	lang, ok := ga.GetLanguage("Vagrantfile")
	assert.True(t, ok)
	assert.Equal(t, "Ruby", lang)

	lang, ok = ga.GetLanguage("MyApp.podspec")
	assert.True(t, ok)
	assert.Equal(t, "Ruby", lang)
}
