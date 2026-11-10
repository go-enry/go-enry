package enry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmacsModelineFirstDelimiter(t *testing.T) {
	for _, test := range []struct {
		name, content string
		expected      []string
	}{
		{"fontset", `fontset = "-*-courier-medium-r-*-*-12-*-*-*-m-*-*-*"`, nil},
		{"short form", "// -*- c++ -*-", []string{"C++"}},
		{"compact form", "// -*-c++-*-", []string{"C++"}},
		{"explicit mode", "// -*- coding: utf-8; mode: c++; -*-", []string{"C++"}},
		{"first delimiter on line", "// -*- ruby -*- -*- python -*-", []string{"Ruby"}},
		{"unknown first mode", "// -*- unknown-mode -*- -*- python -*-", nil},
		{"last matching line", "// -*- ruby -*-\n// -*- python -*-", []string{"Python"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, GetLanguagesByEmacsModeline("", []byte(test.content), nil))
		})
	}
}

func TestGtkRCFontset(t *testing.T) {
	content := []byte(`fontset = "-*-courier-medium-r-*-*-12-*-*-*-m-*-*-*"`)
	assert.Equal(t, "GtkRC", GetLanguage("v9t9.gtkrc", content))
}
