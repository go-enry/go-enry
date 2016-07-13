package slinguist

import (
	"bytes"
	"path/filepath"
	"strings"

	"gopkg.in/toqueteos/substring.v1"
)

// From github/linguist.
// curl https://raw.githubusercontent.com/github/linguist/master/lib/linguist/vendor.yml | python -c 'import sys, yaml; l = yaml.load(sys.stdin.read()); print "var skipped = []*regexp.Regexp{\n" + "\n".join(["\tregexp.MustCompile(`" + i + "`)," for i in l]) + "\n}"'

var vendorMatchers = substring.Or(
	substring.Suffix(`.framer`),
	substring.SuffixGroup(`.js`,
		substring.After(`cordova-`, substring.Regexp(`\d\.\d(\.\d)?`)),
		substring.After(`d3`, substring.Regexp(`(\.v\d+)?([^.]*)`)),
		substring.After(`jquery-`, substring.Regexp(`\d\.\d+(\.\d+)?`)),
		substring.After(`jquery-ui`, substring.Regexp(`(\-\d\.\d+(\.\d+)?)?(\.\w+)?`)),
		substring.After(`knockout-`, substring.Regexp(`(\d+\.){3}(debug\.)?$`)),
		substring.After(`microsoft`, substring.Regexp(`([Mm]vc)?([Aa]jax|[Vv]alidation)(\.debug)?`)),
		substring.After(`Microsoft`, substring.Regexp(`([Mm]vc)?([Aa]jax|[Vv]alidation)(\.debug)?`)),
		substring.After(`modernizr-`, substring.Regexp(`\d\.\d+(\.\d+)?`)),
		substring.After(`mootools`, substring.Regexp(`([^.]*)\d+\.\d+.\d+([^.]*)`)),
		substring.Has(`angular`),
		substring.Has(`bootstrap`),
		substring.Has(`cordova`),
		substring.Has(`less`),
		substring.Has(`custom.bootstrap`),
		substring.Has(`extjs/`),
		substring.Has(`foundation`),
		substring.Has(`jquery.effects.`),
		substring.Has(`jquery.ui.`),
		substring.Has(`jquery`),
		substring.Has(`modernizr.custom.`),
		substring.Has(`prototype`),
		substring.Has(`backbone`),
		substring.Has(`three`),
		substring.Has(`ember`),
		substring.Has(`babylon`),
		substring.Has(`react`),
		substring.Has(`shBrush`),
		substring.Has(`tiny_mce`),
		substring.Has(`yahoo-`),
		substring.Has(`html5-`),
		substring.Has(`yui`),
		substring.Has(`underscore`),
		substring.Has(`lodash`),
		substring.Has(`lodash.core`),
		substring.Has(`coffee-script`),
		substring.Suffixes(
			`-vsdoc.js`,
			`.intellisense.js`,
			`Chart.js`,
			`ckeditor.js`,
			`controls.js`,
			`dojo.js`,
			`dragdrop.js`,
			`effects.js`,
			`html5shiv.js`,
			`MochiKit.js`,
			`shCore.js`,
			`shLegacy.js`,
			`.min.js`,
			`-min.js`,
		),
	),
	substring.SuffixGroup(`.css`,
		substring.Has(`bootstrap`),
		substring.Has(`custom.bootstrap`),
		substring.Has(`jquery.effects.`),
		substring.Has(`jquery.ui.`),
		substring.Has(`octicons.css`),
		substring.After(`jquery-ui`, substring.Regexp(`(\-\d\.\d+(\.\d+)?)?(\.\w+)?`)),
		substring.Suffixes(
			`animate.css`,
			`bourbon.css`,
			`Bourbon.css`,
			`font-awesome.css`,
			`foundation.css`,
			`import.css`,
			`normalize.css`,
			`.min.css`,
			`-min.css`,
		),
	),
	substring.SuffixGroup(`.scss`,
		substring.Has(`bootstrap`),
		substring.Has(`custom.bootstrap`),
		substring.Has(`sprockets-octicons.scss`),
		substring.Suffixes(
			`animate.scss`,
			`bourbon.scss`,
			`Bourbon.scss`,
			`font-awesome.scss`,
			`foundation.scss`,
			`import.scss`,
			`normalize.scss`,
		),
	),
	substring.SuffixGroup(`.less`,
		substring.Has(`bootstrap`),
		substring.Has(`custom.bootstrap`),
		substring.Suffixes(
			`animate.less`,
			`bourbon.less`,
			`Bourbon.less`,
			`font-awesome.less`,
			`foundation.less`,
			`import.less`,
			`normalize.less`,
		),
	),
	substring.SuffixGroup(`.styl`,
		substring.Has(`bootstrap`),
		substring.Has(`custom.bootstrap`),
		substring.Suffixes(
			`animate.styl`,
			`bourbon.styl`,
			`Bourbon.styl`,
			`font-awesome.styl`,
			`foundation.styl`,
			`import.styl`,
			`normalize.styl`,
		),
	),
	substring.After(`codemirror/`, substring.Or(
		substring.Exact(`lib`),
		substring.Exact(`mode`),
		substring.Exact(`theme`),
		substring.Exact(`addon`),
		substring.Exact(`keymap`),
	)),
	substring.After(`extjs/`, substring.Or(
		substring.Suffixes(`.html`, `.properties`, `.txt`, `.xml`),
		substring.Exact(`.sencha/`),
		substring.Exact(`builds/`),
		substring.Exact(`cmd/`),
		substring.Exact(`docs/`),
		substring.Exact(`examples/`),
		substring.Exact(`locale/`),
		substring.Exact(`packages/`),
		substring.Exact(`plugins/`),
		substring.Exact(`resources/`),
		substring.Exact(`src/`),
		substring.Exact(`welcome/`),
	)),
	substring.After(`tiny_mce/`, substring.Or(
		substring.Exact(`langs`),
		substring.Exact(`plugins`),
		substring.Exact(`themes`),
		substring.Exact(`utils`),
	)),
	substring.Has(`3rd-party/`),
	substring.Has(`3rd_party/`),
	substring.Has(`3rdparty/`),
	substring.Has(`admin_media/`),
	substring.Has(`bower_components/`),
	substring.Has(`cache/`),
	substring.Has(`Dependencies/`),
	substring.Has(`dependencies/`),
	substring.Has(`erlang.mk`),
	substring.Has(`extern/`),
	substring.Has(`external/`),
	substring.Has(`Godeps/_workspace/`),
	substring.Has(`gradle/wrapper/`),
	substring.Has(`MathJax/`),
	substring.Has(`node_modules/`),
	substring.Has(`proguard-rules.pro`),
	substring.Has(`proguard.pro`),
	substring.Has(`Sparkle/`),
	substring.Has(`third-party/`),
	substring.Has(`third_party/`),
	substring.Has(`thirdparty/`),
	substring.Has(`vendor/`),
	substring.Has(`vendors/`),
	substring.Exact(`.osx`),
	substring.Exact(`rebar`),
	substring.Exact(`Vagrantfile`),
	substring.Exact(`waf`),
	substring.Prefixes(
		`.google_apis/`,
		`debian/`,
		`deps/`,
		`inst/extdata/`,
		`packages/`,
		`Packages/`,
		`Pods/`,
		`Samples/`,
		`samples/`,
		`Test/fixture/`,
		`test/fixture/`,
		`Test/fixtures/`,
		`test/fixtures/`,
		`tools/`,
		`vignettes/`,
	),
	substring.Suffixes(
		`.d.ts`,
		`.travis.yml`,
		`activator.bat`,
		`activator`,
		`circle.yml`,
		`config.guess`,
		`config.sub`,
		`configure.ac`,
		`configure`,
		`fabfile.py`,
		`gitattributes`,
		`gitignore`,
		`gitmodules`,
		`gradlew.bat`,
		`gradlew`,
		`run.n`,
		`.DS_Store`,
		`.DS_store`,
		`.dS_Store`,
		`.dS_store`,
		`.Ds_Store`,
		`.Ds_store`,
		`.ds_Store`,
		`.ds_store`,
	),
)

var documentationMatchers = substring.Or(
	substring.Regexp(`^docs?/`),
	substring.Regexp(`(^|/)[Dd]ocumentation/`),
	substring.Regexp(`(^|/)javadoc/`),
	substring.Regexp(`^man/`),
	substring.Regexp(`^[Ee]xamples/`),
	substring.Regexp(`(^|/)CHANGE(S|LOG)?(\.|$)`),
	substring.Regexp(`(^|/)CONTRIBUTING(\.|$)`),
	substring.Regexp(`(^|/)COPYING(\.|$)`),
	substring.Regexp(`(^|/)INSTALL(\.|$)`),
	substring.Regexp(`(^|/)LICEN[CS]E(\.|$)`),
	substring.Regexp(`(^|/)[Ll]icen[cs]e(\.|$)`),
	substring.Regexp(`(^|/)README(\.|$)`),
	substring.Regexp(`(^|/)[Rr]eadme(\.|$)`),
	substring.Regexp(`^[Ss]amples/`),
)

var configurationLanguages = []string{
	"XML", "JSON", "TOML", "YAML", "INI", "SQL",
}

func VendorIndex(path string) int {
	return findIndex(path, vendorMatchers)
}

func DocumentationIndex(path string) int {
	return findIndex(path, documentationMatchers)
}

func findIndex(path string, matchers substring.StringsMatcher) int {
	return matchers.MatchIndex(path)
}

func IsVendor(path string) bool {
	return VendorIndex(path) >= 0
}

func IsDotFile(path string) bool {
	return strings.HasPrefix(filepath.Base(path), ".")
}

func IsDocumentation(path string) bool {
	return DocumentationIndex(path) >= 0
}

const sniffLen = 8000

//IsBinary detects if data is a binary value based on:
//http://git.kernel.org/cgit/git/git.git/tree/xdiff-interface.c?id=HEAD#n198
func IsBinary(data []byte) bool {
	if len(data) > sniffLen {
		data = data[:sniffLen]
	}

	if bytes.IndexByte(data, byte(0)) == -1 {
		return false
	}

	return true
}
