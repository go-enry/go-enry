package enry

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsVendor(t *testing.T) {
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
		assert.Equal(t, is, test.expected, fmt.Sprintf("%v: is = %v, expected: %v", test.name, is, test.expected))
	}
}

func TestIsDocumentation(t *testing.T) {
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
		assert.Equal(t, is, test.expected, fmt.Sprintf("%v: is = %v, expected: %v", test.name, is, test.expected))
	}
}

func TestIsImage(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{name: "TestIsImage_1", path: "invalid.txt", expected: false},
		{name: "TestIsImage_2", path: "image.png", expected: true},
		{name: "TestIsImage_3", path: "image.jpg", expected: true},
		{name: "TestIsImage_4", path: "image.jpeg", expected: true},
		{name: "TestIsImage_5", path: "image.gif", expected: true},
	}

	for _, test := range tests {
		is := IsImage(test.path)
		assert.Equal(t, is, test.expected, fmt.Sprintf("%v: is = %v, expected: %v", test.name, is, test.expected))
	}
}

func TestGetMimeType(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		lang     string
		expected string
	}{
		{name: "TestGetMimeType_1", path: "text.txt", lang: "", expected: "text/plain"},
		{name: "TestGetMimeType_2", path: "file.go", lang: "Go", expected: "text/x-go"},
		{name: "TestGetMimeType_3", path: "image.png", lang: "", expected: "image/png"},
	}

	for _, test := range tests {
		is := GetMIMEType(test.path, test.lang)
		assert.Equal(t, is, test.expected, fmt.Sprintf("%v: is = %v, expected: %v", test.name, is, test.expected))
	}
}

func TestIsConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{name: "TestIsConfiguration_1", path: "foo", expected: false},
		{name: "TestIsConfiguration_2", path: "foo.ini", expected: true},
		{name: "TestIsConfiguration_3", path: "/test/path/foo.json", expected: true},
	}

	for _, test := range tests {
		is := IsConfiguration(test.path)
		assert.Equal(t, is, test.expected, fmt.Sprintf("%v: is = %v, expected: %v", test.name, is, test.expected))
	}
}

func TestIsBinary(t *testing.T) {
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
		assert.Equal(t, is, test.expected, fmt.Sprintf("%v: is = %v, expected: %v", test.name, is, test.expected))
	}
}

func TestIsDotFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{name: "TestIsDotFile_1", path: "foo/bar/./", expected: false},
		{name: "TestIsDotFile_2", path: "./", expected: false},
	}

	for _, test := range tests {
		is := IsDotFile(test.path)
		assert.Equal(t, test.expected, is, fmt.Sprintf("%v: is = %v, expected: %v", test.name, is, test.expected))
	}
}

func TestGetColor(t *testing.T) {
	tests := []struct {
		name     string
		language string
		expected string
	}{
		{name: "TestGetColor_1", language: "Go", expected: "#00ADD8"},
		{name: "TestGetColor_2", language: "SomeRandom", expected: "#cccccc"},
		{name: "TestGetColor_3", language: "HTML", expected: "#e34c26"},
		{name: "TestGetColor_4", language: "HTML+PHP", expected: "#e34c26"},
	}

	for _, test := range tests {
		color := GetColor(test.language)
		assert.Equal(t, test.expected, color, fmt.Sprintf("%v: is = %v, expected: %v", test.name, color, test.expected))
	}
}
