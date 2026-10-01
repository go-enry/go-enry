package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func TestPrintPercentsZeroTotals(t *testing.T) {
	for _, mode := range []string{"byte", "line"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "empty.py"), nil, 0600); err != nil {
				t.Fatal(err)
			}
			var output bytes.Buffer
			printPercents(root, map[string][]string{"Python": {"empty.py"}}, &output, mode)
			if got, want := output.String(), "0.00%\tPython\n"; got != want {
				t.Fatalf("got %q; want %q", got, want)
			}
		})
	}
}

func TestPrintPercentsNonzeroTotals(t *testing.T) {
	root := t.TempDir()
	for name, content := range map[string]string{"a.py": "x", "b.rb": "xxx"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	var output bytes.Buffer
	printPercents(root, map[string][]string{"Python": {"a.py"}, "Ruby": {"b.rb"}}, &output, "byte")
	if got, want := output.String(), "75.00%\tRuby\n25.00%\tPython\n"; got != want {
		t.Fatalf("got %q; want %q", got, want)
	}
}

func TestGetLinesBufferBoundaries(t *testing.T) {
	for _, size := range []int{4095, 4096, 4097, 8192} {
		for _, suffix := range []string{"", "\n", "\r\n", "\nnext"} {
			t.Run(fmt.Sprintf("%d/%q", size, suffix), func(t *testing.T) {
				content := strings.Repeat("x", size) + suffix
				want := 1
				if suffix == "\nnext" {
					want = 2
				}
				path := filepath.Join(t.TempDir(), "long.txt")
				if err := os.WriteFile(path, []byte(content), 0600); err != nil {
					t.Fatal(err)
				}
				for _, input := range [][]byte{[]byte(content), nil} {
					total, code := getLines(path, input)
					if total != want || code != want {
						t.Errorf("got %d/%d; want %d/%d", total, code, want, want)
					}
				}
			})
		}
	}
}
