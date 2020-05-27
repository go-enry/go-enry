package enry

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
func TestIsTestFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{name: "TestPHP_Is", path: "tests/FooTest.php", expected: true},
		{name: "TestPHP_Not", path: "foo/FooTest.php", expected: false},
		{name: "TestJava_Is_1", path: "test/FooTest.java", expected: true},
		{name: "TestJava_Is_2", path: "test/FooTests.java", expected: true},
		{name: "TestJava_Is_3", path: "test/TestFoo.java", expected: true},
		{name: "TestJava_Is_4", path: "test/qux/TestFoo.java", expected: true},
		{name: "TestJava_Not", path: "foo/FooTest.java", expected: false},
		{name: "TestScala_Is_1", path: "test/FooTest.scala", expected: true},
		{name: "TestScala_Is_2", path: "test/FooTests.scala", expected: true},
		{name: "TestScala_Is_3", path: "test/FooSpec.scala", expected: true},
		{name: "TestScala_Is_4", path: "test/qux/FooSpecs.scala", expected: true},
		{name: "TestScala_Not", path: "foo/FooTest.scala", expected: false},
		{name: "TestPython_Is", path: "test_foo.py", expected: true},
		{name: "TestPython_Not", path: "foo_test.py", expected: false},
		{name: "TestGo_Is", path: "foo_test.go", expected: true},
		{name: "TestGo_Not", path: "test_foo.go", expected: false},
		{name: "TestRuby_Is_1", path: "foo_test.rb", expected: true},
		{name: "TestRuby_Is_1", path: "foo_spec.rb", expected: true},
		{name: "TestRuby_Not", path: "foo_specs.rb", expected: false},
		{name: "TestCSharp_Is_1", path: "FooTest.cs", expected: true},
		{name: "TestCSharp_Is_2", path: "foo/FooTests.cs", expected: true},
		{name: "TestCSharp_Not", path: "foo/TestFoo.cs", expected: false},
		{name: "TestJavaScript_Is_1", path: "foo.test.js", expected: true},
		{name: "TestJavaScript_Is_2", path: "foo.spec.js", expected: true},
		{name: "TestJavaScript_Not", path: "footest.js", expected: false},
		{name: "TestTypeScript_Is_1", path: "foo.test.ts", expected: true},
		{name: "TestTypeScript_Is_2", path: "foo.spec.ts", expected: true},
		{name: "TestTypeScript_Not", path: "footest.ts", expected: false},
	}

	for _, test := range tests {
		is := IsTest(test.path)
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

func TestIsGenerated(t *testing.T) {
	testCases := []struct {
		file      string
		load      bool
		generated bool
	}{
		// Xcode project files
		{"Binary/MainMenu.nib", false, true},
		{"Dummy/foo.xcworkspacedata", false, true},
		{"Dummy/foo.xcuserstate", false, true},

		//Cocoapods
		{"Pods/Pods.xcodeproj", false, true},
		{"Pods/SwiftDependency/foo.swift", false, true},
		{"Pods/ObjCDependency/foo.h", false, true},
		{"Pods/ObjCDependency/foo.m", false, true},
		{"Dummy/Pods/Pods.xcodeproj", false, true},
		{"Dummy/Pods/SwiftDependency/foo.swift", false, true},
		{"Dummy/Pods/ObjCDependency/foo.h", false, true},
		{"Dummy/Pods/ObjCDependency/foo.m", false, true},

		//Carthage
		{"Carthage/Build/.Dependency.version", false, true},
		{"Carthage/Build/iOS/Dependency.framework", false, true},
		{"Carthage/Build/Mac/Dependency.framework", false, true},
		{"src/Carthage/Build/.Dependency.version", false, true},
		{"src/Carthage/Build/iOS/Dependency.framework", false, true},
		{"src/Carthage/Build/Mac/Dependency.framework", false, true},

		//Go-specific vendored paths
		{"go/vendor/github.com/foo.go", false, true},
		{"go/vendor/golang.org/src/foo.c", false, true},
		{"go/vendor/gopkg.in/some/nested/path/foo.go", false, true},

		//.NET designer file
		{"Dummy/foo.designer.cs", false, true},
		{"Dummy/foo.Designer.cs", false, true},
		{"Dummy/foo.designer.vb", false, true},
		{"Dummy/foo.Designer.vb", false, true},

		//Composer generated composer.lock file
		{"JSON/composer.lock", false, true},

		//Node modules
		{"Dummy/node_modules/foo.js", false, true},

		//npm shrinkwrap file
		{"Dummy/npm-shrinkwrap.json", false, true},
		{"Dummy/package-lock.json", false, true},
		{"JavaScript/jquery-1.6.1.min.js", true, true},

		//Yarn Plug'n'Play file
		{".pnp.js", false, true},
		{".pnp.cjs", false, true},
		{".pnp.mjs", false, true},

		//Godep saved dependencies
		{"Godeps/Godeps.json", false, true},
		{"Godeps/_workspace/src/github.com/kr/s3/sign.go", false, true},

		//Generated by Zephir
		{"C/exception.zep.c", false, true},
		{"C/exception.zep.h", false, true},
		{"PHP/exception.zep.php", false, true},

		//Minified files
		{"JavaScript/jquery-1.6.1.min.js", true, true},

		//JavaScript with source-maps
		{"JavaScript/namespace.js", true, true},
		{"Generated/inline.js", true, true},

		//CSS with source-maps
		{"Generated/linked.css", true, true},
		{"Generated/inline.css", true, true},

		//Source-map
		{"Data/bootstrap.css.map", true, true},
		{"Generated/linked.css.map", true, true},
		{"Data/sourcemap.v3.map", true, true},
		{"Data/sourcemap.v1.map", true, true},

		//Specflow
		{"Features/BindingCulture.feature.cs", false, true},

		//JFlex
		{"Java/JFlexLexer.java", true, true},

		//GrammarKit
		{"Java/GrammarKit.java", true, true},

		//roxygen2
		{"R/import.Rd", true, true},

		//PostScript
		{"PostScript/lambda.pfa", true, true},

		//Perl ppport.h
		{"Generated/ppport.h", true, true},

		//Graphql Relay
		{"Javascript/__generated__/App_user.graphql.js", false, true},

		//Game Maker Studio 2
		{"JSON/GMS2_Project.yyp", true, true},
		{"JSON/2ea73365-b6f1-4bd1-a454-d57a67e50684.yy", true, true},
		{"Generated/options_main.inherited.yy", true, true},

		//Pipenv
		{"Dummy/Pipfile.lock", false, true},

		//HTML
		{"HTML/attr-swapped.html", true, true},
		{"HTML/extra-attr.html", true, true},
		{"HTML/extra-spaces.html", true, true},
		{"HTML/extra-tags.html", true, true},
		{"HTML/grohtml.html", true, true},
		{"HTML/grohtml.xhtml", true, true},
		{"HTML/makeinfo.html", true, true},
		{"HTML/mandoc.html", true, true},
		{"HTML/node78.html", true, true},
		{"HTML/org-mode.html", true, true},
		{"HTML/quotes-double.html", true, true},
		{"HTML/quotes-none.html", true, true},
		{"HTML/quotes-single.html", true, true},
		{"HTML/uppercase.html", true, true},
		{"HTML/ronn.html", true, true},
		{"HTML/unknown.html", true, false},
		{"HTML/no-content.html", true, false},
		{"HTML/pages.html", true, true},

		//GIMP
		{"C/image.c", true, true},
		{"C/image.h", true, true},

		//Haxe
		{"Generated/Haxe/main.js", true, true},
		{"Generated/Haxe/main.py", true, true},
		{"Generated/Haxe/main.lua", true, true},
		{"Generated/Haxe/Main.cpp", true, true},
		{"Generated/Haxe/Main.h", true, true},
		{"Generated/Haxe/Main.java", true, true},
		{"Generated/Haxe/Main.cs", true, true},
		{"Generated/Haxe/Main.php", true, true},
	}

	for _, tt := range testCases {
		t.Run(tt.file, func(t *testing.T) {
			var content []byte
			if tt.load {
				var err error
				content, err = ioutil.ReadFile(filepath.Join("_testdata", tt.file))
				require.NoError(t, err)
			}

			result := IsGenerated(tt.file, content)
			require.Equal(t, tt.generated, result)
		})
	}
}

func TestFoo(t *testing.T) {
	file := "HTML/uppercase.html"
	content, err := ioutil.ReadFile("_testdata/" + file)
	require.NoError(t, err)
	require.True(t, IsGenerated(file, content))
}
