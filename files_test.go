package slinguist

import (
	"bytes"
	"regexp"

	. "gopkg.in/check.v1"
)

func (s *UtilsSuite) TestIsVendor(c *C) {
	c.Assert(IsVendor("foo/bar"), Equals, false)
	c.Assert(IsVendor("foo/vendor/foo"), Equals, true)
	c.Assert(IsVendor(".travis.yml"), Equals, true)
	c.Assert(IsVendor("foo.framer"), Equals, true)
	c.Assert(IsVendor("foo/bar/coffee-script.js"), Equals, true)
	c.Assert(IsVendor("foo/bar/less.js"), Equals, true)
	c.Assert(IsVendor("foo/bar/underscore.js"), Equals, true)
	c.Assert(IsVendor("foo/bar/lodash.js"), Equals, true)
	c.Assert(IsVendor("foo/bar/lodash.js"), Equals, true)
	c.Assert(IsVendor("foo/bar/lodash.core.js"), Equals, true)
	c.Assert(IsVendor("foo/bar/backbone.js"), Equals, true)
	c.Assert(IsVendor("foo/bar/ember.js"), Equals, true)
	c.Assert(IsVendor("foo/bar/ember.debug.js"), Equals, true)
	c.Assert(IsVendor("foo/bar/ember.prod.js"), Equals, true)
	c.Assert(IsVendor("foo/bar/three.js"), Equals, true)
	c.Assert(IsVendor("foo/bar/babylon.js"), Equals, true)
	c.Assert(IsVendor("foo/bar/babylon.2.3.js"), Equals, true)
	c.Assert(IsVendor("foo/bar/babylon.2.3.core.js"), Equals, true)
	c.Assert(IsVendor("foo/bar/babylon.2.3.max.js"), Equals, true)
	c.Assert(IsVendor("foo/bar/babylon.2.3.noworker.js"), Equals, true)
	c.Assert(IsVendor("foo/bar/html5-3.6-respond-1.4.2.js"), Equals, true)
}

func (s *UtilsSuite) TestIsDocumentation(c *C) {
	c.Assert(IsDocumentation("foo"), Equals, false)
	c.Assert(IsDocumentation("README"), Equals, true)
}

func (s *UtilsSuite) TestIsConfiguration(c *C) {
	c.Assert(IsConfiguration("foo"), Equals, false)
	c.Assert(IsConfiguration("foo.ini"), Equals, true)
	c.Assert(IsConfiguration("foo.json"), Equals, true)
}

func (s *UtilsSuite) TestIsBinary(c *C) {
	c.Assert(IsBinary([]byte("foo")), Equals, false)

	binary := []byte{0}
	c.Assert(IsBinary(binary), Equals, true)

	binary = bytes.Repeat([]byte{'o'}, 8000)
	binary = append(binary, byte(0))
	c.Assert(IsBinary(binary), Equals, false)
}

const (
	htmlPath = "some/random/dir/file.html"
	jsPath   = "some/random/dir/file.js"
)

func (s *UtilsSuite) BenchmarkVendor(c *C) {
	for i := 0; i < c.N; i++ {
		_ = IsVendor(htmlPath)
	}
}

func (s *UtilsSuite) BenchmarkVendorJS(c *C) {
	for i := 0; i < c.N; i++ {
		_ = IsVendor(jsPath)
	}
}

var vendorRegexp = []*regexp.Regexp{
	regexp.MustCompile(`(^|/)cache/`),
	regexp.MustCompile(`^[Dd]ependencies/`),
	regexp.MustCompile(`^deps/`),
	regexp.MustCompile(`^tools/`),
	regexp.MustCompile(`(^|/)configure$`),
	regexp.MustCompile(`(^|/)configure.ac$`),
	regexp.MustCompile(`(^|/)config.guess$`),
	regexp.MustCompile(`(^|/)config.sub$`),
	regexp.MustCompile(`node_modules/`),
	regexp.MustCompile(`bower_components/`),
	regexp.MustCompile(`^rebar$`),
	regexp.MustCompile(`erlang.mk`),
	regexp.MustCompile(`Godeps/_workspace/`),
	regexp.MustCompile(`(\.|-)min\.(js|css)$`),
	regexp.MustCompile(`([^\s]*)import\.(css|less|scss|styl)$`),
	regexp.MustCompile(`(^|/)bootstrap([^.]*)\.(js|css|less|scss|styl)$`),
	regexp.MustCompile(`(^|/)custom\.bootstrap([^\s]*)(js|css|less|scss|styl)$`),
	regexp.MustCompile(`(^|/)font-awesome\.(css|less|scss|styl)$`),
	regexp.MustCompile(`(^|/)foundation\.(css|less|scss|styl)$`),
	regexp.MustCompile(`(^|/)normalize\.(css|less|scss|styl)$`),
	regexp.MustCompile(`(^|/)[Bb]ourbon/.*\.(css|less|scss|styl)$`),
	regexp.MustCompile(`(^|/)animate\.(css|less|scss|styl)$`),
	regexp.MustCompile(`third[-_]?party/`),
	regexp.MustCompile(`3rd[-_]?party/`),
	regexp.MustCompile(`vendors?/`),
	regexp.MustCompile(`extern(al)?/`),
	regexp.MustCompile(`^debian/`),
	regexp.MustCompile(`run.n$`),
	regexp.MustCompile(`(^|/)jquery([^.]*)\.js$`),
	regexp.MustCompile(`(^|/)jquery\-\d\.\d+(\.\d+)?\.js$`),
	regexp.MustCompile(`(^|/)jquery\-ui(\-\d\.\d+(\.\d+)?)?(\.\w+)?\.(js|css)$`),
	regexp.MustCompile(`(^|/)jquery\.(ui|effects)\.([^.]*)\.(js|css)$`),
	regexp.MustCompile(`(^|/)prototype(.*)\.js$`),
	regexp.MustCompile(`(^|/)effects\.js$`),
	regexp.MustCompile(`(^|/)controls\.js$`),
	regexp.MustCompile(`(^|/)dragdrop\.js$`),
	regexp.MustCompile(`(.*?)\.d\.ts$`),
	regexp.MustCompile(`(^|/)mootools([^.]*)\d+\.\d+.\d+([^.]*)\.js$`),
	regexp.MustCompile(`(^|/)dojo\.js$`),
	regexp.MustCompile(`(^|/)MochiKit\.js$`),
	regexp.MustCompile(`(^|/)yahoo-([^.]*)\.js$`),
	regexp.MustCompile(`(^|/)yui([^.]*)\.js$`),
	regexp.MustCompile(`(^|/)ckeditor\.js$`),
	regexp.MustCompile(`(^|/)tiny_mce([^.]*)\.js$`),
	regexp.MustCompile(`(^|/)tiny_mce/(langs|plugins|themes|utils)`),
	regexp.MustCompile(`(^|/)MathJax/`),
	regexp.MustCompile(`(^|/)Chart\.js$`),
	regexp.MustCompile(`(^|/)[Cc]ode[Mm]irror/(lib|mode|theme|addon|keymap)`),
	regexp.MustCompile(`(^|/)shBrush([^.]*)\.js$`),
	regexp.MustCompile(`(^|/)shCore\.js$`),
	regexp.MustCompile(`(^|/)shLegacy\.js$`),
	regexp.MustCompile(`(^|/)angular([^.]*)\.js$`),
	regexp.MustCompile(`(^|\/)d3(\.v\d+)?([^.]*)\.js$`),
	regexp.MustCompile(`(^|/)react(-[^.]*)?\.js$`),
	regexp.MustCompile(`(^|/)modernizr\-\d\.\d+(\.\d+)?\.js$`),
	regexp.MustCompile(`(^|/)modernizr\.custom\.\d+\.js$`),
	regexp.MustCompile(`(^|/)knockout-(\d+\.){3}(debug\.)?js$`),
	regexp.MustCompile(`(^|/)admin_media/`),
	regexp.MustCompile(`^fabfile\.py$`),
	regexp.MustCompile(`^waf$`),
	regexp.MustCompile(`^.osx$`),
	regexp.MustCompile(`^Pods/`),
	regexp.MustCompile(`(^|/)Sparkle/`),
	regexp.MustCompile(`(^|/)gradlew$`),
	regexp.MustCompile(`(^|/)gradlew\.bat$`),
	regexp.MustCompile(`(^|/)gradle/wrapper/`),
	regexp.MustCompile(`-vsdoc\.js$`),
	regexp.MustCompile(`\.intellisense\.js$`),
	regexp.MustCompile(`(^|/)jquery([^.]*)\.validate(\.unobtrusive)?\.js$`),
	regexp.MustCompile(`(^|/)jquery([^.]*)\.unobtrusive\-ajax\.js$`),
	regexp.MustCompile(`(^|/)[Mm]icrosoft([Mm]vc)?([Aa]jax|[Vv]alidation)(\.debug)?\.js$`),
	regexp.MustCompile(`^[Pp]ackages\/.+\.\d+\/`),
	regexp.MustCompile(`(^|/)extjs/.*?\.js$`),
	regexp.MustCompile(`(^|/)extjs/.*?\.xml$`),
	regexp.MustCompile(`(^|/)extjs/.*?\.txt$`),
	regexp.MustCompile(`(^|/)extjs/.*?\.html$`),
	regexp.MustCompile(`(^|/)extjs/.*?\.properties$`),
	regexp.MustCompile(`(^|/)extjs/.sencha/`),
	regexp.MustCompile(`(^|/)extjs/docs/`),
	regexp.MustCompile(`(^|/)extjs/builds/`),
	regexp.MustCompile(`(^|/)extjs/cmd/`),
	regexp.MustCompile(`(^|/)extjs/examples/`),
	regexp.MustCompile(`(^|/)extjs/locale/`),
	regexp.MustCompile(`(^|/)extjs/packages/`),
	regexp.MustCompile(`(^|/)extjs/plugins/`),
	regexp.MustCompile(`(^|/)extjs/resources/`),
	regexp.MustCompile(`(^|/)extjs/src/`),
	regexp.MustCompile(`(^|/)extjs/welcome/`),
	regexp.MustCompile(`(^|/)html5shiv\.js$`),
	regexp.MustCompile(`^[Ss]amples/`),
	regexp.MustCompile(`^[Tt]est/fixtures/`),
	regexp.MustCompile(`(^|/)cordova([^.]*)\.js$`),
	regexp.MustCompile(`(^|/)cordova\-\d\.\d(\.\d)?\.js$`),
	regexp.MustCompile(`foundation(\..*)?\.js$`),
	regexp.MustCompile(`^Vagrantfile$`),
	regexp.MustCompile(`.[Dd][Ss]_[Ss]tore$`),
	regexp.MustCompile(`^vignettes/`),
	regexp.MustCompile(`^inst/extdata/`),
	regexp.MustCompile(`octicons.css`),
	regexp.MustCompile(`sprockets-octicons.scss`),
	regexp.MustCompile(`(^|/)activator$`),
	regexp.MustCompile(`(^|/)activator\.bat$`),
	regexp.MustCompile(`proguard.pro`),
	regexp.MustCompile(`proguard-rules.pro`),
	regexp.MustCompile(`gitattributes$`),
	regexp.MustCompile(`gitignore$`),
	regexp.MustCompile(`gitmodules$`),
	regexp.MustCompile(`.travis.yml$`),
	regexp.MustCompile(`circle.yml$`),
}

func isVendorRegexp(s string) bool {
	for _, re := range vendorRegexp {
		found := re.FindStringIndex(s)
		if found != nil {
			return found[1] >= 0
		}
	}
	return false
}

func (s *UtilsSuite) BenchmarkVendorRegexp(c *C) {
	for i := 0; i < c.N; i++ {
		_ = isVendorRegexp(htmlPath)
	}
}

func (s *UtilsSuite) BenchmarkVendorRegexpJS(c *C) {
	for i := 0; i < c.N; i++ {
		_ = isVendorRegexp(htmlPath)
	}
}
