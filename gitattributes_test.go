package enry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			ga := ParseGitAttributes([]byte(tt.content))
			assert.Equal(t, tt.wantRules, len(ga.rules))
		})
	}
}

func TestParseGitAttributesValues(t *testing.T) {
	content := "*.go linguist-vendored linguist-language=Go -linguist-documentation\n"
	ga := ParseGitAttributes([]byte(content))
	require.Len(t, ga.rules, 1)

	rule := ga.rules[0]
	assert.Equal(t, "*.go", rule.pattern)
	assert.Equal(t, "true", rule.attrs["linguist-vendored"])
	assert.Equal(t, "Go", rule.attrs["linguist-language"])
	assert.Equal(t, "false", rule.attrs["linguist-documentation"])
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

		// Trailing slash (directory only) - should not match files
		{"vendor/", "vendor", false},
		{"vendor/", "vendor/foo.go", false},

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
	ga := ParseGitAttributes([]byte("special/lib/* linguist-vendored\nmy/lib/* -linguist-vendored\n"))

	// Explicit vendored
	assert.True(t, ga.IsVendor("special/lib/foo.go"))

	// Explicit not vendored
	assert.False(t, ga.IsVendor("my/lib/bar.go"))

	// No rule - falls back to default
	assert.False(t, ga.IsVendor("src/main.go"))
}

func TestGitAttributesIsDocumentation(t *testing.T) {
	ga := ParseGitAttributes([]byte("docs/* linguist-documentation\napi-docs/* -linguist-documentation\n"))

	assert.True(t, ga.IsDocumentation("docs/guide.md"))
	assert.False(t, ga.IsDocumentation("api-docs/spec.md"))
}

func TestGitAttributesIsGenerated(t *testing.T) {
	ga := ParseGitAttributes([]byte("*.pb.go linguist-generated\nhand-written.pb.go -linguist-generated\n"))

	assert.True(t, ga.IsGenerated("foo.pb.go", nil))

	// Explicit override to not generated - last matching rule wins
	assert.False(t, ga.IsGenerated("hand-written.pb.go", nil))
}

func TestGitAttributesIsDetectable(t *testing.T) {
	ga := ParseGitAttributes([]byte("*.sql linguist-detectable\n"))

	assert.True(t, ga.IsDetectable("schema.sql"))
	assert.False(t, ga.IsDetectable("main.go"))
}

func TestGitAttributesGetLanguage(t *testing.T) {
	ga := ParseGitAttributes([]byte("*.extension linguist-language=Go\n"))

	lang, ok := ga.GetLanguage("foo.extension")
	assert.True(t, ok)
	assert.Equal(t, "Go", lang)

	// No match
	_, ok = ga.GetLanguage("foo.other")
	assert.False(t, ok)
}

func TestGitAttributesGetLanguageAlias(t *testing.T) {
	// "py" and "python" should both resolve to "Python"
	ga := ParseGitAttributes([]byte("*.myext linguist-language=python\n"))

	lang, ok := ga.GetLanguage("test.myext")
	assert.True(t, ok)
	assert.Equal(t, "Python", lang)
}

func TestGitAttributesMultipleRules(t *testing.T) {
	// Last matching rule wins (git precedence)
	content := `*.go linguist-language=Go
special.go linguist-language=Ruby
`
	ga := ParseGitAttributes([]byte(content))

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
	ga := ParseGitAttributes([]byte(content))

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
	ga := ParseGitAttributes([]byte(content))

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
