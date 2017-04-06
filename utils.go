package slinguist

// CODE GENERATED AUTOMATICALLY WITH github.com/src-d/simple-linguist/cli/slinguist-generate
// THIS FILE SHOULD NOT BE EDITED BY HAND
// Extracted from github/linguist commit: dae33dc2b20cddc85d1300435c3be7118a7115a9

import (
	"bytes"
	"path/filepath"
	"strings"

	"gopkg.in/toqueteos/substring.v1"
)

func IsAuxiliaryLanguage(lang string) bool {
	_, ok := auxiliaryLanguages[lang]
	return ok
}

func IsConfiguration(path string) bool {
	lang, _ := GetLanguageByExtension(path)
	_, is := configurationLanguages[lang]

	return is
}

func IsDotFile(path string) bool {
	return strings.HasPrefix(filepath.Base(path), ".")
}

func IsVendor(path string) bool {
	return findIndex(path, vendorMatchers) >= 0
}

func IsDocumentation(path string) bool {
	return findIndex(path, documentationMatchers) >= 0
}

func findIndex(path string, matchers substring.StringsMatcher) int {
	return matchers.MatchIndex(path)
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

var vendorMatchers = substring.Or(
	substring.Regexp(`(^|/)cache/`),
	substring.Regexp(`^[Dd]ependencies/`),
	substring.Regexp(`(^|/)dist/`),
	substring.Regexp(`^deps/`),
	substring.Regexp(`^tools/`),
	substring.Regexp(`(^|/)configure$`),
	substring.Regexp(`(^|/)config.guess$`),
	substring.Regexp(`(^|/)config.sub$`),
	substring.Regexp(`(^|/)aclocal.m4`),
	substring.Regexp(`(^|/)libtool.m4`),
	substring.Regexp(`(^|/)ltoptions.m4`),
	substring.Regexp(`(^|/)ltsugar.m4`),
	substring.Regexp(`(^|/)ltversion.m4`),
	substring.Regexp(`(^|/)lt~obsolete.m4`),
	substring.Regexp(`cpplint.py`),
	substring.Regexp(`node_modules/`),
	substring.Regexp(`bower_components/`),
	substring.Regexp(`^rebar$`),
	substring.Regexp(`erlang.mk`),
	substring.Regexp(`Godeps/_workspace/`),
	substring.Regexp(`.indent.pro`),
	substring.Regexp(`(\.|-)min\.(js|css)$`),
	substring.Regexp(`([^\s]*)import\.(css|less|scss|styl)$`),
	substring.Regexp(`(^|/)bootstrap([^.]*)\.(js|css|less|scss|styl)$`),
	substring.Regexp(`(^|/)custom\.bootstrap([^\s]*)(js|css|less|scss|styl)$`),
	substring.Regexp(`(^|/)font-awesome\.(css|less|scss|styl)$`),
	substring.Regexp(`(^|/)foundation\.(css|less|scss|styl)$`),
	substring.Regexp(`(^|/)normalize\.(css|less|scss|styl)$`),
	substring.Regexp(`(^|/)[Bb]ourbon/.*\.(css|less|scss|styl)$`),
	substring.Regexp(`(^|/)animate\.(css|less|scss|styl)$`),
	substring.Regexp(`third[-_]?party/`),
	substring.Regexp(`3rd[-_]?party/`),
	substring.Regexp(`vendors?/`),
	substring.Regexp(`extern(al)?/`),
	substring.Regexp(`(^|/)[Vv]&#43;endor/`),
	substring.Regexp(`^debian/`),
	substring.Regexp(`run.n$`),
	substring.Regexp(`bootstrap-datepicker/`),
	substring.Regexp(`(^|/)jquery([^.]*)\.js$`),
	substring.Regexp(`(^|/)jquery\-\d\.\d&#43;(\.\d&#43;)?\.js$`),
	substring.Regexp(`(^|/)jquery\-ui(\-\d\.\d&#43;(\.\d&#43;)?)?(\.\w&#43;)?\.(js|css)$`),
	substring.Regexp(`(^|/)jquery\.(ui|effects)\.([^.]*)\.(js|css)$`),
	substring.Regexp(`jquery.fn.gantt.js`),
	substring.Regexp(`jquery.fancybox.(js|css)`),
	substring.Regexp(`fuelux.js`),
	substring.Regexp(`(^|/)jquery\.fileupload(-\w&#43;)?\.js$`),
	substring.Regexp(`(^|/)slick\.\w&#43;.js$`),
	substring.Regexp(`(^|/)Leaflet\.Coordinates-\d&#43;\.\d&#43;\.\d&#43;\.src\.js$`),
	substring.Regexp(`leaflet.draw-src.js`),
	substring.Regexp(`leaflet.draw.css`),
	substring.Regexp(`Control.FullScreen.css`),
	substring.Regexp(`Control.FullScreen.js`),
	substring.Regexp(`leaflet.spin.js`),
	substring.Regexp(`wicket-leaflet.js`),
	substring.Regexp(`.sublime-project`),
	substring.Regexp(`.sublime-workspace`),
	substring.Regexp(`(^|/)prototype(.*)\.js$`),
	substring.Regexp(`(^|/)effects\.js$`),
	substring.Regexp(`(^|/)controls\.js$`),
	substring.Regexp(`(^|/)dragdrop\.js$`),
	substring.Regexp(`(.*?)\.d\.ts$`),
	substring.Regexp(`(^|/)mootools([^.]*)\d&#43;\.\d&#43;.\d&#43;([^.]*)\.js$`),
	substring.Regexp(`(^|/)dojo\.js$`),
	substring.Regexp(`(^|/)MochiKit\.js$`),
	substring.Regexp(`(^|/)yahoo-([^.]*)\.js$`),
	substring.Regexp(`(^|/)yui([^.]*)\.js$`),
	substring.Regexp(`(^|/)ckeditor\.js$`),
	substring.Regexp(`(^|/)tiny_mce([^.]*)\.js$`),
	substring.Regexp(`(^|/)tiny_mce/(langs|plugins|themes|utils)`),
	substring.Regexp(`(^|/)ace-builds/`),
	substring.Regexp(`(^|/)fontello(.*?)\.css$`),
	substring.Regexp(`(^|/)MathJax/`),
	substring.Regexp(`(^|/)Chart\.js$`),
	substring.Regexp(`(^|/)[Cc]ode[Mm]irror/(\d&#43;\.\d&#43;/)?(lib|mode|theme|addon|keymap|demo)`),
	substring.Regexp(`(^|/)shBrush([^.]*)\.js$`),
	substring.Regexp(`(^|/)shCore\.js$`),
	substring.Regexp(`(^|/)shLegacy\.js$`),
	substring.Regexp(`(^|/)angular([^.]*)\.js$`),
	substring.Regexp(`(^|\/)d3(\.v\d&#43;)?([^.]*)\.js$`),
	substring.Regexp(`(^|/)react(-[^.]*)?\.js$`),
	substring.Regexp(`(^|/)modernizr\-\d\.\d&#43;(\.\d&#43;)?\.js$`),
	substring.Regexp(`(^|/)modernizr\.custom\.\d&#43;\.js$`),
	substring.Regexp(`(^|/)knockout-(\d&#43;\.){3}(debug\.)?js$`),
	substring.Regexp(`(^|/)docs?/_?(build|themes?|templates?|static)/`),
	substring.Regexp(`(^|/)admin_media/`),
	substring.Regexp(`(^|/)env/`),
	substring.Regexp(`^fabfile\.py$`),
	substring.Regexp(`^waf$`),
	substring.Regexp(`^.osx$`),
	substring.Regexp(`\.xctemplate/`),
	substring.Regexp(`\.imageset/`),
	substring.Regexp(`^Carthage/`),
	substring.Regexp(`^Pods/`),
	substring.Regexp(`(^|/)Sparkle/`),
	substring.Regexp(`Crashlytics.framework/`),
	substring.Regexp(`Fabric.framework/`),
	substring.Regexp(`BuddyBuildSDK.framework/`),
	substring.Regexp(`Realm.framework`),
	substring.Regexp(`RealmSwift.framework`),
	substring.Regexp(`gitattributes$`),
	substring.Regexp(`gitignore$`),
	substring.Regexp(`gitmodules$`),
	substring.Regexp(`(^|/)gradlew$`),
	substring.Regexp(`(^|/)gradlew\.bat$`),
	substring.Regexp(`(^|/)gradle/wrapper/`),
	substring.Regexp(`-vsdoc\.js$`),
	substring.Regexp(`\.intellisense\.js$`),
	substring.Regexp(`(^|/)jquery([^.]*)\.validate(\.unobtrusive)?\.js$`),
	substring.Regexp(`(^|/)jquery([^.]*)\.unobtrusive\-ajax\.js$`),
	substring.Regexp(`(^|/)[Mm]icrosoft([Mm]vc)?([Aa]jax|[Vv]alidation)(\.debug)?\.js$`),
	substring.Regexp(`^[Pp]ackages\/.&#43;\.\d&#43;\/`),
	substring.Regexp(`(^|/)extjs/.*?\.js$`),
	substring.Regexp(`(^|/)extjs/.*?\.xml$`),
	substring.Regexp(`(^|/)extjs/.*?\.txt$`),
	substring.Regexp(`(^|/)extjs/.*?\.html$`),
	substring.Regexp(`(^|/)extjs/.*?\.properties$`),
	substring.Regexp(`(^|/)extjs/.sencha/`),
	substring.Regexp(`(^|/)extjs/docs/`),
	substring.Regexp(`(^|/)extjs/builds/`),
	substring.Regexp(`(^|/)extjs/cmd/`),
	substring.Regexp(`(^|/)extjs/examples/`),
	substring.Regexp(`(^|/)extjs/locale/`),
	substring.Regexp(`(^|/)extjs/packages/`),
	substring.Regexp(`(^|/)extjs/plugins/`),
	substring.Regexp(`(^|/)extjs/resources/`),
	substring.Regexp(`(^|/)extjs/src/`),
	substring.Regexp(`(^|/)extjs/welcome/`),
	substring.Regexp(`(^|/)html5shiv\.js$`),
	substring.Regexp(`^[Tt]ests?/fixtures/`),
	substring.Regexp(`^[Ss]pecs?/fixtures/`),
	substring.Regexp(`(^|/)cordova([^.]*)\.js$`),
	substring.Regexp(`(^|/)cordova\-\d\.\d(\.\d)?\.js$`),
	substring.Regexp(`foundation(\..*)?\.js$`),
	substring.Regexp(`^Vagrantfile$`),
	substring.Regexp(`.[Dd][Ss]_[Ss]tore$`),
	substring.Regexp(`^vignettes/`),
	substring.Regexp(`^inst/extdata/`),
	substring.Regexp(`octicons.css`),
	substring.Regexp(`sprockets-octicons.scss`),
	substring.Regexp(`(^|/)activator$`),
	substring.Regexp(`(^|/)activator\.bat$`),
	substring.Regexp(`proguard.pro`),
	substring.Regexp(`proguard-rules.pro`),
	substring.Regexp(`^puphpet/`),
	substring.Regexp(`(^|/)\.google_apis/`),
	substring.Regexp(`^Jenkinsfile$`),
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

var configurationLanguages = map[string]bool{
	"XML": true, "JSON": true, "TOML": true, "YAML": true, "INI": true, "SQL": true,
}
