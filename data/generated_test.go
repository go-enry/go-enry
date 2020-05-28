package data

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFirstLineEmptyContent(t *testing.T) {
	require.Nil(t, getFirstLine(nil))
}

func TestForEachLine(t *testing.T) {
	const sample = "foo\nbar\nboomboom\nbleepbloop\n"
	var lines = strings.Split(sample, "\n")

	var result []string
	forEachLine([]byte(sample), func(l []byte) {
		result = append(result, string(l))
	})

	require.Equal(t, lines[:len(lines)-1], result)
}

func TestGetLines(t *testing.T) {
	const sample = "foo\nbar\nboomboom\nbleepbloop\n"

	testCases := []struct {
		lines    int
		expected []string
	}{
		{1, []string{"foo"}},
		{2, []string{"foo", "bar"}},
		{10, []string{"foo", "bar", "boomboom", "bleepbloop"}},
		{-1, []string{"bleepbloop"}},
		{-2, []string{"bleepbloop", "boomboom"}},
		{-10, []string{"bleepbloop", "boomboom", "bar", "foo"}},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprint(tt.lines), func(t *testing.T) {
			lines := getLines([]byte(sample), tt.lines)

			var result = make([]string, len(lines))
			for i, l := range lines {
				result[i] = string(l)
			}

			require.Equal(t, tt.expected, result)
		})
	}
}
