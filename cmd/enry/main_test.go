package main

import (
	"os"
	"path/filepath"
	"testing"

	enry "github.com/go-enry/go-enry/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLines(t *testing.T) {
	tests := []struct {
		content      string
		wantTotal    int
		wantNonBlank int
	}{
		// 0
		{content: "This is one line", wantTotal: 1, wantNonBlank: 1},
		// 1 Test no content
		{content: "", wantTotal: 0, wantNonBlank: 0},
		// 2 A single blank line
		{content: "One blank line\n\nTwo nonblank lines", wantTotal: 3, wantNonBlank: 2},
		// 3 Testing multiple blank lines in a row
		{content: "\n\n", wantTotal: 3, wantNonBlank: 0},
		// 4 '
		{content: "\n\n\n\n", wantTotal: 5, wantNonBlank: 0},
		// 5 Multiple blank lines content on ends
		{content: "content\n\n\n\ncontent", wantTotal: 5, wantNonBlank: 2},
		// 6 Content with blank lines on ends
		{content: "\n\n\ncontent\n\n\n", wantTotal: 7, wantNonBlank: 1},
	}

	for i, test := range tests {
		t.Run("", func(t *testing.T) {
			gotTotal, gotNonBlank := getLines("", []byte(test.content))
			if gotTotal != test.wantTotal || gotNonBlank != test.wantNonBlank {
				t.Errorf("wrong line counts obtained for test case #%d:\n      %7s, %7s\nGOT:   %7d, %7d\nWANT:  %7d, %7d\n", i, "TOTAL", "NON_BLANK",
					gotTotal, gotNonBlank, test.wantTotal, test.wantNonBlank)
			}
		})
	}
}

func TestAnalyzeDirRespectsNestedTrailingSlashVendorOverride(t *testing.T) {
	root := t.TempDir()

	writeFile := func(relativePath string, content string) {
		path := filepath.Join(root, relativePath)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	}

	writeFile(".gitattributes", "vendor/** linguist-vendored\nvendor/github.com/ -linguist-vendored\n")
	writeFile("main.go", "package main\n")
	writeFile("vendor/github.com/pkg/errors/errors.go", "package errors\n")
	writeFile("vendor/acme/lib.go", "package acme\n")

	content, err := os.ReadFile(filepath.Join(root, ".gitattributes"))
	require.NoError(t, err)

	gitAttrs, err := enry.ParseGitAttributes(content)
	require.NoError(t, err)

	out, err := analyzeDir(root, -1, true, gitAttrs)
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{"main.go", "vendor/github.com/pkg/errors/errors.go"}, out["Go"])
}
