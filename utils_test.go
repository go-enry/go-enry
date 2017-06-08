package slinguist

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (s *SimpleLinguistTestSuite) TestIsVendor() {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{name: "TestIsVendor_1", path: "foo/bar", expected: false},
		{name: "TestIsVendor_2", path: "foo/vendor/foo", expected: true},
		{name: "TestIsVendor_3", path: ".sublime-project", expected: true},
		{name: "TestIsVendor_4", path: "leaflet.draw-src.js", expected: true},
		{name: "TestIsVendor_5", path: "foo/bar/MochiKit.js", expected: true},
		{name: "TestIsVendor_6", path: "foo/bar/dojo.js", expected: true},
		{name: "TestIsVendor_7", path: "foo/env/whatever", expected: true},
		{name: "TestIsVendor_8", path: "foo/.imageset/bar", expected: true},
		{name: "TestIsVendor_9", path: "Vagrantfile", expected: true},
	}

	for _, test := range tests {
		is := IsVendor(test.path)
		assert.Equal(s.T(), is, test.expected, fmt.Sprintf("%v: is = %v, expected: %v", test.name, is, test.expected))
	}
}

func (s *SimpleLinguistTestSuite) TestIsDocumentation() {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{name: "TestIsDocumentation_1", path: "foo", expected: false},
		{name: "TestIsDocumentation_2", path: "README", expected: true},
	}

	for _, test := range tests {
		is := IsDocumentation(test.path)
		assert.Equal(s.T(), is, test.expected, fmt.Sprintf("%v: is = %v, expected: %v", test.name, is, test.expected))
	}
}

func (s *SimpleLinguistTestSuite) TestIsConfiguration() {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{name: "TestIsConfiguration_1", path: "foo", expected: false},
		{name: "TestIsConfiguration_2", path: "foo.ini", expected: true},
		{name: "TestIsConfiguration_3", path: "foo.json", expected: true},
	}

	for _, test := range tests {
		is := IsConfiguration(test.path)
		assert.Equal(s.T(), is, test.expected, fmt.Sprintf("%v: is = %v, expected: %v", test.name, is, test.expected))
	}
}

func (s *SimpleLinguistTestSuite) TestIsBinary() {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{name: "TestIsBinary_1", data: []byte("foo"), expected: false},
		{name: "TestIsBinary_2", data: []byte{0}, expected: true},
		{name: "TestIsBinary_3", data: bytes.Repeat([]byte{'o'}, 8000), expected: false},
	}

	for _, test := range tests {
		is := IsBinary(test.data)
		assert.Equal(s.T(), is, test.expected, fmt.Sprintf("%v: is = %v, expected: %v", test.name, is, test.expected))
	}
}

const (
	htmlPath = "some/random/dir/file.html"
	jsPath   = "some/random/dir/file.js"
)

func BenchmarkVendor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = IsVendor(htmlPath)
	}
}

func BenchmarkVendorJS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = IsVendor(jsPath)
	}
}
